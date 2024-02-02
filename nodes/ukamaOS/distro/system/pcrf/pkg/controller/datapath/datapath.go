package datapath

import (
	log "github.com/sirupsen/logrus"
)

type datapath struct {
	ovs     *OvsSwitch
	ueCount int
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

func (d *datapath) AddNewDataPath(ip string) error {
	err := c.ovs.AddUEDataPath()
	if err != nil {
		log.Errorf("Failed to add datapath for UE %s. Error: %v", ip, err.Error())
		return err
	}

	c.ueCount++
	return nil
}

func (d *datapath) DeleteDataPath(ip string) error {
	err := c.ovs.DelteUEDataPath()
	if err != nil {
		log.Errorf("Failed to delete datapath for UE %s. Error: %v", ip, err.Error())
		return err
	}
	return nil
}

func (d *datapath) DataPathCount() error {

	return nil
}

func (d *datapath) DataPathStats() error {

	return nil
}
