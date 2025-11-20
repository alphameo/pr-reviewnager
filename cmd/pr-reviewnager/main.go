package main

import (
	"context"
	"log"
	"os"

	"github.com/alphameo/pr-reviewnager/internal/adapters/api"
	s "github.com/alphameo/pr-reviewnager/internal/application/services"
	ds "github.com/alphameo/pr-reviewnager/internal/domain/services"
	"github.com/alphameo/pr-reviewnager/internal/infrastructure/db/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else {
		port = ":" + port
	}

	ctx := context.Background()
	storage, err := postgres.NewPSQLStorage(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer func() {
		if err := storage.Close(ctx); err != nil {
			log.Fatalf("Error closing storage: %v", err)
		}
	}()

	domainServiceProvider, err := ds.NewDefaultServiceProvider(storage)
	if err != nil {
		log.Fatalf("Failed to create domain service provider: %v", err)
	}

	serviceProvider, err := s.NewDefaultServiceProvider(storage, domainServiceProvider)
	if err != nil {
		log.Fatalf("Failed to create service provider: %v", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	serverImpl, err := api.NewServer(
		serviceProvider.TeamService(),
		serviceProvider.UserService(),
		serviceProvider.PullRequestService(),
	)
	if err != nil {
		log.Fatal("Failed to create server:", err)
	}

	api.RegisterHandlers(e, serverImpl)

	if err := e.Start(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
