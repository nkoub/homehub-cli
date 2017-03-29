package cli

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/jamesnetherton/homehub-cli/cmd"
	"github.com/jamesnetherton/homehub-cli/service"
)

// CLI is a representation of the home hub command line interface
type CLI struct {
	readline *readline.Instance
	commands []cmd.Command
}

// NewCLI creates a new CLI
func NewCLI(commands []cmd.Command, readLine *readline.Instance) (cli *CLI) {
	c := &CLI{
		readLine,
		commands,
	}

	c.readline.Config.AutoComplete = readline.NewPrefixCompleter(readline.PcItemDynamic(c.commandListFunction()))

	return c
}

// Run starts the CLI and processes user input
func (c *CLI) Run() {
	defer c.readline.Close()

	banner()

	for {
		line, err := c.readline.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		if !service.StringIsEmpty(line) {
			commandLine := strings.Split(line, " ")
			commandName := commandLine[0]
			args := append(commandLine[:0], commandLine[1:]...)

			command := c.findCommand(commandName)
			if command != nil {
				command.ExecuteLifecylce(args)
			} else {
				fmt.Println("homehub: Command not found:", strconv.Quote(commandName))
			}
		}
	}
}

func (c *CLI) findCommand(commandName string) (command cmd.Command) {
	for i := 0; i < len(c.commands); i++ {
		if c.commands[i].GetName() == commandName {
			return c.commands[i]
		}
	}
	return nil
}

func (c *CLI) commandListFunction() func(string) []string {
	return func(s string) []string {
		var commands []string
		for i := 0; i < len(c.commands); i++ {
			commands = append(commands, c.commands[i].GetName())
		}
		return commands
	}
}

func banner() {
	fmt.Println(" _   _                           _   _         _")
	fmt.Println("| | | |                         | | | |       | |")
	fmt.Println("| |_| |  ___   _ __ ___    ___  | |_| | _   _ | |__")
	fmt.Println("|  _  | / _ \\ | '_ ` _ \\  / _ \\ |  _  || | | || '_ \\")
	fmt.Println("| | | || (_) || | | | | ||  __/ | | | || |_| || |_) |")
	fmt.Println("\\_| |_/ \\___/ |_| |_| |_| \\___| \\_| |_/ \\__,_||_.__/")
	fmt.Printf("\n=====================================================\n\n")
}