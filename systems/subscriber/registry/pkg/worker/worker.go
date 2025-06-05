/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package worker

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg/client"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg/db"
	simMangerPb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
)

type DeletionWorker struct {
	subscriberRepo    db.SubscriberRepo
	simManagerService client.SimManagerClientProvider
	config            *pkg.DeletionWorkerConfig
	ticker            *time.Ticker
	ctx               context.Context
	cancel            context.CancelFunc
	done              chan struct{}
}

func NewDeletionWorker(
	subscriberRepo db.SubscriberRepo,
	simManagerService client.SimManagerClientProvider,
	config *pkg.DeletionWorkerConfig,
) *DeletionWorker {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &DeletionWorker{
		subscriberRepo:    subscriberRepo,
		simManagerService: simManagerService,
		config:           config,
		ticker:           time.NewTicker(config.CheckInterval),
		ctx:              ctx,
		cancel:           cancel,
		done:             make(chan struct{}),
	}
}

func (dw *DeletionWorker) Start() {
	log.Info("Starting deletion worker routine")
	
	go func() {
		defer close(dw.done)
		defer dw.ticker.Stop()
		
		for {
			select {
			case <-dw.ticker.C:
				dw.checkStuckDeletions()
			case <-dw.ctx.Done():
				log.Info("Deletion worker routine stopping")
				return
			}
		}
	}()
}

func (dw *DeletionWorker) Stop() {
	log.Info("Stopping deletion worker")
	dw.cancel()
	
	select {
	case <-dw.done:
		log.Info("Deletion worker stopped successfully")
	case <-time.After(5 * time.Second):
		log.Warn("Deletion worker stop timeout exceeded")
	}
}

func (dw *DeletionWorker) checkStuckDeletions() {
	threshold := time.Now().Add(-dw.config.DeletionTimeout)
	
	stuckSubscribers, err := dw.subscriberRepo.FindPendingDeletionBefore(threshold)
	if err != nil {
		log.Errorf("Error checking for stuck deletions: %v", err)
		return
	}
	
	if len(stuckSubscribers) == 0 {
		return
	}
	
	log.Infof("Found %d subscribers stuck in pending deletion state", len(stuckSubscribers))
	
	for _, subscriber := range stuckSubscribers {
		if subscriber.DeletionRetryCount >= dw.config.MaxRetries {
			log.Errorf("Subscriber %s has exceeded maximum retry attempts (%d). Manual intervention required.", 
				subscriber.SubscriberId, dw.config.MaxRetries)
			continue
		}
		
		log.Infof("Retrying deletion for subscriber %s (attempt %d/%d)", 
			subscriber.SubscriberId, subscriber.DeletionRetryCount+1, dw.config.MaxRetries)
		
		retryCtx := context.Background()
		go dw.retrySubscriberDeletion(retryCtx, subscriber)
	}
}

func (dw *DeletionWorker) retrySubscriberDeletion(ctx context.Context, subscriber db.Subscriber) {
	err := dw.subscriberRepo.IncrementDeletionRetry(subscriber.SubscriberId)
	if err != nil {
		log.Errorf("Failed to increment retry count for subscriber %s: %v", 
			subscriber.SubscriberId, err)
		return
	}
	
	simManagerClient, err := dw.simManagerService.GetSimManagerService()
	if err != nil {
		log.Errorf("Failed to get SimManagerClient for retry: %v", err)
		return
	}
	
	_, err = simManagerClient.TerminateSimsForSubscriber(ctx, &simMangerPb.TerminateSimsForSubscriberRequest{
		SubscriberId: subscriber.SubscriberId.String(),
	})
	
	if err != nil {
		log.Errorf("Retry failed for subscriber %s: %v", subscriber.SubscriberId, err)
		
		if subscriber.DeletionRetryCount+1 >= dw.config.MaxRetries {
			log.Errorf("Subscriber %s deletion failed after %d attempts. Manual intervention required.", 
				subscriber.SubscriberId, dw.config.MaxRetries)
		}
	} else {
		log.Infof("Retry successful for subscriber %s", subscriber.SubscriberId)
	}
}

