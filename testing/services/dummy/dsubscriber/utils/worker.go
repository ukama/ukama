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
)

func Worker(iccid string, updateChan chan pkg.WMessage, initial pkg.WMessage) {
	mint := 1
	count := 1
	nodeId := initial.NodeId
	cdrClient := initial.CDRClient
	profile := initial.Profile
	expiry := initial.Expiry
	pkgId := initial.PackageId
	expiryDate, _ := time.Parse(time.RFC3339, expiry)
	diff := time.Until(expiryDate)
	totalMints := int(math.Round(diff.Minutes()))

	fmt.Printf("Coroutine %s started with: %d, %s, %s\n", iccid, profile, expiry, pkgId)

	for {
		usage := 0.0
		if mint > totalMints {
			fmt.Printf("Coroutine %s stopping\n", iccid)
			return
		}
		select {
		case msg, ok := <-updateChan:
			if !ok {
				fmt.Printf("Coroutine %s stopping\n", iccid)
				return
			}
			nodeId = msg.NodeId
			profile = msg.Profile
			expiry = msg.Expiry
			expiryDate, _ := time.Parse(time.RFC3339, expiry)
			diff := time.Until(expiryDate)
			totalMints = int(math.Round(diff.Minutes()))
			fmt.Printf("Coroutine %s updated args: %d, %s\n", iccid, profile, expiry)
		default:
		}

		if profile == cenums.PROFILE_MIN {
			usage = MIN + rand.Float64()*(NORMAL-MIN)*0.1
		} else if profile == cenums.PROFILE_NORMAL {
			usage = NORMAL + rand.Float64()*(MAX-NORMAL)*0.1
		} else if profile == cenums.PROFILE_MAX {
			usage = MAX + rand.Float64()*(MAX-NORMAL)*0.1
		}

		if len(iccid) > 4 {
			iccidInImsi := iccid[4:] //TODO: TEMP logic

			err := cdrClient.AddCDR(clients.AddCDRRequest{
				Session:       uint64(count),
				Imsi:          iccidInImsi,
				NodeId:        nodeId,
				Policy:        "policy",
				ApnName:       "apn",
				Ip:            "ip",
				StartTime:     uint64(time.Now().Unix()),
				EndTime:       uint64(time.Now().Unix() + 300),
				LastUpdatedAt: uint64(time.Now().Unix()),
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

		count += 1
		mint += 10
		time.Sleep(10 * time.Minute)
	}
}
