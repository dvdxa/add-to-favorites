package main

import (
	"github.com/dvdxa/add-to-favorites/internal/configs"
	"github.com/dvdxa/add-to-favorites/internal/database/postgres"
	"github.com/dvdxa/add-to-favorites/internal/handlers"
	"github.com/dvdxa/add-to-favorites/internal/repositories"
	"github.com/dvdxa/add-to-favorites/internal/services"
	"github.com/dvdxa/add-to-favorites/pkg/logger"
	"github.com/dvdxa/add-to-favorites/server"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load("env")
	if err != nil {
		log.Fatalf("failed to load env files: %v", err)
	}
	log := logger.GetLogger()
	cfg, err := configs.InitConfig()
	if err != nil {
		log.Fatalf("failed to initialize configs: %v", err)
	}
	pgx, err := postgres.ConnectToPostgres(&cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	repoPort := repositories.NewRepositoryPort(pgx)
	servicePort := services.NewServicePort(repoPort)
	handler := handlers.NewHandler(*log, *servicePort)
	srv := new(server.Server)
	err = srv.Run("0006", handler.InitRoutes())
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
