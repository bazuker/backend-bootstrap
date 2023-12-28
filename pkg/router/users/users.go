package users

import (
	"log"
	"net/http"

	database "github.com/bazuker/backend-bootstrap/pkg/db"
	"github.com/bazuker/backend-bootstrap/pkg/router/helper"
	"github.com/gin-gonic/gin"
)

// HandleUsersMe returns information about the authenticated user.
func HandleUsersMe(c *gin.Context) {
	// Get the database from the context.
	dbContext := c.MustGet(helper.ContextDatabase)
	db := dbContext.(database.Adapter)

	// Find the authenticated user.
	userIDContext := c.MustGet(helper.ContextUserID)
	userID := userIDContext.(string)
	user, err := db.GetUserByID(userID)
	if err != nil {
		log.Printf("failed to get user by ID '%s': %s\n", userID, err.Error())
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			helper.HTTPError{Message: "failed to get user"},
		)
		return
	}

	c.JSON(http.StatusOK, user)
}

// HandleGetUsers returns information about any user if access level is sufficient.
func HandleGetUsers(c *gin.Context) {
	userIDContext := c.MustGet(helper.ContextUserID)
	requestedUserID := c.Param("userid")

	// Check user's access level. Only admins can retrieve information about other users
	userAccessLevelContext := c.MustGet(helper.ContextUserAccessLevel)
	if userAccessLevelContext != database.AccessLevelAdmin &&
		userIDContext.(string) != requestedUserID {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			helper.HTTPError{Message: "insufficient rights"},
		)
		return
	}

	// Get the database from the context.
	dbContext := c.MustGet(helper.ContextDatabase)
	db := dbContext.(database.Adapter)

	// Find the user.
	user, err := db.GetUserByID(requestedUserID)
	if err != nil {
		log.Printf("failed to get user by ID '%s': %s\n", requestedUserID, err.Error())
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			helper.HTTPError{Message: "failed to get user"},
		)
		return
	}

	c.JSON(http.StatusOK, user)
}
