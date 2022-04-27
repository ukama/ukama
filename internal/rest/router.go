package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgtype"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/wI2L/fizz"

	"github.com/ukama/openIoR/services/bootstrap/lookup/cmd/version"
	"github.com/ukama/openIoR/services/bootstrap/lookup/internal"
	"github.com/ukama/openIoR/services/bootstrap/lookup/internal/db"
	"github.com/ukama/openIoR/services/common/rest"
	sr "github.com/ukama/openIoR/services/common/srvcrouter"
	"github.com/ukama/openIoR/services/common/ukama"
)

const NodeIdParamName = "node"

type Router struct {
	fizz     *fizz.Fizz
	port     int
	R        *sr.ServiceRouter
	nodeRepo db.NodeRepo
	orgRepo  db.OrgRepo
}

func NewRouter(config *internal.Config, svcR *sr.ServiceRouter, nodeRepo db.NodeRepo, orgRepo db.OrgRepo, debugMode bool) *Router {

	f := rest.NewFizzRouter(&config.Server, internal.ServiceName, version.Version, internal.IsDebugMode)

	r := &Router{fizz: f,
		port:     config.Server.Port,
		R:        svcR,
		nodeRepo: nodeRepo,
		orgRepo:  orgRepo,
	}

	if !debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r.init()
	return r
}

func (r *Router) Run(close chan error) {
	logrus.Info("Listening on port ", r.port)
	err := r.fizz.Engine().Run(fmt.Sprint(":", r.port))
	if err != nil {
		close <- err
	}
	close <- nil
}

func (r *Router) init() {
	org := r.fizz.Group("/orgs/", "Org", "Organizaton")
	org.POST("", nil, tonic.Handler(r.addOrgHandler, http.StatusOK))
	org.GET("node", nil, tonic.Handler(r.getNodeHandler, http.StatusOK))
	org.POST("node", nil, tonic.Handler(r.postNodeHandler, http.StatusOK))

}

func (r *Router) postNodeHandler(c *gin.Context, req *ReqAddNode) error {
	logrus.Debugf("Received a request to add Node %s to org %s lookingto %s.", req.NodeID, req.OrgName, req.LookingTo)

	id, err := ukama.ValidateNodeId(req.NodeID)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  "error parsing NodeId :" + err.Error(),
		}
	}

	org, err := r.orgRepo.GetByName(req.OrgName)
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  "organization :" + err.Error(),
		}
	}

	err = r.nodeRepo.AddOrUpdate(&db.Node{NodeID: id.StringLowercase(), OrgID: org.ID})
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  "error adding the node mapping :" + err.Error(),
		}
	}

	return nil
}

func (r *Router) getNodeHandler(c *gin.Context, req *ReqGetNode) (*RespGetNode, error) {

	logrus.Debugf("Received a request to read Node %s from org %s lookingFor %s.", req.NodeID, req.OrgName, req.LookingFor)

	id, err := ukama.ValidateNodeId(req.NodeID)
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusBadRequest,
			Message:  "error parsing NodeId :" + err.Error(),
		}
	}

	node, err := r.nodeRepo.Get(id)
	if err != nil {
		return nil, rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  "node :" + err.Error(),
		}
	}

	resp := &RespGetNode{
		NodeId:      id.StringLowercase(),
		Certificate: node.Org.Certificate,
		OrgName:     node.Org.Name,
		Ip:          node.Org.Ip.IPNet.IP.String(),
	}

	return resp, nil
}

func (r *Router) addOrgHandler(c *gin.Context, req *ReqAddOrg) error {
	logrus.Debugf("Received a request to addorg name %s lookingto %s.", req.OrgName, req.LookingTo)

	// var req AddOrgRequest
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	common.ThrowError(c, http.StatusBadRequest, "Error parsing request", err.Error(), err)
	// 	return
	// }

	ip := pgtype.Inet{}
	err := ip.Set(req.Ip + "/32")
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusNotFound,
			Message:  "Error parsing IP :" + err.Error(),
		}
	}

	err = r.orgRepo.Upsert(&db.Org{Name: req.OrgName, Certificate: req.Certificate, Ip: ip})
	if err != nil {
		return rest.HttpError{
			HttpCode: http.StatusInternalServerError,
			Message:  "Error adding org :" + err.Error(),
		}
	}

	return nil
}
