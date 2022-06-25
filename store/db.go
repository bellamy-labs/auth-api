package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/bellamy-labs/auth-api/config"
	"github.com/bellamy-labs/auth-api/ent"
)

func InitDB() error {
	client, err := ent.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True", config.Cfg.DBUser, config.Cfg.DBPassword, config.Cfg.DBHost, config.Cfg.DBPort, config.Cfg.DBName))
	if err != nil {
		return errors.New("failed opening connection to sqlite: " + err.Error())
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		return errors.New("failed creating schema resources: " + err.Error())
	}

	return nil
}
