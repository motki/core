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
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"

	_ "github.com/motki/motkid/http/middleware"
	_ "github.com/motki/motkid/http/route"

	"github.com/motki/motkid/db"
	"github.com/motki/motkid/eveapi"
	"github.com/motki/motkid/evecentral"
	"github.com/motki/motkid/evedb"
	"github.com/motki/motkid/http"
	"github.com/motki/motkid/http/auth"
	"github.com/motki/motkid/http/session"
	"github.com/motki/motkid/http/template"
	"github.com/motki/motkid/log"
	"github.com/motki/motkid/mail"
	"github.com/motki/motkid/model"
	"github.com/motki/motkid/worker"

	"os/signal"
	"syscall"
	"time"

	modaccount "github.com/motki/motkid/http/module/account"
	modauth "github.com/motki/motkid/http/module/auth"
	modhome "github.com/motki/motkid/http/module/home"
	modindustry "github.com/motki/motkid/http/module/industry"
	modmarket "github.com/motki/motkid/http/module/market"
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

		Model: mdl,

		EveCentral: ec,
		EveDB:      edb,
		EveAPI:     api,
	}, nil
}

type AbortFunc func()

func (env *Env) BlockUntilAbortWith(abort chan os.Signal, fns ...AbortFunc) {
	signal.Notify(abort, syscall.SIGINT, syscall.SIGTERM)
	select {
	case s := <-abort:
		env.Logger.Warnf("app: signal %+v received, shutting down...", s)
		ct := make(chan struct{}, 0)
		wg := &sync.WaitGroup{}
		for _, fn := range fns {
			wg.Add(1)
			go func(fn AbortFunc) {
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

func (env *Env) BlockUntilAbort(abort chan os.Signal) {
	env.BlockUntilAbortWith(abort, env.abortFunc())
}

func (env *Env) abortFunc() AbortFunc {
	return func() {
		if err := env.Scheduler.Shutdown(); err != nil {
			env.Logger.Warnf("app: error shutting down scheduler: %s", err.Error())
		}
	}
}

type WebEnv struct {
	*Env

	Mailer    *mail.Sender
	Sessions  session.Manager
	Templates template.Renderer
	Auth      auth.Manager
	Web       *http.Server
}

func NewWebEnv(conf *Config) (*WebEnv, error) {
	env, err := NewEnv(conf)
	if err != nil {
		return nil, err
	}
	mailer := mail.NewSender(conf.Mail, env.Logger)
	mailer.DoNotSend, err = mail.NewModelList(env.Model, "unsubscribe")
	if err != nil {
		return nil, errors.Wrap(err, "app: unable to init 'unsubscribe' list")
	}
	sessions := session.NewManager(conf.HTTP.Session, env.Logger)
	templates := template.NewRenderer(conf.HTTP.Templating, env.Logger)
	authManager := auth.NewManager(
		auth.NewFormLoginAuthenticator(env.Model, env.Logger, "/login/begin"),
		auth.NewEveAPIAuthorizer(env.Model, env.EveAPI, env.Logger),
		sessions,
	)
	srv, err := http.New(conf.HTTP, env.Logger)
	if err != nil {
		return nil, errors.Wrap(err, "app: unable to initialize web server")
	}
	err = srv.Register(
		modauth.New(sessions, authManager, templates, env.Model, env.Scheduler, mailer, env.Logger),
		modhome.New(sessions, templates, mailer, env.Logger),
		modmarket.New(authManager, templates, env.Model, env.EveDB, env.Logger),
		modaccount.New(authManager, templates, env.Model, env.EveDB, env.Logger),
		modindustry.New(authManager, templates, env.Model, env.EveDB, env.Logger),
	)
	if err != nil {
		return nil, errors.Wrap(err, "app: error registering http modules")
	}

	return &WebEnv{
		Env: env,

		Mailer:    mailer,
		Sessions:  sessions,
		Templates: templates,
		Auth:      authManager,
		Web:       srv,
	}, nil
}

func (webEnv *WebEnv) BlockUntilAbort(abort chan os.Signal) {
	webEnv.BlockUntilAbortWith(abort, webEnv.Env.abortFunc(), webEnv.abortFunc())
}

func (webEnv *WebEnv) abortFunc() AbortFunc {
	return func() {
		if err := webEnv.Web.Shutdown(); err != nil {
			webEnv.Logger.Warnf("app: error shutting down web server: %s", err.Error())
		}
	}
}
