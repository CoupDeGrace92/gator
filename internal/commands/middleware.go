package cmnd

import (
	"fmt"
	"github.com/CoupDeGrace92/gator/internal/database"
	"github.com/CoupDeGrace92/gator/internal/config"
	"context"
	"Errors"
	"database/sql"
)

func middlewareLoggedIn(handler func() error) func(*config.State,cmd Command) error{
	return func(s config.State, cmd Command){
		user, err := s.Db.GetUser(context.Background(), s.CfgPoint.CurrentUser)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows){
				err = fmt.Errorf("User not found")
				return err
			}
			err := fmt.Errorf("Error with GetUser: %v", userErr)
			return err
		}
		return handler(s, cmd, user)
	}
}