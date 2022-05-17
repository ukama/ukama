package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/services/cloud/hss/pb/gen"
	"github.com/ukama/ukama/services/common/msgbus"
	"github.com/wagslane/go-rabbitmq"
	"google.golang.org/protobuf/encoding/protojson"
)

type DeviceFeederReqGenerator struct {
	publisher *rabbitmq.Publisher
}

type DevicesUpdateRequest struct {
	Target     string `json:"target"` // Target devices in form of "organization.network.device-id". Device id and network could be wildcarded
	HttpMethod string `json:"httpMethod"`
	Path       string `json:"path"`
	Body       string `json:"body"`
}

func NewDeviceFeederReqGenerator(amqpConnStr string) (*DeviceFeederReqGenerator, error) {

	publisher, err := rabbitmq.NewPublisher(amqpConnStr, rabbitmq.Config{})
	if err != nil {
		return nil, err
	}

	return &DeviceFeederReqGenerator{
		publisher: publisher,
	}, nil
}

func (d DeviceFeederReqGenerator) ImsiAdded(org string, imsi *pb.ImsiRecord) {
	body, err := protojson.Marshal(imsi)
	if err != nil {
		log.Error("Failed to marshal imsi record: ", err)
		return
	}

	d.sendMessage(org, "/hss/subscriber/", http.MethodPost, body)
}

func (d DeviceFeederReqGenerator) ImsiUpdated(org string, imsi *pb.ImsiRecord) {
	body, err := protojson.Marshal(imsi)
	if err != nil {
		log.Error("Failed to marshal imsi record: ", err)
		return
	}

	d.sendMessage(org, "/hss/subscriber/", http.MethodPost, body)
}

func (d DeviceFeederReqGenerator) ImsiDeleted(org string, imsi string) {
	delBody :=
		fmt.Sprintf(`{
		"subscriber_info":{ "imsi": "%s" } 
	}`, imsi)

	d.sendMessage(org, "/hss/subscriber/", http.MethodDelete, []byte(delBody))
}

func (d DeviceFeederReqGenerator) GutiAdded(org string, imsi string, guti *pb.Guti) {
	gs, err := json.Marshal(guti)
	if err != nil {
		log.Error("Failed to marshal guti record: ", err)
		return
	}

	body := fmt.Sprintf(`
	{	
		"imsi": "%s",
		%s
	}`, imsi, string(gs))
	if err != nil {
		log.Error("Failed to marshal imsi record: ", err)
	}

	d.sendMessage(org, "/hss/guti/", http.MethodPut, []byte(body))
}

func (d DeviceFeederReqGenerator) TaiUpdated(org string, tai *pb.UpdateTaiRequest) {
	body := fmt.Sprintf(`
	{	
		"imsi": "%s",
		"tai": {
			"plmn_id": "%s",
			"tac": %d
		}
	}`, tai.Imsi, tai.PlmnId, tai.Tac)

	d.sendMessage(org, "/hss/tai/", http.MethodPut, []byte(body))
}

func (d DeviceFeederReqGenerator) sendMessage(org string, path string, method string, body []byte) {
	log.Infof("Sending message to device feeder.  Org: %s, path: %s", org, path)
	req := DevicesUpdateRequest{
		Target:     org + ".*",
		Path:       path,
		HttpMethod: method,
		Body:       string(body),
	}

	js, err := json.Marshal(req)
	if err != nil {
		log.Error("Failed to marshal request: ", err)
	}

	err = d.publisher.Publish(js, []string{string(msgbus.DeviceFeederRequestRoutingKey)},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange(msgbus.DefaultExchange))
	if err != nil {
		log.Error("Failed to publish message: ", err)
	}
}
