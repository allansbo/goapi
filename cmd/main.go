package main

import (
	"github.com/allansbo/goapi/internal/app/server"
	"github.com/allansbo/goapi/internal/config"
	"github.com/allansbo/goapi/internal/pkg/logs"
)

var cfg *config.EnvConfig

func init() {
	logs.ConfigLog()

	var err error
	cfg, err = config.LoadEnvConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	server.Initialize(cfg)
}
