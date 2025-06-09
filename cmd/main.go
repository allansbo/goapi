package main

import (
	"fmt"
	"github.com/allansbo/goapi/internal/app/server"
	"github.com/allansbo/goapi/internal/config"
	"github.com/allansbo/goapi/internal/domain/usecase"
	"github.com/allansbo/goapi/internal/pkg/logs"
	"github.com/allansbo/goapi/internal/provider/db"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type locationService struct {
	cfg        *config.EnvConfig
	repository db.Repository
	server     *server.AppServer
	quit       chan os.Signal
}

var service locationService

func init() {
	logs.ConfigLog()

	var err error
	service.cfg, err = config.LoadEnvConfig()
	if err != nil {
		slog.Error("error on loading environment", "error", err.Error())
		panic(err)
	}

	service.repository = db.NewMongoDBRepository(service.cfg)
	if err := service.repository.Ping(); err != nil {
		slog.Error("error on handling mongodb", "error", err.Error())
		panic(err)
	}

	usecase.LoadLocationUseCase(service.repository)
}

//	@title			Location API
//	@version		1.0
//	@description	API to manage locations from vehicles
//	@host			localhost:8080
//	@BasePatch		/api/v1/
func main() {
	service.quit = make(chan os.Signal, 1)
	signal.Notify(service.quit, syscall.SIGTERM, syscall.SIGINT)
	go service.shutdown()

	service.server = server.NewAppServer(service.cfg.AppPort)
	service.server.Start()
}

func (s *locationService) shutdown() {
	<-s.quit
	fmt.Println("\nClosing tasks. Please wait.")
	slog.Info("Shutdown routine, closing tasks.")

	slog.Info("Closing Context")
	if s.repository != nil {
		s.repository.Stop()
	}

	if s.server != nil {
		if err := s.server.FiberApp.Shutdown(); err != nil {
			slog.Error("error on shutting down fiber app", "error", err.Error())
		}
	}

	fmt.Println("Tasks closed")
	slog.Info("Tasks closed")
	time.Sleep(time.Millisecond * 200)
	fmt.Println("System off")
	slog.Info("System off")
	os.Exit(0)
}
