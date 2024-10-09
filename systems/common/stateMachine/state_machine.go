/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package statemachine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Event struct {
	Name       string    `json:"name"`
	Timestamp  time.Time `json:"timestamp"`
	OldState   string    `json:"old_state"`
	NewState   string    `json:"new_state"`
	IsSubstate bool      `json:"is_substate"`
}

type TransitionCallback func(event Event)

type Transition struct {
	ToState string   `json:"to_state"`
	Trigger []string `json:"trigger"`
}

type SubState struct {
	Events      []string               `json:"events"`
	Transitions map[string]Transition `json:"transition"`
}

func (ss *SubState) UnmarshalJSON(data []byte) error {
	aux := struct {
		Events      []string     `json:"events"`
		Transitions []Transition `json:"transition"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	ss.Events = aux.Events
	ss.Transitions = make(map[string]Transition)
	for _, t := range aux.Transitions {
		for _, trigger := range t.Trigger {
			ss.Transitions[trigger] = t
		}
	}
	return nil
}

type State struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Events      []string              `json:"events"`
	Transitions map[string]Transition `json:"transition"`
	SubState    *SubState             `json:"substate,omitempty"`
	OnEnter     func()                `json:"-"`
	OnExit      func()                `json:"-"`
}

func (s *State) UnmarshalJSON(data []byte) error {
	aux := struct {
		Name        string       `json:"name"`
		Description string       `json:"description"`
		Events      []string     `json:"events"`
		Transitions []Transition `json:"transition"`
		SubState    *SubState    `json:"substate,omitempty"`
	}{
		Transitions: make([]Transition, 0),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	s.Name = aux.Name
	s.Description = aux.Description
	s.Events = aux.Events
	s.Transitions = make(map[string]Transition)
	for _, t := range aux.Transitions {
		for _, trigger := range t.Trigger {
			s.Transitions[trigger] = t
		}
	}
	s.SubState = aux.SubState
	return nil
}

type StateMachineConfig struct {
	Version string           `json:"version"`
	Entity  string           `json:"entity"`
	File    string           `json:"file"`
	States  map[string]State `json:"states"`
}

func (smc *StateMachineConfig) UnmarshalJSON(data []byte) error {
	aux := struct {
		Version string  `json:"version"`
		Entity  string  `json:"entity"`
		File    string  `json:"file"`
		States  []State `json:"states"` // Keep as an array
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	smc.Version = aux.Version
	smc.Entity = aux.Entity
	smc.File = aux.File
	smc.States = make(map[string]State)
	for _, state := range aux.States {
		smc.States[state.Name] = state
	}
	return nil
}

type StateMachine struct {
	handler TransitionCallback
	mu      sync.Mutex
}

type StateMachineInstance struct {
	InstanceID      string
	CurrentState    string
	CurrentSubstate string
	Config          StateMachineConfig
	StateMachine    *StateMachine
}

func NewStateMachine(handler TransitionCallback) *StateMachine {
	return &StateMachine{handler: handler}
}

func (sm *StateMachine) NewInstance(configFile, instanceID, initialState string) (*StateMachineInstance, error) {
	config, err := LoadConfig(configFile)
	if err != nil {
		return nil, err
	}

	if _, exists := config.States[initialState]; !exists {
		return nil, fmt.Errorf("invalid initial state: %s", initialState)
	}

	instance := &StateMachineInstance{
		InstanceID:      instanceID,
		CurrentState:    initialState,
		CurrentSubstate: "",
		Config:          config,
		StateMachine:    sm,
	}

	return instance, nil
}

// Transition processes a state transition for the instance.
func (instance *StateMachineInstance) Transition(eventName string) {
	instance.StateMachine.mu.Lock()
	defer instance.StateMachine.mu.Unlock()

	oldState := instance.CurrentState
	currentState, exists := instance.Config.States[instance.CurrentState]

	if !exists {
		log.Infof("Current state not found: %s\n", instance.CurrentState)
		return
	}

	if currentState.OnExit != nil {
		currentState.OnExit()
	}

	log.Infof("Current state: %s\n", currentState.Name)

	if transition, exists := currentState.Transitions[eventName]; exists {
		instance.CurrentState = transition.ToState
		event := Event{
			Name:       eventName,
			Timestamp:  time.Now(),
			OldState:   oldState,
			NewState:   instance.CurrentState,
			IsSubstate: false,
		}

		if instance.StateMachine.handler != nil {
			instance.StateMachine.handler(event)
		}

		if newState, exists := instance.Config.States[instance.CurrentState]; exists && newState.OnEnter != nil {
			newState.OnEnter()
		}
	} else if currentState.SubState != nil {
		if subTransition, exists := currentState.SubState.Transitions[eventName]; exists {
			instance.CurrentSubstate = subTransition.ToState
			event := Event{
				Name:       eventName,
				Timestamp:  time.Now(),
				OldState:   oldState,
				NewState:   instance.CurrentSubstate,
				IsSubstate: true,
			}

			if instance.StateMachine.handler != nil {
				instance.StateMachine.handler(event)
			}

			if newSubState, exists := instance.Config.States[instance.CurrentSubstate]; exists && newSubState.OnEnter != nil {
				newSubState.OnEnter()
			}
		} else {
			log.Infof("No substate transition found for event: %s in state: %s\n", eventName, instance.CurrentState)
		}
	} else {
		log.Infof("No transition found for event: %s in state: %s\n", eventName, instance.CurrentState)
	}
}

func LoadConfig(configFile string) (StateMachineConfig, error) {
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