package main

import (
	"log"

	handler "github.com/Mamvriyskiy/lab3-template/src/gateway/handler"
	services "github.com/Mamvriyskiy/lab3-template/src/gateway/services"
	server "github.com/Mamvriyskiy/lab3-template/src/server"
	redis "github.com/Mamvriyskiy/lab3-template/src/gateway/rollback"
)

func main() {
	redis.InitRedis()

	services := services.NewServices()
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run("8080", handlers.InitRouters()); err != nil {
		log.Fatal("Failed to start server: ", err)
		return
	}
}
