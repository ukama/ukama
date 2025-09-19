/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package utils

import (
	"fmt"
	"io"
	"os"
	"time"

	fakeit "github.com/brianvoe/gofakeit/v7"
	"github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/testing/common/amqp"
	conf "github.com/ukama/ukama/testing/common/config"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	defaultDuration = 5 * time.Second
)

func getRoutingKey(orgName string) msgbus.RoutingKeyBuilder {
	return msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem("messaging").SetOrgName(orgName).SetService("mesh")
}

func Run(amqpConf conf.Queue, route string, payload *anypb.Any, out io.Writer) error {

	aClient := amqp.NewAmqpClient(amqpConf, defaultDuration)

	logrus.Infof("Publishing event to %s with payload %v", route, payload)

	respData, err := aClient.PublishMessage(amqpConf.Vhost, amqpConf.Exchange, route, payload)
	if err != nil {
		return fmt.Errorf("failled to publish event: %w", err)
	}

	outputBuf, err := amqp.Serialize(respData, "json")
	if err != nil {
		return fmt.Errorf("error while serializing output data: %w", err)
	}

	_, err = fmt.Fprint(out, string(outputBuf))
	if err != nil {
		return fmt.Errorf("error while writting output: %w", err)
	}

	return nil
}

func PushNodeOnline(orgName, nodeId string, m mb.MsgBusServiceClient) {
	route := getRoutingKey(orgName).SetAction("online").SetObject("node").MustBuild()

	evt := &epb.NodeOnlineEvent{
		NodeId:       nodeId,
		NodeIp:       fakeit.IPv4Address(),
		MeshIp:       fakeit.IPv4Address(),
		MeshHostName: fakeit.DomainName(),
	}

	logrus.Infof("Pushing NodeOnline event for node %s, rpc %+v", nodeId, evt)

	if err := m.PublishRequest(route, evt); err != nil {
		logrus.Errorf("Failed to publish %s event. Error %s", route, err.Error())
	}
}

func PushNodeOnlineViaREST(amqpConf conf.Queue, org, nodeId string, m mb.MsgBusServiceClient) {
	route := getRoutingKey(org).SetAction("online").SetObject("node").MustBuild()

	msg := epb.NodeOnlineEvent{
		NodeId:       nodeId,
		NodeIp:       fakeit.IPv4Address(),
		MeshIp:       fakeit.IPv4Address(),
		MeshHostName: fakeit.DomainName(),
	}

	anyMsg, err := anypb.New(&msg)
	if err != nil {
		logrus.Errorf("Failed to convert message to Any: %v", err)
	}

	err = Run(amqpConf, route, anyMsg, os.Stdout)
	if err != nil {
		logrus.Errorf("Failed to publish %s event. Error %s", route, err.Error())
	}

	logrus.Infof("Successfully published NodeOnline event for node %s", nodeId)
}

func PushNodeReset(orgName, nodeId string, m mb.MsgBusServiceClient) {
	route := getRoutingKey(orgName).SetAction("reset").SetObject("node").MustBuild()
	evt := &epb.NodeOnlineEvent{
		NodeId:       nodeId,
		NodeIp:       fakeit.IPv4Address(),
		MeshIp:       fakeit.IPv4Address(),
		MeshHostName: fakeit.DomainName(),
	}
	if err := m.PublishRequest(route, evt); err != nil {
		logrus.Errorf("Failed to publish %s event. Error %s", route, err.Error())
	}
}

func PushNodeResetViaREST(amqpConf conf.Queue, org, nodeId string, m mb.MsgBusServiceClient) {
	route := getRoutingKey(org).SetAction("reset").SetObject("node").MustBuild()

	msg := epb.NodeOnlineEvent{
		NodeId:       nodeId,
		NodeIp:       fakeit.IPv4Address(),
		MeshIp:       fakeit.IPv4Address(),
		MeshHostName: fakeit.DomainName(),
	}

	anyMsg, err := anypb.New(&msg)
	if err != nil {
		logrus.Errorf("Failed to convert message to Any: %v", err)
	}

	err = Run(amqpConf, route, anyMsg, os.Stdout)
	if err != nil {
		logrus.Errorf("Failed to publish %s event. Error %s", route, err.Error())
	}

	logrus.Infof("Successfully published NodeOnline event for node %s", nodeId)
}

func PushNodeRFOn(orgName, nodeId string, m mb.MsgBusServiceClient) {
	route := getRoutingKey(orgName).SetAction("rfon").SetObject("node").MustBuild()
	evt := &epb.NodeOnlineEvent{
		NodeId:       nodeId,
		NodeIp:       fakeit.IPv4Address(),
		MeshIp:       fakeit.IPv4Address(),
		MeshHostName: fakeit.DomainName(),
	}
	if err := m.PublishRequest(route, evt); err != nil {
		logrus.Errorf("Failed to publish %s event. Error %s", route, err.Error())
	}
}

func PushNodeRFOnViaREST(amqpConf conf.Queue, org, nodeId string, m mb.MsgBusServiceClient) {
	route := getRoutingKey(org).SetAction("rfon").SetObject("node").MustBuild()

	msg := epb.NodeOnlineEvent{
		NodeId:       nodeId,
		NodeIp:       fakeit.IPv4Address(),
		MeshIp:       fakeit.IPv4Address(),
		MeshHostName: fakeit.DomainName(),
	}

	anyMsg, err := anypb.New(&msg)
	if err != nil {
		logrus.Errorf("Failed to convert message to Any: %v", err)
	}

	err = Run(amqpConf, route, anyMsg, os.Stdout)
	if err != nil {
		logrus.Errorf("Failed to publish %s event. Error %s", route, err.Error())
	}

	logrus.Infof("Successfully published NodeOnline event for node %s", nodeId)
}

func PushNodeRFOff(orgName, nodeId string, m mb.MsgBusServiceClient) {
	route := getRoutingKey(orgName).SetAction("rfoff").SetObject("node").MustBuild()
	evt := &epb.NodeOnlineEvent{
		NodeId:       nodeId,
		NodeIp:       fakeit.IPv4Address(),
		MeshIp:       fakeit.IPv4Address(),
		MeshHostName: fakeit.DomainName(),
	}
	if err := m.PublishRequest(route, evt); err != nil {
		logrus.Errorf("Failed to publish %s event. Error %s", route, err.Error())
	}
}

func PushNodeRFOffViaREST(amqpConf conf.Queue, org, nodeId string, m mb.MsgBusServiceClient) {
	route := getRoutingKey(org).SetAction("rfoff").SetObject("node").MustBuild()

	msg := epb.NodeOnlineEvent{
		NodeId:       nodeId,
		NodeIp:       fakeit.IPv4Address(),
		MeshIp:       fakeit.IPv4Address(),
		MeshHostName: fakeit.DomainName(),
	}

	anyMsg, err := anypb.New(&msg)
	if err != nil {
		logrus.Errorf("Failed to convert message to Any: %v", err)
	}

	err = Run(amqpConf, route, anyMsg, os.Stdout)
	if err != nil {
		logrus.Errorf("Failed to publish %s event. Error %s", route, err.Error())
	}

	logrus.Infof("Successfully published NodeOnline event for node %s", nodeId)
}

func PushNodeOff(orgName, nodeId string, m mb.MsgBusServiceClient) {
	route := getRoutingKey(orgName).SetAction("offline").SetObject("node").MustBuild()
	evt := &epb.NodeOfflineEvent{
		NodeId: nodeId,
	}
	if err := m.PublishRequest(route, evt); err != nil {
		logrus.Errorf("Failed to publish %s event. Error %s", route, err.Error())
	}
}

func PushNodeOffViaREST(amqpConf conf.Queue, org, nodeId string, m mb.MsgBusServiceClient) {
	route := getRoutingKey(org).SetAction("offline").SetObject("node").MustBuild()

	msg := epb.NodeOfflineEvent{
		NodeId: nodeId,
	}

	anyMsg, err := anypb.New(&msg)
	if err != nil {
		logrus.Errorf("Failed to convert message to Any: %v", err)
	}

	err = Run(amqpConf, route, anyMsg, os.Stdout)
	if err != nil {
		logrus.Errorf("Failed to publish %s event. Error %s", route, err.Error())
	}

	logrus.Infof("Successfully published NodeOnline event for node %s", nodeId)
}
