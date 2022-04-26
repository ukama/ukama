package iface

import (
	"fmt"
	"lwm2m-gateway/pkg/lwm2m"
	stat "lwm2m-gateway/specs/common/spec"
	spec "lwm2m-gateway/specs/lwm2mIface/spec"
	"os"

	"github.com/ukama/ukamaX/common/msgbus"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

var ifMsgClient msgbus.IMsgBus
var msgHandlerName string = "Lwm2mGateway"
var evtHandlerPrefix string = "EventHandler_"

// MsgBus Config
type MsgBusConfig struct {
	Uri string
}

var MsgBusConf = MsgBusConfig{
	Uri: getEnv("RABBIT_URI", "amqp://guest:guest@localhost:5672/"),
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

// Initialize Iface for receiving request messages
func Start() {
	log.Infof("Iface:: Connecting LwM2M Gateway to MessageBus.\n")
	ifMsgClient = &msgbus.MsgClient{}
	ifMsgClient.ConnectToBroker(MsgBusConf.Uri)
	err := ifMsgClient.Subscribe(msgbus.LwM2MQ.Queue, msgbus.LwM2MQ.Exchange, msgbus.LwM2MQ.ExchangeType, msgbus.LwM2MQ.ReqRountingKeys, msgHandlerName, IfMsgHandlerCB)
	//err := ifMsgClient.SubscribeToQueue(msgbus.LwM2MQ.Queue, msgHandlerName, IfMsgHandlerCB)
	failOnError(err, "Could not start subscribe to "+msgbus.LwM2MQ.Exchange+msgbus.LwM2MQ.ExchangeType)

	//Register callback for Events
	lwm2m.RegisterCallbackForEvents(PublishMsgOnQueue)
}

// Initialize Iface for receiving request messages
func Stop() {
	log.Infof("Iface:: Closing LwM2M Gateway connection to MessageBus.\n")
	if ifMsgClient != nil {
		ifMsgClient.Close()
	}
}

// Publish message to msgbus
func PublishMsg(message proto.Message) {
	// Marshal
	data, err := proto.Marshal(message)
	if err != nil {
		log.Errorf("Iface:: fail marshal: %s", err.Error())
		return
	}
	log.Debugf("Iface:: Proto data for message is %+v", data)

	// Publish a message
	err = ifMsgClient.Publish(data, msgbus.LwM2MQ.Queue, msgbus.LwM2MQ.Exchange, msgbus.LwM2MQ.RespRountingKeys[0], msgbus.LwM2MQ.ExchangeType)
	if err != nil {
		log.Errorf(err.Error())
	}
}

// Publish message to msgbus Queue
func PublishMsgOnQueue(message proto.Message, uuid string) {
	// Marshal
	data, err := proto.Marshal(message)
	if err != nil {
		log.Errorf("Iface:: fail marshal: %s", err.Error())
		return
	}

	queueName := evtHandlerPrefix + uuid
	log.Debugf("Iface:: publishing on Qeueue %s. Proto data for message is %+v ", queueName, data)

	// Publish a message
	// Publish a message
	err = ifMsgClient.Publish(data, queueName, msgbus.LwM2MQ.Exchange, msgbus.EventDeviceCreate, msgbus.LwM2MQ.ExchangeType)
	if err != nil {
		log.Errorf(err.Error())
	}
}

//Unmarshal Request Message
func unmarshalRequestMsg(d amqp.Delivery) (*spec.Lwm2MConfigReqMsg, error) {
	// unmarshal
	reqMsg := &spec.Lwm2MConfigReqMsg{}
	err := proto.Unmarshal(d.Body, reqMsg)
	if err != nil {
		log.Errorf("Iface:: Fail unmarshal: %s", d.Body)
		return nil, err
	}

	log.Debugf("Iface:: LWM2M Received request msg: %v", reqMsg)
	log.Debugf("Iface:: LWM2M Received request msg: Object %d, Instance %d and Resource %d", reqMsg.Uri.Object, reqMsg.Uri.Instance, reqMsg.Uri.Resource)
	return reqMsg, nil
}

//Handling incoming requests
func IfMsgHandlerCB(d amqp.Delivery, ch chan<- bool) {
	var respCode stat.StatusCode
	respMsg := &spec.Lwm2MConfigRespMsg{}
	done := true

	// Unmarshal
	reqMsg, err := unmarshalRequestMsg(d)
	if err != nil {
		log.Errorf("Iface:: Failed to decode Update request for the config. Error:: %+v.", err)
		return
	}

	// Filter Read/Write configs.
	switch msgbus.RoutingKey(d.RoutingKey) {
	case msgbus.RequestDeviceUpdateConfig:

		respCode = lwm2m.WriteConfig(reqMsg)

		// Preparing update response
		respMsg = &spec.Lwm2MConfigRespMsg{
			Device:     reqMsg.Device,
			Uri:        reqMsg.Uri,
			CfgData:    reqMsg.CfgData,
			StatusCode: respCode,
		}

	case msgbus.RequestDeviceReadConfig:
		var readCfg string
		respCode := lwm2m.ReadConfig(reqMsg, &readCfg)
		// Preparing read response
		respMsg = &spec.Lwm2MConfigRespMsg{
			Device:     reqMsg.Device,
			Uri:        reqMsg.Uri,
			CfgData:    readCfg,
			StatusCode: respCode,
		}

	case msgbus.CommandDeviceExecuteResource:
		respCode = lwm2m.ExecCommand(reqMsg)

		// Preparing read response
		respMsg = &spec.Lwm2MConfigRespMsg{
			Device:     reqMsg.Device,
			Uri:        reqMsg.Uri,
			CfgData:    reqMsg.CfgData,
			StatusCode: respCode,
		}

	case msgbus.RequestDeviceSetobserveConfig:
		respCode = lwm2m.ObserveConfig(reqMsg)

		// Preparing read response
		respMsg = &spec.Lwm2MConfigRespMsg{
			Device:     reqMsg.Device,
			Uri:        reqMsg.Uri,
			CfgData:    reqMsg.CfgData,
			StatusCode: respCode,
		}

	case msgbus.RequestDeviceCancelobserveConfig:
		respCode = lwm2m.CancelObservationOnConfig(reqMsg)

		// Preparing read response
		respMsg = &spec.Lwm2MConfigRespMsg{
			Device:     reqMsg.Device,
			Uri:        reqMsg.Uri,
			CfgData:    reqMsg.CfgData,
			StatusCode: respCode,
		}

	default:
		respCode = stat.StatusCode_ERR_NOT_IMPLEMENTED
		log.Errorf("Iface:: Unknown Routing key received %s.", d.RoutingKey)
	}

	log.Debugf("Iface:: Response to Controller:: Response Message :: %v", respMsg)
	log.Debugf("Iface:: Response to Controller:: Response code :: %v", respCode)

	// If request RPC request publish the RPC repsonse to client.
	if d.ReplyTo != "" {

		// marshal
		data, err := proto.Marshal(respMsg)
		if err != nil {
			log.Errorf("Iface:: fail marshal: %s", err.Error())
			return
		}

		// Publish a message
		err = ifMsgClient.PublishRPCResponse(data, d.CorrelationId, msgbus.RoutingKey(d.ReplyTo))
		if err != nil {
			log.Errorf(err.Error())
		}

	}

	// Respond to PublishRPCRequest
	ch <- done
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Errorf("Iface:: %s: %s", msg, err)
		panic(fmt.Sprintf("Iface:: %s: %s", msg, err))
	}
}
