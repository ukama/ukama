/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package msgbus

// deprecated routing keys
const (
	RequestDeviceUpdateConfig                    RoutingKey = "REQUEST.DEVICE.UPDATE.CONFIG"
	ResponseDeviceUpdateConfig                   RoutingKey = "RESPONSE.DEVICE.UPDATE.CONFIG"
	NotificationGitServerCreate                  RoutingKey = "NOTIFICATION.GITSERVER.CREATE.*"
	RequestDeviceCreate                          RoutingKey = "REQUEST.DEVICE.CREATE.*"
	ResponseDeviceCreate                         RoutingKey = "RESPONSE.DEVICE.CREATE.*"
	RequestDeviceDelete                          RoutingKey = "REQUEST.DEVICE.DELETE.*"
	ResponseDeviceDelete                         RoutingKey = "RESPONSE.DEVICE.DELETE.*"
	RequestDeviceReadConfig                      RoutingKey = "REQUEST.DEVICE.READ.CONFIG"
	ResponseDeviceReadConfig                     RoutingKey = "RESPONSE.DEVICE.READ.CONFIG"
	CommandControllerExecuteReloadMetric         RoutingKey = "CMD.CONTROLLER.EXEC.RELOAD_METRIC"
	ResponseCommandControllerExecuteReloadMetric RoutingKey = "CMD.CONTROLLER.EXEC.RELOAD_METRIC"
	RequestDeviceSetobserveConfig                RoutingKey = "REQUEST.DEVICE.OBSERVE.CONFIG"
	ResponseDeviceSetobserveConfig               RoutingKey = "RESPONSE.DEVICE.OBSERVE.CONFIG"
	RequestDeviceCancelobserveConfig             RoutingKey = "REQUEST.DEVICE.CANCEL.CONFIG"
	ResponseDeviceCancelobserveConfig            RoutingKey = "RESPONSE.DEVICE.CANCEL.CONFIG"
	CommandDeviceExecuteResource                 RoutingKey = "CMD.DEVICE.EXEC.RESOURCE"
	ResponseDeviceExecuteResource                RoutingKey = "RESPONSE.DEVICE.EXEC.RESOURCE"
	EventDeviceCreate                            RoutingKey = "EVENT.DEVICE.CREATE.*"
	EventVirtNodeUpdateStatus                    RoutingKey = "EVENT.VIRTNODE.UPDATE.STATUS"
)

// actual routing keys
const (
	NodeConnectedRoutingKey     RoutingKey = "event.device.mesh.link.connect"
	UserRegisteredRoutingKey      RoutingKey = "event.cloud.identity.user.create"
	NodeFeederRequestRoutingKey RoutingKey = "request.cloud.node-feeder"
	OrgCreatedRoutingKey          RoutingKey = "event.cloud.org.org.created"
	OrgDeletedRoutingKey          RoutingKey = "event.cloud.org.org.deleted"
	NodeUpdatedRoutingKey         RoutingKey = "event.cloud.node.node.updated"

	DefaultExchange = "amq.topic"
)

type NodeUpdateBody struct {
	NodeId string `json:"nodeId"`
	State  string `json:"state"`
	Name   string `json:"name"`
}

type OrgCreatedBody struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

var DeviceQ = MsgBusQConfig{
	Exchange:     "DEVICE_EXCHANGE",
	Queue:        "DEVICE_HANDLE_QUEUE",
	ExchangeType: "topic",
	ReqRountingKeys: []RoutingKey{
		RequestDeviceCreate, RequestDeviceDelete,
		RequestDeviceReadConfig, RequestDeviceUpdateConfig,
		CommandControllerExecuteReloadMetric, CommandDeviceExecuteResource,
		RequestDeviceSetobserveConfig, RequestDeviceCancelobserveConfig,
	},
	RespRountingKeys: []RoutingKey{
		ResponseDeviceCreate, ResponseDeviceDelete,
		ResponseDeviceUpdateConfig, ResponseDeviceReadConfig,
		ResponseCommandControllerExecuteReloadMetric, ResponseDeviceExecuteResource,
		ResponseDeviceSetobserveConfig, ResponseDeviceCancelobserveConfig,
	},
}

var GNotifyQ = MsgBusQConfig{
	Exchange:         "DEVICE_EXCHANGE",
	Queue:            "GNOTIFY_QUEUE",
	ExchangeType:     "topic",
	ReqRountingKeys:  []RoutingKey{NotificationGitServerCreate},
	RespRountingKeys: []RoutingKey{},
}

var LwM2MQ = MsgBusQConfig{
	Exchange:     "LWM2M_EXCHANGE",
	Queue:        "LWM2M_QUEUE",
	ExchangeType: "topic",
	ReqRountingKeys: []RoutingKey{RequestDeviceReadConfig, RequestDeviceUpdateConfig,
		CommandDeviceExecuteResource, RequestDeviceSetobserveConfig,
		RequestDeviceCancelobserveConfig},
	RespRountingKeys: []RoutingKey{
		ResponseDeviceUpdateConfig, ResponseDeviceReadConfig,
		ResponseDeviceExecuteResource, ResponseDeviceSetobserveConfig,
		ResponseDeviceCancelobserveConfig,
	},
}
