/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package node

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ukama/ukama/systems/common/rest/client"

	log "github.com/sirupsen/logrus"
)

const ControllerEndpoint = "/v1/controller"

type RestartNodeResponse struct {
	OperationId string `json:"operation_id"`
	ResourceKey string `json:"resource_key"`
	Status      string `json:"status"`
}

type ToggleSwitchPortRequest struct {
	Status bool  `json:"status"`
	Port   int32 `json:"port"`
}

type ToggleSwitchPortResponse struct {
	OperationId string `json:"operation_id"`
	ResourceKey string `json:"resource_key"`
	Status      string `json:"status"`
}

type ToggleRadioResponse struct {
	OperationId string `json:"operation_id"`
	ResourceKey string `json:"resource_key"`
	Status      string `json:"status"`
}

type ToggleServiceResponse struct {
	OperationId string `json:"operation_id"`
	ResourceKey string `json:"resource_key"`
	Status      string `json:"status"`
}

type ConfigNodeResponse struct {
	OperationId string `json:"operation_id"`
	ResourceKey string `json:"resource_key"`
	Status      string `json:"status"`
}

type DeleteNodeConfigResponse struct {
	OperationId string `json:"operation_id"`
	ResourceKey string `json:"resource_key"`
	Status      string `json:"status"`
}

type NodeControllerClient interface {
	RestartNode(nodeId string) (*RestartNodeResponse, error)
	ToggleSwitchPort(nodeId string, req ToggleSwitchPortRequest) (*ToggleSwitchPortResponse, error)
	ToggleRadio(nodeId, state string) (*ToggleRadioResponse, error)
	ToggleService(nodeId, state string) (*ToggleServiceResponse, error)
	ConfigNode(nodeId string) (*ConfigNodeResponse, error)
	DeleteNodeConfig(nodeId string) (*DeleteNodeConfigResponse, error)
}

type nodeControllerClient struct {
	u *url.URL
	R *client.Resty
}

func NewNodeControllerClient(h string, options ...client.Option) *nodeControllerClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &nodeControllerClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (c *nodeControllerClient) RestartNode(nodeId string) (*RestartNodeResponse, error) {
	log.Debugf("Restarting node: %v", nodeId)

	resp, err := c.R.Post(c.u.String()+ControllerEndpoint+"/nodes/"+nodeId+"/restart", nil)
	if err != nil {
		log.Errorf("RestartNode failure. error: %s", err.Error())

		return nil, fmt.Errorf("RestartNode failure: %w", err)
	}

	out := &RestartNodeResponse{}

	err = json.Unmarshal(resp.Body(), out)
	if err != nil {
		log.Tracef("Failed to deserialize restart node response. Error message is: %s", err.Error())

		return nil, fmt.Errorf("restart node response deserialization failure: %w", err)
	}

	log.Infof("RestartNode: %+v", out)

	return out, nil
}

func (c *nodeControllerClient) ToggleSwitchPort(nodeId string, req ToggleSwitchPortRequest) (*ToggleSwitchPortResponse, error) {
	log.Debugf("Toggling switch port for node %v: %v", nodeId, req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %w", err)
	}

	resp, err := c.R.Post(c.u.String()+ControllerEndpoint+"/nodes/"+nodeId+"/switch-port", b)
	if err != nil {
		log.Errorf("ToggleSwitchPort failure. error: %s", err.Error())

		return nil, fmt.Errorf("ToggleSwitchPort failure: %w", err)
	}

	out := &ToggleSwitchPortResponse{}

	err = json.Unmarshal(resp.Body(), out)
	if err != nil {
		log.Tracef("Failed to deserialize toggle switch port response. Error message is: %s", err.Error())

		return nil, fmt.Errorf("toggle switch port response deserialization failure: %w", err)
	}

	log.Infof("ToggleSwitchPort: %+v", out)

	return out, nil
}

func (c *nodeControllerClient) ToggleRadio(nodeId, state string) (*ToggleRadioResponse, error) {
	log.Debugf("Toggling radio for node %v to %v", nodeId, state)

	resp, err := c.R.Post(c.u.String()+ControllerEndpoint+"/nodes/"+nodeId+"/radio/"+state, nil)
	if err != nil {
		log.Errorf("ToggleRadio failure. error: %s", err.Error())

		return nil, fmt.Errorf("ToggleRadio failure: %w", err)
	}

	out := &ToggleRadioResponse{}

	err = json.Unmarshal(resp.Body(), out)
	if err != nil {
		log.Tracef("Failed to deserialize toggle radio response. Error message is: %s", err.Error())

		return nil, fmt.Errorf("toggle radio response deserialization failure: %w", err)
	}

	log.Infof("ToggleRadio: %+v", out)

	return out, nil
}

func (c *nodeControllerClient) ToggleService(nodeId, state string) (*ToggleServiceResponse, error) {
	log.Debugf("Toggling service for node %v to %v", nodeId, state)

	resp, err := c.R.Post(c.u.String()+ControllerEndpoint+"/nodes/"+nodeId+"/service/"+state, nil)
	if err != nil {
		log.Errorf("ToggleService failure. error: %s", err.Error())

		return nil, fmt.Errorf("ToggleService failure: %w", err)
	}

	out := &ToggleServiceResponse{}

	err = json.Unmarshal(resp.Body(), out)
	if err != nil {
		log.Tracef("Failed to deserialize toggle service response. Error message is: %s", err.Error())

		return nil, fmt.Errorf("toggle service response deserialization failure: %w", err)
	}

	log.Infof("ToggleService: %+v", out)

	return out, nil
}

func (c *nodeControllerClient) ConfigNode(nodeId string) (*ConfigNodeResponse, error) {
	log.Debugf("Sending config to node: %v", nodeId)

	resp, err := c.R.Post(c.u.String()+ControllerEndpoint+"/nodes/"+nodeId+"/config", nil)
	if err != nil {
		log.Errorf("ConfigNode failure. error: %s", err.Error())

		return nil, fmt.Errorf("ConfigNode failure: %w", err)
	}

	out := &ConfigNodeResponse{}

	err = json.Unmarshal(resp.Body(), out)
	if err != nil {
		log.Tracef("Failed to deserialize config node response. Error message is: %s", err.Error())

		return nil, fmt.Errorf("config node response deserialization failure: %w", err)
	}

	log.Infof("ConfigNode: %+v", out)

	return out, nil
}

func (c *nodeControllerClient) DeleteNodeConfig(nodeId string) (*DeleteNodeConfigResponse, error) {
	log.Debugf("Deleting config from node: %v", nodeId)

	resp, err := c.R.Delete(c.u.String() + ControllerEndpoint + "/nodes/" + nodeId + "/config")
	if err != nil {
		log.Errorf("DeleteNodeConfig failure. error: %s", err.Error())

		return nil, fmt.Errorf("DeleteNodeConfig failure: %w", err)
	}

	out := &DeleteNodeConfigResponse{}

	err = json.Unmarshal(resp.Body(), out)
	if err != nil {
		log.Tracef("Failed to deserialize delete node config response. Error message is: %s", err.Error())

		return nil, fmt.Errorf("delete node config response deserialization failure: %w", err)
	}

	log.Infof("DeleteNodeConfig: %+v", out)

	return out, nil
}
