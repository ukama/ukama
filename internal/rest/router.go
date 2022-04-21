package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgtype"
	"github.com/sirupsen/logrus"

	"github.com/ukama/openIoR/services/bootstrap/lookup/internal/db"
	common "github.com/ukama/openIoR/services/common/rest"
	sr "github.com/ukama/openIoR/services/common/srvcrouter"
	"github.com/ukama/openIoR/services/common/ukama"
)

const NodeIdParamName = "node"
const orgNameParamName = "org"
const requestPUTorPOST = "looking_to"
const requestGET = "looking_for"

type Router struct {
	gin      *gin.Engine
	R        *sr.ServiceRouter
	nodeRepo db.NodeRepo
	orgRepo  db.OrgRepo
}

func NewRouter(svcR *sr.ServiceRouter, nodeRepo db.NodeRepo, orgRepo db.OrgRepo, debugMode bool) *Router {
	r := &Router{R: svcR, nodeRepo: nodeRepo, orgRepo: orgRepo}
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

	org := rt.gin.Group("/orgs/")
	{
		org.POST("", rt.addOrgHandler)
		org.GET("node", rt.getDeviceHandler)
		org.POST("node", rt.postDeviceHandler)
	}
}

func (rt *Router) pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (rt *Router) postDeviceHandler(c *gin.Context) {

	orgName := c.Query(orgNameParamName)
	lookingTo := c.Query(requestPUTorPOST)
	nodeId := c.Query(NodeIdParamName)

	logrus.Debugf("Received a request to add device %s to org %s lookingto %s.", nodeId, orgName, lookingTo)

	id, err := ukama.ValidateNodeId(nodeId)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorMessage{
			Message: "Error parsing NodeID",
			Details: err.Error(),
		})
		return
	}

	org, err := rt.orgRepo.GetByName(orgName)
	if err != nil {
		common.SendErrorResponseFromGet(c, "organisation", err)
		return
	}

	err = rt.nodeRepo.AddOrUpdate(&db.Node{NodeID: id.StringLowercase(), OrgID: org.ID})
	if err != nil {
		common.ThrowError(c, http.StatusInternalServerError, "Error adding the node mapping", err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Mapping added or updated"})
}

func (rt *Router) getDeviceHandler(c *gin.Context) {

	orgName := c.Query(orgNameParamName)
	lookingFor := c.Query(requestGET)
	nodeId := c.Query(NodeIdParamName)

	logrus.Debugf("Received a request to read device %s from org %s lookingFor %s.", nodeId, orgName, lookingFor)

	id, err := ukama.ValidateNodeId(nodeId)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorMessage{
			Message: "Error parsing NodeID",
			Details: err.Error(),
		})
		return
	}

	node, err := rt.nodeRepo.Get(id)
	if err != nil {
		common.SendErrorResponseFromGet(c, "node", err)
		return
	}

	resp := GetDeviceResponse{
		NodeId:      id.StringLowercase(),
		Certificate: node.Org.Certificate,
		OrgName:     node.Org.Name,
		Ip:          node.Org.Ip.IPNet.IP.String(),
	}

	c.JSON(http.StatusOK, resp)
}

func (rt *Router) addOrgHandler(c *gin.Context) {
	name := c.Query(orgNameParamName)
	lookingTo := c.Query(requestPUTorPOST)

	logrus.Debugf("Received a request to addorg name %s lookingto %s.", name, lookingTo)

	var req AddOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ThrowError(c, http.StatusBadRequest, "Error parsing request", err.Error(), err)
		return
	}

	ip := pgtype.Inet{}
	err := ip.Set(req.Ip + "/32")
	if err != nil {
		common.ThrowError(c, http.StatusBadRequest, "Error parsing IP", err.Error(), err)
		return
	}

	err = rt.orgRepo.Upsert(&db.Org{Name: name, Certificate: req.Certificate, Ip: ip})
	if err != nil {
		common.ThrowError(c, http.StatusBadRequest, "Error parsing request", err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Organisation added or updated"})
}
