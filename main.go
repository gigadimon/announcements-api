package main

import (
	"announce-api/router"
	"announce-api/server"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	srv := server.NewServer(router.Init())

	srv.Run()
}
