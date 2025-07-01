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

const SubscriberEndpoint = "/v1/subscriber"

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
