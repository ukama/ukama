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

	"github.com/ukama/ukama/systems/common/rest/client"

	log "github.com/sirupsen/logrus"
)

const (
	UkamaSimsEndpoint = "/v1/asr"
)

type UkamaAgentClient interface {
	BindSim(iccid string) (*UkamaSimInfo, error)
	GetSimInfo(iccid string) (*UkamaSimInfo, error)
	GetUsages(iccid, cdrType, from, to, region string) (map[string]any, map[string]any, error)
	ActivateSim(req client.AgentRequestData) error
	DeactivateSim(req client.AgentRequestData) error
	UpdatePackage(req client.AgentRequestData) error
	TerminateSim(iccid string) error
}

type ukamaAgentClient struct {
	u *url.URL
	R *client.Resty
}

func NewUkamaAgentClient(h string) *ukamaAgentClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse  %s url. Error: %s", h, err.Error())
	}

	return &ukamaAgentClient{
		u: u,
		R: client.NewResty(),
	}
}

// Bind sim is a no-op for ukama sims for now
func (o *ukamaAgentClient) BindSim(iccid string) (*UkamaSimInfo, error) {
	return &UkamaSimInfo{}, nil
}

func (o *ukamaAgentClient) GetSimInfo(iccid string) (*UkamaSimInfo, error) {
	log.Debugf("Getting ukama sim info: %v", iccid)

	sim := UkamaSim{}

	resp, err := o.R.Get(o.u.String() + UkamaSimsEndpoint + "/" + iccid)
	if err != nil {
		log.Errorf("GetSimInfo failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetSimInfo failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &sim)
	if err != nil {
		log.Tracef("Failed to deserialize ukama sim info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("ukama sim info deserialization failure: %w", err)
	}

	log.Infof("Ukama Sim Info: %+v", sim)

	return sim.Record, nil
}

func (o *ukamaAgentClient) GetUsages(iccid, cdrType, from, to, region string) (map[string]any, map[string]any, error) {
	log.Debugf("Getting ukama sim usages: %v", iccid)

	usage := UkamaSimUsage{}

	resp, err := o.R.Get(o.u.String() + UkamaSimsEndpoint + "/usage/" +
		fmt.Sprintf("?iccid=%s&cdr_type=%s&from=%s&to=%s&region=%s", iccid, cdrType, from, to, region))
	if err != nil {
		log.Errorf("GetSim usages failure. error: %s", err.Error())

		return nil, nil, fmt.Errorf("GetSim usages failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &usage)
	if err != nil {
		log.Tracef("Failed to deserialize ukama sim info. Error message is: %s", err.Error())

		return nil, nil, fmt.Errorf("ukama sim info deserialization failure: %w", err)
	}

	log.Infof("ukama data usage (of type %T): %+v", usage.Usage, usage.Usage)
	log.Infof("ukama data cost (of type %T): %+v", usage.Cost, usage.Cost)

	return usage.Usage, usage.Cost, nil
}

func (o *ukamaAgentClient) ActivateSim(req client.AgentRequestData) error {
	log.Debugf("Activating ukama sim: %v", req.Iccid)

	_, err := o.R.C.R().SetBody(req).Put(o.u.String() + UkamaSimsEndpoint + "/" + req.Iccid)
	if err != nil {
		log.Errorf("ActivateSim failure. error: %s", err.Error())

		return fmt.Errorf("ActivateSim failure: %w", err)
	}

	return nil
}

func (o *ukamaAgentClient) DeactivateSim(req client.AgentRequestData) error {
	log.Debugf("Deactivating ukama sim: %v", req.Iccid)

	_, err := o.R.C.R().SetBody(req).Delete(o.u.String() + UkamaSimsEndpoint + "/" + req.Iccid)
	if err != nil {
		log.Errorf("DeactivateSim failure. error: %s", err.Error())

		return fmt.Errorf("DeactivateSim failure: %w", err)
	}

	return nil
}

func (o *ukamaAgentClient) UpdatePackage(req client.AgentRequestData) error {
	log.Debugf("Updating ukama sim's pacakge: %v", req.Iccid)

	_, err := o.R.C.R().SetBody(req).Patch(o.u.String() + UkamaSimsEndpoint + "/" + req.Iccid)
	if err != nil {
		log.Errorf("Update sim's package failure. error: %s", err.Error())

		return fmt.Errorf("update sim's package failure: %w", err)
	}

	return nil
}

func (o *ukamaAgentClient) TerminateSim(iccid string) error {
	log.Debugf("Terminating ukama sim: %v", iccid)

	return nil
}

type UkamaSim struct {
	Record *UkamaSimInfo `json:"sim"`
}

type UkamaSimUsage struct {
	Usage map[string]any `json:"usage"`
	Cost  map[string]any `json:"cost"`
}

type UkamaSimInfo struct {
	Imsi        string `protobuf:"bytes,1,opt,name=Imsi,json=imsi,proto3" json:"Imsi,omitempty"`
	SimId       string `protobuf:"bytes,2,opt,name=SimId,json=simId,proto3" json:"SimId,omitempty"`
	Iccid       string `protobuf:"bytes,3,opt,name=Iccid,json=iccid,proto3" json:"Iccid,omitempty"`
	Key         []byte `protobuf:"bytes,4,opt,name=Key,json=k,proto3" json:"Key,omitempty"`
	Op          []byte `protobuf:"bytes,5,opt,name=Op,json=op,proto3" json:"Op,omitempty"`
	Amf         []byte `protobuf:"bytes,6,opt,name=Amf,json=amf,proto3" json:"Amf,omitempty"`
	Apn         *Apn   `protobuf:"bytes,7,opt,name=Apn,json=apn,proto3" json:"Apn,omitempty"`
	AlgoType    uint32 `protobuf:"varint,8,opt,name=AlgoType,json=algo_type,proto3" json:"AlgoType,omitempty"`
	UeDlAmbrBps uint32 `protobuf:"varint,9,opt,name=UeDlAmbrBps,json=ue_dl_ambr_bps,proto3" json:"UeDlAmbrBps,omitempty"`
	UeUlAmbrBps uint32 `protobuf:"varint,10,opt,name=UeUlAmbrBps,json=ue_ul_ambr_bps,proto3" json:"UeUlAmbrBps,omitempty"`
	Sqn         uint64 `protobuf:"varint,11,opt,name=Sqn,json=sqn,proto3" json:"Sqn,omitempty"`
	CsgIdPrsent bool   `protobuf:"varint,12,opt,name=CsgIdPrsent,json=csg_id_prsent,proto3" json:"CsgIdPrsent,omitempty"`
	CsgId       uint32 `protobuf:"varint,13,opt,name=CsgId,json=csg_id,proto3" json:"CsgId,omitempty"`
	PackageId   string `protobuf:"bytes,14,opt,name=PackageId,json=package_id,proto3" json:"PackageId,omitempty"`
}

type Apn struct {
	Name string `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
}
