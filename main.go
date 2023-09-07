package main

import (
	"context"

	"log/slog"

	"github.com/vod/api"
	conn "github.com/vod/config/database"
	db "github.com/vod/db/sqlc"
	util "github.com/vod/utils"
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
