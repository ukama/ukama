package datapath

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

type OvsSwitch struct {
	bridgeName string
	Ip         string
	netType    string
	ofActor    *ofctrl.OfActor
	ctrler     *ofctrl.Controller
}

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

// NewOvsSwitch Creates a new OVS switch instance
func NewOvsSwitch(bridgeName, localIP, netType string) (*OvsSwitch, error) {

	sw := new(OvsSwitch)

	//Create a controller
	sw.ofActor = new(OfActor)
	sw.ctrler = ofctrl.NewController(sw.ofActor)
	sw.bridgeName = bridgeName

	log.Infof("wait for 2sec for ovs bridge ovsbr12 to get created..")
	time.Sleep(2 * time.Second)
	go sw.ctrler.Connect("/var/run/openvswitch/ovsbr12.mgmt")

	//wait for 8sec and see if switch connects
	time.Sleep(8 * time.Second)
	if !sw.ofActor.isSwitchConnected {
		log.Fatalf("%s switch did not connect within 20sec", bridgeName)
	}

	// Create initial tables
	sw.ofActor.inputTable = sw.ofActor.Switch.DefaultTable()
	if sw.ofActor.inputTable == nil {
		log.Fatalf("Failed to get input Table")
		return nil, fmt.Errorf("failed to get input Table for switch")
	}
	log.Infof("Switch connected. Creating tables..")

	return sw, nil
}

func (o *OvsSwitch) DeleteMeter(id uint32) error {
	meterMod := openflow15.NewMeterMod()
	meterMod.MeterId = id
	meterMod.Command = openflow15.MC_DELETE
	return o.ofActor.Switch.Send(meterMod)
}

func (o *OvsSwitch) AddMeter(id, rate, burstSize uint32) error {
	meter := ofctrl.NewMeter(id, ofctrl.MeterKbps, o.ofActor.Switch)

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
		return err
	}
	return nil
}

func (o *OvsSwitch) CreateMetersForUE(rxMeter, txMeter, rxRate, txRate, burstSize uint32) error {
	err := o.AddMeter(rxMeter, rxRate, burstSize)
	if err != nil {
		log.Errorf("Failed to create RX meter.Error: %v", err)
		return err
	}

	err = o.AddMeter(txMeter, txRate, burstSize)
	if err != nil {
		log.Errorf("Failed to create TX meter.Error: %v", err)
		o.DeleteMeter(rxMeter)
		return err
	}
	return nil
}

func (o *OvsSwitch) DeleteMetersForUE(rxMeter, txMeter uint32) error {
	err := o.DeleteMeter(rxMeter)
	if err != nil {
		log.Errorf("Failed to delete RX meter.Error: %v", err)
		return err
	}

	err = o.DeleteMeter(txMeter)
	if err != nil {
		log.Errorf("Failed to delete TX meter.Error: %v", err)
		return err
	}
	return nil
}

func addActionsToFlow(f *ofctrl.Flow, meter uint32) *ofctrl.Flow {
	/* Add actions */
	m := ofctrl.NewMeterAction(meter)
	rNormalAction := ofctrl.NewOutputNormal()
	f.ApplyAction(m)
	f.ApplyAction(rNormalAction)
	return f
}

func (o *OvsSwitch) createTxFlow(ip *net.IP) (*ofctrl.Flow, error) {
	f, err := o.ofActor.inputTable.NewFlow(ofctrl.FlowMatch{
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

func (o *OvsSwitch) createRxFlow(ip *net.IP) (*ofctrl.Flow, error) {
	f, err := o.ofActor.inputTable.NewFlow(ofctrl.FlowMatch{
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

func (o *OvsSwitch) updateFlowForUE(a *OfActor, ipString string, rxMeter, txMeter uint32, oprationType int) error {

	ip := net.ParseIP(ipString)
	if ip == nil {
		log.Errorf("Invalid IP address")
		return fmt.Errorf("invalid ip %s", ipString)
	}

	/* Add RX flow */
	rxF, err := o.createRxFlow(&ip)
	if err != nil {
		log.Errorf("Failed to create RX flow for the UE %s with meter id %d. Error %s", ipString, rxMeter, err.Error())
		return err
	}

	rxF = addActionsToFlow(rxF, rxMeter)

	/* Submit flow */
	err = rxF.Send(oprationType)
	if err != nil {
		log.Errorf("Failed to submit RX flow for the UE %s with meter id %d. Error %s", ipString, rxMeter, err.Error())
		return err
	}
	rxF.UpdateInstallStatus(true)

	/* Add TX flow */
	txF, err := o.createTxFlow(&ip)
	if err != nil {
		log.Errorf("Failed to create TX flow for the UE %s with meter id %d. Error %s", ipString, txMeter, err.Error())
		return err
	}

	/* Add actions */
	txF = addActionsToFlow(txF, txMeter)

	/* Submit flow */
	err = txF.Send(oprationType)
	if err != nil {
		log.Errorf("Failed to submit RX flow for the UE %s with meter id %d. Error %s", ipString, txMeter, err.Error())
		return err
	}
	rxF.UpdateInstallStatus(true)
	return nil
}

func (o *OvsSwitch) AddFlowForUE(ipString string, rxMeter, txMeter uint32) error {
	err := o.updateFlowForUE(ipString, rxMeter, txMeter, openflow15.FC_ADD)
	if err != nil {
		log.Errorf("failed to add flow for UE %s. Error: %s", ipString, err.Error())
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

func (o *OvsSwitch) deleteFlowfromSwitch(ip net.IP, dp DataPath) error {
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
	err := o.ofActor.Send(flow)
	if err != nil {
		log.Errorf("failed to delete flow for UE %v. Error: %s", ip, err.Error())
		return err
	}
	return nil
}

func (o *OvsSwitch) deleteFlowFromTable(ip net.IP, dp DataPath) error {

	// Delete flow from the table
	f := new(ofctrl.Flow)
	f.Table = o.ofActor.inputTable
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

	flowKey := getFlowKey(f.Match)

	err := o.ofActor.inputTable.DeleteFlow(flowKey)
	if err != nil {
		log.Errorf("Failed to remove flow for UE %v from the table.Error: %s", ip, err.Error())
		return err
	}

	return nil
}

func (o *OvsSwitch) deleteFlowForTXPath(ip net.IP) error {

	/* Delete flow from switch */
	err := o.deleteFlowfromSwitch(ip, TX_PATH)
	if err != nil {
		log.Errorf("Deleting TX Path for UE %v flow from switch failed.Error %s", ip, err.Error())
		return err
	}

	/* Delete flow from table */
	err = o.deleteFlowFromTable(ip, TX_PATH)
	if err != nil {
		log.Errorf("Deleting TX Path for UE %v flow from switch failed.Error %s", ip, err.Error())
		return err
	}

	return nil

}

func (o *OvsSwitch) deleteFlowForRXPath(ip net.IP) error {

	/* Delete flow from switch */
	err := o.deleteFlowfromSwitch(ip, RX_PATH)
	if err != nil {
		log.Errorf("Deleting RX Path for UE %v flow from switch failed.Error %s", ip, err.Error())
		return err
	}

	/* Delete flow from table */
	err = o.deleteFlowFromTable(ip, RX_PATH)
	if err != nil {
		log.Errorf("Deleting RX Path for UE %v flow from switch failed.Error %s", ip, err.Error())
		return err
	}

	return nil
}

func (o *OvsSwitch) DeleteFlowForUE(ipString string) error {
	ip := net.ParseIP(ipString)
	if ip == nil {
		log.Errorf("Invalid IP address %s", ipString)
		return fmt.Errorf("invalid ip %s", ipString)
	}

	err := o.deleteFlowForTXPath(ip)
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

func (o *OvsSwitch) AddUEDataPath(ipString string, rxMeter, txMeter, rxRate, txRate, burstSize uint32) error {

	/* Create Meters */
	err := o.CreateMetersForUE(rxMeter, txMeter, rxRate, txRate, burstSize)
	if err != nil {
		log.Errorf("Failed to create meters for UE %s. Error: %v", ipString, err)
		return err
	}

	/* Add Flows */
	err = o.AddFlowForUE(ipString, rxMeter, txMeter)
	if err != nil {
		log.Errorf("Failed to create flows for UE %s. Error: %v", ipString, err)
		o.DeleteMetersForUE(rxMeter, txMeter)
		return err
	}

	return nil
}

func (o *OvsSwitch) DelteUEDataPath(ipString string, rxMeter, txMeter uint32) error {

	/* Delete Flows */
	err := o.DeleteFlowForUE(ipString)
	if err != nil {
		log.Errorf("Failed to delete flows for UE %s. Error: %v", ipString, err)
		o.DeleteMetersForUE(rxMeter, txMeter)
		return err
	}

	/* Delete Meters */
	err = o.DeleteMetersForUE(rxMeter, txMeter)
	if err != nil {
		log.Errorf("Failed to delete meters for UE %s. Error: %v", ipString, err)
		return err
	}

	return nil
}
