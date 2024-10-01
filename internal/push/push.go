/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package push

import (
	"fmt"
	"io"
	"time"

	"github.com/ukama/msgcli/util"
)

const (
	defaultURI      = "http://localhost:15672"
	defaultUsr      = "guest"
	defaultPwd      = "guest"
	defaultVhost    = "%2F"
	defaultExchange = "amq.topic"
	defaultDuration = 5 * time.Second
)

func Run(org, route, msg string, out io.Writer, cfg *util.Config) error {
	// fmt.Printf("Push called with args: %q %q %q\n", org, route, msg)

	routingKey, payload, err := prepareEvent(org, route, msg)
	if err != nil {
		return fmt.Errorf("failled to prepare event: %w", err)
	}

	aClient := NewAmqpClient(defaultURI, defaultUsr, defaultPwd, defaultDuration)

	respData, err := aClient.PublishMessage(defaultVhost, defaultExchange, routingKey, payload)
	if err != nil {
		return fmt.Errorf("failled to publish event: %w", err)
	}

	outputBuf, err := util.Serialize(respData, cfg.OutputFormat)
	if err != nil {
		return fmt.Errorf("error while serializing output data: %w", err)
	}

	_, err = fmt.Fprint(out, outputBuf)
	if err != nil {
		return fmt.Errorf("error while writting output: %w", err)
	}

	return nil
}
