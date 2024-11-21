/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"context"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/node/state/pb/gen"
	"google.golang.org/grpc"
)

type State struct {
	conn    *grpc.ClientConn
	client  pb.StateServiceClient
	timeout time.Duration
	host    string
}

func NewState(stateHost string, timeout time.Duration) *State {

	conn, err := grpc.NewClient(stateHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewStateServiceClient(conn)

	return &State{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    stateHost,
	}
}

func NewStateFromClient(mClient pb.StateServiceClient) *State {
	return &State{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (r *State) Close() {
	r.conn.Close()
}

func (r *State) GetStates(nodeId string) (*pb.GetStatesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.GetStates(ctx, &pb.GetStatesRequest{NodeId: nodeId})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *State) GetStatesHistory(nodeId string, pageSize int32, pageNumber int32, startTime, endTime string) (*pb.GetStatesHistoryResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.GetStatesHistory(ctx, &pb.GetStatesHistoryRequest{
		NodeId:     nodeId,
		PageSize:   pageSize,
		PageNumber: pageNumber,
		StartTime:  startTime,
		EndTime:    endTime,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *State) EnforeTransition(nodeId string, event string) (*pb.EnforceStateTransitionResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.client.EnforceStateTransition(ctx, &pb.EnforceStateTransitionRequest{
		NodeId: nodeId,
		Event:  event,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
