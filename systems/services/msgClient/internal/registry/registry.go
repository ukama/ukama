/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package registry

import (
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Topology interface {
	EnsureQueue(queue string, routes []string) error
	DeleteQueue(queue string) error
	SetMaxLength(queue string, max int64) error
	OwnedQueues() ([]string, error)
}

type Listener interface {
	Stop()
}

type StartListener func(service, uri string, onUnreachable func()) (Listener, error)

type Config struct {
	RefreshInterval time.Duration
	MissedRefreshes uint32
	QueueMaxLength  int64
	FlushAfter      time.Duration
}

type service struct {
	mu       sync.Mutex
	name     string
	uuid     string
	instance string
	uri      string
	routes   []string
	ready    bool
	listener Listener
	lastSeen time.Time
}

// Registry is msgclient's in-memory view of registered services; RabbitMQ holds the durable state.
type Registry struct {
	mu       sync.Mutex
	services map[string]*service
	uuids    map[string]*service
	owned    map[string]time.Time

	topo  Topology
	start StartListener
	cfg   Config
	now   func() time.Time
}

func New(topo Topology, start StartListener, cfg Config) *Registry {
	return &Registry{
		services: make(map[string]*service),
		uuids:    make(map[string]*service),
		owned:    make(map[string]time.Time),
		topo:     topo,
		start:    start,
		cfg:      cfg,
		now:      time.Now,
	}
}

// LoadOwned picks up queues created before a restart; their flush timer starts now.
func (r *Registry) LoadOwned() error {
	queues, err := r.topo.OwnedQueues()
	if err != nil {
		return err
	}

	now := r.now()
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, q := range queues {
		r.owned[q] = now
	}

	log.Infof("Found %d service queues", len(queues))
	return nil
}

func (r *Registry) Register(name, system, instance, uri string, routes []string) (string, error) {
	routes = normalize(routes)
	s := r.service(name, system)

	s.mu.Lock()
	defer s.mu.Unlock()

	now := r.now()
	r.mu.Lock()
	s.lastSeen = now
	s.instance = instance
	if len(routes) > 0 {
		r.owned[name] = now
	}
	r.mu.Unlock()

	if s.uri != uri {
		s.stopListener()
		s.uri = uri
	}

	if !s.ready || !slices.Equal(s.routes, routes) {
		switch {
		case len(routes) > 0:
			if err := r.topo.EnsureQueue(name, routes); err != nil {
				return "", err
			}
			if err := r.topo.SetMaxLength(name, r.cfg.QueueMaxLength); err != nil {
				log.Warnf("[%s] Failed to set queue size limit. Error %s", name, err.Error())
			}
		case len(s.routes) > 0:
			if err := r.topo.EnsureQueue(name, nil); err != nil {
				return "", err
			}
			s.stopListener()
		}

		s.routes = routes
		s.ready = true
	}

	if len(routes) > 0 && s.listener == nil {
		if err := r.startListener(s); err != nil {
			return "", err
		}
	}

	return s.uuid, nil
}

func (r *Registry) Resume(id string) error {
	s := r.byUuid(id)
	if s == nil {
		return fmt.Errorf("service %s not registered", id)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ready && len(s.routes) > 0 && s.listener == nil {
		return r.startListener(s)
	}

	return nil
}

func (r *Registry) Pause(id string) error {
	s := r.byUuid(id)
	if s == nil {
		return fmt.Errorf("service %s not registered", id)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.stopListener()

	return nil
}

// Unregister removes the service and deletes its queue.
func (r *Registry) Unregister(id string) error {
	s := r.byUuid(id)
	if s == nil {
		return fmt.Errorf("service %s not registered", id)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.stopListener()

	if err := r.topo.DeleteQueue(s.name); err != nil {
		return err
	}

	r.mu.Lock()
	delete(r.owned, s.name)
	r.mu.Unlock()

	s.routes = nil
	s.ready = false

	log.Infof("[%s] Unregistered, queue deleted", s.name)
	return nil
}

// Lookup returns the service name and instance for a uuid, if registered since msgclient started.
func (r *Registry) Lookup(id string) (string, string, bool) {
	s := r.byUuid(id)
	if s == nil {
		return "", "", false
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	return s.name, s.instance, true
}

// Run pauses services whose registration lapsed and flushes queues unregistered for too long.
func (r *Registry) Run(stop <-chan struct{}) {
	t := time.NewTicker(r.cfg.RefreshInterval)
	defer t.Stop()

	for {
		select {
		case <-stop:
			return
		case <-t.C:
			r.sweep()
		}
	}
}

func (r *Registry) Close() {
	r.mu.Lock()
	services := make([]*service, 0, len(r.services))
	for _, s := range r.services {
		services = append(services, s)
	}
	r.mu.Unlock()

	for _, s := range services {
		s.mu.Lock()
		s.stopListener()
		s.mu.Unlock()
	}
}

func (r *Registry) sweep() {
	now := r.now()
	lease := time.Duration(r.cfg.MissedRefreshes) * r.cfg.RefreshInterval

	r.mu.Lock()
	var lapsed []*service
	for _, s := range r.services {
		if now.Sub(s.lastSeen) > lease {
			lapsed = append(lapsed, s)
		}
	}

	var flush []string
	if r.cfg.FlushAfter > 0 {
		for q, seen := range r.owned {
			if now.Sub(seen) > r.cfg.FlushAfter {
				flush = append(flush, q)
			}
		}
	}
	r.mu.Unlock()

	for _, s := range lapsed {
		go r.expire(s)
	}

	for _, q := range flush {
		r.flush(q, now)
	}
}

func (r *Registry) expire(s *service) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.listener != nil {
		log.Warnf("[%s] Registration lapsed, pausing delivery", s.name)
		s.stopListener()
	}
}

func (r *Registry) flush(queue string, now time.Time) {
	r.mu.Lock()
	s := r.services[queue]
	r.mu.Unlock()

	if s != nil {
		s.mu.Lock()
		defer s.mu.Unlock()
	}

	r.mu.Lock()
	seen, ok := r.owned[queue]
	r.mu.Unlock()
	if !ok || now.Sub(seen) <= r.cfg.FlushAfter {
		return
	}

	if s != nil {
		s.stopListener()
	}

	if err := r.topo.DeleteQueue(queue); err != nil {
		log.Errorf("[%s] Failed to flush queue. Error %s", queue, err.Error())
		return
	}

	r.mu.Lock()
	delete(r.owned, queue)
	r.mu.Unlock()

	if s != nil {
		s.routes = nil
		s.ready = false
	}

	log.Warnf("[%s] Queue flushed, service not registered for %s", queue, r.cfg.FlushAfter)
}

// startListener expects s.mu held.
func (r *Registry) startListener(s *service) error {
	var l Listener
	l, err := r.start(s.name, s.uri, func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		if s.listener == l {
			s.stopListener()
		}
	})
	if err != nil {
		return err
	}

	s.listener = l
	return nil
}

// stopListener expects s.mu held.
func (s *service) stopListener() {
	if s.listener != nil {
		s.listener.Stop()
		s.listener = nil
	}
}

func (r *Registry) service(name, system string) *service {
	r.mu.Lock()
	defer r.mu.Unlock()

	s, ok := r.services[name]
	if !ok {
		s = &service{
			name: name,
			uuid: uuid.NewSHA1(uuid.NameSpaceOID, []byte(system+"/"+name)).String(),
		}
		r.services[name] = s
		r.uuids[s.uuid] = s
	}

	return s
}

func (r *Registry) byUuid(id string) *service {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.uuids[id]
}

func normalize(routes []string) []string {
	out := slices.Clone(routes)
	slices.Sort(out)
	return slices.Compact(out)
}
