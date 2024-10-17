/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package statemachine

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNewStateMachine(t *testing.T) {
	handler := func(event Event) {}
	sm := NewStateMachine(handler)
	if sm == nil {
		t.Error("NewStateMachine returned nil")
	}
}

func TestNewInstance(t *testing.T) {
	// Create a temporary config file
	config := `{
		"version": "1.0",
		"entity": "test",
		"file": "test.json",
		"states": [
			{
				"name": "initial",
				"description": "Initial state",
				"events": ["start"],
				"transition": [
					{
						"to_state": "running",
						"trigger": ["start"]
					}
				]
			},
			{
				"name": "running",
				"description": "Running state",
				"events": ["stop"],
				"transition": [
					{
						"to_state": "stopped",
						"trigger": ["stop"]
					}
				]
			},
			{
				"name": "stopped",
				"description": "Stopped state",
				"events": ["start"],
				"transition": [
					{
						"to_state": "running",
						"trigger": ["start"]
					}
				]
			}
		]
	}`

	tmpfile, err := ioutil.TempFile("", "test_config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(config)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	sm := NewStateMachine(nil)
	instance, err := sm.NewInstance(tmpfile.Name(), "test-instance", "initial")
	if err != nil {
		t.Errorf("NewInstance failed: %v", err)
	}

	if instance.CurrentState != "initial" {
		t.Errorf("Expected initial state to be 'initial', got '%s'", instance.CurrentState)
	}
}

func TestTransition(t *testing.T) {
	config := StateMachineConfig{
		Version: "1.0",
		Entity:  "test",
		File:    "test.json",
		States: map[string]State{
			"initial": {
				Name:        "initial",
				Description: "Initial state",
				Events:      []string{"start"},
				Transitions: map[string]Transition{
					"start": {ToState: "running", Trigger: []string{"start"}},
				},
			},
			"running": {
				Name:        "running",
				Description: "Running state",
				Events:      []string{"stop"},
				Transitions: map[string]Transition{
					"stop": {ToState: "stopped", Trigger: []string{"stop"}},
				},
			},
			"stopped": {
				Name:        "stopped",
				Description: "Stopped state",
				Events:      []string{"start"},
				Transitions: map[string]Transition{
					"start": {ToState: "running", Trigger: []string{"start"}},
				},
			},
		},
	}

	var lastEvent Event
	handler := func(event Event) {
		lastEvent = event
	}

	sm := NewStateMachine(handler)
	instance := &StateMachineInstance{
		InstanceID:   "test-instance",
		CurrentState: "initial",
		Config:       config,
		StateMachine: sm,
	}

	// Test valid transition
	err := instance.Transition("start")
	if err != nil {
		t.Errorf("Transition failed: %v", err)
	}
	if instance.CurrentState != "running" {
		t.Errorf("Expected state to be 'running', got '%s'", instance.CurrentState)
	}
	if lastEvent.Name != "start" || lastEvent.OldState != "initial" || lastEvent.NewState != "running" {
		t.Errorf("Unexpected event: %+v", lastEvent)
	}

	// Test invalid transition
	err = instance.Transition("invalid")
	if err == nil {
		t.Error("Expected error for invalid transition, got nil")
	}
}

func TestLoadConfig(t *testing.T) {
	config := `{
		"version": "1.0",
		"entity": "test",
		"file": "test.json",
		"states": [
			{
				"name": "initial",
				"description": "Initial state",
				"events": ["start"],
				"transition": [
					{
						"to_state": "running",
						"trigger": ["start"]
					}
				]
			}
		]
	}`

	tmpfile, err := ioutil.TempFile("", "test_config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(config)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	loadedConfig, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Errorf("LoadConfig failed: %v", err)
	}

	if loadedConfig.Version != "1.0" || loadedConfig.Entity != "test" || loadedConfig.File != "test.json" {
		t.Errorf("Unexpected config values: %+v", loadedConfig)
	}

	if _, exists := loadedConfig.States["initial"]; !exists {
		t.Error("Expected 'initial' state not found in loaded config")
	}
}

func TestSubStateTransition(t *testing.T) {
	config := StateMachineConfig{
		Version: "1.0",
		Entity:  "test",
		File:    "test.json",
		States: map[string]State{
			"main": {
				Name:        "main",
				Description: "Main state with substates",
				Events:      []string{"sub1", "sub2"},
				SubState: &SubState{
					Events: []string{"sub1", "sub2"},
					Transitions: map[string]Transition{
						"sub1": {ToState: "substate1", Trigger: []string{"sub1"}},
						"sub2": {ToState: "substate2", Trigger: []string{"sub2"}},
					},
				},
			},
		},
	}

	var lastEvent Event
	handler := func(event Event) {
		lastEvent = event
	}

	sm := NewStateMachine(handler)
	instance := &StateMachineInstance{
		InstanceID:   "test-instance",
		CurrentState: "main",
		Config:       config,
		StateMachine: sm,
	}

	// Test substate transition
	err := instance.Transition("sub1")
	if err != nil {
		t.Errorf("Substate transition failed: %v", err)
	}
	if instance.CurrentSubstate != "substate1" {
		t.Errorf("Expected substate to be 'substate1', got '%s'", instance.CurrentSubstate)
	}
	if lastEvent.Name != "sub1" || lastEvent.OldState != "main" || lastEvent.NewSubstate != "substate1" || !lastEvent.IsSubstate {
		t.Errorf("Unexpected event for substate transition: %+v", lastEvent)
	}
}