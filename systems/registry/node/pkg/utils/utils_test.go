/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package utils_test

import (
	"errors"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/tj/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/node/pkg/db"
	"github.com/ukama/ukama/systems/registry/node/pkg/utils"

	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
)

func TestRepoErrorToGrpc(t *testing.T) {
	cases := []struct {
		name string
		err  error
		code codes.Code
	}{
		{"StatusErrorKeepsCode", status.Error(codes.FailedPrecondition, "not assigned"),
			codes.FailedPrecondition},
		{"RecordNotFound", gorm.ErrRecordNotFound, codes.NotFound},
		{"DuplicateKey", &pgconn.PgError{Code: "23505"}, codes.AlreadyExists},
		{"OtherError", errors.New("boom"), codes.Internal},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.code, status.Code(utils.RepoErrorToGrpc(tc.err, "node")))
		})
	}
}

func TestIsValidNodeConnectivity(t *testing.T) {
	for _, v := range []string{"Online", "offline", "UNDEFINED"} {
		assert.True(t, utils.IsValidNodeConnectivity(v), v)
	}

	for _, v := range []string{"", "bogus", "1", "high"} {
		assert.False(t, utils.IsValidNodeConnectivity(v), v)
	}
}

func TestProcessNodeDuplErrors(t *testing.T) {
	err := utils.ProcessNodeDuplErrors(&pgconn.PgError{Code: "23505"}, "uk-node")

	assert.Equal(t, codes.AlreadyExists, status.Code(err))
	assert.Equal(t, codes.Internal,
		status.Code(utils.ProcessNodeDuplErrors(errors.New("boom"), "uk-node")))
}

func TestInvalidNodeIDError(t *testing.T) {
	err := utils.InvalidNodeIDError("bad-id", errors.New("too short"))

	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Contains(t, status.Convert(err).Message(), "bad-id")
	assert.Contains(t, status.Convert(err).Message(), "too short")
}

func TestDbNodeToPbNode(t *testing.T) {
	towerId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE).StringLowercase()
	ampId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_AMPNODE).StringLowercase()
	siteId := uuid.NewV4()

	t.Run("WithoutSite", func(t *testing.T) {
		n := utils.DbNodeToPbNode(&db.Node{
			Id:   towerId,
			Name: "tower-1",
			Type: ukama.NODE_ID_TYPE_TOWERNODE,
			Status: db.NodeStatus{Connectivity: ukama.NodeConnectivityOnline,
				State: ukama.NodeStateReady},
		})

		assert.Equal(t, towerId, n.Id)
		assert.Equal(t, cpb.NodeConnectivity_Online, n.Status.Connectivity)
		assert.Equal(t, cpb.NodeState_Ready, n.Status.State)
		assert.Nil(t, n.Site)
		assert.Nil(t, n.Attached)
	})

	t.Run("WithSiteAndAttached", func(t *testing.T) {
		n := utils.DbNodeToPbNode(&db.Node{
			Id:       towerId,
			Site:     db.Site{NodeId: towerId, SiteId: siteId},
			Attached: []*db.Node{{Id: ampId}},
		})

		assert.Equal(t, siteId.String(), n.Site.SiteId)
		assert.Len(t, n.Attached, 1)
		assert.Equal(t, ampId, n.Attached[0].Id)
	})

	t.Run("ListKeepsOrder", func(t *testing.T) {
		nodes := utils.DbNodesToPbNodes([]db.Node{{Id: towerId}, {Id: ampId}})

		assert.Equal(t, towerId, nodes[0].Id)
		assert.Equal(t, ampId, nodes[1].Id)
	})
}
