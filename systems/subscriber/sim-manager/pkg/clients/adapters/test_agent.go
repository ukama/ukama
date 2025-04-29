/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package adapters

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ukama/ukama/systems/common/rest/client"

	pb "github.com/ukama/ukama/systems/subscriber/test-agent/pb/gen"
)

type TestAgentAdapter struct {
	conn   *grpc.ClientConn
	host   string
	client pb.TestAgentServiceClient
}

func NewTestAgentAdapter(testAgentHost string, timeout time.Duration) (*TestAgentAdapter, error) {

	testAgentConn, err := grpc.NewClient(testAgentHost, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return &TestAgentAdapter{
		conn:   testAgentConn,
		host:   testAgentHost,
		client: pb.NewTestAgentServiceClient(testAgentConn)}, nil
}

func (t *TestAgentAdapter) BindSim(ctx context.Context, req client.AgentRequestData) (any, error) {
	return t.client.BindSim(ctx, &pb.BindSimRequest{Iccid: req.Iccid})
}

func (t *TestAgentAdapter) GetSim(ctx context.Context, iccid string) (any, error) {
	return t.client.GetSim(ctx, &pb.GetSimRequest{Iccid: iccid})
}

func (t *TestAgentAdapter) GetUsages(ctx context.Context, iccid, cdrType, from, to, region string) (any, any, error) {
	return nil, nil, nil
}

func (t *TestAgentAdapter) ActivateSim(ctx context.Context, req client.AgentRequestData) error {
	_, err := t.client.ActivateSim(ctx, &pb.ActivateSimRequest{Iccid: req.Iccid})

	return err
}

func (t *TestAgentAdapter) DeactivateSim(ctx context.Context, req client.AgentRequestData) error {
	_, err := t.client.DeactivateSim(ctx, &pb.DeactivateSimRequest{Iccid: req.Iccid})

	return err
}

func (t *TestAgentAdapter) UpdatePackage(ctx context.Context, req client.AgentRequestData) error {
	// think of how to use ctx with restclient
	return nil
}

func (t *TestAgentAdapter) TerminateSim(ctx context.Context, iccid string) error {
	_, err := t.client.TerminateSim(ctx, &pb.TerminateSimRequest{Iccid: iccid})

	return err
}

func (t *TestAgentAdapter) Close() {
	_ = t.conn.Close()
}
