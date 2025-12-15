package config

import (
	"os"
	"fmt"
	"encoding/json"
)

type Config struct {
	DbUrl        string `json:"db_url"`
	CurrentUser  string `json: "current_user_name"`
}

func (C Config) SetUser(name string) {
	C.CurrentUser = name
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

