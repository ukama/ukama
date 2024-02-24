package datapath

import (
	log "github.com/sirupsen/logrus"
)

type dataPath struct {
	ovs     *OvsSwitch
	ueCount uint32
}

type DataPath interface {
	AddNewDataPath(ip string, rxMeter, txMeter, rxRate, txRate, burstSize uint32, rxCookie, txCookie uint64) error
	DeleteDataPath(ip string, rxMeter, txMeter uint32) error
	DataPathCount() uint32
	DataPathStats(rxCookieID, txCookieID uint64) (uint64, uint64, uint64, uint64, error)
}

func InitDataPath(name, ip, netType, mgmt string) (*dataPath, error) {
	var err error
	d := &dataPath{ueCount: 0}
	d.ovs, err = NewOvsSwitch(name, ip, netType, mgmt)
	if err != nil {
		log.Errorf("error connecting bridge %s at %s .Error: %v", name, ip, err)
		return nil, err
	}
	return d, nil
}

func (d *dataPath) AddNewDataPath(ip string, rxMeter, txMeter, rxRate, txRate, burstSize uint32, rxCookie, txCookie uint64) error {
	err := d.ovs.AddUEDataPath(ip, rxMeter, txMeter, rxRate, txRate, burstSize, rxCookie, txCookie)
	if err != nil {
		log.Errorf("Failed to add datapath for UE %s. Error: %v", ip, err.Error())
		return err
	}

	d.ueCount++
	return nil
}

func (d *dataPath) DeleteDataPath(ip string, rxMeter, txMeter uint32) error {
	err := d.ovs.DeleteUEDataPath(ip, rxMeter, txMeter)
	if err != nil {
		log.Errorf("Failed to delete datapath for UE %s. Error: %v", ip, err.Error())
		return err
	}
	d.ueCount--
	return nil
}

func (d *dataPath) DataPathCount() uint32 {
	return d.ueCount
}

func (d *dataPath) DataPathStats(rxCookieID, txCookieID uint64) (uint64, uint64, uint64, uint64, error) {

	rxBC, rxPC, txBC, txPC, err := d.ovs.DataPathUEStats(rxCookieID, txCookieID)
	if err != nil {
		log.Errorf("Error getting UE pathstats %s", err.Error())
		return 0, 0, 0, 0, err
	}
	return rxBC, rxPC, txBC, txPC, nil
}
