package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/motki/motkid/cli/text"
	"github.com/motki/motkid/log"
	"github.com/peterh/liner"
)

// A Command is a console application.
type Command interface {
	// Description returns a ~40 character sentence describing the Command.
	Description() string

	// Prefixes defines the names that the Command will be invoked under.
	//
	// This is a slice to allow for alternatives and shorthands to be considered
	// as a prefix for the Command.
	Prefixes() []string

	// Handle executes the given subcmd, if any, with the given arguments.
	Handle(subcmd string, args ...string)

	// PrintHelp prints a helpful overview describing the Command its subcommands.
	PrintHelp()
}

// A Server handles command-line requests.
type Server struct {
	*liner.State
	logger log.Logger

	abort chan struct{}

	// commands contains all the commands registered with the p.
	commands []Command

	// commandLookup acts as a lookup table, pairing a Command with each prefix it specifies.
	commandLookup map[string]Command
}

// NewServer initializes a new CLI server.
func NewServer(logger log.Logger) *Server {
	return &Server{
		State:  liner.NewLiner(),
		logger: logger,

		abort: make(chan struct{}, 1),

		commands:      []Command{},
		commandLookup: map[string]Command{},
	}
}

func (srv *Server) SetCommands(commands ...Command) {
	srv.commands = append(commands, quitCommand{srv}, helpCommand{srv})
	cmdNames := []string{}
	for _, cmd := range srv.commands {
		for _, prefix := range cmd.Prefixes() {
			cmdNames = append(cmdNames, prefix)
			srv.commandLookup[prefix] = cmd
		}
	}
	srv.SetCompleter(func(line string) []string {
		res := []string{}
		for _, v := range cmdNames {
			if strings.HasPrefix(v, line) {
				res = append(res, v)
			}
		}
		return res
	})
}

// LoopCLI starts an endless loop to perform commands read from stdin.
//
// This function is intended to be started in a goroutine.
func (srv *Server) LoopCLI() {
	for {
		err := func() error {
			cmd, err := srv.Prompt("> ")
			if err != nil {
				if err == liner.ErrPromptAborted {
					return err
				}
				if err == io.EOF {
					err = errors.New("unexpected EOF")
					fmt.Println()
				}
				srv.logger.Debugf("unable to read command: %s", err.Error())
				return nil
			}
			srv.AppendHistory(cmd)
			parts := strings.Split(cmd, " ")
			if len(parts) < 1 {
				srv.PrintHelp()
				return nil
			}
			if cmd, ok := srv.commandLookup[parts[0]]; ok {
				var subcmd string
				var args []string
				if len(parts) > 1 {
					subcmd = parts[1]
					args = parts[2:]
				}
				if subcmd == "help" {
					cmd.PrintHelp()
					return nil
				}
				cmd.Handle(subcmd, args...)
			} else {
				fmt.Println("Unknown Command:", parts[0])
				srv.PrintHelp()
			}
			return nil
		}()
		if err == liner.ErrPromptAborted {
			srv.abort <- struct{}{}
		}
		select {
		case <-srv.abort:
			return
		default:
			// no op
		}
	}
}

// PrintHelp prints the application-level help text.
func (srv *Server) PrintHelp() {
	fmt.Println(`motki is a command-line utility for interacting with a motkid installation.

Commands:`)
	for _, cmd := range srv.commands {
		for _, prefix := range cmd.Prefixes() {
			fmt.Printf("  %s %s\n", text.PadTextRight(prefix, 15), cmd.Description())
			break
		}
	}
	fmt.Println()
	fmt.Println(`More information about a particular Command can be shown by running`)
	fmt.Println(`  help <Command>`)
	fmt.Println()
}

// quitCommand handles exiting the application on Command.
type quitCommand struct {
	env *Server
}

func (c quitCommand) Prefixes() []string {
	return []string{"quit", "exit", "\\q", "q"}
}

func (c quitCommand) Handle(subcmd string, args ...string) {
	c.env.abort <- struct{}{}
}

func (c quitCommand) Description() string {
	return "Quits the application."
}

func (c quitCommand) PrintHelp() {
	fmt.Printf(`Command "quit" exits the application.

Aliases for quit:
	quit
	q
	exit
	\q

%s`, text.WrapText(`Additionally, the program can be exited by sending a SIGINT or SIGKILL signal, for example by pressing CTRL+C.`, text.StandardTerminalWidthInChars))
}

// helpCommand handles printing help information for all registered commands.
type helpCommand struct {
	env *Server
}

func (c helpCommand) Prefixes() []string {
	return []string{"help"}
}

func (c helpCommand) Handle(subcmd string, args ...string) {
	if len(subcmd) > 0 {
		if cmd, ok := c.env.commandLookup[subcmd]; ok {
			fmt.Println()
			cmd.PrintHelp()
			fmt.Println()
			return
		}
	}
	c.PrintHelp()
}

func (c helpCommand) Description() string {
	return "Displays this help text."
}

func (c helpCommand) PrintHelp() {
	c.env.PrintHelp()
}
