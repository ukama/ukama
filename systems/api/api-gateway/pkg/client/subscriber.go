package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

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

type AddSubscriberRequest struct {
	OrgId                 string `json:"org_id" validate:"required"`
	NetworkId             string `json:"network_id" validate:"required"`
	FirstName             string `json:"first_name,omitempty"`
	LastName              string `json:"last_name,omitempty"`
	Email                 string `json:"email,omitempty"`
	PhoneNumber           string `json:"phone_number,omitempty"`
	Address               string `json:"address,omitempty"`
	Dob                   string `json:"date_of_birth,omitempty"`
	ProofOfIdentification string `json:"proof_of_identification,omitempty"`
	IdSerial              string `json:"id_serial,omitempty"`
}

type SubscriberClient interface {
	Get(Id string) (*SubscriberInfo, error)
	Add(req AddSubscriberRequest) (*SubscriberInfo, error)
}

type subscriberClient struct {
	u *url.URL
	R *Resty
}

func NewSubscriberClient(h string) *subscriberClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse  %s url. Error: %s", h, err.Error())
	}

	return &subscriberClient{
		u: u,
		R: NewResty(),
	}
}

func (s *subscriberClient) Add(req AddSubscriberRequest) (*SubscriberInfo, error) {
	log.Debugf("Adding subscriber: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	subscriber := Subscriber{}

	resp, err := s.R.Post(s.u.String()+SubscriberEndpoint, b)
	if err != nil {
		log.Errorf("AddSubscriber failure. error: %s", err.Error())

		return nil, fmt.Errorf("AddSubscriber failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &subscriber)
	if err != nil {
		log.Tracef("Failed to deserialize subscriber info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("subscriber info deserailization failure: %w", err)
	}

	log.Infof("Subscriber Info: %+v", subscriber.SubscriberInfo)

	return subscriber.SubscriberInfo, nil
}

func (s *subscriberClient) Get(id string) (*SubscriberInfo, error) {
	log.Debugf("Getting subscriber: %v", id)

	subscriber := Subscriber{}

	resp, err := s.R.Get(s.u.String() + SubscriberEndpoint + "/" + id)
	if err != nil {
		log.Errorf("GetSubscriber failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetSubscriber failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &subscriber)
	if err != nil {
		log.Tracef("Failed to deserialize subscriber info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("subscriber info deserailization failure: %w", err)
	}

	log.Infof("Subscriber Info: %+v", subscriber.SubscriberInfo)

	return subscriber.SubscriberInfo, nil
}
