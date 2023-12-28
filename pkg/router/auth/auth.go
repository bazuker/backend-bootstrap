package auth

import (
	"github.com/bazuker/backend-bootstrap/pkg/router/helper"
	"log"
	"net/http"
	"time"

	"github.com/akyoto/cache"
	"github.com/bazuker/backend-bootstrap/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// HandleAuthGoogleInitiation handles the initiation of Google authentication.
// "redirect_url" must be passed via query to redirect users after authentication is complete.
func HandleAuthGoogleInitiation(c *gin.Context) {
	redirectURL := c.Query("redirect_url")
	if len(redirectURL) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.HTTPError{
			Message: "'redirect_url' is missing in the request",
		})
		return
	}

	state := oauthGoogleLogin(c.Writer, c.Request)

	sessionCache := c.MustGet(helper.CacheContext).(*cache.Cache)
	sessionCache.Set(state, redirectURL, time.Hour)
}

func HandleAuthGoogleCallback(c *gin.Context) {
	googleUser, state, err := oauthGoogleCallback(c.Writer, c.Request)
	if err != nil {
		log.Println("failed to process google auth callback:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	database := c.MustGet(helper.DatabaseContext).(db.Adapter)
	// Try to find the user by email/
	user, err := database.GetUserByEmail(googleUser.Email)
	if err != nil {
		if err != db.ErrNotFound {
			// Something went wrong.
			log.Println("failed to get user by email:", err)
			return
		}

		user = &db.User{
			ID:            uuid.NewString(),
			FirstName:     googleUser.GivenName,
			LastName:      googleUser.FamilyName,
			Email:         googleUser.Email,
			VerifiedEmail: true, // Verified because this is Google SSO.
		}

		err = database.CreateUser(user)
		if err != nil {
			log.Println("failed to create user:", err)
			return
		}
	}

	sessionCache := c.MustGet(helper.CacheContext).(*cache.Cache)
	redirectURL, ok := sessionCache.Get(state)
	if !ok {
		log.Println("error, state", state, "is missing")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	http.Redirect(c.Writer, c.Request, redirectURL.(string), http.StatusTemporaryRedirect)
}
