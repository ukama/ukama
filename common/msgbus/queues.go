package msgbus

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

var GNotifyQ = MsgBusQConfig{
	Exchange:         "DEVICE_EXCHANGE",
	Queue:            "GNOTIFY_QUEUE",
	ExchangeType:     "topic",
	ReqRountingKeys:  []RoutingKeyType{NotificationGitServerCreate},
	RespRountingKeys: []RoutingKeyType{},
}

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
