package main

import (
	"announce-api/db"
	"announce-api/handlers"
	"announce-api/router"
	"announce-api/server"
	"announce-api/services"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	dbConfig := db.Config{
		Driver:   "postgres",
		Host:     os.Getenv("PG_HOST"),
		Port:     os.Getenv("PG_PORT"),
		DBName:   "announcements",
		User:     "postgres",
		Password: os.Getenv("PG_PASSWD"),
		SSLMode:  "disable",
	}

	dbClient, err := db.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Connection to database failed: %s", err.Error())
	}

	minioClient, err := services.NewMinioStorage()
	if err != nil {
		log.Fatalf("Connection to object storage failed: %s", err.Error())
	}

	service := services.Init(services.InitParams{DBClient: dbClient, ObjectStorageClient: minioClient})

	handler := handlers.Init(service)

	srv := server.NewServer(router.Init(handler))

	srv.Run()
}
