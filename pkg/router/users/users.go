package users

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	database "github.com/bazuker/backend-bootstrap/pkg/db"
	"github.com/bazuker/backend-bootstrap/pkg/fileStore"
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
			helper.HTTPMessage{Message: "failed to get user"},
		)
		return
	}

	c.JSON(http.StatusOK, user)
}

// HandleUsersMePhoto handles user photo upload.
func HandleUsersMePhoto(c *gin.Context) {
	// Get the file from the request.
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Println("failed to get the form file:", err)
		c.JSON(
			http.StatusBadRequest,
			helper.HTTPMessage{
				Message: "failed to get the form file",
			},
		)
		return
	}

	// Check if the file extension is appropriate.
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	log.Println(ext)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		c.JSON(
			http.StatusBadRequest,
			helper.HTTPMessage{
				Message: "invalid image format. Only JPEG and PNG are supported",
			},
		)
		return
	}

	// Open the file.
	file, err := fileHeader.Open()
	defer file.Close()
	if err != nil {
		log.Println("failed to open the form file:", err)
		c.JSON(
			http.StatusBadRequest,
			helper.HTTPMessage{
				Message: "failed to open the form file",
			},
		)
		return
	}

	// Read the file.
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, file)
	if err != nil {
		log.Println("failed to read the form file:", err)
		c.JSON(
			http.StatusInternalServerError,
			helper.HTTPMessage{
				Message: "failed to read the form file",
			},
		)
		return
	}

	// Save (or upload) to the file store.
	userIDContext := c.MustGet(helper.ContextUserID)
	userID := userIDContext.(string)
	fileStoreContext := c.MustGet(helper.ContextFileStore)
	fs := fileStoreContext.(fileStore.FileStore)
	objectKey := fmt.Sprintf("%s-photo%s", userID, ext)
	err = fs.PutObject(buf.Bytes(), objectKey)
	if err != nil {
		log.Println("failed to save the file to the filestore:", err)
		c.JSON(
			http.StatusInternalServerError,
			helper.HTTPMessage{
				Message: "failed to save the file to the filestore",
			},
		)
		return
	}

	// Update user's photo in the database.
	dbContext := c.MustGet(helper.ContextDatabase)
	db := dbContext.(database.Adapter)
	err = db.UpdateUserPhotoURL(userID, objectKey)
	if err != nil {
		log.Printf("failed to update user '%s' photo URL: %s\n", userID, err.Error())
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			helper.HTTPMessage{Message: "failed to update user photo URL"},
		)
		return
	}

	c.JSON(http.StatusOK, helper.HTTPMessage{Message: "ok"})
}

// HandleUsersMeDeletePhoto handles user photo deletion.
func HandleUsersMeDeletePhoto(c *gin.Context) {
	userIDContext := c.MustGet(helper.ContextUserID)
	userID := userIDContext.(string)

	// Get the database from the context.
	dbContext := c.MustGet(helper.ContextDatabase)
	db := dbContext.(database.Adapter)

	// Find the user.
	user, err := db.GetUserByID(userID)
	if err != nil {
		log.Printf("failed to get user by ID '%s': %s\n", userID, err.Error())
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			helper.HTTPMessage{Message: "failed to get user"},
		)
		return
	}

	// No photo - nothing to delete.
	if len(user.PhotoURL) == 0 {
		c.JSON(http.StatusOK, helper.HTTPMessage{Message: "ok"})
		return
	}

	// Delete the photo from the file store.
	fileStoreContext := c.MustGet(helper.ContextFileStore)
	fs := fileStoreContext.(fileStore.FileStore)
	err = fs.DeleteObject(user.PhotoURL)
	if err != nil {
		log.Printf("failed to delete user photo '%s': %s\n", user.PhotoURL, err.Error())
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			helper.HTTPMessage{Message: "failed to delete photo"},
		)
		return
	}

	// Delete the photo URL from the database.
	err = db.UpdateUserPhotoURL(userID, "")
	if err != nil {
		log.Printf("failed to update user '%s' photo URL: %s\n", userID, err.Error())
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			helper.HTTPMessage{Message: "failed to update user photo URL"},
		)
		return
	}

	c.JSON(http.StatusOK, helper.HTTPMessage{Message: "ok"})
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
			helper.HTTPMessage{Message: "insufficient rights"},
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
			helper.HTTPMessage{Message: "failed to get user"},
		)
		return
	}

	c.JSON(http.StatusOK, user)
}
