package test

import (
	"context"
	"strconv"
)

type WorkflowExitHandlerFxn func(ctx context.Context, w *Workflow) error
type WorkflowSetupHandlerFxn func(ctx context.Context, w *Workflow) error
type WorkflowStateCheckHandlerFxn func(ctx context.Context, w *Workflow) (bool, error)

type ExitHandlerFxn func(ctx context.Context, t *TestCase) error
type SetupHandlerFxn func(ctx context.Context, t *TestCase) error
type StateHanderFxn func(ctx context.Context, t *TestCase) (bool, error)
type TestFxn func(ctx context.Context, t *TestCase) error

type StateType uint8

const (
	StateTypeUnknown   StateType = iota
	StateTypeWaiting             = 1
	StateTypeUnderTest           = 2
	StateTypePass                = 3
	StateTypeFail                = 4
	StateTypeInvalid             = 5
	StateTypeTested              = 6
)

func (s StateType) String() string {
	t := map[StateType]string{0: "unknown", 1: "waiting", 2: "under_test", 3: "pass", 4: "fail", 5: "invalid", 6: "tested"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseType(value string) StateType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return StateType(i)
	}

	t := map[string]StateType{"unknown": 0, "waiting": 1, "under_test": 2, "pass": 3, "fail": 4, "invalid": 5, "tested": 6}

	v, ok := t[value]
	if !ok {
		return StateType(0)
	}

	return StateType(v)
}
