package users

import (
	database "github.com/bazuker/backend-bootstrap/pkg/db"
	"github.com/bazuker/backend-bootstrap/pkg/router/helper"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// HandleUsersMe returns information about the authenticated user.
func HandleUsersMe(c *gin.Context) {
	userIDContext := c.MustGet(helper.ContextUserID)
	dbContext := c.MustGet(helper.ContextDatabase)
	db := dbContext.(database.Adapter)

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
