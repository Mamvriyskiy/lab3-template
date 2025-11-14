package main

import (
	"log"

	"github.com/joho/godotenv"

	handler "github.com/Mamvriyskiy/lab3-template/src/bonus/handler"
	repo "github.com/Mamvriyskiy/lab3-template/src/bonus/repository"
	services "github.com/Mamvriyskiy/lab3-template/src/bonus/services"
	server "github.com/Mamvriyskiy/lab3-template/src/server"
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
		DBName:   "privileges",
		SSLMode:  "disable",
	})

	if err != nil {
		log.Fatal("Error connect db:", err.Error())
		return
	}

	repos := repo.NewRepository(db)
	service := services.NewServices(repos)
	handlers := handler.NewHandler(service)

	srv := new(server.Server)
	if err := srv.Run("8050", handlers.InitRouters()); err != nil {
		log.Fatal("Failed to start server: ", err)
		return
	}
}
