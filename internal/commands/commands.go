package cmnd

import(
	"fmt"
	"github.com/CoupDeGrace92/gator/internal/config"
	"github.com/CoupDeGrace92/gator/internal/database"
	"time"
	"github.com/google/uuid"
	"context"
	"errors"
	"database/sql"
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
	_, userErr := s.Db.GetUser(context.Background(), cmd.Args[0])
	if userErr != nil {
		if errors.Is(userErr, sql.ErrNoRows){
			err := fmt.Errorf("User not found")
			return err
		}
		err := fmt.Errorf("Error with GetUser: %v", userErr)
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

func HandlerRegister(s *config.State, cmd Command) error {
	if len(cmd.Args) !=1 {
		err := fmt.Errorf("Error: register expects a single name argument, %v arguments found", len(cmd.Args))
		return err
	}
	//Check to make sure the name does not exist in the database:
	_, userErr := s.Db.GetUser(context.Background(), cmd.Args[0])
	if userErr == nil{
		err := fmt.Errorf("Another user found with name %s", cmd.Args[0])
		return err
	} else if !errors.Is(userErr, sql.ErrNoRows){
		e := fmt.Errorf("Error: Non-ErrNoRows error in get users: \n%v", userErr)
		return e
	}

	//HERE WE CREATE ARE ARG PARAMS STRUCT:
	var argParams database.CreateUserParams
	argParams.ID = uuid.New()
	now := time.Now()
	argParams.CreatedAt = now
	argParams.UpdatedAt = now
	argParams.Name = cmd.Args[0]
	addedUser, err := s.Db.CreateUser(context.Background(), argParams)
	if err != nil{
		e := fmt.Errorf("Error in creating user: %v",err)
		return e
	}
	s.CfgPoint.SetUser(cmd.Args[0]) 
	fmt.Printf("User was created:\n	ID: %v\n	CreatedAt: %v\n 	UpdatedAt: %v\n 	Name: %v\n",addedUser.ID, addedUser.CreatedAt, addedUser.UpdatedAt, addedUser.Name)
	return nil
}

func HandlerReset(s *config.State,cmd Command) error {
	err := s.Db.DeleteUsers(context.Background())
	if err != nil{
		err = fmt.Errorf("Error in dropping the users table: %v\n", err)
		return err
	}
	fmt.Println("Users table reset")
	return nil
}