package pkg

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
)

type ConfigData struct {
	FileName  string `json:"fileName"`
	App       string `json:"app"`
	Version   string `json:"version"`
	Data      []byte `json:"data"`
	Reason    int    `json:"reason"`
	Timestamp uint32 `json:"timestamp"`
	FileCount int    `json:"file_count"`
}

func NodeFeederPublishMessage(c *Config, k string, m mb.MsgBusServiceClient) error {

	configReq := ConfigData{
		FileName:  "abcd", /* filename with path */
		App:       "configd",
		Version:   "abcdef",
		Reason:    1,
		Timestamp: 1000,
		FileCount: 2,
		Data:      []byte("{ \"name\": \"config\", \"value\": \"0.0.1\"}"),
	}

	jReq, err := json.Marshal(configReq)
	if err != nil {
		log.Errorf("Failed to marshal configdata %+v. Errors %s", configReq, err.Error())
		return err
	}

	msg := &pb.NodeFeederMessage{
		Target:     "ukamaorg" + "." + "network" + "." + "site" + "." + "uk-000000-hnode-00-0000",
		HTTPMethod: "POST",
		Path:       "/v1/configd/config",
		Msg:        jReq,
	}

	err = m.PublishRequest(k, msg)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", msg, k, err.Error())
		return err
	}
	return nil
}
