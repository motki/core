package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/motki/motkid/cli"
	"github.com/motki/motkid/cli/auth"
	"github.com/motki/motkid/cli/command"
	"github.com/motki/motkid/cli/text"
)

// A CLIEnv wraps an *Env, providing CLI specific facilities.
type CLIEnv struct {
	*Env
	CLI      *cli.Server
	Prompter *cli.Prompter

	// historyPath is the path to the CLI history file.
	historyPath string

	// abort is used to shutdown the program.
	//
	// Write any os.Signal to this channel and the program will attempt to
	// shutdown gracefully.
	//
	// The implication is os.Exit() should not be used; stick solely to exiting
	// by writing to this channel to exit.
	abort chan os.Signal
}

// CLIConfig wraps a *Config and contains optional credentials.
type CLIConfig struct {
	*Config

	username string
	password string
}

// NewCLIConfig creates a new CLI-specific configuration using the given conf.
func NewCLIConfig(appConf *Config) CLIConfig {
	return CLIConfig{Config: appConf}
}

// WithCredentials returns a copy of the CLIConfig with the given credentials.
func (c CLIConfig) WithCredentials(username, password string) CLIConfig {
	return CLIConfig{c.Config, username, password}
}

// NewCLIEnv initializes a new CLI environment.
//
// If the given CLIConfig contains a username or password, authentication
// will be attempted. If authentication fails, an error is returned.
func NewCLIEnv(conf CLIConfig, historyPath string) (*CLIEnv, error) {
	appEnv, err := NewEnv(conf.Config)
	if err != nil {
		return nil, err
	}
	if !filepath.IsAbs(historyPath) {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		historyPath = filepath.Join(cwd, historyPath)
	}
	srv := cli.NewServer(appEnv.Logger)
	prompter := cli.NewPrompter(srv, appEnv.EveDB, appEnv.Logger)
	sess := auth.NewSession(appEnv.Client, appEnv.Model, appEnv.EveAPI, appEnv.Logger)
	if conf.username != "" || conf.password != "" {
		if _, err := sess.Authenticate(conf.username, conf.password); err != nil {
			return nil, err
		} else {
			fmt.Println("Welcome, " + text.Boldf(conf.username) + "!")
		}
	} else {
		fmt.Printf("Welcome to the %s command line interface.\n", text.Boldf("motki"))
		fmt.Println()
		fmt.Printf("Enter \"%s\" in the prompt for detailed help information.\n", text.Boldf("help"))
	}
	srv.SetCommands(
		command.NewEVETypesCommand(prompter),
		command.NewProductCommand(appEnv.Client, prompter, appEnv.EveAPI, appEnv.Logger))
	if f, err := os.Open(historyPath); err == nil {
		srv.ReadHistory(f)
		f.Close()
	}
	srv.SetCtrlCAborts(true)
	env := &CLIEnv{
		Env:      appEnv,
		CLI:      srv,
		Prompter: prompter,

		historyPath: historyPath,

		abort: make(chan os.Signal, 1),
	}

	return env, nil
}

func (env *CLIEnv) LoopCLI() {
	env.CLI.LoopCLI()
	env.abort <- os.Interrupt
}

// BlockUntilAbort will block until it receives the abort signal.
//
// This function attempts to perform a graceful shutdown, shutting
// down all related services and doing whatever clean up processes are
// necessary.
func (env *CLIEnv) BlockUntilAbort(abort chan os.Signal) {
	if abort == nil {
		abort = env.abort
	}
	abortFuncs := append(
		append([]abortFunc{}, env.Env.abortFuncs()...),
		func() {
			if f, err := os.Create(env.historyPath); err == nil {
				env.CLI.WriteHistory(f)
				f.Close()
			} else {
				env.Logger.Warnf("unable to write CLI history: %s", err.Error())
			}
			env.CLI.Close()
		})
	env.BlockUntilAbortWith(abort, abortFuncs...)
}
