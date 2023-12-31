package auth

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/akyoto/cache"
	"github.com/bazuker/backend-bootstrap/pkg/db"
	"github.com/bazuker/backend-bootstrap/pkg/manager/helper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// HandleAuthGoogleInitiation handles the initiation of Google authentication.
// "redirect_url" must be passed via query to redirect users after authentication is complete.
func HandleAuthGoogleInitiation(c *gin.Context) {
	redirectURL := c.Query("redirect_url")
	if len(redirectURL) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.HTTPMessage{
			Message: "'redirect_url' is missing in the request",
		})
		return
	}

	state := oauthGoogleLogin(c.Writer, c.Request)

	sessionCache := c.MustGet(helper.ContextCache).(*cache.Cache)
	sessionCache.Set(state, redirectURL, time.Hour)
}

// HandleAuthGoogleCallback handles the callback from Google and if successful, redirects the user to 'redirect_url'
// An 'access_token' will be added to the query of the 'redirect_url'.
func HandleAuthGoogleCallback(c *gin.Context) {
	googleUser, state, err := oauthGoogleCallback(c.Writer, c.Request)
	if err != nil {
		log.Println("failed to process google auth callback:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	database := c.MustGet(helper.ContextDatabase).(db.Adapter)
	// Try to find the user by email/
	user, err := database.GetUserByEmail(googleUser.Email)
	if err != nil {
		if err != db.ErrNotFound {
			// Something went wrong.
			log.Println("failed to get user by email:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		user = db.User{
			ID:            uuid.NewString(),
			FirstName:     googleUser.GivenName,
			LastName:      googleUser.FamilyName,
			Email:         googleUser.Email,
			VerifiedEmail: true, // Verified because this is Google SSO.
			AccessLevel:   db.AccessLevelBasic,
		}

		err = database.CreateUser(&user)
		if err != nil {
			log.Println("failed to create user:", err)
			return
		}

		log.Printf("created a new user with ID '%s'\n", user.ID)
	}

	sessionCache := c.MustGet(helper.ContextCache).(*cache.Cache)
	redirectURL, ok := sessionCache.Get(state)
	if !ok {
		log.Println("failed to get state from the cache", state, "is missing")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	sessionCache.Delete(state)

	redirectURLStr := redirectURL.(string)
	u, err := url.Parse(redirectURLStr)
	if err != nil {
		log.Println("failed to parse redirect_url:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Generate an access token.
	accessToken := helper.GenerateRandomString(32)
	query := u.Query()
	query.Set("access_token", accessToken)
	u.RawQuery = query.Encode()
	// Create a user session.
	sessionData := helper.SessionData{
		UserID:      user.ID,
		AccessLevel: user.AccessLevel,
	}
	sessionCache.Set(accessToken, sessionData, time.Hour*24)

	http.Redirect(c.Writer, c.Request, u.String(), http.StatusTemporaryRedirect)
}

func CheckAuthenticationMiddleware(c *gin.Context) {
	accessToken := c.GetHeader("Access-Token")
	if len(accessToken) == 0 {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			helper.HTTPMessage{Message: "'Access-Token' header is missing"},
		)
		return
	}

	// Load information about the user from the cache.
	sessionCache := c.MustGet(helper.ContextCache).(*cache.Cache)
	session, ok := sessionCache.Get(accessToken)
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			helper.HTTPMessage{Message: "no access"},
		)
		return
	}
	sessionData := session.(helper.SessionData)

	// Store the relevant information in the context for other handlers to use.
	c.Set(helper.ContextUserID, sessionData.UserID)
	c.Set(helper.ContextUserAccessLevel, sessionData.AccessLevel)

	c.Next()
}
