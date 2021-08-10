package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

func ThrowError(c *gin.Context, status int, message string, details string, err error) {
	c.JSON(status, ErrorMessage{
		Message: message,
		Details: details,
	})
	logrus.Errorf("Message: %s. Error: %s", message, err)
}

func SendErrorResponseFromGet(c *gin.Context, entityType string, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ThrowError(c, http.StatusNotFound, entityType+" not found", "", err)
	} else {
		ThrowError(c, http.StatusInternalServerError, "Error finding the "+entityType, "", err)
	}
}

func GetUuidFromPath(c *gin.Context, uuidKey string) (id uuid.UUID, isValid bool) {
	uuidStr := c.Param(uuidKey)
	id, err := uuid.FromString(uuidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{
			Message: "Error parsing UUID",
			Details: err.Error(),
		})
		return uuid.UUID{}, false
	}
	return id, true
}
