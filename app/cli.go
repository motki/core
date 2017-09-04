package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/motki/motkid/cli"
	"github.com/motki/motkid/cli/command"
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
	// Write any os.Signal to this channel and the program attempt to shutdown
	// gracefully.
	//
	// The implication is os.Exit() should not be used; stick solely to exiting
	// by writing to this channel to exit.
	abort chan os.Signal
}

// NewCLIEnv initializes a new CLI environment.
func NewCLIEnv(conf *Config, historyPath string) (*CLIEnv, error) {
	appEnv, err := NewEnv(conf)
	if err != nil {
		return nil, errors.Wrap(err, "app: unable to initialize command line environment")
	}
	if !filepath.IsAbs(historyPath) {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, errors.Wrap(err, "app: unable to initialize command line environment")
		}
		historyPath = filepath.Join(cwd, historyPath)
	}
	srv := cli.NewServer(appEnv.Logger)
	prompter := cli.NewPrompter(srv, appEnv.EveDB, appEnv.Logger)
	srv.SetCommands(
		command.NewEVETypesCommand(prompter),
		command.NewProductCommand(prompter, appEnv.EveDB, appEnv.Model, appEnv.Logger))
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
	env.BlockUntilAbortWith(
		abort,
		func() {
			fmt.Println("Exiting.")
			if err := env.Scheduler.Shutdown(); err != nil {
				env.Logger.Warnf("app: error shutting down scheduler: %s", err.Error())
			}
		},
		func() {
			if f, err := os.Create(env.historyPath); err == nil {
				env.CLI.WriteHistory(f)
				f.Close()
			} else {
				env.Logger.Warnf("unable to write CLI history: %s", err.Error())
			}
			env.CLI.Close()
		})
}
