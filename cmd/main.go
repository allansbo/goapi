package main

import (
	"github.com/allansbo/goapi/internal/app/server"
	"github.com/allansbo/goapi/internal/config"
	"log"
)

var cfg *config.EnvConfig

func init() {
	var err error
	cfg, err = config.LoadEnvConfig()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	server.RunServer(cfg)
}
