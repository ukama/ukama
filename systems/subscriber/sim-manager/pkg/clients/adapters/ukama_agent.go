/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package adapters

import (
	"encoding/json"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest/client"
)

const (
	UkamaOperatorSimsEndpoint = "/v1/asr"
)

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

// Bind sim is a no-op for operator sims for now
func (o *ukamaAgentClient) BindSim(iccid string) (*OperatorSimInfo, error) {
	return &OperatorSimInfo{}, nil
}

func (o *ukamaAgentClient) GetSimInfo(iccid string) (*OperatorSimInfo, error) {
	log.Debugf("Getting operator sim info: %v", iccid)

	sim := UkamaSimInfo{}

	resp, err := o.R.Get(o.u.String() + UkamaOperatorSimsEndpoint + "/" + iccid)
	if err != nil {
		log.Errorf("GetSimInfo failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetSimInfo failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &sim)
	if err != nil {
		log.Tracef("Failed to deserialize operator sim info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("operator sim info deserailization failure: %w", err)
	}

	log.Infof("Operator Ukama Sim Info: %+v", sim)

	return &OperatorSimInfo{
		Iccid: sim.Iccid,
		Imsi:  sim.Imsi,
	}, nil
}

func (o *ukamaAgentClient) GetUsages(iccid, cdrType, from, to, region string) (map[string]any, map[string]any, error) {
	log.Debugf("Getting operator sim info: %v", iccid)

	usage := OperatorUsage{}

	resp, err := o.R.Get(o.u.String() + UkamaOperatorSimsEndpoint + "/usage/" +
		fmt.Sprintf("?iccid=%s&cdr_type=%s&from=%s&to=%s&region=%s", iccid, cdrType, from, to, region))
	if err != nil {
		log.Errorf("GetSimInfo failure. error: %s", err.Error())

		return nil, nil, fmt.Errorf("GetSimInfo failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &usage)
	if err != nil {
		log.Tracef("Failed to deserialize operator sim info. Error message is: %s", err.Error())

		return nil, nil, fmt.Errorf("operator sim info deserailization failure: %w", err)
	}

	log.Infof("Operator data usage (of type %T): %+v", usage.Usage, usage.Usage)
	log.Infof("Operator data cost (of type %T): %+v", usage.Cost, usage.Cost)

	return usage.Usage, usage.Cost, nil
}

func (o *ukamaAgentClient) ActivateSim(req ReqData) error {
	log.Debugf("Activationg operator sim: %v", req.Iccid)

	_, err := o.R.C.R().SetBody(req).Put(o.u.String() + UkamaOperatorSimsEndpoint + "/" + req.Iccid)
	if err != nil {
		log.Errorf("ActivateSim failure. error: %s", err.Error())

		return fmt.Errorf("ActivateSim failure: %w", err)
	}

	return nil
}

func (o *ukamaAgentClient) DeactivateSim(req ReqData) error {
	log.Debugf("Deactivating operator sim: %v", req.Iccid)

	_, err := o.R.C.R().SetBody(req).Delete(o.u.String() + UkamaOperatorSimsEndpoint + "/" + req.Iccid)
	if err != nil {
		log.Errorf("DeactivateSim failure. error: %s", err.Error())

		return fmt.Errorf("DeactivateSim failure: %w", err)
	}

	return nil
}

func (o *ukamaAgentClient) UpdatePackage(req ReqData) error {
	log.Debugf("Deactivating operator sim: %v", req.Iccid)

	_, err := o.R.C.R().SetBody(req).Patch(o.u.String() + UkamaOperatorSimsEndpoint + "/" + req.Iccid)
	if err != nil {
		log.Errorf("DeactivateSim failure. error: %s", err.Error())

		return fmt.Errorf("DeactivateSim failure: %w", err)
	}

	return nil
}

func (o *ukamaAgentClient) TerminateSim(iccid string) error {
	log.Debugf("Terminating operator sim: %v", iccid)
	return nil
}

type OperatorSimInfo struct {
	Iccid string `json:"iccid"`
	Imsi  string `json:"imsi"`
}

type OperatorSim struct {
	SimInfo *OperatorSimInfo `json:"sim"`
}

type OperatorUsage struct {
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
