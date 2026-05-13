/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package parser

import (
	"encoding/json"
	"time"
)

type HealthPayload struct {
	SchemaVersion string            `json:"schemaVersion"`
	NodeID        string            `json:"nodeId"`
	NodeType      string            `json:"nodeType"`
	ReportedAt    time.Time         `json:"reportedAt"`
	Capabilities  []string          `json:"capabilities"`
	System        HealthSystem      `json:"system"`
	Interfaces    HealthInterfaces  `json:"interfaces"`
	Apps          []HealthApp       `json:"apps"`
	Events        []json.RawMessage `json:"events"`
}

func ParseHealthPayload(raw json.RawMessage) (*HealthPayload, error) {
	if len(raw) == 0 {
		raw = []byte("{}")
	}
	var p HealthPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

type HealthSystem struct {
	UptimeSec int64           `json:"uptimeSec"`
	Starter   StarterStatus   `json:"starter"`
	Power     SystemPowerInfo `json:"power"`
}

type StarterStatus struct {
	Available          bool   `json:"available"`
	State              string `json:"state"`
	UpdateInProgress   bool   `json:"updateInProgress"`
	SwitchRequested    bool   `json:"switchRequested"`
	TerminateRequested bool   `json:"terminateRequested"`
	ExitCode           int    `json:"exitCode"`
}

type SystemPowerInfo struct {
	Available    bool    `json:"available"`
	Ok           bool    `json:"ok"`
	Board        string  `json:"board"`
	Reason       string  `json:"reason"`
	TotalWatts   float64 `json:"totalWatts"`
	TemperatureC float64 `json:"temperatureC"`
}

type HealthInterfaces struct {
	Cellular   *CellularInterface       `json:"cellular,omitempty"`
	Radio      *RadioInterface          `json:"radio,omitempty"`
	GPS        *GPSInterface            `json:"gps,omitempty"`
	Backhaul   *BackhaulInterface       `json:"backhaul,omitempty"`
	FEM        *FEMInterface            `json:"fem,omitempty"`
	Switch     *SwitchInterface         `json:"switch,omitempty"`
	Controller *NodeControllerInterface `json:"controller,omitempty"`
}

type CellularInterface struct {
	Available bool   `json:"available"`
	Error     string `json:"error"`
}

type RadioInterface struct {
	Available bool   `json:"available"`
	State     string `json:"state"`
}

type GPSInterface struct {
	Available   bool      `json:"available"`
	Lock        bool      `json:"lock"`
	Coordinates string    `json:"coordinates"`
	Time        time.Time `json:"time"`
}

type BackhaulInterface struct {
	Available  bool    `json:"available"`
	State      string  `json:"state"`
	LinkGuess  string  `json:"linkGuess"`
	Confidence float64 `json:"confidence"`
}

type FEMInterface struct {
	Available bool      `json:"available"`
	FEMs      []*FEMUnit `json:"fems"`
}

type FEMUnit struct {
	Unit    int     `json:"unit"`
	Present bool    `json:"present"`
	GPIO    FEMGPIO `json:"gpio"`
}

type FEMGPIO struct {
	TxRfEnable   bool `json:"txRfEnable"`
	RxRfEnable   bool `json:"rxRfEnable"`
	PaVdsEnable  bool `json:"paVdsEnable"`
	RfPalEnable  bool `json:"rfPalEnable"`
	Vds28vEnable bool `json:"vds28vEnable"`
	PsuPgood     bool `json:"psuPgood"`
}

type SwitchInterface struct {
	Available       bool         `json:"available"`
	Reachable       bool         `json:"reachable"`
	State           string       `json:"state"`
	Model           string       `json:"model"`
	SoftwareVersion string       `json:"softwareVersion"`
	PortCount       int          `json:"portCount"`
	Policy          SwitchPolicy `json:"policy"`
	Ports           []SwitchPort `json:"ports"`
}

type SwitchPolicy struct {
	State  string `json:"state"`
	Hash   string `json:"hash"`
	Source string `json:"source"`
	Error  string `json:"error"`
}

type SwitchPort struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Present        bool    `json:"present"`
	AdminState     string  `json:"adminState"`
	LinkState      string  `json:"linkState"`
	PoeState       string  `json:"poeState"`
	PoeOperational bool    `json:"poeOperational"`
	SpeedBps       int64   `json:"speedBps"`
	PowerWatts     float64 `json:"powerWatts"`
	Fault          string  `json:"fault"`
}

type NodeControllerInterface struct {
	Available        bool                     `json:"available"`
	CommOk           bool                     `json:"commOk"`
	ChargeState      string                   `json:"chargeState"`
	ErrorCode        int                      `json:"errorCode"`
	Error            string                   `json:"error"`
	ActiveAlarmCount int                      `json:"activeAlarmCount"`
	Solar            ControllerSolarMetrics   `json:"solar"`
	Battery          ControllerBatteryMetrics `json:"battery"`
	Load             ControllerLoadMetrics    `json:"load"`
}

type ControllerSolarMetrics struct {
	VoltageV float64 `json:"voltageV"`
	CurrentA float64 `json:"currentA"`
	PowerW   float64 `json:"powerW"`
}

type ControllerBatteryMetrics struct {
	VoltageV float64 `json:"voltageV"`
	CurrentA float64 `json:"currentA"`
	SocPct   int     `json:"socPct"`
}

type ControllerLoadMetrics struct {
	OutputOn bool    `json:"outputOn"`
	CurrentA float64 `json:"currentA"`
}

type HealthApp struct {
	Space     string       `json:"space"`
	Name      string       `json:"name"`
	Tag       string       `json:"tag"`
	Version   string       `json:"version"`
	State     string       `json:"state"`
	PID       int          `json:"pid"`
	Resources AppResources `json:"resources"`
}

type AppResources struct {
	CPUPercent     float64 `json:"cpuPercent"`
	MemoryRssKb    float64   `json:"memoryRssKb"`
	DiskReadBytes  float64   `json:"diskReadBytes"`
	DiskWriteBytes float64   `json:"diskWriteBytes"`
}
