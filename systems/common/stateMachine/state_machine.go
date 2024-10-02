/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package stateMachine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

type Event string
type StateID string

type TransitionCallback func(oldState, event, newState string)

type Transition struct {
	ToState string   `json:"to_state"`
	Trigger []string `json:"trigger"`
}

type SubState struct {
	Events      []string               `json:"events"`
	Transitions map[string]Transition `json:"transition"`
}
type State struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Events      []string              `json:"events"`
	Transitions map[string]Transition `json:"transition"`
	SubState    *SubState             `json:"substate,omitempty"`
}

type StateMachineConfig struct {
	Version string           `json:"version"`
	Entity  string           `json:"entity"`
	File    string           `json:"file"`
	States  map[string]State `json:"states"`
}
type StateMachine struct {
	handler TransitionCallback
}

type StateMachineInstance struct {
	InstanceID    string
	CurrentState  string
	Config        StateMachineConfig
	StateMachine  *StateMachine
}

func (ss *SubState) UnmarshalJSON(data []byte) error {
	type Alias SubState
	aux := &struct {
		Transitions []Transition `json:"transition"`
		*Alias
	}{
		Alias: (*Alias)(ss),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	ss.Transitions = make(map[string]Transition)
	for _, t := range aux.Transitions {
		for _, trigger := range t.Trigger {
			ss.Transitions[trigger] = t
		}
	}
	return nil
}

func (s *State) UnmarshalJSON(data []byte) error {
	type Alias State
	aux := &struct {
		Transitions []Transition `json:"transition"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	s.Transitions = make(map[string]Transition)
	for _, t := range aux.Transitions {
		for _, trigger := range t.Trigger {
			s.Transitions[trigger] = t
		}
	}
	return nil
}


func (smc *StateMachineConfig) UnmarshalJSON(data []byte) error {
	type Alias StateMachineConfig
	aux := &struct {
		States []State `json:"states"`
		*Alias
	}{
		Alias: (*Alias)(smc),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	smc.States = make(map[string]State)
	for _, state := range aux.States {
		smc.States[state.Name] = state
	}
	return nil
}


func NewStateMachine(handler TransitionCallback) *StateMachine {
	return &StateMachine{handler: handler}
}

func (sm *StateMachine) NewInstance(configFile, instanceID, initialState string) (*StateMachineInstance, error) {
	config, err := loadConfig(configFile)
	if err != nil {
		return nil, err
	}

	if _, exists := config.States[initialState]; !exists {
		log.WithError(fmt.Errorf("invalid initial state: %s", initialState)).Error("Invalid initial state provided")
		return nil, fmt.Errorf("invalid initial state: %s", initialState)
	}

	return &StateMachineInstance{
		InstanceID:    instanceID,
		CurrentState:  initialState,
		Config:        config,
		StateMachine:  sm,
	}, nil
}

func loadConfig(configFile string) (StateMachineConfig, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return StateMachineConfig{}, fmt.Errorf("error reading config file: %v", err)
	}

	var config StateMachineConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return StateMachineConfig{}, fmt.Errorf("error parsing JSON: %v", err)
	}

	return config, nil
}

func (instance *StateMachineInstance) Transition(event string) error {
	oldState := instance.CurrentState
	currentState, exists := instance.Config.States[instance.CurrentState]
	if !exists {
		return fmt.Errorf("current state not found: %s", instance.CurrentState)
	}

	log.Infof("Current state: %s\n", currentState.Name)

	if transition, exists := currentState.Transitions[event]; exists {
		log.Infof("Found transition for event '%s' in state '%s' -> transitioning to state '%s'\n", event, currentState.Name, transition.ToState)

		instance.CurrentState = transition.ToState
		log.Infof("State updated to: %s\n", instance.CurrentState)

		if instance.StateMachine.handler != nil {
			instance.StateMachine.handler(oldState, event, instance.CurrentState)
		}
		return nil
	}

	return fmt.Errorf("no transition found for event: %s in state: %s", event, instance.CurrentState)
}

