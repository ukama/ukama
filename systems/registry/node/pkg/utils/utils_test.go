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

	"github.com/tj/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/node/pkg/db"
	"github.com/ukama/ukama/systems/registry/node/pkg/utils"

	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
)

func TestIsValidNodeConnectivity(t *testing.T) {
	for _, v := range []string{"Online", "offline", "UNDEFINED"} {
		assert.True(t, utils.IsValidNodeConnectivity(v), v)
	}

	for _, v := range []string{"", "bogus", "1", "high"} {
		assert.False(t, utils.IsValidNodeConnectivity(v), v)
	}
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
