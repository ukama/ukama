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
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

func (sm *SimManager) AllocateSim(req *pb.AllocateSimRequest) (*pb.AllocateSimResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.AllocateSim(ctx, req)
}

func (sm *SimManager) GetSim(simId string) (*pb.GetSimResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.GetSim(ctx, &pb.GetSimRequest{SimId: simId})
}

func (sm *SimManager) ListSims(iccid, imsi, subscriberId, networkId, simType, simStatus string, trafficPolicy uint32,
	isPhysical, sort bool, count uint32) (*pb.ListSimsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
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

// Deprecated: Use pkg.client.SimManager.ListSims with subscriberId as filtering param instead.
func (sm *SimManager) GetSimsBySub(subscriberId string) (*pb.GetSimsBySubscriberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.GetSimsBySubscriber(ctx, &pb.GetSimsBySubscriberRequest{SubscriberId: subscriberId})
}

// Deprecated: Use pkg.client.SimManager.ListSims with networkId as filtering param instead.
func (sm *SimManager) GetSimsByNetwork(networkId string) (*pb.GetSimsByNetworkResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.GetSimsByNetwork(ctx, &pb.GetSimsByNetworkRequest{NetworkId: networkId})
}

func (sm *SimManager) ToggleSimStatus(simId string, status string) (*pb.ToggleSimStatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.ToggleSimStatus(ctx, &pb.ToggleSimStatusRequest{SimId: simId, Status: status})
}

func (sm *SimManager) AddPackageToSim(req *pb.AddPackageRequest) (*pb.AddPackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.AddPackageForSim(ctx, req)
}

func (sm *SimManager) RemovePackageForSim(req *pb.RemovePackageRequest) (*pb.RemovePackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.RemovePackageForSim(ctx, req)
}

func (sm *SimManager) TerminateSim(simId string) (*pb.TerminateSimResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.TerminateSim(ctx, &pb.TerminateSimRequest{SimId: simId})
}

func (sm *SimManager) ListPackagesForSim(simId, dataPlanId, fromStartDate, toStartDate, fromEndDate,
	toEndDate string, isActive, asExpired, sort bool, count uint32) (*pb.ListPackagesForSimResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.ListPackagesForSim(ctx, &pb.ListPackagesForSimRequest{
		SimId:         simId,
		DataPlanId:    dataPlanId,
		FromStartDate: fromStartDate,
		ToStartDate:   toStartDate,
		FromEndDate:   fromEndDate,
		ToEndDate:     toEndDate,
		IsActive:      isActive,
		AsExpired:     asExpired,
		Sort:          sort,
		Count:         count,
	})
}

// Deprecated: Use pkg.client.SimManager.ListPackagesForSim with simId as filtering param instead.
func (sm *SimManager) GetPackagesForSim(simId string) (*pb.GetPackagesForSimResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.GetPackagesForSim(ctx, &pb.GetPackagesForSimRequest{SimId: simId})
}

func (sm *SimManager) SetActivePackageForSim(req *pb.SetActivePackageRequest) (*pb.SetActivePackageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	return sm.client.SetActivePackageForSim(ctx, req)
}

func (sm *SimManager) GetUsages(simId, simType, cdrType, from, to, region string) (*pb.UsageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
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
