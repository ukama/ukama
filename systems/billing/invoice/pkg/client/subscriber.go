package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
)

const SubscriberEndpoint = "/v1/subscribers"

type SubscriberInfo struct {
	SubscriberId          uuid.UUID `json:"subscriber_id,omitempty"`
	OrgId                 uuid.UUID `json:"org_id,omitempty"`
	NetworkId             uuid.UUID `json:"network_id,omitempty"`
	FirstName             string    `json:"first_name,omitempty"`
	LastName              string    `json:"last_name,omitempty"`
	Email                 string    `json:"email,omitempty"`
	PhoneNumber           string    `json:"phone_number,omitempty"`
	Address               string    `json:"address,omitempty"`
	Dob                   string    `json:"date_of_birth,omitempty"`
	ProofOfIdentification string    `json:"proof_of_identification,omitempty"`
	IdSerial              string    `json:"id_serial,omitempty"`
	CreatedAt             time.Time `json:"created_at,omitempty"`
}

type Subscriber struct {
	SubscriberInfo *SubscriberInfo `json:"subscriber"`
}

type SubscriberClient interface {
	Get(Id string) (*SubscriberInfo, error)
}

type subscriberClient struct {
	R *rest.RestClient
}

func NewSubscriberClient(h string, debug bool) *subscriberClient {
	f, err := rest.NewRestClient(h, debug)
	if err != nil {
		log.Fatalf("Can't connect to %s url. Error %s", h, err.Error())
	}

	return &subscriberClient{
		R: f,
	}
}

func (s *subscriberClient) Get(id string) (*SubscriberInfo, error) {
	errStatus := &rest.ErrorMessage{}

	subscriber := Subscriber{}

	resp, err := s.R.C.R().
		SetError(errStatus).
		Get(s.R.URL.String() + SubscriberEndpoint + "/" + id)

	if err != nil {
		log.Errorf("Failed to send api request to subscriber/registry. Error %s", err.Error())

		return nil, fmt.Errorf("api request to subscriber system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to fetch subscriber info. HTTP resp code %d and Error message is %s",
			resp.StatusCode(), errStatus.Message)

		return nil, fmt.Errorf("Subscriber Info failure %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &subscriber)
	if err != nil {
		log.Tracef("Failed to deserialize subscriber info. Error message is %s",
			err.Error())

		return nil, fmt.Errorf("Subscriber info deserailization failure: %w", err)
	}

	log.Infof("Subscriber Info: %+v", subscriber)

	return subscriber.SubscriberInfo, nil
}
