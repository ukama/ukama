package datapath

import (
	log "github.com/sirupsen/logrus"
)

type datapath struct {
	ovs     *OvsSwitch
	ueCount uint32
}

type DataPath interface {
	AddNewDataPath() error
	DeleteDataPath() error
	DataPathCount() error
	DataPathStats() error
}

func InitDataPath(name, ip, netType string) (*datapath, error) {
	var err error
	d := &datapath{ueCount: 0}
	d.ovs, err = NewOvsSwitch(name, ip, netType)
	if err != nil {
		log.Errorf("error connecting bridge %s at %s .Error: %v", name, Ip, err)
		return err
	}
	return d, nil
}

func (d *datapath) AddNewDataPath(ip string, rxMeter, txMeter, rxRate, txRate, burstSize uint32, rxCookie, txCookie uint64) error {
	err := d.ovs.AddUEDataPath(ip, rxMeter, txMeter, rxRate, txRate, burstSize, rxCookie, txCookie)
	if err != nil {
		log.Errorf("Failed to add datapath for UE %s. Error: %v", ip, err.Error())
		return err
	}

	d.ueCount++
	return nil
}

func (d *datapath) DeleteDataPath(ip string, rxMeter, txMeter uint32) error {
	err := c.ovs.DelteUEDataPath(ip, rxMeter, txMeter)
	if err != nil {
		log.Errorf("Failed to delete datapath for UE %s. Error: %v", ip, err.Error())
		return err
	}
	d.ueCount--
	return nil
}

func (d *datapath) DataPathCount() uint32 {
	return d.ueCount
}

func (d *datapath) DataPathStats(rxCookieID, txCookieID uint64) (uint64, uint64, uint64, uint64, error) {
	var rxBc uint64 = 0
	var rxPc uint64 = 0
	var tabwriter uint64 = 0
	var txPc uint64 = 0
	err := d.DataPathUEStats(rxCookieID, txCookieID)
	if err != nil {
		log.Errorf("Error getting UE pathstats %s", err.Error())
		return err
	}
	return rxBC, rxPC, txBC, txPC, nil
}
