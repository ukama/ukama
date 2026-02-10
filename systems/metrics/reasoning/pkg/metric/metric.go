/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package metric

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/errors"
)


func formatMetricsResponse(metricName string, w io.Writer, b io.ReadCloser) error {
	bytes, err := io.ReadAll(b)
	if err != nil {
		log.Errorf("Failed to read prometheus response for %s Error: %v", metricName, err)
		return err
	}

	rmap := map[string]interface{}{}
	err = json.Unmarshal([]byte(bytes), &rmap)
	if err != nil {
		log.Errorf("Failed to unmarshal prometheus response for %s Error: %v", metricName, err)
		return err
	}
	rmap["Name"] = metricName

	rb, err := json.Marshal(rmap)
	if err != nil {
		log.Errorf("Failed to marshal prometheus response for %s Error: %v", metricName, err)
		return err
	}

	n, err := w.Write(rb)
	if err != nil {
		log.Errorf("Failed to add prometheus response to ws response for %s Error: %v", metricName, err)
		return err
	}

	log.Infof("Updated %d bytes of response: %s", n, string(rb))
	return nil
}

func ProcessPromRequest(ctx context.Context, metricName string, url string, data url.Values, w io.Writer, formatting bool) (httpStatusCode int, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "failed to create request")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	log.Infof("Request is: %v Body %+v", req, data.Encode())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "failed to execute request")
	}
	log.Infof("Response Body %+v", res.Body)
	if formatting {
		err = formatMetricsResponse(metricName, w, res.Body)
		if err != nil {
			return http.StatusInternalServerError, errors.Wrap(err, "failed to format response")
		}
	} else {
		_, err = io.Copy(w, res.Body)
		if err != nil {
			return http.StatusInternalServerError, errors.Wrap(err, "failed to copy response")
		}
	}

	err = res.Body.Close()
	if err != nil {
		log.Warnf("failed to properly close response body. Error: %v", err)
	}

	return res.StatusCode, nil
}