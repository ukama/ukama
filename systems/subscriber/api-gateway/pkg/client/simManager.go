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
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	cclient "github.com/ukama/ukama/systems/common/rest/client"
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
)

type SimManager struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.SimManagerServiceClient
	host    string
}

func NewSimManager(host string, timeout time.Duration) *SimManager {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()),
		ugrpc.TracingDialOption())
	if err != nil {
		log.Fatalf("Failed to connect to Sim Manager Service: %v", err)
	}
	client := pb.NewSimManagerServiceClient(conn)

	return &SimManager{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewSimManagerFromClient(SimManagerClient pb.SimManagerServiceClient) *SimManager {
	return &SimManager{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  SimManagerClient,
	}
}

func (sm *SimManager) Close() {
	if sm.conn != nil {
		if err := sm.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close Sim Manager Service connection: %v", err)
		}
	}
}

func (sm *SimManager) AllocateSim(ctx context.Context, req *pb.AllocateSimRequest) (*pb.SimResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sm.timeout)
	defer cancel()

	return sm.client.AllocateSim(ctx, req)
}

func (sm *SimManager) GetSim(ctx context.Context, simId string) (*pb.SimResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sm.timeout)
	defer cancel()

	return sm.client.GetSim(ctx, &pb.SimRequest{SimId: simId})
}

func (sm *SimManager) ListSims(ctx context.Context, iccid, imsi, subscriberId, networkId, simType, simStatus string, trafficPolicy uint32,
	isPhysical, sort bool, count uint32) (*pb.ListSimsResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sm.timeout)
	defer cancel()

	return sm.client.ListSims(ctx, &pb.ListSimsRequest{
		Iccid:         iccid,
		Imsi:          imsi,
		SubscriberId:  subscriberId,
		NetworkId:     networkId,
		SimType:       simType,
		SimStatus:     simStatus,
		TrafficPolicy: trafficPolicy,
		IsPhysical:    isPhysical,
		Sort:          sort,
		Count:         count,
	})
}

func (sm *SimManager) ToggleSimServiceStatus(ctx context.Context, simId string, status string) (*pb.ToggleSimServiceStatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sm.timeout)
	defer cancel()

	return sm.client.ToggleSimServiceStatus(ctx, &pb.ToggleSimServiceStatusRequest{SimId: simId, Status: status})
}

func (sm *SimManager) AddPackageToSim(ctx context.Context, req *pb.AddPackageRequest) (*pb.PackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sm.timeout)
	defer cancel()

	return sm.client.AddPackageForSim(ctx, req)
}

func (sm *SimManager) ListPackagesForSim(ctx context.Context, simId, dataPlanId, fromStartDate, toStartDate, fromEndDate,
	toEndDate string, isCurrentlyInUse, isExpired, sort bool, count uint32) (*pb.ListPackagesForSimResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sm.timeout)
	defer cancel()

	return sm.client.ListPackagesForSim(ctx, &pb.ListPackagesForSimRequest{
		SimId:            simId,
		DataPlanId:       dataPlanId,
		FromStartDate:    fromStartDate,
		ToStartDate:      toStartDate,
		FromEndDate:      fromEndDate,
		ToEndDate:        toEndDate,
		IsCurrentlyInUse: isCurrentlyInUse,
		IsExpired:        isExpired,
		Sort:             sort,
		Count:            count,
	})
}

func (sm *SimManager) RemovePackageForSim(ctx context.Context, req *pb.PackageRequest) (*pb.PackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sm.timeout)
	defer cancel()

	return sm.client.RemovePackageForSim(ctx, req)
}

func (sm *SimManager) SetPackageInUseForSim(ctx context.Context, req *pb.PackageRequest) (*pb.PackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sm.timeout)
	defer cancel()

	return sm.client.SetPackageInUseForSim(ctx, req)
}

func (sm *SimManager) UnsetPackageInUseForSim(ctx context.Context, req *pb.PackageRequest) (*pb.PackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sm.timeout)
	defer cancel()

	return sm.client.UnsetPackageInUseForSim(ctx, req)
}

func (sm *SimManager) TerminateSim(ctx context.Context, simId string) (*pb.TerminateSimResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sm.timeout)
	defer cancel()

	return sm.client.TerminateSim(ctx, &pb.SimRequest{SimId: simId})
}

func (sm *SimManager) GetUsages(ctx context.Context, simId, simType, cdrType, from, to, region string) (*pb.UsageResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sm.timeout)
	defer cancel()

	resp, err := sm.client.GetUsages(ctx,
		&pb.UsageRequest{
			SimId:   simId,
			SimType: simType,
			Type:    cdrType,
			From:    from,
			To:      to,
			Region:  region,
		})

	if err != nil {
		return nil, cclient.HandleRestErrorStatus(err)
	}

	return resp, nil
}

func (sm *SimManager) GetSimToken(ctx context.Context, iccid string) (*pb.SimTokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sm.timeout)
	defer cancel()

	return sm.client.GenerateSimToken(ctx,
		&pb.SimTokenRequest{
			Iccid: iccid,
		})
}

// Deprecated: Use pkg.client.SimManager.ListPackagesForSim with simId as filtering param instead.
func (sm *SimManager) GetPackagesForSim(ctx context.Context, simId string) (*pb.GetPackagesForSimResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sm.timeout)
	defer cancel()

	return sm.client.GetPackagesForSim(ctx, &pb.GetPackagesForSimRequest{SimId: simId})
}

// Deprecated: Use pkg.client.SimManager.ListSims with subscriberId as filtering param instead.
func (sm *SimManager) GetSimsBySub(ctx context.Context, subscriberId string) (*pb.GetSimsBySubscriberResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sm.timeout)
	defer cancel()

	return sm.client.GetSimsBySubscriber(ctx, &pb.GetSimsBySubscriberRequest{SubscriberId: subscriberId})
}

// Deprecated: Use pkg.client.SimManager.ListSims with networkId as filtering param instead.
func (sm *SimManager) GetSimsByNetwork(ctx context.Context, networkId string) (*pb.GetSimsByNetworkResponse, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sm.timeout)
	defer cancel()

	return sm.client.GetSimsByNetwork(ctx, &pb.GetSimsByNetworkRequest{NetworkId: networkId})
}
