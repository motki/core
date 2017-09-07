package app

import (
	"os"

	"github.com/pkg/errors"

	"github.com/motki/motkid/http"
	"github.com/motki/motkid/http/auth"
	"github.com/motki/motkid/http/session"
	"github.com/motki/motkid/http/template"
	"github.com/motki/motkid/mail"

	modaccount "github.com/motki/motkid/http/module/account"
	modassets "github.com/motki/motkid/http/module/assets"
	modauth "github.com/motki/motkid/http/module/auth"
	modhome "github.com/motki/motkid/http/module/home"
	modindustry "github.com/motki/motkid/http/module/industry"
	modmarket "github.com/motki/motkid/http/module/market"
)

// A WebEnv wraps a regular Env, providing web and mail servers.
type WebEnv struct {
	*Env

	Mailer    *mail.Sender
	Sessions  session.Manager
	Templates template.Renderer
	Auth      auth.Manager
	Web       *http.Server

	unexported struct{}
}

// NewWebEnv creates a new web environment using the given configuration.
//
// This function will initialize a regular Env before it initializes the
// web and mail server related functionality.
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
		return nil, errors.Wrap(err, "app: unable to initialize web environment")
	}
	err = srv.Register(
		modassets.New(),
		modauth.New(sessions, authManager, templates, env.Model, env.Scheduler, mailer, env.Logger),
		modhome.New(sessions, templates, mailer, env.Logger),
		modmarket.New(authManager, templates, env.Model, env.EveDB, env.Logger),
		modaccount.New(authManager, templates, env.Model, env.EveDB, env.Logger),
		modindustry.New(authManager, templates, env.Model, env.EveDB, env.Logger),
	)
	if err != nil {
		return nil, errors.Wrap(err, "app: unable to initialize web environment")
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

// BlockUntilAbort will block until it receives the abort signal.
//
// This function performs the default shutdown procedure when it receives
// an abort signal.
//
// See BlockUntilAbortWith for more details.
func (webEnv *WebEnv) BlockUntilAbort(abort chan os.Signal) {
	webEnv.BlockUntilAbortWith(abort, webEnv.Env.abortFunc(), webEnv.abortFunc())
}

// abortFunc returns a function to be called when the application is
// shutting down.
func (webEnv *WebEnv) abortFunc() abortFunc {
	return func() {
		if err := webEnv.Web.Shutdown(); err != nil {
			webEnv.Logger.Warnf("app: error shutting down web server: %s", err.Error())
		}
	}
}
