// Package app contains functionality related to creating an integrated
// environment with all the necessary dependencies.
//
// The goal with this package is to provide a single, reusable base for
// interacting with a given motkid installation.
//
// This package imports every other motkid package. As such, it cannot be
// imported from the "library" portion of the project. It is intended to be
// used from an external package, such as is done in the motkid and motki
// commands.
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

	_ "github.com/motki/motkid/cli"
	_ "github.com/motki/motkid/cli/auth"
	_ "github.com/motki/motkid/cli/command"
	_ "github.com/motki/motkid/cli/text"
	"github.com/motki/motkid/db"
	"github.com/motki/motkid/eveapi"
	"github.com/motki/motkid/evecentral"
	"github.com/motki/motkid/evedb"
	"github.com/motki/motkid/http"
	_ "github.com/motki/motkid/http/auth"
	_ "github.com/motki/motkid/http/middleware"
	_ "github.com/motki/motkid/http/route"
	_ "github.com/motki/motkid/http/session"
	_ "github.com/motki/motkid/http/template"
	"github.com/motki/motkid/log"
	"github.com/motki/motkid/mail"
	"github.com/motki/motkid/model"
	"github.com/motki/motkid/worker"

	_ "github.com/motki/motkid/http/module/account"
	_ "github.com/motki/motkid/http/module/assets"
	_ "github.com/motki/motkid/http/module/auth"
	_ "github.com/motki/motkid/http/module/home"
	_ "github.com/motki/motkid/http/module/industry"
	_ "github.com/motki/motkid/http/module/market"
)

// Config represents a fully configured motkid installation.
type Config struct {
	Logging  log.Config    `toml:"logging"`
	Database db.Config     `toml:"db"`
	HTTP     http.Config   `toml:"http"`
	Mail     mail.Config   `toml:"mail"`
	EVEAPI   eveapi.Config `toml:"eveapi"`
}

// NewConfig loads a TOML configuration from the given path.
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

// An Env is a fully integrated environment.
//
// This struct contains all the core services needed by motkid, but
// does not contain any web or mail server related services.
type Env struct {
	conf *Config

	Logger    log.Logger
	DB        *db.ConnPool
	Scheduler *worker.Scheduler
	Model     *model.Manager

	EveCentral *evecentral.EveCentral
	EveDB      *evedb.EveDB
	EveAPI     *eveapi.EveAPI
}

// NewEnv creates an Env using the given configuration.
func NewEnv(conf *Config) (*Env, error) {
	logger := log.New(conf.Logging)
	pool, err := db.New(conf.Database, logger)
	if err != nil {
		return nil, errors.Wrap(err, "app: unable to initialize db connection pool")
	}
	work := worker.New(logger)

	ec := evecentral.New()
	edb := evedb.New(pool)
	api := eveapi.New(conf.EVEAPI, logger)
	mdl := model.NewManager(pool, edb, api, ec)

	return &Env{
		conf: conf,

		Logger:    logger,
		DB:        pool,
		Scheduler: work,
		Model:     mdl,

		EveCentral: ec,
		EveDB:      edb,
		EveAPI:     api,
	}, nil
}

// abortFunc is a simple function intended to be called prior to application exit.
type abortFunc func()

// BlockUntilAbortWith will block until it receives the abort signal.
//
// This function attempts to perform a graceful shutdown, shutting
// down all services and doing whatever clean up processes are necessary.
//
// Each pre-exit task exists in the form of an abortFunc.
//
// Note that each abortFunc is run concurrently and there is a finite amount
// of time for them to return before the application exits anyway.
func (env *Env) BlockUntilAbortWith(abort chan os.Signal, fns ...abortFunc) {
	signal.Notify(abort, syscall.SIGINT, syscall.SIGTERM)
	select {
	case s := <-abort:
		env.Logger.Warnf("app: signal %+v received, shutting down...", s)
		ct := make(chan struct{}, 0)
		wg := &sync.WaitGroup{}
		for _, fn := range fns {
			wg.Add(1)
			go func(fn abortFunc) {
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
}

// BlockUntilAbort will block until it receives the abort signal.
//
// This function performs the default shutdown procedure when it receives
// an abort signal.
//
// See BlockUntilAbortWith for more details.
func (env *Env) BlockUntilAbort(abort chan os.Signal) {
	env.BlockUntilAbortWith(abort, env.abortFunc())
}

// abortFunc returns a function to be called when the application is
// shutting down.
func (env *Env) abortFunc() abortFunc {
	return func() {
		if err := env.Scheduler.Shutdown(); err != nil {
			env.Logger.Warnf("app: error shutting down scheduler: %s", err.Error())
		}
	}
}
