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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/errors"
)

type PrometheusPayload struct {
	Metrics []string `json:"metrics"`
	Start string `json:"start"`
	End string `json:"end"`
	Step string `json:"step"`
	Filters []Filter `json:"filters"`
	Operation string `json:"operation"`
}

type Filter struct {
	Key string `json:"key"`
	Value string `json:"value"`
}

type PrometheusRequestData struct {
	Url     string `json:"url"`
	Query   url.Values `json:"query"`
	Payload PrometheusPayload `json:"payload"`
}

func getFiltersQuery(filters []Filter) string {
	var fq []string
	for _, f := range filters {
		fq = append(fq, fmt.Sprintf("%s='%s'", f.Key, f.Value))
	}
	return strings.ReplaceAll(strings.Join(fq, ","), "'", "\"")
}

func GetPrometheusRequestUrl(
	baseUrl string,
	pp PrometheusPayload,
) PrometheusRequestData {

	queries := make([]string, 0)

	u := fmt.Sprintf("%s/api/v1/query_range", strings.TrimSuffix(baseUrl, "/"))

	data := url.Values{}
	data.Set("start", pp.Start)
	data.Set("end", pp.End)
	data.Set("step", pp.Step)

	filtersQuery := getFiltersQuery(pp.Filters)

	for _, metric := range pp.Metrics {
		queries = append(
			queries,
			fmt.Sprintf(`%s{%s}`, metric, filtersQuery),
		)
	}
	query := fmt.Sprintf("(%s)", strings.Join(queries, " or "))

	if pp.Operation != "" {
		data.Set("query", fmt.Sprintf("%s(%s)", pp.Operation, query))
	} else {
		data.Set("query", query)
	}

	return PrometheusRequestData{
		Url: u,
		Query: data,
		Payload: pp,
	}
}

func formatMetricsResponse(w io.Writer, b io.ReadCloser, prd PrometheusRequestData) error {
	for _, metric := range prd.Payload.Metrics {
		bytes, err := io.ReadAll(b)
		if err != nil {
			log.Errorf("Failed to read prometheus response for %s Error: %v", metric, err)
			return err
		}

		rmap := map[string]interface{}{}
		err = json.Unmarshal([]byte(bytes), &rmap)
		if err != nil {
			log.Errorf("Failed to unmarshal prometheus response for %s Error: %v", metric, err)
			return err
		}
		rmap["Name"] = metric

		rb, err := json.Marshal(rmap)
		if err != nil {
			log.Errorf("Failed to marshal prometheus response for %s Error: %v", metric, err)
			return err
		}

		n, err := w.Write(rb)
		if err != nil {
			log.Errorf("Failed to add prometheus response to ws response for %s Error: %v", metric, err)
			return err
		}

		log.Infof("Updated %d bytes of response: %s", n, string(rb))
		
	}

	
	return nil
}

func ProcessPromRequest(ctx context.Context, prd PrometheusRequestData, w io.Writer, formatting bool) (httpStatusCode int, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, prd.Url, strings.NewReader(prd.Query.Encode()))
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "failed to create request")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(prd.Query.Encode())))

	log.Infof("Request is: %v Body %+v", req, prd.Query.Encode())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "failed to execute request")
	}
	log.Infof("Response Body %+v", res.Body)
	if formatting {
		err = formatMetricsResponse(w, res.Body, prd)
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