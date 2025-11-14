package main

import (
	"log"

	handler "github.com/Mamvriyskiy/lab3-template/src/flight/handler"
	repo "github.com/Mamvriyskiy/lab3-template/src/flight/repository"
	services "github.com/Mamvriyskiy/lab3-template/src/flight/services"
	server "github.com/Mamvriyskiy/lab3-template/src/server"
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
		DBName:   "flights",
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
	if err := srv.Run("8060", handlers.InitRouters()); err != nil {
		log.Fatal("Error occurred while running http server:", err.Error())
		return
	}
}
