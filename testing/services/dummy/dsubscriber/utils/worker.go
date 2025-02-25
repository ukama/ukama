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
	"time"

	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/pkg"
)

func Worker(id string, updateChan chan pkg.WMessage, initial pkg.WMessage) {
	count := 1.0
	profile := initial.Profile
	scenario := initial.Scenario

	fmt.Printf("Coroutine %s started with: %d, %s\n", id, profile, scenario)

	for {
		count += 0.1
		time.Sleep(1 * time.Second)
		select {
		case msg := <-updateChan:
			profile = msg.Profile
			scenario = msg.Scenario
			fmt.Printf("Coroutine %s updated args: %d, %s\n", id, profile, scenario)
		default:
		}

		fmt.Printf("Coroutine %s working with: %d, %s\n", id, profile, scenario)
	}
}
