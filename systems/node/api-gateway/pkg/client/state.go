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
 
 type NodeState struct {
	 conn    *grpc.ClientConn
	 client  pb.StateServiceClient
	 timeout time.Duration
	 host    string
 }
 
 func NewNodeState(nodeStateHost string, timeout time.Duration) *NodeState {
	 conn, err := grpc.NewClient(nodeStateHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	 if err != nil {
		 logrus.Fatalf("did not connect: %v", err)
	 }
	 client := pb.NewStateServiceClient(conn)
 
	 return &NodeState{
		 conn:    conn,
		 client:  client,
		 timeout: timeout,
		 host:    nodeStateHost,
	 }
 }
 
 func NewNodeStateFromClient(mClient pb.StateServiceClient) *NodeState {
	 return &NodeState{
		 host:    "localhost",
		 timeout: 1 * time.Second,
		 conn:    nil,
		 client:  mClient,
	 }
 }
 
 func (r *NodeState) Close() {
	 r.conn.Close()
 }
 
 func (r *NodeState) GetByNodeId(nodeId string) (*pb.GetByNodeIdResponse, error) {
	 ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	 defer cancel()
 
	 res, err := r.client.GetByNodeId(ctx, &pb.GetByNodeIdRequest{NodeId: nodeId})
	 if err != nil {
		 return nil, err
	 }
 
	 return res, nil
 }

 
 func (r *NodeState) ListAll() (*pb.ListAllResponse, error) {
	 ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	 defer cancel()
 
	 req := &pb.ListAllRequest{}
 
	 res, err := r.client.ListAll(ctx, req)
	 if err != nil {
		 return nil, err
	 }
 
	 return res, nil
 }
 
 func (r *NodeState) GetStateHistory(nodeId string) (*pb.GetStateHistoryResponse, error) {
	 ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	 defer cancel()
 
	 req := &pb.GetStateHistoryRequest{
		 NodeId: nodeId,
	 }
 
	 res, err := r.client.GetStateHistory(ctx, req)
	 if err != nil {
		 return nil, err
	 }
 
	 return res, nil
 }