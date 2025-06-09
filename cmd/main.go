package main

import (
	"github.com/allansbo/goapi/internal/app/server"
	"github.com/allansbo/goapi/internal/config"
	"github.com/allansbo/goapi/internal/domain/usecase"
	"github.com/allansbo/goapi/internal/pkg/logs"
	"github.com/allansbo/goapi/internal/provider/db"
	"log/slog"
)

var (
	cfg        *config.EnvConfig
	repository db.Repository
	err        error
)

func init() {
	logs.ConfigLog()

	cfg, err = config.LoadEnvConfig()
	if err != nil {
		slog.Error("error on loading environment", "error", err.Error())
		panic(err)
	}

	repository = db.NewMongoDBRepository(cfg)
	if err := repository.Ping(); err != nil {
		slog.Error("error on handling mongodb", "error", err.Error())
		panic(err)
	}

	usecase.LoadLocationUseCase(repository)
}

func main() {
	server.Initialize(cfg.AppPort)
}
