package pkg

import (
	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"google.golang.org/protobuf/types/known/anypb"
)

func NodeFeederPublishMessage(c *Config, k string, m mb.MsgBusServiceClient) error {

	configReq := &pb.Config{
		Filename: "abcd", /* filename with path */
		App:      "configd",
		Data:     []byte("{ \"name\": \"config\", \"value\": \"0.0.1\"}"),
	}

	anyMsg, err := anypb.New(configReq)
	if err != nil {
		log.Errorf("failed to create message: %v", err)
		return err
	}

	msg := &pb.NodeFeederMessage{
		Target:     "ukamaorg" + "." + "network" + "." + "site" + "." + "uk-000000-hnode-00-0000",
		HTTPMethod: "POST",
		Path:       "/v1/configd/config",
		Msg:        anyMsg,
	}

	err = m.PublishRequest(k, msg)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", msg, k, err.Error())
		return err
	}
	return nil
}
