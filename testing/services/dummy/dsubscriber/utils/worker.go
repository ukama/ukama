/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package utils

import (
	"fmt"
	"math"
	"math/rand/v2"
	"time"

	agent "github.com/ukama/ukama/systems/common/rest/client/ukamaagent"
	cenums "github.com/ukama/ukama/testing/common/enums"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/clients"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/pkg"
)

func Worker(iccid string, updateChan chan pkg.WMessage, initial pkg.WMessage, rc pkg.RoutineConfig) {
	interval := rc.Interval
	count := 1
	nodeId := initial.NodeId
	cdrClient := initial.CDRClient
	agent := initial.Agent
	profile := initial.Profile
	expiry := initial.Expiry
	status := initial.Status
	scenario := initial.Scenario
	imsi := initial.Imsi

	fmt.Printf("Coroutine %s started with: %d, %s\n", iccid, profile, expiry)

	runLogic(iccid, nodeId, imsi, profile, scenario, cdrClient, count, interval, rc, agent)

	ticker := time.NewTicker(time.Duration(rc.Interval) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-updateChan:
			if !ok {
				fmt.Printf("Coroutine %s stopping for status inactive\n", iccid)
				return
			}
			if msg.NodeId != "" {
				nodeId = msg.NodeId
			}
			profile = msg.Profile
			scenario = msg.Scenario
			expiry = msg.Expiry
			status = msg.Status
			imsi = msg.Imsi
			fmt.Printf("Coroutine %s updated args: %d, %s\n", iccid, profile, expiry)
		case <-ticker.C:
			if status {
				runLogic(iccid, nodeId, imsi, profile, scenario, cdrClient, count, interval, rc, agent)
				count += 1
				interval += rc.Interval
			}
		}
		if status {
			expiryDate, _ := time.Parse(time.RFC3339, expiry)
			diff := time.Until(expiryDate)
			totalMints := uint64(math.Round(diff.Minutes()))
			if interval > totalMints {
				fmt.Printf("Coroutine %s stopping for limit reach\n", iccid)
				return
			}
		}
	}
}

func runLogic(iccid, nodeId, imsi string, profile cenums.Profile, scenario cenums.SCENARIOS, cdrClient clients.CDRClient, count int, interval uint64, rc pkg.RoutineConfig, ua agent.UkamaAgentClient) {
	usage := 0.0
	switch profile {
	case cenums.PROFILE_MIN:
		usage = rc.Min + rand.Float64()*(rc.Normal-rc.Min)*0.1
	case cenums.PROFILE_NORMAL:
		usage = rc.Normal + rand.Float64()*(rc.Max-rc.Normal)*0.1
	case cenums.PROFILE_MAX:
		usage = rc.Max + rand.Float64()*(rc.Max-rc.Normal)*0.1
	}

	res, err := ua.GetSimInfo(iccid)
	if err != nil {
		fmt.Printf("Coroutine PostCDR for IMSI: %s error: %v\n", imsi, err)
		return
	}
	start := time.Now()
	end := start.Add(time.Duration(rc.Interval*60) * time.Second)
	fmt.Printf("Coroutine PostCDR for IMSI: %s , Start: %d, End: %d\n", imsi, start.Unix(), end.Unix())
	err = cdrClient.AddCDR(clients.AddCDRRequest{
		Session:       uint64(count),
		Imsi:          imsi,
		NodeId:        nodeId,
		Policy:        res.Policy.UUID,
		ApnName:       "apn",
		Ip:            "ip",
		StartTime:     uint64(start.Unix()),
		EndTime:       uint64(end.Unix()),
		LastUpdatedAt: uint64(start.Unix()),
		TxBytes:       uint64(usage * 1024 * 1024),
		RxBytes:       uint64(usage * 1024 * 1024),
		TotalBytes:    uint64(usage * 1024 * 1024),
	})

	if err != nil {
		fmt.Printf("Coroutine PostCDR for IMSI: %s error: %v\n", imsi, err)
	}

}
