package main

import (
	"context"

	"log/slog"

	"github.com/tiktok/api"
	conn "github.com/tiktok/config/database"
	db "github.com/tiktok/db/sqlc"
	util "github.com/tiktok/utils"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		slog.Info("cannot load config", err)
	}

	dbConn, err := conn.NewPostgres(context.Background(), config.DB_URI)
	if err != nil {
		slog.Info("cannot connect to db", err)
	}
	store := db.New(dbConn.DB)

	server, err := api.NewServer(store, config)
	if err != nil {
		slog.Info("cannot create server")
	}

	err = server.Start(":8080")
	if err != nil {
		slog.Info("=======", err)
	}

}
