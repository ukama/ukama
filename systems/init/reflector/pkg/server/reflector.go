/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/init/reflector/pb/gen"
	"github.com/ukama/ukama/systems/init/reflector/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReflectorServer struct {
	pb.UnimplementedReflectorServiceServer
	config *pkg.Config
}

func NewReflectorServer(config *pkg.Config) *ReflectorServer {
	return &ReflectorServer{
		config: config,
	}
}

func (s *ReflectorServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	log.Infof("Ping request received")

	return &pb.PingResponse{
		Message: fmt.Sprintf("OK ts=%d", time.Now().UnixMilli()),
	}, nil
}

func (s *ReflectorServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Infof("Get request received")

	return &pb.GetResponse{
		ReflectorNearUrl: s.config.BaseURL(),
		ReflectorFarUrl:  s.config.BaseURL(),
		Version:          "ukama-reflector-1",
	}, nil
}

func (s *ReflectorServer) Download(ctx context.Context, req *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	log.Infof("Download request received")

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	if req.Bytes <= 0 {
		return nil, status.Error(codes.InvalidArgument, "bytes must be > 0")
	}

	if req.Bytes > s.config.ServiceConfig.MaxDownloadBytes {
		return nil, status.Errorf(codes.ResourceExhausted, "bytes exceeds max download bytes (%d)", s.config.ServiceConfig.MaxDownloadBytes)
	}

	payload := make([]byte, req.Bytes)

	if req.ChunkDelayMs > 0 && req.ChunkBytes > 0 {
		chunks := req.Bytes / int64(req.ChunkBytes)
		if req.Bytes%int64(req.ChunkBytes) != 0 {
			chunks++
		}
		if chunks > 0 {
			time.Sleep(time.Duration(chunks*int64(req.ChunkDelayMs)) * time.Millisecond)
		}
	}

	return &pb.DownloadResponse{
		Payload: payload,
	}, nil
}

func (s *ReflectorServer) Upload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	log.Infof("Upload request received")

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	payload := req.GetPayload()
	if int64(len(payload)) > s.config.ServiceConfig.MaxUploadBytes {
		return nil, status.Errorf(codes.ResourceExhausted, "payload exceeds max upload bytes (%d)", s.config.ServiceConfig.MaxUploadBytes)
	}

	sum := sha256.Sum256(payload)
	nowTs := time.Now().Unix()

	return &pb.UploadResponse{
		Ok:            true,
		BytesReceived: int64(len(payload)),
		Sha256:        hex.EncodeToString(sum[:]),
		Ts:            nowTs,
	}, nil
}