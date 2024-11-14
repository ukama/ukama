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
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Event struct {
	Name        string    `json:"name"`
	Timestamp   time.Time `json:"timestamp"`
	InstanceID  string    `json:"instance_id"`
	OldState    string    `json:"old_state"`
	NewState    string    `json:"new_state"`
	OldSubstate string    `json:"old_substate"`
	NewSubstate string    `json:"new_substate"`
}

type TransitionCallback func(event Event)

type Transition struct {
	ToState string   `json:"to_state"`
	Trigger []string `json:"trigger"`
}

type SubState struct {
	Events      []string              `json:"events"`
	Transitions map[string]Transition `json:"transition"`
}

type State struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Events      []string              `json:"events"`
	Transitions map[string]Transition `json:"transition"`
	SubState    *SubState             `json:"substate,omitempty"`
	OnEnter     func() error          `json:"-"`
	OnExit      func() error          `json:"-"`
}

type StateMachineConfig struct {
	Version string           `json:"version"`
	Entity  string           `json:"entity"`
	File    string           `json:"file"`
	States  map[string]State `json:"states"`
}

type configCache struct {
	configs map[string]StateMachineConfig
	mu      sync.RWMutex
}

var (
	cache     *configCache
	cacheOnce sync.Once
)

func getConfigCache() *configCache {
	cacheOnce.Do(func() {
		cache = &configCache{
			configs: make(map[string]StateMachineConfig),
		}
	})
	return cache
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

func (smc *StateMachineConfig) UnmarshalJSON(data []byte) error {
	aux := struct {
		Version string  `json:"version"`
		Entity  string  `json:"entity"`
		File    string  `json:"file"`
		States  []State `json:"states"`
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
	mu      sync.RWMutex
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

func validateStateTransitions(config StateMachineConfig) error {
	for stateName, state := range config.States {
		for trigger, transition := range state.Transitions {
			if _, exists := config.States[transition.ToState]; !exists {
				return fmt.Errorf("state '%s' has transition to non-existent state '%s' on trigger '%s'",
					stateName, transition.ToState, trigger)
			}
		}

		if state.SubState != nil {
			for trigger := range state.SubState.Transitions {
				if trigger == "" {
					return fmt.Errorf("empty trigger found in substate transitions for state '%s'", stateName)
				}
			}
		}
	}
	return nil
}

func validateEventConsistency(config StateMachineConfig) error {
	for stateName, state := range config.States {
		declaredEvents := make(map[string]bool)
		for _, event := range state.Events {
			declaredEvents[event] = true
		}
		if state.SubState != nil {
			for _, event := range state.SubState.Events {
				declaredEvents[event] = true
			}
		}

		for trigger := range state.Transitions {
			if !declaredEvents[trigger] && trigger != "enter" {
				return fmt.Errorf("state '%s' uses undeclared event '%s' in transition",
					stateName, trigger)
			}
		}

		if state.SubState != nil {
			for trigger := range state.SubState.Transitions {
				if !declaredEvents[trigger] && trigger != "enter" {
					return fmt.Errorf("state '%s' substate uses undeclared event '%s' in transition",
						stateName, trigger)
				}
			}
		}
	}
	return nil
}

func validateEventUniqueness(config StateMachineConfig) error {
	for stateName, state := range config.States {
		eventCount := make(map[string]int)
		for _, event := range state.Events {
			eventCount[event]++
			if eventCount[event] > 1 {
				return fmt.Errorf("duplicate event '%s' found in state '%s'",
					event, stateName)
			}
		}

		if state.SubState != nil {
			substateEventCount := make(map[string]int)
			for _, event := range state.SubState.Events {
				substateEventCount[event]++
				if substateEventCount[event] > 1 {
					return fmt.Errorf("duplicate event '%s' found in substate of state '%s'",
						event, stateName)
				}
			}
		}
	}
	return nil
}

func (sm *StateMachine) NewInstance(configFile, instanceID, initialState string) (*StateMachineInstance, error) {
	config, err := LoadConfig(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
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

func (instance *StateMachineInstance) Transition(eventName string) error {
	instance.StateMachine.mu.Lock()
	defer instance.StateMachine.mu.Unlock()

	oldState := instance.CurrentState
	oldSubstate := instance.CurrentSubstate

	currentState, exists := instance.Config.States[instance.CurrentState]
	if !exists {
		return fmt.Errorf("current state not found: %s", instance.CurrentState)
	}

	log.Infof("Processing event %s in state %s (substate: %s)",
		eventName, currentState.Name, instance.CurrentSubstate)

	newMainState := instance.CurrentState
	newSubState := instance.CurrentSubstate

	mainStateTransition, hasMainTransition := currentState.Transitions[eventName]

	var substateTransition Transition
	var hasSubstateTransition bool
	if currentState.SubState != nil {
		substateTransition, hasSubstateTransition = currentState.SubState.Transitions[eventName]
	}

	if hasMainTransition {
		if currentState.OnExit != nil {
			if err := currentState.OnExit(); err != nil {
				return fmt.Errorf("error in OnExit for state %s: %w", instance.CurrentState, err)
			}
		}
		newMainState = mainStateTransition.ToState
		log.Infof("Main state transition: %s -> %s", currentState.Name, newMainState)
	}

	if hasSubstateTransition {
		newSubState = substateTransition.ToState
		log.Infof("Substate transition: %s -> %s", instance.CurrentSubstate, newSubState)
	}

	if newMainState != instance.CurrentState {
		newState, exists := instance.Config.States[newMainState]
		if !exists {
			return fmt.Errorf("new state not found: %s", newMainState)
		}

		if newState.OnEnter != nil {
			if err := newState.OnEnter(); err != nil {
				return fmt.Errorf("error in OnEnter for state %s: %w", newMainState, err)
			}
		}

		if newState.SubState == nil {
			newSubState = ""
		}
	}

	instance.CurrentState = newMainState
	instance.CurrentSubstate = newSubState

	if instance.StateMachine.handler != nil {
		instance.StateMachine.handler(Event{
			Name:        eventName,
			Timestamp:   time.Now(),
			InstanceID:  instance.InstanceID,
			OldState:    oldState,
			NewState:    instance.CurrentState,
			OldSubstate: oldSubstate,
			NewSubstate: instance.CurrentSubstate,
		})
	}

	return nil
}

func LoadConfig(configFile string) (StateMachineConfig, error) {
	cache := getConfigCache()

	cache.mu.RLock()
	if config, exists := cache.configs[configFile]; exists {
		cache.mu.RUnlock()
		return config, nil
	}
	cache.mu.RUnlock()

	cache.mu.Lock()
	defer cache.mu.Unlock()

	if config, exists := cache.configs[configFile]; exists {
		return config, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return StateMachineConfig{}, fmt.Errorf("error reading config file: %v", err)
	}

	var config StateMachineConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return StateMachineConfig{}, fmt.Errorf("error parsing JSON: %v", err)
	}

	if err := validateStateTransitions(config); err != nil {
		return StateMachineConfig{}, fmt.Errorf("invalid state transitions: %v", err)
	}

	if err := validateEventConsistency(config); err != nil {
		return StateMachineConfig{}, fmt.Errorf("event consistency error: %v", err)
	}

	if err := validateEventUniqueness(config); err != nil {
		return StateMachineConfig{}, fmt.Errorf("event uniqueness error: %v", err)
	}

	cache.configs[configFile] = config
	return config, nil
}
