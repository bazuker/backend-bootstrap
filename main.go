package main

import (
	"log"
	"time"

	"github.com/akyoto/cache"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/bazuker/backend-bootstrap/pkg/db/dynamodb"
	"github.com/bazuker/backend-bootstrap/pkg/fileStore/s3"
	"github.com/bazuker/backend-bootstrap/pkg/router"
)

func main() {
	log.Println("Initializing")

	// Initialize AWS stuff.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := dynamodb.New(dynamodb.Config{
		AWSSession:     sess,
		UsersTableName: "backend-bootstrap-users",
	})
	fs := s3.New(s3.Config{
		AWSSession: sess,
		Bucket:     "backend-bootstrap-storage",
	})

	// Initialize the router.
	r := router.New(router.Config{
		Address:             ":9999",
		MaxUploadFilesizeMB: 16,
		Cache:               cache.New(time.Minute * 5),
		DB:                  db,
		FileStore:           fs,
	})

	log.Println("Running")

	// Start the HTTP server using the router.
	if err := r.Run(); err != nil {
		log.Fatalln(err)
	}
}
