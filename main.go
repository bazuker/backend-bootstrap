package main

import (
	"log"
	"time"

	"github.com/akyoto/cache"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/bazuker/backend-bootstrap/pkg/db/dynamodb"
	"github.com/bazuker/backend-bootstrap/pkg/filestore/s3"
	"github.com/bazuker/backend-bootstrap/pkg/manager"
)

func main() {
	log.Println("Initializing")

	/*
		Ready for the cloud? Get your environment variables ready!
		export AWS_REGION=""
		export AWS_ACCESS_KEY_ID=""
		export AWS_SECRET_ACCESS_KEY=""
		export GOOGLE_OAUTH_CLIENT_ID=""
		export GOOGLE_OAUTH_CLIENT_SECRET=""
		export GOOGLE_OAUTH_REDIRECT_URL=""
	*/

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

	// Initialize the manager.
	m := manager.New(manager.Config{
		ServerAddress:             ":9999",
		ServerMaxUploadFilesizeMB: 16,
		Cache:                     cache.New(time.Minute * 5),
		DB:                        db,
		FileStore:                 fs,
	})

	/*
		// Not ready for the cloud yet or just want to test locally?
		// No problemo!
		// import localFS "github.com/bazuker/backend-bootstrap/pkg/filestore/local"
		// import localDB "github.com/bazuker/backend-bootstrap/pkg/db/local"
		// Initialize the local database and file storage.
		db, err := localDB.New(localDB.Config{
			// The directory must already exist. Only the file will be created.
			Filename: "localdata/database.json",
		})
		if err != nil {
			log.Println("failed to create local database:", err)
			return
		}
		fs := localFS.New(localFS.Config{
			// The directory must already exist.
			Directory: "localdata/",
		})
		// Initialize the manager.
		m := manager.New(manager.Config{
			ServerAddress:             ":9999",
			ServerMaxUploadFilesizeMB: 16,
			Cache:                     cache.New(time.Minute * 5),
			DB:                        db,
			FileStore:                 fs,
		})
	*/

	log.Println("Running")

	// Start the HTTP server using the router.
	if err := m.Start(); err != nil {
		log.Fatalln(err)
	}
}
