package main

import _ "github.com/lib/pq"

import (
	"github.com/CoupDeGrace92/gator/internal/config"
	"github.com/CoupDeGrace92/gator/internal/commands"
	"os"
	"github.com/CoupDeGrace92/gator/internal/database"
	"database/sql"
	"fmt"
)

func main() {
	conf := config.Read()

	db, err := sql.Open("postgres", conf.DbUrl)
	if err != nil {
		fmt.Printf("Error openning connection to %v: \n%v\n", conf.DbUrl, err)
		os.Exit(1)
	}

	dbQueries := database.New(db)
	
	state := &config.State{
		CfgPoint: &conf,
		Db:       dbQueries,
	}

	commands := cmnd.Commands{FuncHandlers: make(map[string]func(*config.State, cmnd.Command)error)}
	commands.Register("login", cmnd.HandlerLogin)
	commands.Register("register", cmnd.HandlerRegister)
	commands.Register("reset", cmnd.HandlerReset)
	cliCommands := os.Args
	if len(cliCommands)<2{
		fmt.Printf("Error: expecting at least one command line argument in addition to program name\n")
		os.Exit(1)
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
	err=commands.Run(state, command)
	if err != nil {
		fmt.Printf("Error encountered running %s:\n	-%v\nTerminating process now\n",commandName, err)
		os.Exit(1)
	}
	os.Exit(0)
}