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
	"time"

	"github.com/ukama/ukama/systems/common/rest/client"

	log "github.com/sirupsen/logrus"
)

const HealthEndpoint = "/v1/health"

type listInterfacesResponse struct {
	Interfaces InterfaceInfo `json:"interfaces,omitempty"`
}

type InterfaceInfo struct {
	Cellular   *CellularInterfaceInfo       `json:"cellular,omitempty"`
	Radio      *RadioInterfaceInfo          `json:"radio,omitempty"`
	Gps        *GPSInterfaceInfo            `json:"gps,omitempty"`
	Backhaul   *BackhaulInterfaceInfo       `json:"backhaul,omitempty"`
	Fem        *FEMInterfaceInfo            `json:"fem,omitempty"`
	Switch     *SwitchInterfaceInfo         `json:"switch,omitempty"`
	Controller *NodeControllerInterfaceInfo `json:"controller,omitempty"`
}

type CellularInterfaceInfo struct {
	Available bool   `json:"available,omitempty"`
	Error     string `json:"error,omitempty"`
}

type RadioInterfaceInfo struct {
	Available bool   `json:"available,omitempty"`
	State     string `json:"state,omitempty"`
}

type GPSInterfaceInfo struct {
	Available   bool      `json:"available,omitempty"`
	Lock        bool      `json:"lock,omitempty"`
	Coordinates string    `json:"coordinates"`
	Time        time.Time `json:"time,omitempty"`
}

type BackhaulInterfaceInfo struct {
	Available  bool    `json:"available,omitempty"`
	State      string  `json:"state,omitempty"`
	LinkGuess  string  `json:"linkGuess,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
}

type FEMInterfaceInfo struct {
	Available bool          `json:"available,omitempty"`
	Fems      []FEMUnitInfo `json:"fems,omitempty"`
}

type FEMUnitInfo struct {
	Unit    int32        `json:"unit,omitempty"`
	Present bool         `json:"present,omitempty"`
	Gpio    *FEMGPIOInfo `json:"gpio,omitempty"`
}

type FEMGPIOInfo struct {
	TxRfEnable   bool `json:"txRfEnable,omitempty"`
	RxRfEnable   bool `json:"rxRfEnable,omitempty"`
	PaVdsEnable  bool `json:"paVdsEnable,omitempty"`
	RfPalEnable  bool `json:"rfPalEnable,omitempty"`
	Vds28vEnable bool `json:"vds28vEnable,omitempty"`
	PsuPgood     bool `json:"psuPgood,omitempty"`
}

type SwitchInterfaceInfo struct {
	Available       bool                       `json:"available,omitempty"`
	Reachable       bool                       `json:"reachable,omitempty"`
	State           string                     `json:"state,omitempty"`
	Model           string                     `json:"model,omitempty"`
	SoftwareVersion string                     `json:"softwareVersion,omitempty"`
	PortCount       int32                      `json:"portCount,omitempty"`
	Policy          *SwitchInterfacePolicyInfo `json:"policy,omitempty"`
	Ports           []SwitchPortInfo           `json:"ports,omitempty"`
}

type SwitchInterfacePolicyInfo struct {
	State  string `json:"state,omitempty"`
	Hash   string `json:"hash,omitempty"`
	Source string `json:"source,omitempty"`
	Error  string `json:"error,omitempty"`
}

type SwitchPortInfo struct {
	Id             int64   `json:"id,omitempty"`
	Name           string  `json:"name,omitempty"`
	Present        bool    `json:"present,omitempty"`
	AdminState     string  `json:"adminState,omitempty"`
	LinkState      string  `json:"linkState,omitempty"`
	PoeState       string  `json:"poeState,omitempty"`
	PoeOperational bool    `json:"poeOperational,omitempty"`
	SpeedBps       int64   `json:"speedBps,omitempty"`
	PowerWatts     float64 `json:"powerWatts,omitempty"`
	Fault          string  `json:"fault,omitempty"`
}

type NodeControllerInterfaceInfo struct {
	Available        bool                          `json:"available,omitempty"`
	CommOk           bool                          `json:"commOk,omitempty"`
	ChargeState      string                        `json:"chargeState,omitempty"`
	ErrorCode        int32                         `json:"errorCode,omitempty"`
	Error            string                        `json:"error,omitempty"`
	ActiveAlarmCount int32                         `json:"activeAlarmCount,omitempty"`
	Solar            *ControllerSolarMetricsInfo   `json:"solar,omitempty"`
	Battery          *ControllerBatteryMetricsInfo `json:"battery,omitempty"`
	Load             *ControllerLoadMetricsInfo    `json:"load,omitempty"`
}

type ControllerSolarMetricsInfo struct {
	VoltageV float64 `json:"voltageV,omitempty"`
	CurrentA float64 `json:"currentA,omitempty"`
	PowerW   float64 `json:"powerW,omitempty"`
}

type ControllerBatteryMetricsInfo struct {
	VoltageV float64 `json:"voltageV,omitempty"`
	CurrentA float64 `json:"currentA,omitempty"`
	SocPct   int32   `json:"socPct,omitempty"`
}

type ControllerLoadMetricsInfo struct {
	OutputOn bool    `json:"outputOn,omitempty"`
	CurrentA float64 `json:"currentA,omitempty"`
}

type NodeHealthClient interface {
	GetInterfaces(interfaceName, nodeId, reportId string) (InterfaceInfo, error)
}

type nodeHealthClient struct {
	u *url.URL
	R *client.Resty
}

func NewNodeHealthClient(h string, options ...client.Option) *nodeHealthClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &nodeHealthClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (h *nodeHealthClient) GetInterfaces(interfaceName, nodeId, reportId string) (InterfaceInfo, error) {
	log.Debugf("Getting interfaces: interfaceName=%q nodeId=%q reportId=%q", interfaceName, nodeId, reportId)

	q := url.Values{}
	if reportId != "" {
		q.Set("reportId", reportId)
	}
	if nodeId != "" {
		q.Set("nodeId", nodeId)
	}
	if interfaceName != "" {
		q.Set("interfaceName", interfaceName)
	}

	resp, err := h.R.GetWithQuery(h.u.String()+HealthEndpoint+"/interfaces", q.Encode())
	if err != nil {
		log.Errorf("GetInterfaces failure. error: %s", err.Error())

		return InterfaceInfo{}, fmt.Errorf("GetInterfaces failure: %w", err)
	}

	var out listInterfacesResponse
	err = json.Unmarshal(resp.Body(), &out)
	if err != nil {
		log.Tracef("Failed to deserialize interfaces. Error message is: %s", err.Error())

		return InterfaceInfo{}, fmt.Errorf("interfaces deserialization failure: %w", err)
	}

	log.Infof("Interfaces: %+v", out.Interfaces)

	return out.Interfaces, nil
}
