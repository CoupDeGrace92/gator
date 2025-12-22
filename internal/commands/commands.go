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
	if len(cmd.Args)!=1 {
		err := fmt.Errorf("agg expects one arguments, found %v\n", len(cmd.Args))
		return err
	}
	//The argument passed should be in a form time.ParseDuration recognizes
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		err = fmt.Errorf("Error parsing time: %v", err)
	}
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err := ScrapeFeeds(s, context.Background())
		if err != nil{
			return err
		}
	}
	return nil
}

func HandlerAddFeed(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) != 2{
		err := fmt.Errorf("Error: Add feed expects 2 args, recieved %v\n", len(cmd.Args))
		return err
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
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
		return err
	}
	fmt.Printf("Created feed succesfully:\n	ID: %v\n	CreatedAt: %v\n	UpdatedAt: %v\n	Name: %s\n	Url: %s\n	UserId: %v\n",	addedFeed.ID, addedFeed.CreatedAt, addedFeed.UpdatedAt, addedFeed.Name, addedFeed.Url, addedFeed.UserID)
	cmd.Args[0] = cmd.Args[1]
	cmd.Args = cmd.Args[:1]
	err = HandlerFollow(s, cmd, user)
	if err != nil{
		err = fmt.Errorf("Error following created feed: %v\n", err)
		return err
	}
	return nil
}

func HandlerFeeds(s *config.State, cmd Command) error {
	if len(cmd.Args) != 0{
		err := fmt.Errorf("Error: feeds expects no arguments, recieved %v", len(cmd.Args))
		return err
	}
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		err = fmt.Errorf("Error getting feeds from db: %v", err)
		return err
	}

	if len(feeds) == 0 {
		err = fmt.Errorf("No feeds in database")
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("Feed: %s URL: %s User: %s\n", feed.FeedName, feed.FeedUrl, feed.User)
	}

	return nil
}

func HandlerFollow(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1{
		err := fmt.Errorf("Error: follow expects one argument, recieved %v", len(cmd.Args))
		return err
	}
	url := cmd.Args[0]
	feedID, err:= s.Db.GetFeedIds(context.Background(), url)
	if err != nil {
		err := fmt.Errorf("Error with getting feed name from the db: %v", err)
		return err
	}

	var argParams database.CreateFeedFollowsParams
	argParams.ID = uuid.New()
	now := time.Now()
	argParams.CreatedAt = now
	argParams.UpdatedAt = now
	argParams.UserID = user.ID//this is a uuid.UUID object, not a string but is handled similarly to strings
	argParams.FeedID = feedID //this is a uuid.UUID object...

	feedFollows, err := s.Db.CreateFeedFollows(context.Background(), argParams)
	if err != nil {
		err = fmt.Errorf("Error with creating feed follows in the db: %v", err)
		return err
	}

	fmt.Printf("Feed Name: %s\n Current User: %s\n", feedFollows.FeedName, feedFollows.User)

	return nil
}

func HandlerFollowing(s *config.State, cmd Command) error {
	if len(cmd.Args)!=0{
		err := fmt.Errorf("Error: following expects no args, recieved %v", len(cmd.Args))
		return err
	}
	feedFollowSlice, err := s.Db.GetFeedFollowsForUser(context.Background(), s.CfgPoint.CurrentUser)
	if err != nil{
		err = fmt.Errorf("Error with getting feed follows from db layar: %v\n", err)
		return err
	}
	if len(feedFollowSlice)==0{
		fmt.Println("No feeds followed")
		return nil
	}
	fmt.Println("Following feeds:")
	for _, feed := range feedFollowSlice {
		fmt.Printf("	%s\n", feed.FeedName)
	}
	return nil
}

func HandlerUnfollow(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args)!=1{
		err:= fmt.Errorf("Error: Unfollow expects 1 args, recieved %v\n", len(cmd.Args))
		return err
	}
	// We have GetFeedIds(ctx, url) uuid  AND GetUser(ctx, name) User  where User.ID is the uuid for user id
	//name is in s.CfgPoint.CurrentUser
	userId, err := s.Db.GetUser(context.Background(), s.CfgPoint.CurrentUser)
	if err != nil{
		err = fmt.Errorf("Error getting user: %v", err)
		return err
	}
	feedId, err := s.Db.GetFeedIds(context.Background(), cmd.Args[0])
	if err != nil{
		err = fmt.Errorf("Error getting feedId: %v\n", err)
		return err 
	}
	var unfollowPs database.UnfollowParams
	unfollowPs.UserID = userId.ID
	unfollowPs.FeedID = feedId
	err = s.Db.Unfollow(context.Background(), unfollowPs)
	if err != nil {
		err = fmt.Errorf("Error in unfollowing %v: %v\n", cmd.Args[0], err)
		return err
	}
	fmt.Printf("Succesfully unfollowed feed from %v\n", cmd.Args[0])
	return nil
}