/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package broker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"

	mb "github.com/ukama/ukama/systems/common/msgbus"
)

// Every service queue is bound here so msgclient can find the queues it owns after a restart.
const ownedExchange = "msgclient.owned"

type Broker struct {
	uri      string
	vhost    string
	mgmtUri  string
	user     string
	password string
	http     *http.Client

	Listen  *rabbitmq.Conn
	Publish *rabbitmq.Conn
}

// New blocks until both RabbitMQ connections are up.
func New(uri, mgmtUri, user, password string) (*Broker, error) {
	u, err := amqp.ParseURI(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid queue uri: %w", err)
	}

	b := &Broker{
		uri:      uri,
		vhost:    u.Vhost,
		mgmtUri:  strings.TrimSuffix(mgmtUri, "/"),
		user:     user,
		password: password,
		http:     &http.Client{Timeout: 10 * time.Second},
	}

	b.Listen = dial(uri, "listen")
	b.Publish = dial(uri, "publish")

	return b, nil
}

func dial(uri, name string) *rabbitmq.Conn {
	for {
		conn, err := rabbitmq.NewConn(uri, rabbitmq.WithConnectionOptionsLogging)
		if err == nil {
			return conn
		}
		log.Warnf("Failed to open %s connection to RabbitMQ at %s. Error %s. Retrying.", name, mb.RemovePassFromConnection(uri), err.Error())
		time.Sleep(2 * time.Second)
	}
}

func (b *Broker) Close() {
	_ = b.Listen.Close()
	_ = b.Publish.Close()
}

// EnsureQueue declares the service queue and makes its bindings on the event exchange match routes.
func (b *Broker) EnsureQueue(queue string, routes []string) error {
	current, err := b.queueRoutes(queue)
	if err != nil {
		return err
	}

	want := make(map[string]bool, len(routes))
	for _, r := range routes {
		want[r] = true
	}

	return b.withChannel(func(ch *amqp.Channel) error {
		if _, err := ch.QueueDeclare(queue, true, false, false, false, nil); err != nil {
			return fmt.Errorf("declare queue %s: %w", queue, err)
		}

		if err := ch.ExchangeDeclare(ownedExchange, "fanout", true, false, false, false, nil); err != nil {
			return fmt.Errorf("declare exchange %s: %w", ownedExchange, err)
		}

		if err := ch.QueueBind(queue, "", ownedExchange, false, nil); err != nil {
			return fmt.Errorf("bind queue %s to %s: %w", queue, ownedExchange, err)
		}

		for _, r := range routes {
			if current[r] {
				continue
			}
			if err := ch.QueueBind(queue, r, mb.DefaultExchange, false, nil); err != nil {
				return fmt.Errorf("bind queue %s to %s: %w", queue, r, err)
			}
		}

		for r := range current {
			if want[r] {
				continue
			}
			if err := ch.QueueUnbind(queue, r, mb.DefaultExchange, nil); err != nil {
				return fmt.Errorf("unbind queue %s from %s: %w", queue, r, err)
			}
			log.Infof("Unbound route %s from queue %s", r, queue)
		}

		return nil
	})
}

func (b *Broker) DeleteQueue(queue string) error {
	err := b.withChannel(func(ch *amqp.Channel) error {
		_, err := ch.QueueDelete(queue, false, false, false)
		return err
	})
	if err != nil {
		return fmt.Errorf("delete queue %s: %w", queue, err)
	}

	return b.SetMaxLength(queue, 0)
}

// SetMaxLength caps the queue through an operator policy; RabbitMQ drops the oldest event when full. 0 removes the cap.
func (b *Broker) SetMaxLength(queue string, max int64) error {
	path := "/api/operator-policies/" + url.PathEscape(b.vhost) + "/" + url.PathEscape("msgclient-"+queue)

	if max <= 0 {
		return b.mgmt(http.MethodDelete, path, nil, nil)
	}

	policy := map[string]interface{}{
		"pattern":    "^" + regexp.QuoteMeta(queue) + "$",
		"apply-to":   "queues",
		"priority":   0,
		"definition": map[string]interface{}{"max-length": max},
	}

	return b.mgmt(http.MethodPut, path, policy, nil)
}

// OwnedQueues lists the service queues msgclient created.
func (b *Broker) OwnedQueues() ([]string, error) {
	var bindings []struct {
		Destination string `json:"destination"`
	}

	path := "/api/exchanges/" + url.PathEscape(b.vhost) + "/" + ownedExchange + "/bindings/source"
	if err := b.mgmt(http.MethodGet, path, nil, &bindings); err != nil {
		return nil, err
	}

	queues := make([]string, 0, len(bindings))
	for _, bd := range bindings {
		queues = append(queues, bd.Destination)
	}

	return queues, nil
}

func (b *Broker) queueRoutes(queue string) (map[string]bool, error) {
	var bindings []struct {
		Source     string `json:"source"`
		RoutingKey string `json:"routing_key"`
	}

	path := "/api/queues/" + url.PathEscape(b.vhost) + "/" + url.PathEscape(queue) + "/bindings"
	if err := b.mgmt(http.MethodGet, path, nil, &bindings); err != nil {
		return nil, err
	}

	routes := make(map[string]bool)
	for _, bd := range bindings {
		if bd.Source == mb.DefaultExchange {
			routes[bd.RoutingKey] = true
		}
	}

	return routes, nil
}

func (b *Broker) withChannel(fn func(ch *amqp.Channel) error) error {
	conn, err := amqp.Dial(b.uri)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer func() { _ = ch.Close() }()

	return fn(ch)
}

// mgmt calls the RabbitMQ management API; a 404 counts as success with nothing returned.
func (b *Broker) mgmt(method, path string, body, out interface{}) error {
	var reader *strings.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = strings.NewReader(string(buf))
	} else {
		reader = strings.NewReader("")
	}

	req, err := http.NewRequest(method, b.mgmtUri+path, reader)
	if err != nil {
		return err
	}
	req.SetBasicAuth(b.user, b.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.http.Do(req)
	if err != nil {
		return fmt.Errorf("rabbitmq management %s %s: %w", method, path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("rabbitmq management %s %s: %s", method, path, resp.Status)
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}

	return nil
}
