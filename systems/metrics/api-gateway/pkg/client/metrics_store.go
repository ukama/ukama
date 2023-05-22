package client

import (
	log "github.com/sirupsen/logrus"
	rc "github.com/ukama/ukama/systems/common/rest"
)

type MetricsStore struct {
	conn *rc.RestClient
}

func NewMetricsStore(host string, debug bool) (*MetricsStore, error) {

	c, err := rc.NewRestClient(host, debug)
	if err != nil {
		log.Errorf("Failed to create a client to  metrics store %s. Error %s", host, err.Error())
		return nil, err
	}
	return &MetricsStore{
		conn: c,
	}, nil
}
