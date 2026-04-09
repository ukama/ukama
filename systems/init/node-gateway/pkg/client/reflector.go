/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/init/reflector/pb/gen"
)
 
 type ReflectorEP interface {
	Ping(req *pb.PingRequest) (*pb.PingResponse, error)
	Get(req *pb.GetRequest) (*pb.GetResponse, error)
	Download(req *pb.DownloadRequest) (*pb.DownloadResponse, error)
	Upload(req *pb.UploadRequest) (*pb.UploadResponse, error)
}
 
 type Reflector struct {
	conn    *grpc.ClientConn
	client  pb.ReflectorServiceClient
	timeout time.Duration
	host    string
}
 
func NewReflector(reflectorHost string, timeout time.Duration) *Reflector {
 
	conn, err := grpc.NewClient(reflectorHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Reflector Service:  %v", err)
	}
	client := pb.NewReflectorServiceClient(conn)
 
	return &Reflector{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    reflectorHost,
	}
}
 
func NewReflectorFromClient(mClient pb.ReflectorServiceClient) *Reflector {
	return &Reflector{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
 	}
}
 
func (r *Reflector) Close() {
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close Reflector Service connection: %v", err)
		}
	}
}

func (r *Reflector) Ping(req *pb.PingRequest) (*pb.PingResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.Ping(ctx, req)
}

func (r *Reflector) Get(req *pb.GetRequest) (*pb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.client.Get(ctx, req)
}

func (r *Reflector) Download(req *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.client.Download(ctx, req)
}

func (r *Reflector) Upload(req *pb.UploadRequest) (*pb.UploadResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.client.Upload(ctx, req)
}