package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgtype"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"ukamaX/lookup/internal/db"
)

type Router struct {
	gin *gin.Engine

	nodeRepo db.NodeRepo
	orgRepo  db.OrgRepo
}

func NewRouter(nodeRepo db.NodeRepo, orgRepo db.OrgRepo, debugMode bool) *Router {
	r := &Router{nodeRepo: nodeRepo, orgRepo: orgRepo}
	if !debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r.init()
	return r
}

func (rt *Router) Run() {
	err := rt.gin.Run()
	if err != nil {
		panic(err)
	}
}

func (rt *Router) init() {
	rt.gin = gin.Default()

	rt.gin.GET("/ping", rt.pingHandler)
	rt.gin.GET("/devices/:uuid", rt.getDeviceHandler)
	rt.gin.POST("/devices/:uuid", rt.postDeviceHandler)
	rt.gin.POST("/orgs/:name", rt.addOrgHandler)
}

func (rt *Router) pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (rt *Router) postDeviceHandler(c *gin.Context) {
	id, isValid := rt.getUuidFromPath(c)
	if !isValid {
		return
	}

	var req DeviceMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		throwError(c, http.StatusBadRequest, "Error parsing request", err.Error(), err)
		return
	}
	if req.Org == "" {
		throwError(c, http.StatusBadRequest, "Organisation field is empty", "", nil)
		return
	}

	org, err := rt.orgRepo.GetByName(req.Org)
	if err != nil {
		rt.sendErrorResponseFromGet(c, "organisation", err)
		return
	}

	err = rt.nodeRepo.AddOrUpdate(&db.Node{UUID: id, OrgID: org.ID})
	if err != nil {
		throwError(c, http.StatusInternalServerError, "Error adding the node mapping", err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Mapping added or updated"})
}

func (rt *Router) getDeviceHandler(c *gin.Context) {
	id, isValid := rt.getUuidFromPath(c)
	if !isValid {
		return
	}

	node, err := rt.nodeRepo.Get(id)
	if err != nil {
		rt.sendErrorResponseFromGet(c, "node", err)
		return
	}

	resp := GetDeviceResponse{
		Uuid:        id,
		Certificate: node.Org.Certificate,
		OrgName:     node.Org.Name,
		Ip:          node.Org.Ip.IPNet.IP.String(),
	}

	c.JSON(http.StatusOK, resp)
}

func (rt *Router) addOrgHandler(c *gin.Context) {
	name := c.Param("name")
	var req AddOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		throwError(c, http.StatusBadRequest, "Error parsing request", err.Error(), err)
		return
	}
	ip := pgtype.Inet{}
	err := ip.Set(req.Ip + "/32")
	if err != nil {
		throwError(c, http.StatusBadRequest, "Error parsing IP", err.Error(), err)
		return
	}
	err = rt.orgRepo.Upsert(&db.Org{Name: name, Certificate: req.Certificate, Ip: ip})
	if err != nil {
		throwError(c, http.StatusBadRequest, "Error parsing request", err.Error(), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Organisation added or updated"})
}

func (rt *Router) getUuidFromPath(c *gin.Context) (id uuid.UUID, isValid bool) {
	uuidStr := c.Param("uuid")
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

func throwError(c *gin.Context, status int, message string, details string, err error) {
	c.JSON(status, ErrorMessage{
		Message: message,
		Details: details,
	})
	logrus.Errorf("Message: %s. Error: %s", message, err)
}

func (rt *Router) sendErrorResponseFromGet(c *gin.Context, entityType string, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		throwError(c, http.StatusNotFound, entityType+" not found", "", err)
	} else {
		throwError(c, http.StatusInternalServerError, "Error finding the "+entityType, "", err)
	}
}
