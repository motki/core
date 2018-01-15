// Package app contains functionality related to creating an integrated
// environment with all the necessary dependencies.
//
// The goal with this package is to provide a single, reusable base for
// interacting with the various functions provided by the motki library.
//
// This package imports every other motki library package. As such, it cannot
// be imported from the "library" portion of the project. It is intended to be
// used from an external package.
//
// This package provides two types of environments: client-only, and
// client/server.
//
// Client-only environments, represented by ClientEnv, only contain client-
// side logic and require a remote motki grpc server to provide data.
//
// Client/server environments, represented by Env, contain both the client-
// side services and server-side services, including a grpc server.
package app

import (
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"google.golang.org/grpc/test/bufconn"

	"github.com/motki/motki/db"
	"github.com/motki/motki/eveapi"
	"github.com/motki/motki/evedb"
	"github.com/motki/motki/evemarketer"
	"github.com/motki/motki/log"
	"github.com/motki/motki/model"
	"github.com/motki/motki/proto"
	"github.com/motki/motki/proto/client"
	"github.com/motki/motki/proto/server"
	"github.com/motki/motki/worker"
)

// Config represents the configuration of an Env or ClientEnv.
//
// Note that the Database and EVEAPI Config structs may not be populated
// if the intent is to use the Config for creating a ClientEnv.
type Config struct {
	Logging  log.Config    `toml:"logging"`
	Database db.Config     `toml:"db"`
	EVEAPI   eveapi.Config `toml:"eveapi"`
	Backend  proto.Config  `toml:"backend"`
}

// NewConfigFromTOMLFile loads a TOML configuration from the given path.
func NewConfigFromTOMLFile(tomlPath string) (*Config, error) {
	if !filepath.IsAbs(tomlPath) {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		tomlPath = filepath.Join(cwd, tomlPath)
	}
	c, err := ioutil.ReadFile(tomlPath)
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	_, err = toml.Decode(string(c), conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// A ClientEnv is an environment without any server-side components.
//
// For a ClientEnv to function, it must connect to a remote motki grpc server.
type ClientEnv struct {
	conf *Config

	Logger    log.Logger
	Scheduler *worker.Scheduler
	Client    client.Client

	signals chan os.Signal
}

// NewClientEnv creates a ClientEnv using the given configuration.
// A ClientEnv will not have an associated gRPC server, nor any database,
// or eveapi, etc.
func NewClientEnv(conf *Config) (*ClientEnv, error) {
	if conf.Backend.Kind == proto.BackendLocalGRPC {
		return nil, errors.New("app: cannot create client-only env with local grpc backend")
	}
	logger := log.New(conf.Logging)
	work := worker.New(logger)
	cl, err := client.New(conf.Backend, logger)
	if err != nil {
		return nil, errors.Wrap(err, "app: unable to init grpc client")
	}
	return &ClientEnv{
		conf: conf,

		Logger:    logger,
		Scheduler: work,

		Client: cl,
	}, nil
}

// ShutdownFunc is a simple function intended to be called prior to application
// exit.
type ShutdownFunc func()

// BlockUntilSignalWith will block until it receives the signals signal.
//
// This function attempts to perform a graceful shutdown, shutting
// down all services and doing whatever clean up processes are necessary.
//
// Each pre-exit task exists in the form of a ShutdownFunc. Each ShutdownFunc
// is run concurrently and there is a finite amount of time for them to return
// before the application exits anyway.
func (env *ClientEnv) BlockUntilSignalWith(signals chan os.Signal, fns ...ShutdownFunc) {
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	s := <-signals
	env.Logger.Warnf("app: signal %+v received, shutting down...", s)
	ct := make(chan struct{}, 0)
	wg := &sync.WaitGroup{}
	for _, fn := range fns {
		wg.Add(1)
		go func(fn ShutdownFunc) {
			fn()
			wg.Done()
		}(fn)
	}
	go func() {
		wg.Wait()
		close(ct)
	}()
	t := time.NewTimer(5 * time.Second)
	select {
	case <-t.C:
		env.Logger.Warnf("app: timeout waiting for services to shutdown")
		os.Exit(1)

	case <-ct:
		env.Logger.Debugf("app: graceful shutdown complete; exiting")
		os.Exit(0)
	}
}

// BlockUntilSignal will block until it receives a signal.
//
// This function performs the default shutdown procedure when it receives
// a signal.
//
// See BlockUntilSignalWith for more details.
func (env *ClientEnv) BlockUntilSignal(abort chan os.Signal) {
	env.signals = abort
	env.BlockUntilSignalWith(abort, env.shutdownFuncs()...)
}

// Shutdown begins a graceful shutdown process.
func (env *ClientEnv) Shutdown() {
	env.signals <- os.Interrupt
}

// shutdownFuncs returns a list of functions to be called when the application
// is shutting down.
func (env *ClientEnv) shutdownFuncs() []ShutdownFunc {
	return []ShutdownFunc{
		func() {
			if err := env.Scheduler.Shutdown(); err != nil {
				env.Logger.Warnf("app: error shutting down scheduler: %s", err.Error())
			}
		}}
}

// An Env is a fully integrated environment.
//
// This struct contains all the client and server services provided by the
// motki library. It does not contain any web or mail server related services.
type Env struct {
	*ClientEnv

	DB    *db.ConnPool
	Model *model.Manager

	EveCentral *evemarketer.EveMarketer
	EveDB      *evedb.EveDB
	EveAPI     *eveapi.EveAPI

	// GRPC application server.
	Server server.Server

	// Prevent external packages from constructing the struct themselves.
	//
	// Use NewEnv or NewClientEnv to create a new MOTKI environment.
	unexported interface{}
}

// NewEnv creates an Env using the given configuration.
func NewEnv(conf *Config) (*Env, error) {
	logger := log.New(conf.Logging)
	pool, err := db.New(conf.Database, logger)
	if err != nil {
		return nil, errors.Wrap(err, "app: unable to initialize db connection pool")
	}
	work := worker.New(logger)

	ec := evemarketer.New()
	edb := evedb.New(pool)
	api := eveapi.New(conf.EVEAPI, logger)
	mdl := model.NewManager(pool, edb, api, ec)

	if conf.Backend.Kind == proto.BackendLocalGRPC {
		conf.Backend.LocalGRPC.Listener = bufconn.Listen(1024)
	}
	cl, err := client.New(conf.Backend, logger)
	if err != nil {
		return nil, errors.Wrap(err, "app: unable to init grpc client")
	}
	srv, err := server.New(conf.Backend, mdl, edb, api, logger)
	if err != nil {
		return nil, errors.Wrap(err, "app: unable to init grpc server")
	}

	// Start serving gRPC immediately.
	err = srv.Serve()
	if err != nil {
		return nil, errors.Wrap(err, "app: unable to start grpc server")
	}

	return &Env{
		ClientEnv: &ClientEnv{
			conf: conf,

			Logger:    logger,
			Scheduler: work,

			Client: cl,
		},

		DB:     pool,
		Model:  mdl,
		Server: srv,

		EveCentral: ec,
		EveDB:      edb,
		EveAPI:     api,
	}, nil
}

// BlockUntilSignal will block until it receives a signal.
//
// This function performs the default shutdown procedure when it receives
// a signal.
//
// See BlockUntilSignalWith for more details.
func (env *Env) BlockUntilSignal(signals chan os.Signal) {
	env.signals = signals
	env.BlockUntilSignalWith(signals, append([]ShutdownFunc{
		func() {
			if env.Server == nil {
				return
			}
			if err := env.Server.Shutdown(); err != nil {
				env.Logger.Warnf("app: error shutting down grpc server: %s", err.Error())
			}
		}},
		env.shutdownFuncs()...)...)
}
