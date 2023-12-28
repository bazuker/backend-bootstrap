package manager

import (
	"net/http"

	"github.com/akyoto/cache"
	"github.com/bazuker/backend-bootstrap/pkg/db"
	"github.com/bazuker/backend-bootstrap/pkg/fileStore"
	authHandlers "github.com/bazuker/backend-bootstrap/pkg/manager/auth"
	"github.com/bazuker/backend-bootstrap/pkg/manager/helper"
	usersHandlers "github.com/bazuker/backend-bootstrap/pkg/manager/users"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Manager is a smart HTTP server and router that handles requests routing
// and provides useful context to handlers.
type Manager struct {
	router *gin.Engine
	cfg    Config
}

type Config struct {
	// ServerAddress is server HTTP address
	ServerAddress string
	// ServerMaxUploadFilesizeMB maximum upload filesize in megabytes.
	ServerMaxUploadFilesizeMB int64
	// ServerCORS is cross-origin resource sharing configuration
	ServerCORS *cors.Config
	// DB is a database adapter.
	DB db.Adapter
	// Cache is the session cache.
	Cache *cache.Cache
	// FileStore is a file storage provider.
	FileStore fileStore.FileStore
}

func New(cfg Config) *Manager {
	if cfg.ServerMaxUploadFilesizeMB < 1 {
		cfg.ServerMaxUploadFilesizeMB = 10
	}
	if cfg.ServerCORS == nil {
		corsCfg := cors.DefaultConfig()
		corsCfg.AllowOrigins = []string{"*"}
		corsCfg.AllowHeaders = []string{"*"}
		cfg.ServerCORS = &corsCfg
	}
	return &Manager{router: gin.Default(), cfg: cfg}
}

func (r *Manager) Start() error {
	r.router.MaxMultipartMemory = r.cfg.ServerMaxUploadFilesizeMB << 20

	r.router.Use(cors.New(*r.cfg.ServerCORS))

	api := r.router.Group("/api")
	api.Use(contextMiddleware(r.cfg.DB, r.cfg.Cache, r.cfg.FileStore))

	v1 := api.Group("/v1")

	/* Authentication */
	auth := v1.Group("/auth")
	// Route to initiate Google authentication.
	// e.g. https://example.com/api/v1/auth/google
	auth.Match([]string{http.MethodGet, http.MethodPost}, "/google", authHandlers.HandleAuthGoogleInitiation)
	// Route to handle Google authentication callback.
	auth.Match([]string{http.MethodGet, http.MethodPost}, "/google/callback", authHandlers.HandleAuthGoogleCallback)

	/* Users */
	users := v1.Group("/users")
	// Authentication check middleware will verify that the user is authentication
	// i.e. has 'Access-Token' header with a token that exists in the session cache.
	// and also create 'ContextUserID' for convenience.
	users.Use(authHandlers.CheckAuthenticationMiddleware)
	// Protected route that returns information about the authenticated user.
	// e.g. https://example.com/api/v1/users/me
	users.GET("/me", usersHandlers.HandleUsersMe)
	// Protected route that allows users to upload profile photos.
	users.POST("/me/photo", usersHandlers.HandleUsersMePhoto)
	// Protected route that allows users to delete profile photo.
	users.DELETE("/me/photo", usersHandlers.HandleUsersMeDeletePhoto)
	// Protected route that returns information about a user.
	// Users with basic access can only get information about themselves.
	// Users with admin access can get information about any user.
	users.GET("/:userid", usersHandlers.HandleGetUsers)

	return r.router.Run(r.cfg.ServerAddress)
}

// contextMiddleware sets additional useful context to be used by other handlers.
func contextMiddleware(
	adapter db.Adapter,
	sessionCache *cache.Cache,
	fileStore fileStore.FileStore,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(helper.ContextDatabase, adapter)
		c.Set(helper.ContextCache, sessionCache)
		c.Set(helper.ContextFileStore, fileStore)
		c.Next()
	}
}
