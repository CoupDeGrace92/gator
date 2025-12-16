package main

import (
	"github.com/CoupDeGrace92/gator/internal/config"
	"os"
	"log"
)

func main() {
	conf := config.Read()
	state := &config.State{CfgPoint: &conf}
	commands := config.Commands{FuncHandlers: make(map[string]func(*config.State, config.Command)error)}
	commands.Register("login", config.HandlerLogin)
	cliCommands := os.Args
	if len(cliCommands)<2{
		log.Fatal("Error: expecting at least one command line argument in addition to program name\n")
	}
	commandName := cliCommands[1]
	var args []string
	if len(cliCommands) > 2{
		args = cliCommands[2:]
	} else {
		args = []string{}
	}
	command := config.Command{
		Name: commandName,
		Args: args,
	}
	err:=commands.Run(state, command)
	if err != nil {
		log.Fatalf("Error encountered running %s:\n	-%v\nTerminating process now\n",commandName, err)
	}
	return
}