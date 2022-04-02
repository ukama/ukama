package rest

import (
	"errors"
	"net/http"

	"net/http"

	"github.com/ukama/openIoR/services/common/ukama"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
		ThrowError(c, http.StatusInternalServerError, "Error getting the "+entityType, "", err)
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

// GetNodeIdFromPath extracts nodeId from url parameter with nodeIdKey name and returns the value as NodeId object
// in case of parsing error or missing parameter it will use context to return error message and http.StatusBadRequest code
func GetNodeIdFromPath(c *gin.Context, nodeIdKey string) (id ukama.NodeID, isValid bool) {
	nodeId := c.Param(nodeIdKey)
	id, err := ukama.ValidateNodeId(nodeId)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{
			Message: "Error parsing NodeID",
			Details: err.Error(),
		})
		return "", false
	}
	return id, true
}
