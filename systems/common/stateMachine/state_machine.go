package stateMachine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type Transition struct {
	Name           string   `json:"name"`
	Trigger        []string `json:"trigger"`
	ExpectedEvents []string `json:"expectedEvents,omitempty"`
	Timeout        int      `json:"timeout,omitempty"`
}

type SubState struct {
	Events     []string     `json:"events,omitempty"`
	Transition []Transition `json:"transition"`
}

type State struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Events      []string    `json:"events,omitempty"`
	Transition  []Transition `json:"transition"`
	SubState    *SubState   `json:"substate,omitempty"`
}

type StateMachineConfig struct {
	Version string  `json:"version"`
	Entity  string  `json:"entity"`
	File    string  `json:"file"`
	States  []State `json:"states"`
}

type StateMachine struct {
	config StateMachineConfig
}

func NewStateMachine(configPath string) (*StateMachine, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config StateMachineConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return &StateMachine{config: config}, nil
}

func (sm *StateMachine) GetNextState(currentState string, events []string) (string, error) {
	var stateObj *State
	for _, state := range sm.config.States {
		if strings.EqualFold(state.Name, currentState) {
			stateObj = &state
			break
		}
	}

	if stateObj == nil {
		validStates := make([]string, len(sm.config.States))
		for i, state := range sm.config.States {
			validStates[i] = state.Name
		}
		return "", fmt.Errorf("invalid current state: %s. Valid states are: %s", currentState, strings.Join(validStates, ", "))
	}

	nextState := stateObj.Name
	for _, event := range events {
		transitionFound := false

		// Check main state transitions
		for _, transition := range stateObj.Transition {
			if contains(transition.Trigger, event) {
				nextState = transition.Name
				transitionFound = true
				break
			}
		}

		// If no main state transition found, check substate transitions
		if !transitionFound && stateObj.SubState != nil {
			for _, transition := range stateObj.SubState.Transition {
				if contains(transition.Trigger, event) {
					transitionFound = true
					break
				}
			}
		}

		if transitionFound {
			// Update stateObj for the next iteration if state changed
			for i := range sm.config.States {
				if sm.config.States[i].Name == nextState {
					stateObj = &sm.config.States[i]
					break
				}
			}
		} else {
			fmt.Printf("No transition found for event: %s\n", event)
		}
	}

	return nextState, nil
}

// Helper function to check if a slice contains a certain item
func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

