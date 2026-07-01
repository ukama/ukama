/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package datapath

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"antrea.io/libOpenflow/openflow15"
	"antrea.io/libOpenflow/util"
	"antrea.io/ofnet/ofctrl"

	log "github.com/sirupsen/logrus"
)

type UeDataPath int

const (
	RX_PATH UeDataPath = iota
	TX_PATH
)

type OfActor struct {
	Switch            *ofctrl.OFSwitch
	isSwitchConnected bool
	inputTable        *ofctrl.Table
	connectedCount    int
	tlvTableStatus    *ofctrl.TLVTableStatus
	tlvMapCh          chan struct{}
	mu                sync.RWMutex
}

type OvsSwitch struct {
	bridgeName       string
	Ip               string
	managementSocket string
	ofActor          *OfActor
	ctrler           *ofctrl.Controller
}

func (o *OfActor) PacketRcvd(sw *ofctrl.OFSwitch, packet *ofctrl.PacketIn) {
	log.Printf("OF: Received packet: %+v", packet.Data)
}

func (o *OfActor) SwitchConnected(sw *ofctrl.OFSwitch) {
	log.Printf("OF: Switch connected: %v", sw.DPID())

	o.mu.Lock()
	defer o.mu.Unlock()

	o.Switch = sw
	o.isSwitchConnected = true
	o.connectedCount++
}

func (o *OfActor) MultipartReply(sw *ofctrl.OFSwitch, rep *openflow15.MultipartReply) {
}

func (o *OfActor) SwitchDisconnected(sw *ofctrl.OFSwitch) {
	log.Printf("OF: Switch disconnected: %v", sw.DPID())

	o.mu.Lock()
	defer o.mu.Unlock()

	o.isSwitchConnected = false
}

func (o *OfActor) TLVMapReplyRcvd(sw *ofctrl.OFSwitch, tlvTableStatus *ofctrl.TLVTableStatus) {
	log.Printf("OF: Receive TLVMapTable reply: %s", tlvTableStatus)

	o.mu.Lock()
	defer o.mu.Unlock()

	o.tlvTableStatus = tlvTableStatus
	if o.tlvMapCh != nil {
		close(o.tlvMapCh)
		o.tlvMapCh = nil
	}
}

func (o *OfActor) FlowGraphEnabledOnSwitch() bool {
	return true
}

func (o *OfActor) TLVMapEnabledOnSwitch() bool {
	return true
}

func (o *OfActor) connected() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.isSwitchConnected
}

func (o *OfActor) count() int {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.connectedCount
}

func (o *OfActor) switchHandle() *ofctrl.OFSwitch {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.Switch
}

func (o *OfActor) tableHandle() *ofctrl.Table {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.inputTable
}

func (o *OfActor) setInputTable(t *ofctrl.Table) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.inputTable = t
}

func TryConnect(c *ofctrl.Controller, sock string) {
	if err := c.Connect(sock); err != nil {
		log.Errorf("Failed to connect OpenFlow controller socket %s: %v",
			sock, err)
	}
}

// NewOvsSwitch creates a new OVS switch instance.
//
// mgmtPath is the OVS management directory, not the full socket path.
// Example:
//
//	bridgeName = br0
//	mgmtPath   = /usr/local/var/run/openvswitch
//	socket     = /usr/local/var/run/openvswitch/br0.mgmt
func NewOvsSwitch(bridgeName, localIP, netType, mgmtPath string) (*OvsSwitch, error) {
	var ofsw *ofctrl.OFSwitch

	sw := &OvsSwitch{}

	sw.ofActor = &OfActor{}
	sw.ctrler = ofctrl.NewController(sw.ofActor)
	sw.bridgeName = bridgeName
	sw.Ip = localIP
	sw.managementSocket = fmt.Sprintf("%s/%s.mgmt", mgmtPath, bridgeName)

	go TryConnect(sw.ctrler, sw.managementSocket)

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if sw.ofActor.connected() {
			break
		}
		time.Sleep(250 * time.Millisecond)
	}

	if !sw.ofActor.connected() {
		return nil, fmt.Errorf("%s switch did not connect within 10 sec using %s",
			bridgeName, sw.managementSocket)
	}

	ofsw = sw.ofActor.switchHandle()
	if ofsw == nil {
		return nil, fmt.Errorf("switch handle is nil after connect")
	}

	sw.ofActor.setInputTable(ofsw.DefaultTable())
	if sw.ofActor.tableHandle() == nil {
		return nil, fmt.Errorf("failed to get input table for switch")
	}

	ofsw.EnableMonitor()

	log.Infof("Switch connected bridge=%s socket=%s",
		sw.bridgeName, sw.managementSocket)

	return sw, nil
}

func (o *OvsSwitch) Status() Status {
	if o == nil || o.ofActor == nil {
		return Status{}
	}

	return Status{
		Bridge:           o.bridgeName,
		ManagementSocket: o.managementSocket,
		Connected:        o.ofActor.connected(),
		ConnectedCount:   o.ofActor.count(),
	}
}

func (o *OvsSwitch) switchHandle() (*ofctrl.OFSwitch, error) {
	sw := o.ofActor.switchHandle()
	if sw == nil || !o.ofActor.connected() {
		return nil, fmt.Errorf("OVS switch is not connected")
	}

	return sw, nil
}

func (o *OvsSwitch) inputTable() (*ofctrl.Table, error) {
	t := o.ofActor.tableHandle()
	if t == nil || !o.ofActor.connected() {
		return nil, fmt.Errorf("OVS input table is not ready")
	}

	return t, nil
}

func (o *OvsSwitch) DeleteMeter(id uint32) error {
	var sw *ofctrl.OFSwitch
	var err error

	sw, err = o.switchHandle()
	if err != nil {
		return err
	}

	meterMod := openflow15.NewMeterMod()
	meterMod.MeterId = id
	meterMod.Command = openflow15.MC_DELETE

	return sw.Send(meterMod)
}

func (o *OvsSwitch) AddMeter(id, rate, burstSize uint32) error {
	var sw *ofctrl.OFSwitch
	var err error
	var mb util.Message

	sw, err = o.switchHandle()
	if err != nil {
		return err
	}

	meter := ofctrl.NewMeter(id, ofctrl.MeterKbps, sw)
	mbDrop := new(openflow15.MeterBandDrop)
	meterBandHeader := *openflow15.NewMeterBandHeader()

	meterBandHeader.Type = uint16(ofctrl.MeterDrop)
	meterBandHeader.Rate = rate
	meterBandHeader.BurstSize = burstSize

	mbDrop.MeterBandHeader = meterBandHeader
	mb = mbDrop
	meter.AddMeterBand(&mb)

	err = meter.Install()
	if err != nil {
		log.Errorf("failed to install meter id=%d rate=%d burst=%d: %v",
			id, rate, burstSize, err)
		return err
	}

	return nil
}

func (o *OvsSwitch) CreateMetersForUE(rxMeter, txMeter, rxRate, txRate, burstSize uint32) error {
	err := o.AddMeter(rxMeter, rxRate, burstSize)
	if err != nil {
		log.Errorf("Failed to create RX meter. Error: %v", err)
		return err
	}

	err = o.AddMeter(txMeter, txRate, burstSize)
	if err != nil {
		log.Errorf("Failed to create TX meter. Error: %v", err)
		_ = o.DeleteMeter(rxMeter)
		return err
	}

	return nil
}

func (o *OvsSwitch) DeleteMetersForUE(rxMeter, txMeter uint32) error {
	err := o.DeleteMeter(rxMeter)
	if err != nil {
		log.Errorf("Failed to delete RX meter. Error: %v", err)
		return err
	}

	err = o.DeleteMeter(txMeter)
	if err != nil {
		log.Errorf("Failed to delete TX meter. Error: %v", err)
		return err
	}

	return nil
}

func addActionsToFlow(f *ofctrl.Flow, meter uint32) *ofctrl.Flow {
	m := ofctrl.NewMeterAction(meter)
	normal := ofctrl.NewOutputNormal()

	f.ApplyAction(m)
	f.ApplyAction(normal)

	return f
}

func parseIPv4(ipString string) (net.IP, error) {
	ip := net.ParseIP(ipString)
	if ip == nil {
		return nil, fmt.Errorf("invalid ip %s", ipString)
	}

	ip4 := ip.To4()
	if ip4 == nil {
		return nil, fmt.Errorf("invalid ipv4 %s", ipString)
	}

	return ip4, nil
}

func (o *OvsSwitch) createTxFlow(ip *net.IP) (*ofctrl.Flow, error) {
	table, err := o.inputTable()
	if err != nil {
		return nil, err
	}

	f, err := table.NewFlow(ofctrl.FlowMatch{
		Ethertype: 0x0800,
		Priority:  100,
		IpSa:      ip,
	})
	if err != nil {
		log.Errorf("Failed creating TX flow for switch: %v", err)
		return nil, err
	}

	return f, nil
}

func (o *OvsSwitch) createRxFlow(ip *net.IP) (*ofctrl.Flow, error) {
	table, err := o.inputTable()
	if err != nil {
		return nil, err
	}

	f, err := table.NewFlow(ofctrl.FlowMatch{
		Ethertype: 0x0800,
		Priority:  100,
		IpDa:      ip,
	})
	if err != nil {
		log.Errorf("Failed creating RX flow for switch: %v", err)
		return nil, err
	}

	return f, nil
}

func (o *OvsSwitch) updateFlowForUE(ipString string, rxMeter, txMeter uint32, rxCookie, txCookie uint64, operationType int) error {
	ip, err := parseIPv4(ipString)
	if err != nil {
		log.Errorf("Invalid IP address %s", ipString)
		return err
	}

	rxF, err := o.createRxFlow(&ip)
	if err != nil {
		log.Errorf("Failed to create RX flow for UE %s with meter id %d. Error: %s",
			ipString, rxMeter, err.Error())
		return err
	}

	rxF.CookieID = rxCookie
	rxF = addActionsToFlow(rxF, rxMeter)

	err = rxF.Send(operationType)
	if err != nil {
		log.Errorf("Failed to submit RX flow for UE %s with meter id %d. Error: %s",
			ipString, rxMeter, err.Error())
		return err
	}

	rxF.UpdateInstallStatus(true)

	txF, err := o.createTxFlow(&ip)
	if err != nil {
		log.Errorf("Failed to create TX flow for UE %s with meter id %d. Error: %s",
			ipString, txMeter, err.Error())
		return err
	}

	txF.CookieID = txCookie
	txF = addActionsToFlow(txF, txMeter)

	err = txF.Send(operationType)
	if err != nil {
		log.Errorf("Failed to submit TX flow for UE %s with meter id %d. Error: %s",
			ipString, txMeter, err.Error())
		return err
	}

	txF.UpdateInstallStatus(true)

	return nil
}

func (o *OvsSwitch) AddFlowForUE(ipString string, rxMeter, txMeter uint32, rxCookie, txCookie uint64) error {
	err := o.updateFlowForUE(ipString, rxMeter, txMeter, rxCookie, txCookie, openflow15.FC_ADD)
	if err != nil {
		log.Errorf("failed to add flow for UE %s. Error: %s",
			ipString, err.Error())
		return err
	}

	log.Infof("Added flow for UE %s", ipString)
	return nil
}

func getFlowKey(m ofctrl.FlowMatch) string {
	jsonVal, err := json.Marshal(m)
	if err != nil {
		log.Errorf("Error forming flowkey for %+v. Err: %v", m, err)
		return ""
	}

	return string(jsonVal)
}

func (o *OvsSwitch) deleteFlowfromSwitch(ip net.IP, dp UeDataPath) error {
	sw, err := o.switchHandle()
	if err != nil {
		return err
	}

	flow := openflow15.NewFlowMod()
	flow.TableId = 0
	flow.Priority = 100
	flow.Match = *openflow15.NewMatch()
	flow.Command = openflow15.FC_DELETE

	flow.Match.AddField(*openflow15.NewEthTypeField(0x800))

	switch dp {
	case TX_PATH:
		flow.Match.AddField(*openflow15.NewIpv4SrcField(ip, nil))

	case RX_PATH:
		flow.Match.AddField(*openflow15.NewIpv4DstField(ip, nil))
	}

	err = sw.Send(flow)
	if err != nil {
		log.Errorf("failed to delete flow for UE %v. Error: %s",
			ip, err.Error())
		return err
	}

	return nil
}

func (o *OvsSwitch) deleteFlowFromTable(ip net.IP, dp UeDataPath) error {
	table, err := o.inputTable()
	if err != nil {
		return err
	}

	f := new(ofctrl.Flow)
	f.Table = table

	switch dp {
	case TX_PATH:
		f.Match = ofctrl.FlowMatch{
			Ethertype: 0x0800,
			Priority:  100,
			IpSa:      &ip,
		}

	case RX_PATH:
		f.Match = ofctrl.FlowMatch{
			Ethertype: 0x0800,
			Priority:  100,
			IpDa:      &ip,
		}
	}

	f.UpdateInstallStatus(true)

	flowKey := getFlowKey(f.Match)
	if flowKey == "" {
		return fmt.Errorf("empty flow key for UE %v", ip)
	}

	err = table.DeleteFlow(flowKey)
	if err != nil {
		log.Errorf("Failed to remove flow for UE %v from the table. Error: %s",
			ip, err.Error())
		return err
	}

	return nil
}

func (o *OvsSwitch) deleteFlowForTXPath(ip net.IP) error {
	err := o.deleteFlowfromSwitch(ip, TX_PATH)
	if err != nil {
		log.Errorf("Deleting TX path for UE %v flow from switch failed. Error: %s",
			ip, err.Error())
		return err
	}

	err = o.deleteFlowFromTable(ip, TX_PATH)
	if err != nil {
		log.Errorf("Deleting TX path for UE %v flow from table failed. Error: %s",
			ip, err.Error())
		return err
	}

	return nil
}

func (o *OvsSwitch) deleteFlowForRXPath(ip net.IP) error {
	err := o.deleteFlowfromSwitch(ip, RX_PATH)
	if err != nil {
		log.Errorf("Deleting RX path for UE %v flow from switch failed. Error: %s",
			ip, err.Error())
		return err
	}

	err = o.deleteFlowFromTable(ip, RX_PATH)
	if err != nil {
		log.Errorf("Deleting RX path for UE %v flow from table failed. Error: %s",
			ip, err.Error())
		return err
	}

	return nil
}

func (o *OvsSwitch) DeleteFlowForUE(ipString string) error {
	ip, err := parseIPv4(ipString)
	if err != nil {
		log.Errorf("Invalid IP address %s", ipString)
		return err
	}

	err = o.deleteFlowForTXPath(ip)
	if err != nil {
		log.Errorf("Failed to delete TX flow for UE %s", ipString)
		return err
	}

	err = o.deleteFlowForRXPath(ip)
	if err != nil {
		log.Errorf("Failed to delete RX flow for UE %s", ipString)
		return err
	}

	log.Infof("Deleted flow for UE %s", ipString)
	return nil
}

func (o *OvsSwitch) AddUEDataPath(ipString string, rxMeter, txMeter, rxRate, txRate, burstSize uint32, rxCookie, txCookie uint64) error {
	err := o.CreateMetersForUE(rxMeter, txMeter, rxRate, txRate, burstSize)
	if err != nil {
		log.Errorf("Failed to create meters for UE %s. Error: %v",
			ipString, err)
		return err
	}

	err = o.AddFlowForUE(ipString, rxMeter, txMeter, rxCookie, txCookie)
	if err != nil {
		log.Errorf("Failed to create flows for UE %s. Error: %v",
			ipString, err)
		_ = o.DeleteMetersForUE(rxMeter, txMeter)
		return err
	}

	return nil
}

func (o *OvsSwitch) DeleteUEDataPath(ipString string, rxMeter, txMeter uint32) error {
	err := o.DeleteFlowForUE(ipString)
	if err != nil {
		log.Errorf("Failed to delete flows for UE %s. Error: %v",
			ipString, err)
		_ = o.DeleteMetersForUE(rxMeter, txMeter)
		return err
	}

	err = o.DeleteMetersForUE(rxMeter, txMeter)
	if err != nil {
		log.Errorf("Failed to delete meters for UE %s. Error: %v",
			ipString, err)
		return err
	}

	return nil
}

func parseStats(s openflow15.Stats) (uint64, uint64, error) {
	bc := new(openflow15.PBCountStatField)
	pc := new(openflow15.PBCountStatField)
	n := 2

	data, err := s.MarshalBinary()
	if err != nil {
		log.Errorf("Failed to marshal data. Error %s", err.Error())
		return 0, 0, err
	}

	if len(data) < 4 {
		return 0, 0, fmt.Errorf("invalid stats binary length %d", len(data))
	}

	s.Length = binary.BigEndian.Uint16(data[n:])
	n += 2

	for n < int(s.Length) {
		var f util.Message
		var size uint16

		if n+2 >= len(data) {
			break
		}

		switch data[n+2] >> 1 {
		case openflow15.XST_OFB_DURATION:
			f = new(openflow15.TimeStatField)
			size = f.Len()

		case openflow15.XST_OFB_IDLE_TIME:
			f = new(openflow15.TimeStatField)
			size = f.Len()

		case openflow15.XST_OFB_FLOW_COUNT:
			f = new(openflow15.FlowCountStatField)
			size = f.Len()

		case openflow15.XST_OFB_PACKET_COUNT:
			err = pc.UnmarshalBinary(data[n:])
			if err != nil {
				log.Errorf("Failed to unmarshal packet count stats field data %+v",
					data[n:])
				return 0, 0, err
			}
			size = pc.Len()

		case openflow15.XST_OFB_BYTE_COUNT:
			err = bc.UnmarshalBinary(data[n:])
			if err != nil {
				log.Errorf("Failed to unmarshal byte count stats field data %v",
					data[n:])
				return 0, 0, err
			}
			size = bc.Len()

		default:
			return 0, 0, fmt.Errorf("received unknown stats field: %v",
				data[n+2]>>1)
		}

		if size == 0 {
			return 0, 0, fmt.Errorf("invalid zero-size stats field")
		}

		n += int(size)
	}

	return bc.Count, pc.Count, nil
}

func (o *OvsSwitch) dataPathStats(cookieID uint64) (uint64, uint64, error) {
	var stats []*openflow15.FlowDesc
	var err error

	cookieMask := uint64(0xffffffffffffffff)

	sw, err := o.switchHandle()
	if err != nil {
		return 0, 0, err
	}

	stats, err = sw.DumpFlowStats(cookieID, &cookieMask, nil, nil)
	if err != nil {
		log.Errorf("Error getting stats %s", err.Error())
		return 0, 0, err
	}

	if len(stats) == 0 {
		return 0, 0, fmt.Errorf("stats not found for cookie %d (0x%x)",
			cookieID, cookieID)
	}

	bc, pc, err := parseStats(stats[0].Stats)
	if err != nil {
		log.Errorf("Failed to get stats for flow %d (0x%x)",
			cookieID, cookieID)
		return 0, 0, err
	}

	return bc, pc, nil
}

func (o *OvsSwitch) DataPathUEStats(rxCookieID, txCookieID uint64) (uint64, uint64, uint64, uint64, error) {
	rxBC, rxPC, err := o.dataPathStats(rxCookieID)
	if err != nil {
		log.Errorf("Error getting RX path stats %s", err.Error())
		return 0, 0, 0, 0, err
	}

	txBC, txPC, err := o.dataPathStats(txCookieID)
	if err != nil {
		log.Errorf("Error getting TX path stats %s", err.Error())
		return 0, 0, 0, 0, err
	}

	return rxBC, rxPC, txBC, txPC, nil
}
