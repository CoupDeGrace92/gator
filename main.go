package main

import _ "github.com/lib/pq"

import (
	"github.com/CoupDeGrace92/gator/internal/config"
	"github.com/CoupDeGrace92/gator/internal/commands"
	"os"
	"log"
)

func main() {
	conf := config.Read()
	state := &config.State{CfgPoint: &conf}
	commands := cmnd.Commands{FuncHandlers: make(map[string]func(*config.State, cmnd.Command)error)}
	commands.Register("login", cmnd.HandlerLogin)
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
	command := cmnd.Command{
		Name: commandName,
		Args: args,
	}
	err:=commands.Run(state, command)
	if err != nil {
		log.Fatalf("Error encountered running %s:\n	-%v\nTerminating process now\n",commandName, err)
	}
	return
}