package builder

import (
	"github.com/ukama/ukama/services/common/msgbus"
)

// MsgBus Config
type MsgBusConfig struct {
	Uri string
}

// Queue Config
type MsgBusQConfig struct {
	Exchange         string
	Queue            string
	ExchangeType     string
	ReqRountingKeys  []msgbus.RoutingKey
	RespRountingKeys []msgbus.RoutingKey
}

//TODO:: All this config can be moved to toml file.

// Routing Keys
const (
	RequestDeviceUpdateConfig                    msgbus.RoutingKey = "REQUEST.DEVICE.UPDATE.CONFIG"
	ResponseDeviceUpdateConfig                   msgbus.RoutingKey = "RESPONSE.DEVICE.UPDATE.CONFIG"
	NotificationGitServerCreate                  msgbus.RoutingKey = "NOTIFICATION.GITSERVER.CREATE.*"
	RequestDeviceCreate                          msgbus.RoutingKey = "REQUEST.DEVICE.CREATE.*"
	ResponseDeviceCreate                         msgbus.RoutingKey = "RESPONSE.DEVICE.CREATE.*"
	RequestDeviceDelete                          msgbus.RoutingKey = "REQUEST.DEVICE.DELETE.*"
	ResponseDeviceDelete                         msgbus.RoutingKey = "RESPONSE.DEVICE.DELETE.*"
	RequestDeviceReadConfig                      msgbus.RoutingKey = "REQUEST.DEVICE.READ.CONFIG"
	ResponseDeviceReadConfig                     msgbus.RoutingKey = "RESPONSE.DEVICE.READ.CONFIG"
	CommandControllerExecuteReloadMetric         msgbus.RoutingKey = "CMD.CONTROLLER.EXEC.RELOAD_METRIC"
	ResponseCommandControllerExecuteReloadMetric msgbus.RoutingKey = "CMD.CONTROLLER.EXEC.RELOAD_METRIC"
	RequestDeviceSetobserveConfig                msgbus.RoutingKey = "REQUEST.DEVICE.OBSERVE.CONFIG"
	ResponseDeviceSetobserveConfig               msgbus.RoutingKey = "RESPONSE.DEVICE.OBSERVE.CONFIG"
	RequestDeviceCancelobserveConfig             msgbus.RoutingKey = "REQUEST.DEVICE.CANCEL.CONFIG"
	ResponseDeviceCancelobserveConfig            msgbus.RoutingKey = "RESPONSE.DEVICE.CANCEL.CONFIG"
	CommandDeviceExecuteResource                 msgbus.RoutingKey = "CMD.DEVICE.EXEC.RESOURCE"
	ResponseDeviceExecuteResource                msgbus.RoutingKey = "RESPONSE.DEVICE.EXEC.RESOURCE"
	EventDeviceCreate                            msgbus.RoutingKey = "EVENT.DEVICE.CREATE.*"
	EventVirtNodeUpdateStatus                    msgbus.RoutingKey = "EVENT.VIRTNODE.UPDATE.STATUS"
)

// Device Queue Config
var DeviceQ = MsgBusQConfig{
	Exchange:     "DEVICE_EXCHANGE",
	Queue:        "DEVICE_HANDLE_QUEUE",
	ExchangeType: "topic",
	ReqRountingKeys: []msgbus.RoutingKey{
		RequestDeviceCreate, RequestDeviceDelete,
		RequestDeviceReadConfig, RequestDeviceUpdateConfig,
		CommandControllerExecuteReloadMetric, CommandDeviceExecuteResource,
		RequestDeviceSetobserveConfig, RequestDeviceCancelobserveConfig,
	},
	RespRountingKeys: []msgbus.RoutingKey{
		ResponseDeviceCreate, ResponseDeviceDelete,
		ResponseDeviceUpdateConfig, ResponseDeviceReadConfig,
		ResponseCommandControllerExecuteReloadMetric, ResponseDeviceExecuteResource,
		ResponseDeviceSetobserveConfig, ResponseDeviceCancelobserveConfig, EventDeviceCreate,
	},
}

//TODO:: May be change rotingkeys to [MessageType].[OperationType].[OpernadType].[ResourceType]
var GNotifyQ = MsgBusQConfig{
	Exchange:         "DEVICE_EXCHANGE",
	Queue:            "GNOTIFY_QUEUE",
	ExchangeType:     "topic",
	ReqRountingKeys:  []msgbus.RoutingKey{NotificationGitServerCreate},
	RespRountingKeys: []msgbus.RoutingKey{},
}

// LwM2M Queue Config
var LwM2MQ = MsgBusQConfig{
	Exchange:     "LWM2M_EXCHANGE",
	Queue:        "LWM2M_QUEUE",
	ExchangeType: "topic",
	ReqRountingKeys: []msgbus.RoutingKey{RequestDeviceReadConfig, RequestDeviceUpdateConfig,
		CommandDeviceExecuteResource, RequestDeviceSetobserveConfig,
		RequestDeviceCancelobserveConfig},
	RespRountingKeys: []msgbus.RoutingKey{
		ResponseDeviceUpdateConfig, ResponseDeviceReadConfig,
		ResponseDeviceExecuteResource, ResponseDeviceSetobserveConfig,
		ResponseDeviceCancelobserveConfig,
	},
}
