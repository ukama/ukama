/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package queue

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/services/msgClient/internal/db"
	"google.golang.org/protobuf/proto"
)

type QueuePublisher struct {
	q          string
	name       string
	exchange   string `default:"amq.topic"`
	instanceId string
	pub        mb.QPub
}

func NewQueuePublisher(s db.Service) (*QueuePublisher, error) {

	pub, err := mb.NewQPub(s.MsgBusUri, s.Name, s.Exchange, s.InstanceId)
	if err != nil {
		log.Errorf("Failed to create publisher. Error: %s", err.Error())
		return nil, err
	}

	qp := &QueuePublisher{
		q:          s.PublQueue,
		name:       s.Name,
		instanceId: s.InstanceId,
		pub:        pub,
		exchange:   s.Exchange,
	}

	return qp, nil
}

func (p *QueuePublisher) Publish(key string, payload proto.Message) error {

	err := make(chan error, 1)
	go func(err chan error) {
		e := p.pub.PublishProto(payload, key)
		if e != nil {
			log.Errorf("Failed to publish message. Error %s", e.Error())
			err <- e
		}

		log.Debugf("Publishing: \n Service: %s InstanceId: %s Queue: %s Key: %s \n Message: \n %+v", p.name, p.instanceId, p.q, key, payload)
		err <- nil
	}(err)

	select {
	case ret := <-err:
		if ret != nil {
			return ret
		}
	case <-time.After(2 * time.Second):
		return fmt.Errorf("timout while publishing message for Service %s InstanceId %s Key %s", p.name, p.instanceId, key)
	}

	return nil
}

func (p *QueuePublisher) Close() error {

	err := p.pub.Close()
	if err != nil {
		log.Errorf("Closing publisher for Service: %s InstanceId: %s failed. Error: %s", p.name, p.instanceId, err.Error())
		return err
	}

	log.Infof("Closed publisher for Service: %s InstanceId: %s", p.name, p.instanceId)
	return err
}
