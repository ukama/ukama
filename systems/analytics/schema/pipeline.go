/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package schema

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

// Internal pipeline routing key templates (8-segment repo convention:
// <type>.<source>.<scope>.<org>.<system>.<service>.<object>.<action>).
// The ledger remains the source of truth; these events are the fast path.
const (
	WindowReadyRoute = "event.cloud.local.{{ .Org}}.analytics.ingest.window.ready"
	KpiComputedRoute = "event.cloud.local.{{ .Org}}.analytics.analysis.kpi.computed"
)

// WindowReady is the payload of the window.ready event.
type WindowReady struct {
	DatasetKey string
	WindowID   int64
	OrgName    string
}

// KpiComputed is the payload of the kpi.computed event.
type KpiComputed struct {
	KpiKey   string
	WindowID int64
	OrgName  string
}

// Payloads travel as structpb.Struct (a protobuf well-known type) so no
// custom proto codegen is needed and any consumer can decode them.

func (w WindowReady) ToStruct() (*structpb.Struct, error) {
	return structpb.NewStruct(map[string]interface{}{
		"dataset_key": w.DatasetKey,
		"window_id":   strconv.FormatInt(w.WindowID, 10),
		"org":         w.OrgName,
	})
}

func (k KpiComputed) ToStruct() (*structpb.Struct, error) {
	return structpb.NewStruct(map[string]interface{}{
		"kpi_key":   k.KpiKey,
		"window_id": strconv.FormatInt(k.WindowID, 10),
		"org":       k.OrgName,
	})
}

func UnmarshalWindowReady(msg *anypb.Any) (*WindowReady, error) {
	s := &structpb.Struct{}
	if err := msg.UnmarshalTo(s); err != nil {
		return nil, fmt.Errorf("unmarshal window.ready: %w", err)
	}

	f := s.GetFields()
	id, err := strconv.ParseInt(f["window_id"].GetStringValue(), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("window.ready window_id: %w", err)
	}

	return &WindowReady{
		DatasetKey: f["dataset_key"].GetStringValue(),
		WindowID:   id,
		OrgName:    f["org"].GetStringValue(),
	}, nil
}

func UnmarshalKpiComputed(msg *anypb.Any) (*KpiComputed, error) {
	s := &structpb.Struct{}
	if err := msg.UnmarshalTo(s); err != nil {
		return nil, fmt.Errorf("unmarshal kpi.computed: %w", err)
	}

	f := s.GetFields()
	id, err := strconv.ParseInt(f["window_id"].GetStringValue(), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("kpi.computed window_id: %w", err)
	}

	return &KpiComputed{
		KpiKey:   f["kpi_key"].GetStringValue(),
		WindowID: id,
		OrgName:  f["org"].GetStringValue(),
	}, nil
}

// CanonicalScope renders a scope map as deterministic JSON (sorted keys) so
// it can participate in unique indexes.
func CanonicalScope(scope map[string]string) string {
	if len(scope) == 0 {
		return "{}"
	}

	keys := make([]string, 0, len(scope))
	for k := range scope {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	ordered := make(map[string]string, len(scope))
	for _, k := range keys {
		ordered[k] = scope[k]
	}

	b, _ := json.Marshal(ordered) // map marshal sorts keys in Go

	return string(b)
}

// ParseScope decodes a canonical scope JSON string back into a map.
func ParseScope(scope string) map[string]string {
	out := map[string]string{}
	_ = json.Unmarshal([]byte(scope), &out)

	return out
}
