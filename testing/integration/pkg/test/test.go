/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package test

import (
	"context"
	"strconv"
	"testing"
)

type WorkflowExitHandlerFxn func(ctx context.Context, w *Workflow) error
type WorkflowSetupHandlerFxn func(t *testing.T, ctx context.Context, w *Workflow) error
type WorkflowStateCheckHandlerFxn func(ctx context.Context, w *Workflow) (bool, error)

type ExitHandlerFxn func(ctx context.Context, tc *TestCase) error
type SetupHandlerFxn func(t *testing.T, ctx context.Context, tc *TestCase) error
type StateHanderFxn func(ctx context.Context, tc *TestCase) (bool, error)
type TestFxn func(ctx context.Context, tc *TestCase) error

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
