/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	stm "github.com/ukama/ukama/systems/common/stateMachine"
)

const TimeoutEventName = "timeout"

func (n *StateEventServer) applyTimeout(ctx context.Context, instance *stm.StateMachineInstance,
	nodeId string, enteredAt, now time.Time) (bool, error) {
	prevState := instance.CurrentState

	moved, err := instance.TimeoutTransition(enteredAt, now)
	if err != nil {
		return false, fmt.Errorf("failed to apply timeout for node %s in state %s: %w",
			nodeId, prevState, err)
	}

	if !moved {
		return false, nil
	}

	log.Infof("Node %s timed out in state %s after %s, advancing to %s",
		nodeId, prevState, now.Sub(enteredAt), instance.CurrentState)

	return n.persistTransition(ctx, instance, nodeId, prevState, TimeoutEventName)
}

func (n *StateEventServer) RunTimeouts(ctx context.Context, now time.Time) (int, error) {
	states, err := n.s.ListLatestStates()
	if err != nil {
		return 0, fmt.Errorf("failed to list latest node states: %w", err)
	}

	advanced := 0

	for i := range states {
		state := states[i]

		if state.NodeId == "" {
			continue
		}

		substate := ""
		if len(state.SubState) > 0 {
			substate = state.SubState[len(state.SubState)-1]
		}

		moved, err := n.runNodeTimeout(ctx, state.NodeId, state.CurrentState.String(),
			substate, state.CreatedAt, now)
		if err != nil {
			log.Errorf("Failed to run timeout for node %s: %v", state.NodeId, err)

			continue
		}

		if moved {
			advanced++
		}
	}

	return advanced, nil
}

func (n *StateEventServer) runNodeTimeout(ctx context.Context, nodeId, storedState, storedSubstate string,
	enteredAt, now time.Time) (bool, error) {
	mutexValue, _ := n.processingMutex.LoadOrStore(nodeId, &sync.Mutex{})
	mutex := mutexValue.(*sync.Mutex)

	mutex.Lock()
	defer mutex.Unlock()

	instance, err := n.getOrCreateInstance(nodeId, storedState, storedSubstate)
	if err != nil {
		return false, fmt.Errorf("failed to create state machine instance: %w", err)
	}

	return n.applyTimeout(ctx, instance, nodeId, enteredAt, now)
}

func (n *StateEventServer) StartTimeoutWorker(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		log.Infof("Node state timeout worker disabled")

		return
	}

	log.Infof("Starting node state timeout worker with interval %s", interval)

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Infof("Stopping node state timeout worker")

				return
			case <-ticker.C:
				advanced, err := n.RunTimeouts(ctx, time.Now().UTC())
				if err != nil {
					log.Errorf("Node state timeout sweep failed: %v", err)

					continue
				}

				if advanced > 0 {
					log.Infof("Node state timeout sweep advanced %d node(s)", advanced)
				}
			}
		}
	}()
}
