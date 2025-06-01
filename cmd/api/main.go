package main

import (
	"go-manage-hex/cmd/config"
	"go-manage-hex/internal/infrastructure/http/server"
	"log"
)

func main() {

	config.LoadEnv()

	router := server.StartServer()

	if routerErr := router.Run(config.Port); routerErr != nil {
		log.Fatal(routerErr)
	}
}
