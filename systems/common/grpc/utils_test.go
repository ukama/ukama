/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package grpc

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/store"

	pgconnv5 "github.com/jackc/pgx/v5/pgconn"
)

func TestStoreErrorToGrpc(t *testing.T) {
	wrap := func(err error) error {
		return fmt.Errorf("node uk-1: %w", err)
	}

	cases := []struct {
		name string
		err  error
		code codes.Code
		msg  string
	}{
		{"NotFound", wrap(store.ErrNotFound), codes.NotFound, "node uk-1: not found"},
		{"AlreadyExists", wrap(store.ErrAlreadyExists), codes.AlreadyExists,
			"node uk-1: already exists"},
		{"Invalid", wrap(store.ErrInvalid), codes.InvalidArgument, "node uk-1: invalid"},
		{"Precondition", wrap(store.ErrPrecondition), codes.FailedPrecondition,
			"node uk-1: precondition failed"},
		{"GormNotFound", wrap(gorm.ErrRecordNotFound), codes.NotFound,
			"node record not found"},
		{"DuplicateKey", wrap(&pgconn.PgError{Code: "23505"}), codes.AlreadyExists,
			"node already exists"},
		{"DuplicateKeyV5", wrap(&pgconnv5.PgError{Code: "23505"}), codes.AlreadyExists,
			"node already exists"},
		{"Canceled", wrap(context.Canceled), codes.Canceled, "node uk-1: context canceled"},
		{"DeadlineExceeded", wrap(context.DeadlineExceeded), codes.DeadlineExceeded,
			"node uk-1: context deadline exceeded"},
		{"Other", errors.New("boom"), codes.Internal, "boom"},
		{"StatusIsNotPassedThrough", status.Error(codes.FailedPrecondition, "x"),
			codes.Internal, "rpc error: code = FailedPrecondition desc = x"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := status.Convert(StoreErrorToGrpc(tc.err, "node"))

			assert.Equal(t, tc.code, s.Code())
			assert.Equal(t, tc.msg, s.Message())
		})
	}

	t.Run("Nil", func(t *testing.T) {
		assert.NoError(t, StoreErrorToGrpc(nil, "node"))
	})
}
