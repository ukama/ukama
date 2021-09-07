package msgbus

import (
	"os"
)

type RoutingKeyType string

//Servcie Config
type Config struct {
}

// MsgBus Config
type MsgBusConfig struct {
	Uri string
}

// Queue Config
type MsgBusQConfig struct {
	Exchange         string
	Queue            string
	ExchangeType     string
	ReqRountingKeys  []RoutingKeyType
	RespRountingKeys []RoutingKeyType
}

var MsgBusConf = MsgBusConfig{
	Uri: getEnv("RABBIT_URI", "amqp://guest:guest@localhost:5672/"),
}

const (
	RequestDeviceUpdateConfig                    RoutingKeyType = "REQUEST.DEVICE.UPDATE.CONFIG"
	ResponseDeviceUpdateConfig                   RoutingKeyType = "RESPONSE.DEVICE.UPDATE.CONFIG"
	NotificationGitServerCreate                  RoutingKeyType = "NOTIFICATION.GITSERVER.CREATE.*"
	RequestDeviceCreate                          RoutingKeyType = "REQUEST.DEVICE.CREATE.*"
	ResponseDeviceCreate                         RoutingKeyType = "RESPONSE.DEVICE.CREATE.*"
	RequestDeviceDelete                          RoutingKeyType = "REQUEST.DEVICE.DELETE.*"
	ResponseDeviceDelete                         RoutingKeyType = "RESPONSE.DEVICE.DELETE.*"
	RequestDeviceReadConfig                      RoutingKeyType = "REQUEST.DEVICE.READ.CONFIG"
	ResponseDeviceReadConfig                     RoutingKeyType = "RESPONSE.DEVICE.READ.CONFIG"
	CommandControllerExecuteReloadMetric         RoutingKeyType = "CMD.CONTROLLER.EXEC.RELOAD_METRIC"
	ResponseCommandControllerExecuteReloadMetric RoutingKeyType = "CMD.CONTROLLER.EXEC.RELOAD_METRIC"
	RequestDeviceSetobserveConfig                RoutingKeyType = "REQUEST.DEVICE.OBSERVE.CONFIG"
	ResponseDeviceSetobserveConfig               RoutingKeyType = "RESPONSE.DEVICE.OBSERVE.CONFIG"
	RequestDeviceCancelobserveConfig             RoutingKeyType = "REQUEST.DEVICE.CANCEL.CONFIG"
	ResponseDeviceCancelobserveConfig            RoutingKeyType = "RESPONSE.DEVICE.CANCEL.CONFIG"
	CommandDeviceExecuteResource                 RoutingKeyType = "CMD.DEVICE.EXEC.RESOURCE"
	ResponseDeviceExecuteResource                RoutingKeyType = "RESPONSE.DEVICE.EXEC.RESOURCE"
	EventDeviceCreate                            RoutingKeyType = "EVENT.DEVICE.CREATE.*"
)

var DeviceQ = MsgBusQConfig{
	Exchange:     "DEVICE_EXCHANGE",
	Queue:        "DEVICE_HANDLE_QUEUE",
	ExchangeType: "topic",
	ReqRountingKeys: []RoutingKeyType{
		RequestDeviceCreate, RequestDeviceDelete,
		RequestDeviceReadConfig, RequestDeviceUpdateConfig,
		CommandControllerExecuteReloadMetric, CommandDeviceExecuteResource,
		RequestDeviceSetobserveConfig, RequestDeviceCancelobserveConfig,
	},
	RespRountingKeys: []RoutingKeyType{
		ResponseDeviceCreate, ResponseDeviceDelete,
		ResponseDeviceUpdateConfig, ResponseDeviceReadConfig,
		ResponseCommandControllerExecuteReloadMetric, ResponseDeviceExecuteResource,
		ResponseDeviceSetobserveConfig, ResponseDeviceCancelobserveConfig,
	},
}

//TODO:: May be change rotingkeys to [MessageType].[OperationType].[OpernadType].[ResourceType]
var GNotifyQ = MsgBusQConfig{
	Exchange:         "DEVICE_EXCHANGE",
	Queue:            "GNOTIFY_QUEUE",
	ExchangeType:     "topic",
	ReqRountingKeys:  []RoutingKeyType{NotificationGitServerCreate},
	RespRountingKeys: []RoutingKeyType{},
}

// var LwM2MQ = MsgBusQConfig{
// 	Exchange:         "DEVICE_EXCHANGE",
// 	Queue:            "LWM2M_QUEUE",
// 	ExchangeType:     "topic",
// 	ReqRountingKeys:  []string{"REQUEST.DEVICE.UPDATE.CONFIG", "REQUEST.DEVICE.READ.CONFIG"},
// 	RespRountingKeys: []string{"RESPONSE.DEVICE.UPDATE.CONFIG", "RESPONSE.DEVICE.UPDATE.CONFIG"},
// }

//TODO:: May be change rotingkeys to [MessageType].[OperationType].[OpernadType].[ResourceType]
var LwM2MQ = MsgBusQConfig{
	Exchange:     "LWM2M_EXCHANGE",
	Queue:        "LWM2M_QUEUE",
	ExchangeType: "topic",
	ReqRountingKeys: []RoutingKeyType{RequestDeviceReadConfig, RequestDeviceUpdateConfig,
		CommandDeviceExecuteResource, RequestDeviceSetobserveConfig,
		RequestDeviceCancelobserveConfig},
	RespRountingKeys: []RoutingKeyType{
		ResponseDeviceUpdateConfig, ResponseDeviceReadConfig,
		ResponseDeviceExecuteResource, ResponseDeviceSetobserveConfig,
		ResponseDeviceCancelobserveConfig,
	},
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
