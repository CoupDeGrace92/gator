package cmnd

import(
	"fmt"
	"github.com/CoupDeGrace92/gator/internal/config"
	"github.com/CoupDeGrace92/gator/internal/database"
	"github.com/CoupDeGrace92/gator/internal/web"
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

func HandlerGetUsers(s *config.State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		err = fmt.Errorf("Error getting users: %v", err)
		return err
	}
	for _, user := range users {
		if user == s.CfgPoint.CurrentUser {
			fmt.Printf("* %s (current)\n", user)
			continue
		}
		fmt.Printf("* %s\n", user)
	}
	return nil
}

func HandlerAggregate(s *config.State, cmd Command) error {
	if len(cmd.Args)!=0 {
		err := fmt.Errorf("agg expects no arguments, found %v\n", len(cmd.Args))
		return err
	}
	feed, err:= web.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil{
		err = fmt.Errorf("Error fetching feed: %v \n", err)
		return err
	}
	fmt.Println(*feed)
	return nil
}

func HandlerAddFeed(s *config.State, cmd Command) error {
	if len(cmd.Args) != 2{
		err := fmt.Errorf("Error: Add feed expects 2 args, recieved %v\n", len(cmd.Args))
		return err
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	user, err := s.Db.GetUser(context.Background(), s.CfgPoint.CurrentUser)
	if err != nil{
		err = fmt.Errorf("Error getting current user: %v\n", err)
		return err
	}
	user_id := user.ID
	
	//Now we need to create the argParams:
	var argParams database.CreateFeedParams
	argParams.ID = uuid.New()
	now := time.Now()
	argParams.CreatedAt = now
	argParams.UpdatedAt = now
	argParams.Name = name
	argParams.Url = url
	argParams.UserID = user_id

	addedFeed, err := s.Db.CreateFeed(context.Background(), argParams)
	if err != nil{
		err = fmt.Errorf("Error creating feed: %v\n", err)
	}
	fmt.Printf("Created feed succesfully:\n	ID: %v\n	CreatedAt: %v\n	UpdatedAt: %v\n	Name: %s\n	Url: %s\n	UserId: %v\n",	addedFeed.ID, addedFeed.CreatedAt, addedFeed.UpdatedAt, addedFeed.Name, addedFeed.Url, addedFeed.UserID)
	return nil
}