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

	"github.com/ukama/ukama/utils/msgcli/util"
)

const (
	defaultDuration = 5 * time.Second
)

func Run(org, scope, route, msg string, out io.Writer, cfg *util.PushConfig) error {
	routingKey, payload, err := prepareEvent(org, scope, route, msg)
	if err != nil {
		return fmt.Errorf("failled to prepare event: %w", err)
	}

	aClient := NewAmqpClient(cfg.ClusterURL, cfg.ClusterUsr, cfg.ClusterPswd, defaultDuration)

	respData, err := aClient.PublishMessage(cfg.Vhost, cfg.Exchange, routingKey, payload)
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
