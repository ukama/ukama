/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pkg

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/ukama"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Device5xxServerError struct {
	error
}

type Device4xxServerError struct {
	error
}

type RequestExecutor interface {
	Execute(req *cpb.NodeFeederMessage) error
}

type requestExecutor struct {
	nodeResolver   NodeIpResolver
	devicePort     int
	timeoutSeconds int
}

func NewRequestExecutor(deviceNet NodeIpResolver, devicePort int, timeoutSeconds int) *requestExecutor {
	return &requestExecutor{nodeResolver: deviceNet, devicePort: devicePort, timeoutSeconds: timeoutSeconds}
}

func (e *requestExecutor) Execute(req *cpb.NodeFeederMessage) error {
	segs := strings.Split(req.Target, ".")
	if len(segs) != 4 {
		return fmt.Errorf("invalid target format")
	}
	nodeId, err := ukama.ValidateNodeId(segs[3])
	if err != nil {
		return errors.Wrap(err, "invalid node id format")
	}
	ip, err := e.nodeResolver.Resolve(nodeId)
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
			logrus.Warningf("node %s not found", nodeId)
			return nil
		}

		return errors.Wrap(err, "error resolving device ip for device "+segs[1])
	}
	logrus.Infof("resolved device ip: %s", ip)

	strUrl := fmt.Sprintf("http://%s/%s", ip, strings.Trim(req.Path, "/"))
	if e.devicePort != 0 {
		strUrl = fmt.Sprintf("http://%s:%d/%s", ip, e.devicePort, strings.Trim(req.Path, "/"))
	}
	u, err := url.Parse(strUrl)
	if err != nil {
		return errors.Wrap(err, "malformed url")
	}

	c := http.Client{
		Timeout: time.Duration(e.timeoutSeconds) * time.Second,
	}

	httpReq := http.Request{
		Body:   io.NopCloser(bytes.NewReader((req.GetMsg()))),
		Header: map[string][]string{"Content-Type": {"application/json"}},
		Method: req.HTTPMethod,
		URL:    u,
	}
	logrus.Infof("sending request %+v to %s ", httpReq, u.String())

	resp, err := c.Do(&httpReq)
	if err != nil {
		return errors.Wrap(err, "error sending request")
	}

	logrus.Infof("Response status: %d", resp.StatusCode)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Warning("error reading response body ", err)
	} else {
		logrus.Debugf("Response body: %s", string(b))
	}

	// request with >500 status code considered as filed
	if resp.StatusCode >= 500 {
		return Device5xxServerError{
			fmt.Errorf("server error: %d", resp.StatusCode),
		}
	} else if resp.StatusCode >= 400 {
		return Device4xxServerError{
			fmt.Errorf("server error: %d", resp.StatusCode),
		}
	}

	if err != nil {
		return errors.Wrap(err, "error sending request")
	}

	return nil
}
