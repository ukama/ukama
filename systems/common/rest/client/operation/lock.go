/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package operation

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ukama/ukama/systems/common/rest/client"
)

var ErrBusy = errors.New("resource busy")

type BusyError struct {
	ResourceKey string
	Holder      *OperationInfo
}

func (e *BusyError) Error() string {
	if e.Holder == nil {
		return fmt.Sprintf("%s is locked by an active operation", e.ResourceKey)
	}
	return fmt.Sprintf("%s is busy with %s operation %s", e.ResourceKey,
		e.Holder.Type, e.Holder.Id)
}

func (e *BusyError) Is(target error) bool {
	return target == ErrBusy
}

type LockRequest struct {
	Type         string
	System       string
	RequestedBy  string
	LeaseSeconds uint32
	Targets      []string
	Siblings     []string
}

func IsTerminal(status string) bool {
	switch status {
	case StatusSuccess, StatusFailed, StatusTimeout, StatusCancelled:
		return true
	}
	return false
}

func LockThenCheck(m ManagerClient, req LockRequest) ([]*OperationInfo, error) {
	ops := make([]*OperationInfo, 0, len(req.Targets))
	targets := make(map[string]bool, len(req.Targets))

	for _, key := range req.Targets {
		if targets[key] {
			continue
		}
		targets[key] = true

		resp, err := m.Start(StartRequest{
			Type:         req.Type,
			System:       req.System,
			ResourceKey:  key,
			RequestedBy:  req.RequestedBy,
			LeaseSeconds: req.LeaseSeconds,
		})
		if err != nil {
			if isConflict(err) {
				err = &BusyError{ResourceKey: key}
			}
			return nil, release(m, ops, req.RequestedBy, err)
		}
		ops = append(ops, resp.Operation)
	}

	for _, key := range req.Siblings {
		if targets[key] {
			continue
		}
		holder, err := m.GetByResource(key)
		if err != nil {
			return nil, release(m, ops, req.RequestedBy,
				fmt.Errorf("read lock for %s: %w", key, err))
		}
		if holder != nil && !IsTerminal(holder.Status) {
			return nil, release(m, ops, req.RequestedBy,
				&BusyError{ResourceKey: key, Holder: holder})
		}
	}

	return ops, nil
}

func release(m ManagerClient, ops []*OperationInfo, actor string, cause error) error {
	errs := []error{cause}
	for _, op := range ops {
		if _, err := m.Fail(op.Id, op.FencingToken, actor,
			"not started: "+cause.Error()); err != nil {
			errs = append(errs, fmt.Errorf("release %s: %w", op.ResourceKey, err))
		}
	}
	return errors.Join(errs...)
}

func isConflict(err error) bool {
	var status *client.ErrorStatus
	return errors.As(err, &status) && status.StatusCode == http.StatusConflict
}
