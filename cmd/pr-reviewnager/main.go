package main

import (
	"context"
	"log"

	"github.com/alphameo/pr-reviewnager/internal/adapters/api"
	s "github.com/alphameo/pr-reviewnager/internal/application/services"
	ds "github.com/alphameo/pr-reviewnager/internal/domain/services"
	"github.com/alphameo/pr-reviewnager/internal/infrastructure/db/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	dsn := "host=localhost user=postgres password='' dbname=pr-reviewnager port=5432 sslmode=disable"
	port := ":8080"

	ctx := context.Background()
	conn, err := postgres.NewConnection(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	queries := postgres.NewQueries(conn)

	teamRepo, err := postgres.NewTeamRepository(queries, conn)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
	userRepo, err := postgres.NewUserRepository(queries)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
	prRepo, err := postgres.NewPullRequestRepository(queries, conn)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}

	prDomainServ, err := ds.NewDefaultPullRequestDomainService(userRepo, prRepo, teamRepo)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}

	teamServ, err := s.NewDefaultTeamService(teamRepo, userRepo)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
	userServ, err := s.NewDefaulUserService(userRepo)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
	prServ, err := s.NewDefaultPullRequestService(prDomainServ, prRepo)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	serverImpl, err := api.NewServer(
		teamServ,
		userServ,
		prServ,
	)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}

	api.RegisterHandlers(e, serverImpl)

	if err := e.Start(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
