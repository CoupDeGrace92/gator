package config

import (
	"os"
	"fmt"
	"encoding/json"
	"github.com/CoupDeGrace92/gator/internal/database"
)

type Config struct {
	DbUrl        string `json:"db_url"`
	CurrentUser  string `json: "current_user_name"`
}

type State struct {
	CfgPoint     *Config
	Db			 *database.Queries
}

type Command struct {
	Name	string
	Args	[]string
}

type Commands struct {
	FuncHandlers  map[string]func(*State, Command)error
}

func (C *Commands) Run(S *State, Cmd Command)error {
	f, ok := C.FuncHandlers[Cmd.Name]
	if !ok {
		err := fmt.Errorf("Command %s not found", Cmd.Name)
		return err
	}
	err := f(S, Cmd)
	if err != nil {
		err = fmt.Errorf("Error running %s: %v",Cmd.Name, err)
		return err
	}
	return nil
}

func (C *Commands) Register(name string, f func(*State, Command)error) {
	C.FuncHandlers[name] = f
	return
}

func (C Config) SetUser(Name string) {
	C.CurrentUser = Name
	data, err := json.Marshal(C)
	if err != nil{
		fmt.Printf("Current user name changed, ERROR marshalling the resulting config: %v", err)
		return
	}
	homeDirString, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v", err)
		return
	}
	fullString := homeDirString+"/.gatorconfig.json"
	err = os.WriteFile(fullString, data, 0600)//Permission chosen as rw_______, can also be written in octal as 0o600
	if err != nil {
		fmt.Printf("Error writing file: %v.\nPlease check %s to make sure incomplete write did not take place", err, fullString)
	}
	return
}

func Read() Config {
	homeDirString, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v", err)
		return Config{}
	}
	fullString := homeDirString+"/.gatorconfig.json"
	jsonBytes, err := os.ReadFile(fullString)
	if err != nil {
		fmt.Printf("Error reading file %s: %v",jsonBytes, err)
		return Config{}
	}
	var storageObject Config
	if err := json.Unmarshal(jsonBytes, &storageObject); err != nil{
		fmt.Printf("Error unmarshalling json byte slice")
		return Config{}
	}
	return storageObject
}

func HandlerLogin(s *State,cmd Command) error {
	if len(cmd.Args) != 1 {
		err := fmt.Errorf("Error: login expects 1 argument, %v arguments found", len(cmd.Args))
		return err
	}
	s.CfgPoint.SetUser(cmd.Args[0])
	if checkConfig := Read(); checkConfig.CurrentUser != cmd.Args[0] {
		err:= fmt.Errorf("Error: expected updated username to be %s, was %s",cmd.Args[0],checkConfig.CurrentUser)
		return err
	}
	fmt.Println("Username has been set.")
	return nil
}