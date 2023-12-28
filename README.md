# API Server bootstrap
This is a basic API backend scaffold that can be used to quickly bootstrap your next app. 

## HTTP router
https://github.com/gin-gonic/gin

https://github.com/gin-gonic/examples

## Database
The database is implemented via `Adapter` interface.

`Adapter` acts as the controller layer and is used to control the operations on database.
Currently, [Amazon DynamoDB](https://aws.amazon.com/dynamodb/) is the only implemented adapter.

### DynamoDB setup
Primary index `id`

Secondary index `email` (name `email-index`)

## Sessions and cache
Basic in-memory cache with expiration

https://github.com/akyoto/cache

## Authentication 
Google OAuth 2.0 is conveniently implemented.
Use [Google Console](https://console.cloud.google.com/apis/credentials/oauthclient) to configure OAuth2.0 credentials.