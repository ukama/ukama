package utils

import (
	"log"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

func getRoutingKey(orgName string) msgbus.RoutingKeyBuilder {
	return msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem("messaging").SetOrgName(orgName).SetService("mesh")
}

func PushNodeOnline(orgName, nodeId string, m mb.MsgBusServiceClient) {
	route := getRoutingKey(orgName).SetAction("online").SetObject("node").MustBuild()
	logrus.Infof("Pushing NodeOnline event for node %s", nodeId)
	evt := &epb.NodeOnlineEvent{
		NodeId:       nodeId,
		NodeIp:       gofakeit.IPv4Address(),
		MeshIp:       gofakeit.IPv4Address(),
		MeshHostName: gofakeit.DomainName(),
	}
	if err := m.PublishRequest(route, evt); err != nil {
		log.Fatalf("Failed to publish %s event. Error %s", route, err.Error())
	}
}

func PushNodeReset(orgName, nodeId string, m mb.MsgBusServiceClient) {
	route := getRoutingKey(orgName).SetAction("reset").SetObject("node").MustBuild()
	evt := &epb.NodeOnlineEvent{
		NodeId:       nodeId,
		NodeIp:       gofakeit.IPv4Address(),
		MeshIp:       gofakeit.IPv4Address(),
		MeshHostName: gofakeit.DomainName(),
	}
	if err := m.PublishRequest(route, evt); err != nil {
		log.Fatalf("Failed to publish %s event. Error %s", route, err.Error())
	}
}

func PushNodeRFOn(orgName, nodeId string, m mb.MsgBusServiceClient) {
	route := getRoutingKey(orgName).SetAction("rfon").SetObject("node").MustBuild()
	evt := &epb.NodeOnlineEvent{
		NodeId:       nodeId,
		NodeIp:       gofakeit.IPv4Address(),
		MeshIp:       gofakeit.IPv4Address(),
		MeshHostName: gofakeit.DomainName(),
	}
	if err := m.PublishRequest(route, evt); err != nil {
		log.Fatalf("Failed to publish %s event. Error %s", route, err.Error())
	}
}

func PushNodeRFOff(orgName, nodeId string, m mb.MsgBusServiceClient) {
	route := getRoutingKey(orgName).SetAction("rfoff").SetObject("node").MustBuild()
	evt := &epb.NodeOnlineEvent{
		NodeId:       nodeId,
		NodeIp:       gofakeit.IPv4Address(),
		MeshIp:       gofakeit.IPv4Address(),
		MeshHostName: gofakeit.DomainName(),
	}
	if err := m.PublishRequest(route, evt); err != nil {
		log.Fatalf("Failed to publish %s event. Error %s", route, err.Error())
	}
}

func PushNodeOff(orgName, nodeId string, m mb.MsgBusServiceClient) {
	route := getRoutingKey(orgName).SetAction("off").SetObject("node").MustBuild()
	evt := &epb.NodeOfflineEvent{
		NodeId: nodeId,
	}
	if err := m.PublishRequest(route, evt); err != nil {
		log.Fatalf("Failed to publish %s event. Error %s", route, err.Error())
	}
}
