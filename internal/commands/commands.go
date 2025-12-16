package cmnd

import(
	"fmt"
	"github.com/CoupDeGrace92/gator/internal/config"
)

type Command struct {
	Name	string
	Args	[]string
}

type Commands struct {
	FuncHandlers  map[string]func(*config.State, Command)error
}

func (C *Commands) Run(S *config.State, Cmd Command)error {
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

func (C *Commands) Register(name string, f func(*config.State, Command)error) {
	C.FuncHandlers[name] = f
	return
}

func HandlerLogin(s *config.State,cmd Command) error {
	if len(cmd.Args) != 1 {
		err := fmt.Errorf("Error: login expects 1 argument, %v arguments found", len(cmd.Args))
		return err
	}
	s.CfgPoint.SetUser(cmd.Args[0])
	if checkConfig := config.Read(); checkConfig.CurrentUser != cmd.Args[0] {
		err:= fmt.Errorf("Error: expected updated username to be %s, was %s",cmd.Args[0],checkConfig.CurrentUser)
		return err
	}
	fmt.Println("Username has been set.")
	return nil
}