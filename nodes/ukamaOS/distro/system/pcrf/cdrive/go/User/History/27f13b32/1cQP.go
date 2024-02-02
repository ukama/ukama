package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"antrea.io/libOpenflow/openflow15"
	"antrea.io/libOpenflow/util"
	"antrea.io/ofnet/ofctrl"
	log "github.com/sirupsen/logrus"
)

type DataPath int

const (
	RX_PATH DataPath = 0
	TX_PATH DataPath = 1
)

// OvsSwitch represents on OVS bridge instance
type OvsSwitch struct {
	bridgeName string
	ovsDriver  *ofctrl.OvsDriver
	ofActor    *OfActor
	ctrler     *ofctrl.Controller
}

// type OFBridge struct {
// 	bridgeName string
// 	// Management address
// 	mgmtAddr string
// 	// ofSwitch is the target OFSwitch.
// 	ofSwitch *ofctrl.OFSwitch
// 	// controller helps maintain connections to remote OFSwitch.
// 	controller *ofctrl.Controller
// }

type OfActor struct {
	Switch            *ofctrl.OFSwitch
	isSwitchConnected bool

	inputTable     *ofctrl.Table
	nextTable      *ofctrl.Table
	connectedCount int

	pktInCount     int
	tlvTableStatus *ofctrl.TLVTableStatus
	tlvMapCh       chan struct{}
}

// const (
// 	AddMessage OFOperation = iota
// 	ModifyMessage
// 	DeleteMessage
// )

// type EntryType string
// type OFOperation int

// const (
// 	FlowEntry  EntryType = "FlowEntry"
// 	GroupEntry EntryType = "GroupEntry"
// 	MeterEntry EntryType = "MeterEntry"
// )

// type OFEntry interface {
// 	Add() error
// 	Modify() error
// 	Delete() error
// 	Type() EntryType
// 	// Reset ensures that the entry is "correct" and that the Add /
// 	// Modify / Delete methods can be called on this object. This method
// 	// should be called if a reconnection event happened.
// 	Reset()
// 	// GetBundleMessages returns a slice of ofctrl.OpenFlowModMessage which can be used in Bundle messages. operation
// 	// specifies what operation is expected to be taken on the OFEntry.
// 	GetBundleMessages(operation OFOperation) ([]ofctrl.OpenFlowModMessage, error)
// }

// type Meter interface {
// 	OFEntry
// 	ResetMeterBands() Meter
// 	MeterBand() MeterBandBuilder
// }

// type MeterBandBuilder interface {
// 	MeterType(meterType ofctrl.MeterType) MeterBandBuilder
// 	Rate(rate uint32) MeterBandBuilder
// 	Burst(burst uint32) MeterBandBuilder
// 	PrecLevel(precLevel uint8) MeterBandBuilder
// 	Experimenter(experimenter uint32) MeterBandBuilder
// 	Done() Meter
// }

// type meterBandBuilder struct {
// 	meter           *ofMeter
// 	meterBandHeader *openflow15.MeterBandHeader
// 	prevLevel       uint8
// 	experimenter    uint32
// }

func (o *OfActor) PacketRcvd(sw *ofctrl.OFSwitch, packet *ofctrl.PacketIn) {
	log.Printf("OF: Received packet: %+v", packet.Data)
}

func (o *OfActor) SwitchConnected(sw *ofctrl.OFSwitch) {
	log.Printf("OF: Switch connected: %v", sw.DPID())

	// Store switch for later use
	o.Switch = sw

	o.isSwitchConnected = true
	o.connectedCount += 1
}

func (o *OfActor) MultipartReply(sw *ofctrl.OFSwitch, rep *openflow15.MultipartReply) {
}

func (o *OfActor) SwitchDisconnected(sw *ofctrl.OFSwitch) {
	log.Printf("OF: Switch disconnected: %v", sw.DPID())
	o.isSwitchConnected = false
}

func (o *OfActor) TLVMapReplyRcvd(sw *ofctrl.OFSwitch, tlvTableStatus *ofctrl.TLVTableStatus) {
	log.Printf("OF: Receive TLVMapTable reply: %s", tlvTableStatus)
	o.tlvTableStatus = tlvTableStatus
	if o.tlvMapCh != nil {
		close(o.tlvMapCh)
	}
}

func (o *OfActor) FlowGraphEnabledOnSwitch() bool {
	return true
}

func (o *OfActor) TLVMapEnabledOnSwitch() bool {
	return true
}

// type ofMeter struct {
// 	ofctrl *ofctrl.Meter
// 	bridge *OFBridge
// }

// func (m *ofMeter) Reset() {
// 	m.ofctrl.Switch = m.bridge.ofSwitch
// }

// // Note: use OFSwitch to directly send MeterModification message rather than bundle message is because the
// // current ofnet implementation for OpenFlow bundle does nmain.goot support adding MeterModification.
// func (m *ofMeter) Add() error {
// 	msg := m.ofctrl.GetBundleMessage(openflow15.MC_ADD)
// 	return m.ofctrl.Switch.Send(msg.GetMessage())
// }

// func (m *ofMeter) Modify() error {
// 	msg := m.ofctrl.GetBundleMessage(openflow15.MC_MODIFY)
// 	return m.ofctrl.Switch.Send(msg.GetMessage())
// }

// func (m *ofMeter) Delete() error {
// 	meterMod := openflow15.NewMeterMod()
// 	meterMod.MeterId = m.ofctrl.ID
// 	meterMod.Command = openflow15.MC_DELETE
// 	return m.ofctrl.Switch.Send(meterMod)
// }

// func (m *ofMeter) GetBundleMessages(entryOper OFOperation) ([]ofctrl.OpenFlowModMessage, error) {
// 	var operation int
// 	switch entryOper {
// 	case AddMessage:
// 		operation = openflow15.MC_ADD
// 	case ModifyMessage:
// 		operation = openflow15.MC_MODIFY
// 	case DeleteMessage:
// 		operation = openflow15.MC_DELETE
// 	}
// 	message := m.ofctrl.GetBundleMessage(operation)
// 	return []ofctrl.OpenFlowModMessage{message}, nil
// }

// func (m *ofMeter) ResetMeterBands() Meter {
// 	m.ofctrl.MeterBands = nil
// 	return m
// }

// func (m *ofMeter) Type() EntryType {
// 	return MeterEntry
// }

// func (m *ofMeter) MeterBand() MeterBandBuilder {
// 	return &meterBandBuilder{
// 		meter:           m,
// 		meterBandHeader: openflow15.NewMeterBandHeader(),
// 		prevLevel:       0,
// 		experimenter:    0,
// 	}
// }

// func (m *meterBandBuilder) MeterType(meterType ofctrl.MeterType) MeterBandBuilder {
// 	m.meterBandHeader.Type = uint16(meterType)
// 	return m
// }

// func (m *meterBandBuilder) Rate(rate uint32) MeterBandBuilder {
// 	m.meterBandHeader.Rate = rate
// 	return m
// }

// func (m *meterBandBuilder) Burst(burst uint32) MeterBandBuilder {
// 	m.meterBandHeader.BurstSize = burst
// 	return m
// }

// func (m *meterBandBuilder) PrecLevel(precLevel uint8) MeterBandBuilder {
// 	m.prevLevel = precLevel
// 	return m
// }

// func (m *meterBandBuilder) Experimenter(experimenter uint32) MeterBandBuilder {
// 	m.experimenter = experimenter
// 	return m
// }

// func (m *meterBandBuilder) Done() Meter {
// 	var mb util.Message
// 	switch m.meterBandHeader.Type {
// 	case uint16(ofctrl.MeterDrop):
// 		mbDrop := new(openflow15.MeterBandDrop)
// 		mbDrop.MeterBandHeader = *m.meterBandHeader
// 		mb = mbDrop
// 	case uint16(ofctrl.MeterDSCPRemark):
// 		mbDscp := new(openflow15.MeterBandDSCP)
// 		mbDscp.MeterBandHeader = *m.meterBandHeader
// 		mbDscp.PrecLevel = m.prevLevel
// 		mb = mbDscp
// 	case uint16(ofctrl.MeterExperimenter):
// 		mbExp := new(openflow15.MeterBandExperimenter)
// 		mbExp.MeterBandHeadeDataPathand(&mb)
// 	return m.meter
// }

// NewOvsSwitch Creates a new OVS switch instance
func NewOvsSwitch(bridgeName string) (*OvsSwitch, error) {
	//var err error
	sw := new(OvsSwitch)

	//Create a controller
	sw.ofActor = new(OfActor)
	sw.ctrler = ofctrl.NewController(sw.ofActor)
	sw.bridgeName = bridgeName

	// Create OVS db driver
	//sw.ovsDriver = ofctrl.NewOvsDriver(bridgeName)

	log.Infof("wait for 2sec for ovs bridge %s to get created..", bridgeName)
	time.Sleep(2 * time.Second)
	go sw.ctrler.Connect("/var/run/openvswitch/gtpbr.mgmt")

	//wait for 8sec and see if switch connects
	time.Sleep(8 * time.Second)
	if !sw.ofActor.isSwitchConnected {
		log.Fatalf("%s switch did not connect within 20sec", bridgeName)
	}

	log.Infof("Switch connected. Creating tables..")

	// Create initial tables
	sw.ofActor.inputTable = sw.ofActor.Switch.DefaultTable()
	if sw.ofActor.inputTable == nil {
		log.Fatalf("Failed to get input Table")
		return nil, fmt.Errorf("failed to get input Table for switch")
	}

	// sw.ofActor.nextTable, err = sw.ofActor.Switch.NewTable(1)
	// if err != nil {
	// 	log.Fatalf("Error creating next Table: %v", err)
	// 	return nil, fmt.Errorf("failed to create next Table for switch")
	// }
	log.Infof("Openflow tables created successfully")

	return sw, nil
}

// Delete performs cleanup prior to destruction of the OvsDriver
func (sw *OvsSwitch) Delete() {

	if sw.ovsDriver != nil {
		sw.ovsDriver.Delete()

		// Wait a little for OVS switch to be deleted
		time.Sleep(300 * time.Millisecond)
	}
}

func addMeter(id, rate, burstSize uint32, sw *ofctrl.OFSwitch) (*ofctrl.Meter, error) {
	meter := ofctrl.NewMeter(id, ofctrl.MeterKbps, sw)

	var mb util.Message
	mbDrop := new(openflow15.MeterBandDrop)
	meterBandHeader := *openflow15.NewMeterBandHeader()
	meterBandHeader.Type = uint16(ofctrl.MeterDrop)
	meterBandHeader.Rate = rate
	meterBandHeader.BurstSize = burstSize
	mbDrop.MeterBandHeader = meterBandHeader
	mb = mbDrop
	meter.AddMeterBand(&mb)

	err := meter.Install()
	if err != nil {
		log.Errorf("failed to install.")
		return nil, err
	}
	return meter, nil
}

func deleteMeter(id uint32, sw *ofctrl.OFSwitch) error {
	meterMod := openflow15.NewMeterMod()
	meterMod.MeterId = id
	meterMod.Command = openflow15.MC_DELETE
	return sw.Send(meterMod)
}

/*
	{
				name: "Meter",
				actionFn: func(b Action) FlowBuilder {
					return b.Meter(100)
				},
				expectedActionField: &openflow15.ActionMeter{
					MeterId: 100,
				},
				expectedActionStr: "meter:100",
			},
*/

func addActionsToFlow(f *ofctrl.Flow, meter uint32) *ofctrl.Flow {
	/* Add actions */
	m := ofctrl.NewMeterAction(meter)
	rNormalAction := ofctrl.NewOutputNormal()
	f.ApplyAction(m)
	f.ApplyAction(rNormalAction)
	return f
}

func createTxFlow(a *OfActor, ip *net.IP) (*ofctrl.Flow, error) {
	f, err := a.inputTable.NewFlow(ofctrl.FlowMatch{
		Ethertype: 0x0800,
		Priority:  100,
		IpSa:      ip,
	})
	if err != nil {
		log.Errorf("Failed creating flow for switch")
		return nil, err
	}

	return f, nil
}

func createRxFlow(a *OfActor, ip *net.IP) (*ofctrl.Flow, error) {
	f, err := a.inputTable.NewFlow(ofctrl.FlowMatch{
		Ethertype: 0x0800,
		Priority:  100,
		IpDa:      ip,
	})
	if err != nil {
		log.Errorf("Failed creating flow for switch")
		return nil, err
	}
	return f, nil
}

func updateFlowForUe(a *OfActor, ipString string, meter uint32, oprationType int) error {

	ip := net.ParseIP(ipString)
	if ip == nil {
		log.Errorf("Invalid IP address")
		return fmt.Errorf("invalid ip %s", ipString)
	}

	/* Add RX flow */
	rxF, err := createRxFlow(a, &ip)
	if err != nil {
		log.Errorf("Failed to create RX flow for the UE %s with meter id %d. Error %s", ipString, meter, err.Error())
		return err
	}

	rxF = addActionsToFlow(rxF, meter)

	/* Submit flow */
	err = rxF.Send(oprationType)
	if err != nil {
		log.Errorf("Failed to submit RX flow for the UE %s with meter id %d. Error %s", ipString, meter, err.Error())
		return err
	}
	rxF.UpdateInstallStatus(true)

	/* Add TX flow */
	txF, err := createTxFlow(a, &ip)
	if err != nil {
		log.Errorf("Failed to create TX flow for the UE %s with meter id %d. Error %s", ipString, meter, err.Error())
		return err
	}

	/* Add actions */
	txF = addActionsToFlow(txF, meter)

	/* Submit flow */
	err = txF.Send(oprationType)
	if err != nil {
		log.Errorf("Failed to submit RX flow for the UE %s with meter id %d. Error %s", ipString, meter, err.Error())
		return err
	}
	rxF.UpdateInstallStatus(true)
	return nil
}

func AddFlowForUe(a *OfActor, ipString string, meter uint32) error {
	err := updateFlowForUe(a, ipString, meter, openflow15.FC_ADD)
	if err != nil {
		log.Errorf("failed to add flow for UE %s and meter %d. Error: %s", ipString, meter, err.Error())
		return err
	}

	log.Infof("Added flow for UE %s and meter %d.", ipString, meter)
	return nil
}

func GetFlowKey(m ofctrl.FlowMatch) string {
	jsonVal, err := json.Marshal(m)
	if err != nil {
		log.Errorf("Error forming flowkey for %+v. Err: %v", m, err)
		return ""
	}

	return string(jsonVal)
}

func deleteFlowfromSwitch(a *OfActor, ip net.IP, dp DataPath) error {
	// openflow15 protocol to delete flows from the switch
	flow := openflow15.NewFlowMod()
	flow.TableId = 0
	flow.Priority = 100
	flow.Match = *openflow15.NewMatch()
	flow.Command = openflow15.FC_DELETE
	flow.Match.AddField(*openflow15.NewEthTypeField(0x800))
	if dp == TX_PATH {
		flow.Match.AddField(*openflow15.NewIpv4SrcField(ip, nil))
	} else if dp == RX_PATH {
		flow.Match.AddField(*openflow15.NewIpv4DstField(ip, nil))
	}
	err := a.Switch.Send(flow)
	if err != nil {
		log.Errorf("failed to delete flow for UE %v. Error: %s", ip, err.Error())
		return err
	}
	return nil
}

func deleteFlowFromTable(a *OfActor, ip net.IP, dp DataPath) error {

	// Delete flow from the table
	f := new(ofctrl.Flow)
	f.Table = a.inputTable
	if dp == TX_PATH {
		f.Match = ofctrl.FlowMatch{
			Ethertype: 0x0800,
			Priority:  100,
			IpSa:      &ip,
		}
	} else if dp == RX_PATH {
		f.Match = ofctrl.FlowMatch{
			Ethertype: 0x0800,
			Priority:  100,
			IpDa:      &ip,
		}
	}

	f.UpdateInstallStatus(true)

	flowKey := GetFlowKey(f.Match)

	err := a.inputTable.DeleteFlow(flowKey)
	if err != nil {
		log.Errorf("Failed to remove flow for UE %v from the table.Error: %s", ip, err.Error())
		return err
	}

	return nil
}

func DeleteFlowForTXPath(a *OfActor, ip net.IP) error {

	/* Delete flow from switch */
	err := deleteFlowfromSwitch(a, ip, TX_PATH)
	if err != nil {
		log.Errorf("Deleting TX Path for UE %v flow from switch failed.Error %s", ip, err.Error())
		return err
	}

	/* Delete flow from table */
	err = deleteFlowFromTable(a, ip, TX_PATH)
	if err != nil {
		log.Errorf("Deleting TX Path for UE %v flow from switch failed.Error %s", ip, err.Error())
		return err
	}

	return nil

}

func DeleteFlowForRXPath(a *OfActor, ip net.IP) error {

	/* Delete flow from switch */
	err := deleteFlowfromSwitch(a, ip, RX_PATH)
	if err != nil {
		log.Errorf("Deleting RX Path for UE %v flow from switch failed.Error %s", ip, err.Error())
		return err
	}

	/* Delete flow from table */
	err = deleteFlowFromTable(a, ip, RX_PATH)
	if err != nil {
		log.Errorf("Deleting RX Path for UE %v flow from switch failed.Error %s", ip, err.Error())
		return err
	}

	return nil
}

func DeleteFlowForUe(a *OfActor, ipString string) error {
	ip := net.ParseIP(ipString)
	if ip == nil {
		log.Errorf("Invalid IP address %s", ipString)
		return fmt.Errorf("invalid ip %s", ipString)
	}

	err := DeleteFlowForTXPath(a, ip)
	if err != nil {
		log.Errorf("Failed to delete TX flow for UE %s", ipString)
		return err
	}

	err = DeleteFlowForRXPath(a, ip)
	if err != nil {
		log.Errorf("Failed to delete RX flow for UE %s", ipString)
		return err
	}

	log.Infof("Deleted flow for UE %s", ipString)
	return nil
}

func start() {
	var id uint32 = 3
	// var rate uint32 = 1000
	// var burstSize uint32 = 1500
	log.Infof("Starting ovs controller ")
	sw, err := NewOvsSwitch("gtpbr")
	if err != nil {
		log.Errorf("Error creating switch %s", err.Error())
		return
	}

	// meter, err := addMeter(id, rate, burstSize, sw.ofActor.Switch)
	// if err != nil {
	// 	log.Errorf("Failed to install meter.")
	// }

	// log.Infof("Meter created with id %d", meter.ID)

	// err = deleteMeter(meter.ID, sw.ofActor.Switch)
	// if err != nil {
	// 	log.Errorf("Failed to delete meter.")
	// }

	err = AddFlowForUe(sw.ofActor, "192.168.8.3", id)
	if err != nil {
		log.Errorf("Failed to add  flow.")
	}

	err = DeleteFlowForUe(sw.ofActor, "192.168.8.3")
	if err != nil {
		log.Errorf("Failed to add  flow.")
	}

	// meter := ofctrl.NewMeter(3, ofctrl.MeterKbps, sw.ofActor.Switch)
	// if err != nil {
	// 	log.Errorf("Error creating meter meter %s", err.Error())
	// 	return
	// }

	// var mb util.Message
	// mbDrop := new(openflow15.MeterBandDrop)
	// meterBandHeader := *openflow15.NewMeterBandHeader()
	// meterBandHeader.Type = uint16(ofctrl.MeterDrop)
	// meterBandHeader.Rate = rate
	// meterBandHeader.BurstSize = burstSize
	// mbDrop.MeterBandHeader = meterBandHeader
	// mb = mbDrop
	// meter.AddMeterBand(&mb)

	// err = meter.Install()
	// if err != nil {
	// 	log.Errorf("failed to install.")
	// }

}

func main() {
	start()
}
