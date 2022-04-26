package main

import (
	"fmt"
	"io/ioutil"
	cfg "lwm2m-gateway/pkg/config"
	"lwm2m-gateway/pkg/lwm2m"
	stat "lwm2m-gateway/specs/common/spec"
	spec "lwm2m-gateway/specs/lwm2mIface/spec"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

// Usage
func usageError() {
	fmt.Println("Usage: testService <Command> <ARGS> ")
	fmt.Println("Commands:")
	fmt.Println("\t READ_CONFIG ")
	fmt.Println("\t UPDATE_CONFIG ")
	fmt.Println("\t EXEC")
	fmt.Println("\t HELP")
}

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debugf("Mocking a service..!!")
	log.Debugf("Args length is %d, Args: %+v \n", len(os.Args), os.Args)

	if len(os.Args) < 2 {
		usageError()
		os.Exit(1)
	}

	// Makes sure connection is closed when service exits.
	handleSigterm(func() {
		log.Debugf("ExitingMocking a service..!!")
	})

	// LoadConfig
	err := cfg.LoadConfig("lwm2m-gateway", "json", "configs")
	if err != nil {
		log.Errorf("LwM2MGateway:: LwM2MGateway:: Failed to load config. Err: %s", err.Error())
		os.Exit(1)
	}

	go lwm2m.Receiver()

	switch os.Args[1] {
	case "READ_CONFIG":
		TestReadConfig()
	case "WRITE_CONFIG":
		TestUpdateConfig()
	case "EXEC":
		TestExecCommand()
	case "OBSERVE":
		TestObserve()
	case "CANCEL":
		TestCancel()
	case "HELP":
		usageError()
		os.Exit(1)
	}

	for {
		time.Sleep(10 * time.Second)
	}
}

// Sending Update config request
func TestUpdateConfig() {
	cfgFile := "test/data/json/sid-xx-00-ABCD-0123/UK-0001-HNODE-SA03-1101/device/UK-1001-COM-1101_3328_0.json"
	cfgData := readCfgFile(cfgFile)

	routingKey := "REQUEST.DEVICE.UPDATE.CONFIG"
	reqMsg := &spec.Lwm2MConfigReqMsg{
		Device:  &spec.Device{Name: "ABCDEF", Uuid: "UK-0001-HNODE-SA03-1101"},
		Uri:     &spec.URI{Object: 3328, Instance: 0, Resource: 5821},
		CfgData: cfgData,
	}

	IfMsgTest(routingKey, reqMsg)
}

// Sending read config request
func TestReadConfig() {
	// TODO:: Correct file name
	cfgFile := "test/data/json/sid-xx-00-ABCD-0123/UK-0001-HNODE-SA03-1101/device/UK-1001-COM-1101_3328_0.json"
	cfgData := readCfgFile(cfgFile)

	routingKey := "REQUEST.DEVICE.READ.CONFIG"
	reqMsg := &spec.Lwm2MConfigReqMsg{
		Device:  &spec.Device{Name: "ABCDEF", Uuid: "UK-0001-HNODE-SA03-1101"},
		Uri:     &spec.URI{Object: 34570, Instance: 0, Resource: 0},
		CfgData: cfgData,
	}
	IfMsgTest(routingKey, reqMsg)
}

// Sending read config request
func TestExecCommand() {
	cfgFile := "test/data/json/sid-xx-00-ABCD-0123/UK-0001-HNODE-SA03-1101/device/UK-1001-COM-1101_3328_0.json"
	cfgData := readCfgFile(cfgFile)

	routingKey := "REQUEST.DEVICE.EXEC.COMMAND"
	reqMsg := &spec.Lwm2MConfigReqMsg{
		Device:  &spec.Device{Name: "ABCDEF", Uuid: "UK-0001-HNODE-SA03-1101"},
		Uri:     &spec.URI{Object: 3328, Instance: 0, Resource: 5605},
		CfgData: cfgData,
	}
	IfMsgTest(routingKey, reqMsg)
}

// Sending observation config request
func TestObserve() {
	routingKey := "REQUEST.DEVICE.OBSERVE.CONFIG"
	reqMsg := &spec.Lwm2MConfigReqMsg{
		Device:  &spec.Device{Name: "ABCDEF", Uuid: "UK-0001-HNODE-SA03-1101"},
		Uri:     &spec.URI{Object: 34570, Instance: 0, Resource: 1},
		CfgData: "",
	}
	IfMsgTest(routingKey, reqMsg)
}

// Sending observation cancelation request
func TestCancel() {
	routingKey := "REQUEST.DEVICE.CANCEL.CONFIG"
	reqMsg := &spec.Lwm2MConfigReqMsg{
		Device:  &spec.Device{Name: "ABCDEF", Uuid: "UK-0001-HNODE-SA03-1101"},
		Uri:     &spec.URI{Object: 34570, Instance: 0, Resource: 1},
		CfgData: "",
	}
	IfMsgTest(routingKey, reqMsg)
}

// Mocking  incoming requests
func IfMsgTest(routingKey string, reqMsg *spec.Lwm2MConfigReqMsg) {
	var respCode stat.StatusCode
	respMsg := &spec.Lwm2MConfigRespMsg{}
	log.Debugf("Message received by IfMsgHandlerCB:: %+v", reqMsg)
	// Filter Read/Write configs.
	switch routingKey {
	case "REQUEST.DEVICE.UPDATE.CONFIG":

		respCode = lwm2m.WriteConfig(reqMsg)
		// Preparing update response
		respMsg = &spec.Lwm2MConfigRespMsg{
			Device:     reqMsg.Device,
			Uri:        reqMsg.Uri,
			CfgData:    reqMsg.CfgData,
			StatusCode: respCode,
		}

	case "REQUEST.DEVICE.READ.CONFIG":
		var readCfg string
		rrespCode := lwm2m.ReadConfig(reqMsg, &readCfg)
		// Preparing read response
		respMsg = &spec.Lwm2MConfigRespMsg{
			Device:     reqMsg.Device,
			Uri:        reqMsg.Uri,
			CfgData:    readCfg,
			StatusCode: rrespCode,
		}

	case "REQUEST.DEVICE.EXEC.COMMAND":
		respCode = lwm2m.ExecCommand(reqMsg)
		// Preparing read response
		respMsg = &spec.Lwm2MConfigRespMsg{
			Device:     reqMsg.Device,
			Uri:        reqMsg.Uri,
			CfgData:    reqMsg.CfgData,
			StatusCode: respCode,
		}

	case "REQUEST.DEVICE.OBSERVE.CONFIG":
		respCode = lwm2m.ObserveConfig(reqMsg)

		// Preparing read response
		respMsg = &spec.Lwm2MConfigRespMsg{
			Device:     reqMsg.Device,
			Uri:        reqMsg.Uri,
			CfgData:    reqMsg.CfgData,
			StatusCode: respCode,
		}
	case "REQUEST.DEVICE.CANCEL.CONFIG":
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
		log.Errorf("Unknown Routing key received.")
	}

	log.Debugf("LwM2M RPC Response Message:: %v", respMsg)
	log.Debugf("LwM2M RPC Response Code:: %v", respCode)

	// marshal
	data, err := proto.Marshal(respMsg)
	if err != nil {
		log.Errorf("fail marshal: %s", err.Error())
		return
	}
	log.Debugf(" Marsheld LwM2M RPC Response Message:: %v , Error if any:: %+v", data, err)

}

// Handles Ctrl+C or most other means of "controlled" shutdown gracefully. Invokes the supplied func before exiting.
func handleSigterm(handleExit func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		<-c
		handleExit()
		os.Exit(1)
	}()

}

// Read entire file content, giving us little control but
// making it very simple. No need to close the file.
func readCfgFile(fileP string) string {
	content, err := ioutil.ReadFile(fileP)
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string and print to screen
	str := string(content)
	log.Debugf("Config File : %+v", str)
	return str
}
