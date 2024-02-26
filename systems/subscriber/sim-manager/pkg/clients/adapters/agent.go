/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package adapters

import (
	"context"
	"time"

	"github.com/ukama/ukama/systems/common/ukama"

	log "github.com/sirupsen/logrus"
)

type AgentAdapter interface {
	BindSim(context.Context, string) (any, error)
	GetSim(context.Context, string) (any, error)
	GetUsages(context.Context, string, string, string, string, string) (any, any, error)
	ActivateSim(context.Context, string) error
	DeactivateSim(context.Context, string) error
	TerminateSim(context.Context, string) error
	Close()
}

type AgentFactory interface {
	GetAgentAdapter(ukama.SimType) (AgentAdapter, bool)
}

type agentFactory struct {
	timeout time.Duration
	factory map[ukama.SimType]AgentAdapter
}

func NewAgentFactory(testAgentHost, operatorAgentHost string, timeout time.Duration, debug bool) *agentFactory {
	// we should lookup from provided config to get {realHost, realAgent, timeout} mappings
	// in order to dynamically fill the factory map with available running agents

	// for each {realHost, realAgent, timeout}}
	// agent, err := NewRealAgent(realHost, timeout)
	// handle err
	// factory[SimTypeForAgent] = agent

	// For now we will only use TestAgent for test sim type
	tAgent, err := NewTestAgentAdapter(testAgentHost, timeout)
	if err != nil {
		log.Fatalf("Failed to connect to Agent service at %s. Error: %v", testAgentHost, err)
	}

	// And OperatorAgent for telna sim type
	opAgent, err := NewOperatorAgentAdapter(operatorAgentHost, debug)
	if err != nil {
		log.Fatalf("Failed to connect to Agent service at %s. Error: %v", operatorAgentHost, err)
	}

	var factory = make(map[ukama.SimType]AgentAdapter)

	factory[ukama.SimTypeTest] = tAgent
	factory[ukama.SimTypeOperatorData] = opAgent

	return &agentFactory{
		timeout: timeout,
		factory: factory,
	}
}

func (a *agentFactory) GetAgentAdapter(simType ukama.SimType) (AgentAdapter, bool) {
	agent, ok := a.factory[simType]

	return agent, ok
}

func (a *agentFactory) Close() {
	for _, adapter := range a.factory {
		adapter.Close()
	}
}
