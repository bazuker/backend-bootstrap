# API Server bootstrap
This is a basic API backend scaffold that can be used to quickly bootstrap your next app. 

## HTTP router

https://github.com/gin-gonic/gin

https://github.com/gin-gonic/examples

## Database
The database is implemented via `Adapter` interface.

`Adapter` acts as the controller layer and is used to control the operations on database.
[Amazon AWS DynamoDB](https://aws.amazon.com/dynamodb/) is the only implemented adapter that connects to a real database.
See [dynamodb.go](pkg%2Fdb%2Fdynamodb%2Fdynamodb.go)

Alternatively, for local tests you can use the filesystem adapter. See [local.go](pkg%2Fdb%2Flocal%2Flocal.go)

### DynamoDB setup
Primary index `id`

Secondary index `email` (name `email-index`)

## Sessions and cache
Basic in-memory cache with expiration is implemented.

https://github.com/akyoto/cache

## File storage
The [Amazon AWS S3](https://aws.amazon.com/s3/) is implemented. See [s3.go](pkg%2FfileStore%2Fs3%2Fs3.go)

Alternatively, for local tests you can use the filesystem adapter. See [local.go](pkg%2FfileStore%2Flocal%2Flocal.go)

## Authentication 
Google OAuth 2.0 is conveniently implemented.
Use [Google Console](https://console.cloud.google.com/apis/credentials/oauthclient) to configure OAuth2.0 credentials.