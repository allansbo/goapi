package main

import (
	"github.com/allansbo/goapi/internal/app/server"
	"github.com/allansbo/goapi/internal/domain/usecase"
	"github.com/allansbo/goapi/internal/pkg/logs"
	"log/slog"
)

func init() {
	logs.ConfigLog()

	if err := usecase.LoadAppConfig(); err != nil {
		slog.Error("error on loading environment", "error", err.Error())
		panic(err)
	}

	if err := usecase.LoadDatabaseRepository(); err != nil {
		slog.Error("error on handling mongodb", "error", err.Error())
		panic(err)
	}
}

func main() {
	server.Initialize()
}
