package router

import (
	"net/http"

	"github.com/akyoto/cache"
	"github.com/bazuker/backend-bootstrap/pkg/db"
	authHandlers "github.com/bazuker/backend-bootstrap/pkg/router/auth"
	"github.com/bazuker/backend-bootstrap/pkg/router/helper"
	usersHandlers "github.com/bazuker/backend-bootstrap/pkg/router/users"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Router is a smart HTTP server router that handles requests routing
// and provides useful context to handlers.
type Router struct {
	gin *gin.Engine
	cfg Config
}

type Config struct {
	// HTTP server address
	Address string
	// MaxUploadFilesizeMB maximum upload filesize in megabytes.
	MaxUploadFilesizeMB int64
	// CORS is cross-origin resource sharing configuration
	CORS *cors.Config
	// DB is a database adapter.
	DB db.Adapter
	// Cache is the session cache.
	Cache *cache.Cache
}

func New(cfg Config) *Router {
	if cfg.MaxUploadFilesizeMB < 1 {
		cfg.MaxUploadFilesizeMB = 16
	}
	if cfg.CORS == nil {
		corsCfg := cors.DefaultConfig()
		corsCfg.AllowOrigins = []string{"*"}
		corsCfg.AllowHeaders = []string{"*"}
		cfg.CORS = &corsCfg
	}
	return &Router{cfg: cfg}
}

func (r *Router) Run() error {
	r.gin = gin.Default()
	r.gin.MaxMultipartMemory = r.cfg.MaxUploadFilesizeMB << 20

	r.gin.Use(cors.New(*r.cfg.CORS))

	api := r.gin.Group("/api")
	api.Use(contextMiddleware(r.cfg.DB, r.cfg.Cache))

	v1 := api.Group("/v1")

	auth := v1.Group("/auth")
	// Route to initiate Google authentication.
	// e.g. https://example.com/api/v1/auth/google
	auth.Match([]string{http.MethodGet, http.MethodPost}, "/google", authHandlers.HandleAuthGoogleInitiation)
	// Route to handle Google authentication callback.
	auth.Match([]string{http.MethodGet, http.MethodPost}, "/google/callback", authHandlers.HandleAuthGoogleCallback)

	users := v1.Group("/users")
	users.Use(authHandlers.CheckAuthenticationMiddleware)
	users.GET("/me", usersHandlers.HandleUsersMe)

	return r.gin.Run(r.cfg.Address)
}

// contextMiddleware sets additional useful context to be used by other handlers.
func contextMiddleware(adapter db.Adapter, sessionCache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(helper.ContextDatabase, adapter)
		c.Set(helper.ContextCache, sessionCache)
		c.Next()
	}
}
