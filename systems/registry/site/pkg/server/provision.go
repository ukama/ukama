/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 * Copyright (c) 2026-present, Ukama Inc.
 */
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/registry/site/pkg/db"
	"gorm.io/gorm"
)

const configAttempts = 3
const configTimeout = 60 * time.Second
const configPoll = time.Second

// A lost DELETE is re-sent on this period; in between the worker only polls state.
const cancelResend = 15 * time.Second

// Cancellation that is still unconfirmed this long after the configure
// deadline is abandoned so the operation cannot loop forever.
const cancelGiveUp = 5 * time.Minute

var errNodeOffboarded = errors.New("node offboarded")

// Per-pass cancel window; a variable so tests can shorten it.
var cancelWindow = configTimeout

type provisionStore interface {
	Create(context.Context, *db.Site, []string) (*db.SiteProvision, error)
	Get(context.Context, string) (*db.SiteProvision, error)
	Pending(context.Context) ([]db.SiteProvision, error)
	Save(context.Context, *db.SiteProvision) error
	Run(context.Context, string, func(context.Context, *db.SiteProvision) error) error
}

func (s *SiteServer) StartProvisioning(ctx context.Context, store *gorm.DB, nodeURL string) {
	s.provisions = db.NewProvisionRepo(store)
	s.provisionClient = &nodeProvisionClient{url: nodeURL, http: &http.Client{Timeout: 5 * time.Second}}
	go s.provisionWorker(ctx)
}

func (s *SiteServer) provisionWorker(ctx context.Context) {
	var running sync.Map
	ticker := time.NewTicker(configPoll)
	defer ticker.Stop()
	for {
		ops, err := s.provisions.Pending(ctx)
		if err != nil && ctx.Err() == nil {
			log.Errorf("List pending site configuration: %v", err)
		}
		for _, op := range ops {
			if _, busy := running.LoadOrStore(op.ID, true); busy {
				continue
			}
			go func(id string) {
				defer running.Delete(id)
				err := s.provisions.Run(ctx, id, func(work context.Context, op *db.SiteProvision) error { return s.runProvision(work, op) })
				if err != nil && ctx.Err() == nil {
					log.Errorf("Site configuration %s: %v", id, err)
				}
			}(op.ID)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func attemptNode(op *db.SiteProvision, index int, action string) provisionNode {
	return provisionNode{NodeID: op.Nodes[index], RequestID: fmt.Sprintf("%s-%d-%d", op.ID, op.Attempt, index),
		SiteID: op.ID, NetworkID: op.Site.NetworkId.String(), Action: action}
}

// Each goroutine handles one node; the coordinator owns the group retry.
func (s *SiteServer) waitNodes(ctx context.Context, op *db.SiteProvision, action string, deadline time.Time) error {
	ctx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()
	results := make(chan error, len(op.Nodes))
	for index := range op.Nodes {
		go func(index int) {
			results <- s.waitNode(ctx, attemptNode(op, index, action))
		}(index)
	}
	var first error
	offboarded := false
	for range op.Nodes {
		err := <-results
		if errors.Is(err, errNodeOffboarded) {
			offboarded = true
			continue
		}
		if err != nil && first == nil {
			first = err
		}
	}
	if first == nil && offboarded {
		return errNodeOffboarded
	}
	return first
}

func (s *SiteServer) waitNode(ctx context.Context, node provisionNode) error {
	ticker := time.NewTicker(configPoll)
	defer ticker.Stop()
	action := node.Action
	var dispatched time.Time
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		if action == "cancel" {
			node.Action = "status"
			if time.Since(dispatched) >= cancelResend {
				node.Action, dispatched = "cancel", time.Now()
			}
		}
		result, err := s.provisionClient.Reconcile(ctx, node)
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err == nil {
			if action == "configure" && result.Completed && !result.Cancelled {
				return nil
			}
			if action == "cancel" && result.Cleared {
				return nil
			}
			if action == "cancel" && result.Offboarded {
				return errNodeOffboarded
			}
			if action == "configure" && result.Cancelled {
				return fmt.Errorf("node %s attempt cancelled", node.NodeID)
			}
			if action == "configure" {
				node.Action = "status"
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (s *SiteServer) runProvision(ctx context.Context, op *db.SiteProvision) error {
	switch op.Phase {
	case "configuring":
		err := s.waitNodes(ctx, op, "configure", op.Deadline)
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err == nil {
			op.Phase = "creating"
		} else {
			op.Phase, op.Failure = "cancelling", err.Error()
		}
		return s.provisions.Save(ctx, op)
	case "cancelling":
		// Cleanup has its own retry window and never consumes another attempt.
		err := s.waitNodes(ctx, op, "cancel", time.Now().Add(cancelWindow))
		offboarded := errors.Is(err, errNodeOffboarded)
		if err != nil && !offboarded {
			if ctx.Err() == nil && time.Now().After(op.Deadline.Add(cancelGiveUp)) {
				op.Phase, op.Stop = "failed", true
				op.Failure = fmt.Sprintf("%s; cancel not confirmed: %v", op.Failure, err)
				return s.provisions.Save(ctx, op)
			}
			return err
		}
		if op.Stop || offboarded || op.Attempt == configAttempts {
			op.Phase = "failed"
		} else {
			op.Attempt++
			op.Phase = "configuring"
			op.Deadline = time.Now().UTC().Add(configTimeout)
		}
		return s.provisions.Save(ctx, op)
	case "creating":
		// The rest of Add runs only after the three-node barrier passed.
		err := s.siteRepo.Add(&op.Site, func(_ *db.Site, tx *gorm.DB) error { return db.MarkSitePublishing(tx, op.ID, op.Revision) })
		if err != nil {
			saved, readErr := s.provisions.Get(ctx, op.ID)
			if readErr != nil {
				return readErr
			}
			if saved.Phase != "creating" {
				return nil
			}
			op.Phase, op.Stop, op.Failure = "cancelling", true, err.Error()
			return s.provisions.Save(ctx, op)
		}
		return nil
	case "publishing":
		site, err := s.siteRepo.Get(op.Site.Id)
		if err != nil {
			return err
		}
		op.Site = *site

		if err := s.publishCreatedSite(&op.Site); err != nil {
			return err
		}
		op.Phase = "active"
		if err := s.provisions.Save(ctx, op); err != nil {
			return err
		}
		s.pushSiteCount(op.Site.NetworkId)
	}
	return nil
}

func (s *SiteServer) addProvisionedSite(ctx context.Context, site *db.Site, nodes []string) (*db.Site, error) {
	if s.provisions == nil {
		return nil, fmt.Errorf("site provisioning is not initialized")
	}
	op, err := s.provisions.Create(ctx, site, nodes)
	if err != nil {
		return nil, err
	}
	ticker := time.NewTicker(configPoll)
	defer ticker.Stop()
	for {
		switch op.Phase {
		case "active":
			return &op.Site, nil
		case "failed":
			return nil, fmt.Errorf("site configuration failed: %s", op.Failure)
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}
		op, err = s.provisions.Get(ctx, op.ID)
		if err != nil {
			return nil, err
		}
	}
}
