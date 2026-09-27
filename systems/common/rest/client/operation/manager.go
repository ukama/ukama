/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package operation

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ukama/ukama/systems/common/rest/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const OperationsEndpoint = "/v1/operations"

const (
	StatusPending   = "PENDING"
	StatusRunning   = "RUNNING"
	StatusSuccess   = "SUCCESS"
	StatusFailed    = "FAILED"
	StatusTimeout   = "TIMEOUT"
	StatusCancelled = "CANCELLED"
)

type OperationInfo struct {
	Id             string    `json:"id"`
	Type           string    `json:"type"`
	System         string    `json:"system"`
	Status         string    `json:"status"`
	FencingToken   uint64    `json:"fencing_token"`
	RequestedBy    string    `json:"requested_by,omitempty"`
	IdempotencyKey string    `json:"idempotency_key,omitempty"`
	ResourceKey    string    `json:"resource_key"`
	LeaseExpiresAt time.Time `json:"lease_expires_at"`
	Error          string    `json:"error,omitempty"`
	StartedAt      time.Time `json:"started_at,omitempty"`
	TerminalAt     time.Time `json:"terminal_at,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type StartRequest struct {
	Type           string `json:"type"`
	System         string `json:"system"`
	ResourceKey    string `json:"resource_key"`
	RequestedBy    string `json:"requested_by,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	LeaseSeconds   uint32 `json:"lease_seconds,omitempty"`
	ConflictResourceKeys   []string `json:"conflict_resource_keys,omitempty"`
	AdditionalResourceKeys []string `json:"additional_resource_keys,omitempty"`
}

type StartResponse struct {
	Operation            *OperationInfo `json:"operation,omitempty"`
	ConflictingOperation *OperationInfo `json:"conflicting_operation,omitempty"`
	Operations           []*OperationInfo `json:"operations,omitempty"`
	ConflictsChecked     bool `json:"conflicts_checked,omitempty"`
}

type GetResponse struct {
	Operation *OperationInfo `json:"operation,omitempty"`
}

type MarkRunningRequest struct {
	FencingToken uint64 `json:"fencing_token"`
}

type ForceUnlockRequest struct {
	Actor  string `json:"actor"`
	Reason string `json:"reason"`
}

type CompleteRequest struct {
	Actor  string `json:"actor"`
	Reason string `json:"reason"`
}

type ManagerClient interface {
	Start(StartRequest) (*StartResponse, error)
	Get(id string) (*OperationInfo, error)
	GetByResource(resourceKey string) (*OperationInfo, error)
	MarkRunning(id string, fencingToken uint64) (*OperationInfo, error)
	ForceUnlock(id, actor, reason string) (*OperationInfo, error)
	Complete(id, actor, reason string) (*OperationInfo, error)
}

type managerClient struct {
	u *url.URL
	R *client.Resty
}

func NewManagerClient(host string, options ...client.Option) *managerClient {
	u, err := url.Parse(host)
	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", host, err)
	}
	return &managerClient{u: u, R: client.NewResty(options...)}
}

func (m *managerClient) Start(req StartRequest) (*StartResponse, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error: %w", err)
	}
	resp, err := m.R.Post(m.u.String()+OperationsEndpoint, b)
	if err != nil {
		var responseError *client.ErrorStatus
		if errors.As(err, &responseError) && responseError.StatusCode == http.StatusConflict {
			return nil, status.Error(codes.AlreadyExists, "site resource is locked by an active operation")
		}
		return nil, fmt.Errorf("StartOperation failure: %w", err)
	}
	out := &StartResponse{}
	if uerr := json.Unmarshal(resp.Body(), out); uerr != nil {
		return nil, fmt.Errorf("StartOperation deserialize: %w", uerr)
	}
	if len(req.ConflictResourceKeys) > 0 && !out.ConflictsChecked {
		return nil, fmt.Errorf("operation manager did not confirm site conflict checks; upgrade operation manager and api-gateway")
	}
	return out, nil
}

func (m *managerClient) Get(id string) (*OperationInfo, error) {
	resp, err := m.R.Get(m.u.String() + OperationsEndpoint + "/" + id)
	if err != nil {
		return nil, fmt.Errorf("GetOperation failure: %w", err)
	}
	out := &GetResponse{}
	if uerr := json.Unmarshal(resp.Body(), out); uerr != nil {
		return nil, fmt.Errorf("GetOperation deserialize: %w", uerr)
	}
	return out.Operation, nil
}

func (m *managerClient) GetByResource(resourceKey string) (*OperationInfo, error) {
	resp, err := m.R.Get(m.u.String() + OperationsEndpoint + "?resource_key=" + url.QueryEscape(resourceKey))
	if err != nil {
		return nil, fmt.Errorf("GetByResource failure: %w", err)
	}
	out := &GetResponse{}
	if uerr := json.Unmarshal(resp.Body(), out); uerr != nil {
		return nil, fmt.Errorf("GetByResource deserialize: %w", uerr)
	}
	return out.Operation, nil
}

func (m *managerClient) MarkRunning(id string, fencingToken uint64) (*OperationInfo, error) {
	b, err := json.Marshal(MarkRunningRequest{FencingToken: fencingToken})
	if err != nil {
		return nil, fmt.Errorf("request marshal error: %w", err)
	}
	resp, err := m.R.Post(m.u.String()+OperationsEndpoint+"/"+id+"/run", b)
	if err != nil {
		return nil, fmt.Errorf("MarkRunning failure: %w", err)
	}
	out := &GetResponse{}
	if uerr := json.Unmarshal(resp.Body(), out); uerr != nil {
		return nil, fmt.Errorf("MarkRunning deserialize: %w", uerr)
	}
	return out.Operation, nil
}

func (m *managerClient) ForceUnlock(id, actor, reason string) (*OperationInfo, error) {
	b, err := json.Marshal(ForceUnlockRequest{Actor: actor, Reason: reason})
	if err != nil {
		return nil, fmt.Errorf("request marshal error: %w", err)
	}
	resp, err := m.R.Post(m.u.String()+OperationsEndpoint+"/"+id+"/force-unlock", b)
	if err != nil {
		return nil, fmt.Errorf("ForceUnlock failure: %w", err)
	}
	out := &GetResponse{}
	if uerr := json.Unmarshal(resp.Body(), out); uerr != nil {
		return nil, fmt.Errorf("ForceUnlock deserialize: %w", uerr)
	}
	return out.Operation, nil
}

func (m *managerClient) Complete(id, actor, reason string) (*OperationInfo, error) {
	b, err := json.Marshal(CompleteRequest{Actor: actor, Reason: reason})
	if err != nil {
		return nil, fmt.Errorf("request marshal error: %w", err)
	}
	resp, err := m.R.Post(m.u.String()+OperationsEndpoint+"/"+id+"/complete", b)
	if err != nil {
		return nil, fmt.Errorf("CompleteOperation failure: %w", err)
	}
	out := &GetResponse{}
	if uerr := json.Unmarshal(resp.Body(), out); uerr != nil {
		return nil, fmt.Errorf("CompleteOperation deserialize: %w", uerr)
	}
	return out.Operation, nil
}
