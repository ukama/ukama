/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package subscriber

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
)

const (
	SubscriberEndpoint           = "/v1/subscriber"
	SubscribersByNetworkEndpoint = "/v1/subscribers/networks"
)

type SubscriberInfo struct {
	SubscriberId          uuid.UUID `json:"subscriber_id,omitempty"`
	OrgId                 uuid.UUID `json:"org_id,omitempty"`
	NetworkId             uuid.UUID `json:"network_id,omitempty"`
	Name                  string    `json:"name,omitempty"`
	Email                 string    `json:"email,omitempty"`
	PhoneNumber           string    `json:"phone_number,omitempty"`
	Address               string    `json:"address,omitempty"`
	Dob                   string    `json:"dob,omitempty"`
	ProofOfIdentification string    `json:"proof_of_identification,omitempty"`
	IdSerial              string    `json:"id_serial,omitempty"`
	CreatedAt             time.Time `json:"created_at,omitempty"`
}

type Subscriber struct {
	SubscriberInfo *SubscriberInfo `json:"subscriber"`
}

type SubscriberSimPackage struct {
	Id        string    `json:"id,omitempty"`
	PackageId string    `json:"package_id,omitempty"`
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
	IsActive  bool      `json:"is_active,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type SubscriberNetworkSim struct {
	Id                 string                `json:"id,omitempty"`
	SubscriberId       string                `json:"subscriber_id,omitempty"`
	NetworkId          string                `json:"network_id,omitempty"`
	Package            *SubscriberSimPackage `json:"package,omitempty"`
	Iccid              string                `json:"iccid,omitempty"`
	Msisdn             string                `json:"msisdn,omitempty"`
	Imsi               string                `json:"imsi,omitempty"`
	Type               string                `json:"type,omitempty"`
	Status             string                `json:"status,omitempty"`
	IsPhysical         bool                  `json:"is_physical,omitempty"`
	FirstActivatedOn   time.Time             `json:"first_activated_on,omitempty"`
	LastActivatedOn    time.Time             `json:"last_activated_on,omitempty"`
	ActivationsCount   string                `json:"activations_count,omitempty"`
	DeactivationsCount string                `json:"deactivations_count,omitempty"`
	AllocatedAt        time.Time             `json:"allocated_at,omitempty"`
}

type SubscriberNetworkInfo struct {
	Name                  string                  `json:"name,omitempty"`
	SubscriberId          string                  `json:"subscriber_id,omitempty"`
	NetworkId             string                  `json:"network_id,omitempty"`
	Email                 string                  `json:"email,omitempty"`
	PhoneNumber           string                  `json:"phone_number,omitempty"`
	Address               string                  `json:"address,omitempty"`
	ProofOfIdentification string                  `json:"proof_of_identification,omitempty"`
	CreatedAt             string                  `json:"createdAt,omitempty"`
	DeletedAt             string                  `json:"deletedAt,omitempty"`
	UpdatedAt             string                  `json:"updatedAt,omitempty"`
	Sim                   []*SubscriberNetworkSim `json:"sim,omitempty"`
	Dob                   string                  `json:"dob,omitempty"`
	IdSerial              string                  `json:"id_serial,omitempty"`
	Gender                string                  `json:"gender,omitempty"`
}

type SubscribersByNetworkResponse struct {
	Subscribers []*SubscriberNetworkInfo `json:"subscribers"`
}

type AddSubscriberRequest struct {
	OrgId                 string `json:"org_id" validate:"required"`
	NetworkId             string `json:"network_id" validate:"required"`
	Name                  string `json:"name,omitempty"`
	Email                 string `json:"email,omitempty"`
	PhoneNumber           string `json:"phone_number,omitempty"`
	Address               string `json:"address,omitempty"`
	Dob                   string `json:"dob,omitempty"`
	ProofOfIdentification string `json:"proof_of_identification,omitempty"`
	IdSerial              string `json:"id_serial,omitempty"`
}

type SubscriberClient interface {
	Get(id string) (*SubscriberInfo, error)
	GetByEmail(email string) (*SubscriberInfo, error)
	Add(req AddSubscriberRequest) (*SubscriberInfo, error)
	GetByNetwork(networkId string) ([]*SubscriberNetworkInfo, error)
}

type subscriberClient struct {
	u *url.URL
	R *client.Resty
}

func NewSubscriberClient(h string, options ...client.Option) *subscriberClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &subscriberClient{
		u: u,
		R: client.NewResty(options...),
	}
}

// TODO check upstream returns payload
func (s *subscriberClient) Add(req AddSubscriberRequest) (*SubscriberInfo, error) {
	log.Debugf("Adding subscriber: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %w", err)
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

		return nil, fmt.Errorf("subscriber info deserialization failure: %w", err)
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

		return nil, fmt.Errorf("subscriber info deserialization failure: %w", err)
	}

	log.Infof("Subscriber Info: %+v", subscriber.SubscriberInfo)

	return subscriber.SubscriberInfo, nil
}

func (s *subscriberClient) GetByEmail(email string) (*SubscriberInfo, error) {
	log.Debugf("Getting subscriber: %v", email)

	subscriber := Subscriber{}

	resp, err := s.R.Get(s.u.String() + SubscriberEndpoint + "/email/" + email)
	if err != nil {
		log.Errorf("GetSubscriber failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetSubscriber failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &subscriber)
	if err != nil {
		log.Tracef("Failed to deserialize subscriber info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("subscriber info deserialization failure: %w", err)
	}

	log.Infof("Subscriber Info: %+v", subscriber.SubscriberInfo)

	return subscriber.SubscriberInfo, nil
}

func (s *subscriberClient) GetByNetwork(networkId string) ([]*SubscriberNetworkInfo, error) {
	log.Debugf("Getting subscribers by network: %v", networkId)

	out := SubscribersByNetworkResponse{}

	resp, err := s.R.Get(s.u.String() + SubscribersByNetworkEndpoint + "/" + networkId)
	if err != nil {
		log.Errorf("GetByNetwork failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetByNetwork failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &out)
	if err != nil {
		log.Tracef("Failed to deserialize subscribers by network. Error message is: %s", err.Error())

		return nil, fmt.Errorf("subscribers by network deserialization failure: %w", err)
	}

	log.Infof("Subscribers by network count: %d", len(out.Subscribers))

	return out.Subscribers, nil
}
