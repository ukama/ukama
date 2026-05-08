package adapters

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/node/site-controller/pkg"
)

const SourceSiteController = "site-controller"

type NodeCommandPublisher interface {
	Send(nodeId string, method string, path string, body []byte) error
}

type NodeFeederAdapter struct {
	orgName string
	msgbus  mb.MsgBusServiceClient
	routing msgbus.RoutingKeyBuilder
}

func NewNodeFeederAdapter(orgName string, msgBus mb.MsgBusServiceClient) *NodeFeederAdapter {
	return &NodeFeederAdapter{
		orgName: orgName,
		msgbus:  msgBus,
		routing: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (a *NodeFeederAdapter) Send(nodeId string, method string, path string, body []byte) error {
	route := "request.cloud.local" + "." + a.orgName + "." + pkg.SystemName + "." + pkg.ServiceName + "." + "nodefeeder" + "." + "publish"
	target := fmt.Sprintf("%s...%s", a.orgName, nodeId)
	msg := &epb.NodeFeederMessage{
		Target:     target,
		HttpMethod: method,
		Path:       path,
		Msg:        body,
		NodeId:     nodeId,
	}
	log.Infof("site-controller: publishing node command node=%s method=%s path=%s", nodeId, method, path)
	return a.msgbus.PublishRequest(route, msg)
}

func marshalBody(v interface{}) ([]byte, error) {
	if v == nil {
		return []byte(""), nil
	}
	return json.Marshal(v)
}
