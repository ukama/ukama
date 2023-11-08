/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package billing

import (
	"encoding/json"
	"fmt"
	"net/url"

	bilutil "github.com/ukama/ukama/systems/billing/invoice/pkg/util"
	"github.com/ukama/ukama/testing/integration/pkg/utils"

	log "github.com/sirupsen/logrus"
)

type BillingClient struct {
	u *url.URL
	r utils.Resty
}

func NewBillingClient(h, k string) *BillingClient {
	u, _ := url.Parse(h)

	return &BillingClient{
		u: u,
		r: *utils.NewRestyWithBearer(k),
	}
}

func (s *BillingClient) GetCustomer(custId string) (*bilutil.Customer, error) {
	var rsp = &struct {
		Customer *bilutil.Customer
	}{}

	resp, err := s.r.Get(s.u.String() + "/api/v1/customers/" + custId)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("GetCustomer failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("GetCustomer: response unmarshal error. error: %w", err)
	}

	return rsp.Customer, nil
}

func (s *BillingClient) GetPlan(planCode string) (*bilutil.Plan, error) {
	var rsp = &struct {
		Plan *bilutil.Plan
	}{}

	resp, err := s.r.Get(s.u.String() + "/api/v1/plans/" + planCode)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("GetCustomer failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("GetCustomer: response unmarshal error. error: %w", err)
	}

	return rsp.Plan, nil
}

func (s *BillingClient) GetSubscriptionsByCustomerId(custId string) (*bilutil.Subscription, error) {
	var rsp = &struct {
		Subscriptions []*bilutil.Subscription
	}{}

	resp, err := s.r.Get(s.u.String() + "/api/v1/subscriptions?external_customer_id=" + custId)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("GetSubscriptions failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("GetSubscriptions: response unmarshal error. error: %w", err)
	}

	return rsp.Subscriptions[0], nil
}
