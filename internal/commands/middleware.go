package cmnd

import (
	"fmt"
	"github.com/CoupDeGrace92/gator/internal/database"
	"github.com/CoupDeGrace92/gator/internal/config"
	"context"
)

func MiddlewareLoggedIn(handler func(s *config.State, cmd Command, user database.User) error) func(s *config.State,cmd Command)error{
	return func(s *config.State, cmd Command) error{
		user, err := s.Db.GetUser(context.Background(), s.CfgPoint.CurrentUser)
		if err != nil{
			err = fmt.Errorf("Error getting current user: %v\n", err)
			return err
		}
	return handler(s, cmd, user)
	}
}