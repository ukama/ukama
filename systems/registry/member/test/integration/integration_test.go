/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/ukama/ukama/systems/common/config"
	pb "github.com/ukama/ukama/systems/registry/member/pb/gen"
	"github.com/ukama/ukama/systems/registry/member/pkg/db"

	rconf "github.com/num30/config"
	log "github.com/sirupsen/logrus"
	grpc "google.golang.org/grpc"
)

var tConfig *TestConfig

func init() {
	// load config
	tConfig = &TestConfig{}

	reader := rconf.NewConfReader("integration")

	err := reader.Read(tConfig)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	log.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	log.Infof("Config: %+v\n", tConfig)
}

type TestConfig struct {
	ServiceHost string        `default:"localhost:9090"`
	Queue       *config.Queue `default:"{}"`
	OrgId       string
	OrgName     string
}

func Test_FullFlow(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	log.Infoln("Connecting to member ", tConfig.ServiceHost)

	conn, err := grpc.DialContext(ctx, tConfig.ServiceHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		assert.NoError(t, err, "did not connect: %v", err)

		return
	}

	c := pb.NewMemberServiceClient(conn)

	member := db.Member{
		UserId: uuid.NewV4(),
		Role:   db.Users,
	}

	defer removeMember(t, c, member.UserId)

	var r interface{}

	t.Run("AddMember", func(tt *testing.T) {
		r, err = c.AddMember(ctx, &pb.AddMemberRequest{
			UserUuid: member.UserId.String(),
			Role:     pb.RoleType(db.Users),
		})

		handleResponse(tt, err, r)
	})

	t.Run("GetMember", func(tt *testing.T) {
		r, err = c.GetMember(ctx, &pb.MemberRequest{UserUuid: member.UserId.String()})
		assert.NoError(t, err)
	})

	t.Run("GetMembers", func(tt *testing.T) {
		r, err = c.GetMembers(ctx, &pb.GetMembersRequest{})
		assert.NoError(t, err)
	})

}

func removeMember(t *testing.T, c pb.MemberServiceClient, userId uuid.UUID) {
	t.Helper()

	log.Infoln("Deleting member ", userId.String())

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := c.RemoveMember(ctx, &pb.MemberRequest{UserUuid: userId.String()})
	if err != nil {
		assert.FailNowf(t, "Delete member %s failed: %v\n", userId.String(), err)
	}
}

func handleResponse(t *testing.T, err error, r interface{}) {
	t.Helper()

	log.Printf("Response: %v\n", r)

	if err != nil {
		assert.FailNow(t, "Request failed: %v\n", err)
	}
}
