package main

import (
	"github.com/akyoto/cache"
	"github.com/bazuker/backend-bootstrap/pkg/db/dynamodb"
	"log"
	"time"

	"github.com/bazuker/backend-bootstrap/pkg/router"
)

func main() {
	log.Println("Initializing")

	// Initialize the router.
	r := router.New(router.Config{
		Address:             ":9999",
		MaxUploadFilesizeMB: 16,
		Cache:               cache.New(time.Minute * 5),
		DB: dynamodb.New(dynamodb.Config{
			UsersTableName: "backend-bootstrap-users",
		}),
	})

	log.Println("Running")
	// Start the HTTP server using the router.
	if err := r.Run(); err != nil {
		log.Fatalln(err)
	}
}
