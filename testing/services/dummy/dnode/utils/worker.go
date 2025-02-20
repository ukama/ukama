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

	"github.com/ukama/ukama/testing/services/dummy/dnode/config"
)

func Worker(id string, updateChan chan config.WMessage, initial config.WMessage) {
	currentArg1 := initial.Profile
	currentArg2 := initial.Scenario

	fmt.Printf("Coroutine %s started with: %d, %s\n", id, currentArg1, currentArg2)

	for {
		select {
		case msg := <-updateChan:
			currentArg1 = msg.Profile
			currentArg2 = msg.Scenario
			fmt.Printf("Coroutine %s updated args: %d, %s\n", id, currentArg1, currentArg2)
		default:
			fmt.Printf("Coroutine %s working with: %d, %s\n", id, currentArg1, currentArg2)
			time.Sleep(1 * time.Second)
		}
	}
}
