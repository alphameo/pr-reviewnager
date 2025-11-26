package main

import (
	"context"
	"log"
	"os"

	"github.com/alphameo/pr-reviewnager/internal/adapters/api"
	"github.com/alphameo/pr-reviewnager/internal/cfg"
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
	repoContainer, err := cfg.NewPSQLRepositoryContainer(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to initialize repositories: %v", err)
	}
	defer func() {
		if err := repoContainer.Close(ctx); err != nil {
			log.Fatalf("Error closing storage: %v", err)
		}
	}()

	serviceProvider, err := cfg.NewServiceContainer(repoContainer)
	if err != nil {
		log.Fatalf("Failed to create service provider: %v", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	serverImpl, err := api.NewServer(
		serviceProvider.TeamService,
		serviceProvider.UserService,
		serviceProvider.PullRequestService,
	)
	if err != nil {
		log.Fatal("Failed to create server:", err)
	}

	api.RegisterHandlers(e, serverImpl)

	if err := e.Start(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
