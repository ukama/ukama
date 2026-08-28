/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package ukamaagent

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/validation"

	log "github.com/sirupsen/logrus"
)

const (
	UkamaSimsEndpoint  = "/v1/asr"
	UkamaUsageEndpoint = "/v1/usage"

	defaultStartTime = 1
)

type UkamaAgentClient interface {
	CreateProfile(req client.AgentRequestData) (*UkamaSimInfo, error)
	GetSimInfo(iccid string) (*UkamaSimInfo, error)
	GetUsages(iccid, cdrType, from, to, region string) (map[string]any, map[string]any, error)
	UpdatePackage(req client.AgentRequestData) error
	Update(req client.AgentRequestData) error
	DeleteProfile(req client.AgentRequestData) error
}

type ukamaAgentClient struct {
	u *url.URL
	R *client.Resty
}

func NewUkamaAgentClient(h string) *ukamaAgentClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &ukamaAgentClient{
		u: u,
		R: client.NewResty(),
	}
}

func (u *ukamaAgentClient) CreateProfile(req client.AgentRequestData) (*UkamaSimInfo, error) {
	log.Debugf("Creating profile for ukama sim: %v", req.Iccid)

	_, err := u.R.C.R().SetBody(req).Put(u.u.String() + UkamaSimsEndpoint + "/" + req.Iccid)
	if err != nil {
		log.Errorf("CreateProfile failure. error: %v", err)

		return nil, fmt.Errorf("CreateProfile failure: %w", err)
	}

	return &UkamaSimInfo{}, nil
}

func (u *ukamaAgentClient) GetSimInfo(iccid string) (*UkamaSimInfo, error) {
	log.Debugf("Getting ukama sim info: %v", iccid)

	sim := UkamaSim{}

	resp, err := u.R.Get(u.u.String() + UkamaSimsEndpoint + "/" + iccid)
	if err != nil {
		log.Errorf("GetSimInfo failure. error: %v", err)

		return nil, fmt.Errorf("GetSimInfo failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &sim)
	if err != nil {
		log.Tracef("Failed to deserialize ukama sim info. Error message is: %v", err)

		return nil, fmt.Errorf("ukama sim info deserialization failure: %w", err)
	}

	log.Infof("Ukama Sim Info: %+v", sim)

	return sim.Record, nil
}

func (u *ukamaAgentClient) GetUsages(iccid, cdrType, from, to, region string) (map[string]any, map[string]any, error) {
	log.Debugf("Getting ukama sim usages: %v", iccid)

	if iccid == "" {
		return nil, nil, fmt.Errorf("ukama agent does not support multi-sim usage queries: iccid is required")
	}

	var startTime int64 = defaultStartTime
	var endTime int64 = time.Now().Add(time.Hour).Unix()

	usage := UkamaSimUsage{}

	if from != "" {
		frm, err := validation.FromString(from)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid format for from: %s. Error: %w", from, err)
		}
		startTime = frm.Unix()
	}

	if to != "" {
		t, err := validation.FromString(to)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid format for to: %s. Error: %w", to, err)
		}
		endTime = t.Unix()
	}

	resp, err := u.R.Get(u.u.String() + UkamaUsageEndpoint +
		fmt.Sprintf("/%s?from=%d&to=%d", iccid, startTime, endTime))
	if err != nil {
		log.Errorf("GetSim usages failure. error: %v", err)

		return nil, nil, fmt.Errorf("GetSim usages failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &usage)
	if err != nil {
		log.Tracef("Failed to deserialize ukama sim info. Error message is: %v", err)

		return nil, nil, fmt.Errorf("ukama sim info deserialization failure: %w", err)
	}

	log.Infof("ukama data usage (of type %T): %+v", usage.Usage, usage.Usage)

	return map[string]any{iccid: usage.Usage}, nil, nil
}

func (u *ukamaAgentClient) UpdatePackage(req client.AgentRequestData) error {
	log.Debugf("Updating ukama sim's package: %v", req.Iccid)

	_, err := u.R.C.R().SetBody(req).Patch(u.u.String() + UkamaSimsEndpoint + "/" + req.Iccid + "/package")
	if err != nil {
		log.Errorf("Update sim's package failure. error: %v", err)

		return fmt.Errorf("update sim's package failure: %w", err)
	}

	return nil
}

func (u *ukamaAgentClient) Update(req client.AgentRequestData) error {
	log.Debugf("Updating ukama sim's service status: %v", req.Iccid)

	_, err := u.R.C.R().SetBody(req).Patch(u.u.String() + UkamaSimsEndpoint + "/" + req.Iccid)
	if err != nil {
		log.Errorf("Update sim's service status failure. error: %v", err)

		return fmt.Errorf("update sim's service status failure: %w", err)
	}

	return nil
}

func (u *ukamaAgentClient) DeleteProfile(req client.AgentRequestData) error {
	log.Debugf("Deleting profile for ukama sim: %v", req.Iccid)

	_, err := u.R.C.R().SetBody(req).Delete(u.u.String() + UkamaSimsEndpoint + "/" + req.Iccid)
	if err != nil {
		log.Errorf("DeleteProfile failure. error: %v", err)

		return fmt.Errorf("DeleteProfile failure: %w", err)
	}

	return nil
}

type UkamaSimUsage struct {
	Usage uint64 `json:"usage,string"`
}

type UkamaSim struct {
	Record *UkamaSimInfo `json:"Record"`
}

type UkamaSimInfo struct {
	Imsi              string  `protobuf:"bytes,1,opt,name=Imsi,json=imsi,proto3" json:"Imsi,omitempty"`
	Iccid             string  `protobuf:"bytes,3,opt,name=Iccid,json=iccid,proto3" json:"Iccid,omitempty"`
	Key               []byte  `protobuf:"bytes,4,opt,name=Key,json=k,proto3" json:"Key,omitempty"`
	Op                []byte  `protobuf:"bytes,5,opt,name=Op,json=op,proto3" json:"Op,omitempty"`
	Amf               []byte  `protobuf:"bytes,6,opt,name=Amf,json=amf,proto3" json:"Amf,omitempty"`
	Apn               *Apn    `protobuf:"bytes,7,opt,name=Apn,json=apn,proto3" json:"Apn,omitempty"`
	AlgoType          uint32  `protobuf:"varint,8,opt,name=AlgoType,json=algo_type,proto3" json:"AlgoType,omitempty"`
	UeDlAmbrBps       uint32  `protobuf:"varint,9,opt,name=UeDlAmbrBps,json=ue_dl_ambr_bps,proto3" json:"UeDlAmbrBps,omitempty"`
	UeUlAmbrBps       uint32  `protobuf:"varint,10,opt,name=UeUlAmbrBps,json=ue_ul_ambr_bps,proto3" json:"UeUlAmbrBps,omitempty"`
	Sqn               string  `protobuf:"varint,11,opt,name=Sqn,json=sqn,proto3" json:"Sqn,omitempty"`
	CsgIdPrsent       bool    `protobuf:"varint,12,opt,name=CsgIdPrsent,json=csg_id_prsent,proto3" json:"CsgIdPrsent,omitempty"`
	CsgId             uint32  `protobuf:"varint,13,opt,name=CsgId,json=csg_id,proto3" json:"CsgId,omitempty"`
	PackageId         string  `protobuf:"bytes,14,opt,name=PackageId,json=package_id,proto3" json:"PackageId,omitempty"`
	Policy            *Policy `protobuf:"bytes,17,opt,name=Policy,json=policy,proto3" json:"Policy,omitempty"`
	IsServiceStatusOn bool    `protobuf:"varint,18,opt,name=IsServiceStatusOn,json=is_service_status_on,proto3" json:"IsServiceStatusOn,omitempty"`
}

type Apn struct {
	Name string `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
}

type Policy struct {
	UUID         string `protobuf:"bytes,1,opt,name=UUID,json=uuid,proto3" json:"uuid,omitempty"`
	Burst        string `protobuf:"bytes,2,opt,name=Burst,json=burst,proto3" json:"burst,omitempty"`
	TotalData    string `protobuf:"bytes,3,opt,name=TotalData,json=total_data,proto3" json:"total_data,omitempty"`
	ConsumedData string `protobuf:"bytes,4,opt,name=ConsumedData,json=consumed_data,proto3" json:"consumed_data,omitempty"`
	ULBR         string `protobuf:"bytes,5,opt,name=ULBR,json=ulbr,proto3" json:"ulbr,omitempty"`
	DLBR         string `protobuf:"bytes,6,opt,name=DLBR,json=dlbr,proto3" json:"dlbr,omitempty"`
	StartTime    string `protobuf:"bytes,7,opt,name=StartTime,json=start_time,proto3" json:"start_time,omitempty"`
	EndTime      string `protobuf:"bytes,8,opt,name=EndTime,json=end_time,proto3" json:"end_time,omitempty"`
}
