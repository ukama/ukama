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

	cenums "github.com/ukama/ukama/testing/common/enums"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/clients"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/pkg"
)

const (
	MIN    = 10
	NORMAL = 20
	MAX    = 40
	STEP   = 1
)

func Worker(iccid string, updateChan chan pkg.WMessage, initial pkg.WMessage) {
	mint := uint64(STEP)
	count := 1
	nodeId := initial.NodeId
	cdrClient := initial.CDRClient
	profile := initial.Profile
	expiry := initial.Expiry
	pkgId := initial.PackageId

	fmt.Printf("Coroutine %s started with: %d, %s, %s\n", iccid, profile, expiry, pkgId)

	runLogic(iccid, nodeId, profile, cdrClient, count, mint)

	ticker := time.NewTicker(STEP * time.Minute)
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
			expiry = msg.Expiry
			fmt.Printf("Coroutine %s updated args: %d, %s\n", iccid, profile, expiry)
		case <-ticker.C:
			runLogic(iccid, nodeId, profile, cdrClient, count, mint)
			count += 1
			mint += STEP
		}
		expiryDate, _ := time.Parse(time.RFC3339, expiry)
		diff := time.Until(expiryDate)
		totalMints := uint64(math.Round(diff.Minutes()))
		if mint > totalMints {
			fmt.Printf("Coroutine %s stopping for limit reach\n", iccid)
			return
		}
	}
}

func runLogic(iccid, nodeId string, profile cenums.Profile, cdrClient clients.CDRClient, count int, mint uint64) {
	usage := 0.0
	if profile == cenums.PROFILE_MIN {
		usage = MIN + rand.Float64()*(NORMAL-MIN)*0.1
	} else if profile == cenums.PROFILE_NORMAL {
		usage = NORMAL + rand.Float64()*(MAX-NORMAL)*0.1
	} else if profile == cenums.PROFILE_MAX {
		usage = MAX + rand.Float64()*(MAX-NORMAL)*0.1
	}

	if len(iccid) > 4 {
		iccidInImsi := iccid[4:] //TODO: TEMP logic
		start := time.Now()
		end := start.Add(time.Duration(STEP*60) * time.Second)
		fmt.Printf("Coroutine PostCDR for %s: %d, %d\n", iccid, start.Unix(), end.Unix())
		err := cdrClient.AddCDR(clients.AddCDRRequest{
			Session:       uint64(count),
			Imsi:          iccidInImsi,
			NodeId:        nodeId,
			Policy:        "policy",
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
			fmt.Printf("Coroutine PostCDR for %s error: %v\n", iccid, err)
		}
	} else {
		fmt.Println("String is too short to remove 4 characters")
	}
}
