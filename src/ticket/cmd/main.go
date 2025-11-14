package main

import (
	"log"

	server "github.com/Mamvriyskiy/lab3-template/src/server"
	handler "github.com/Mamvriyskiy/lab3-template/src/ticket/handler"
	repo "github.com/Mamvriyskiy/lab3-template/src/ticket/repository"
	service "github.com/Mamvriyskiy/lab3-template/src/ticket/services"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file found: %v", err)
	}

	db, err := repo.NewPostgresDB(&repo.Config{
		Host:     "postgres",
		Port:     "5432",
		Username: "postgres",
		Password: "postgres",
		DBName:   "tickets",
		SSLMode:  "disable",
	})

	if err != nil {
		return
	}

	repos := repo.NewRepository(db)
	services := service.NewServices(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run("8070", handlers.InitRouters()); err != nil {
		return
	}
}
