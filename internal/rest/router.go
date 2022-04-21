package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgtype"

	"github.com/ukama/openIoR/services/bootstrap/lookup/internal/db"
	common "github.com/ukama/openIoR/services/common/rest"
	sr "github.com/ukama/openIoR/services/common/srvcrouter"
)

const NodeIdParamName = "nodeId"
const orgNameParamName = "org"

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

	org := rt.gin.Group("/orgs/:org")
	{
		org.POST("", rt.addOrgHandler)
		org.GET("devices/:"+NodeIdParamName, rt.getDeviceHandler)
		org.POST("devices/:"+NodeIdParamName, rt.postDeviceHandler)
	}
}

func (rt *Router) pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (rt *Router) postDeviceHandler(c *gin.Context) {
	orgName := c.Param(orgNameParamName)
	id, isValid := common.GetNodeIdFromPath(c, NodeIdParamName)
	if !isValid {
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
	id, isValid := common.GetNodeIdFromPath(c, NodeIdParamName)
	if !isValid {
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
	name := c.Param(orgNameParamName)
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
