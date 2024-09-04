/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package clients_test

// import (
// "bytes"
// "io"
// "net/http"
// "strconv"
// "testing"
// "time"

// "github.com/tj/assert"

// "github.com/ukama/ukama/testing/services/hooks/internal/adapters/mopay"
// )

// const (
// depositId                   = "03cb753f-5e03-4c97-8e47-625115476c72"
// amountCents                 = 200
// currency                    = "USD"
// country                     = "USA"
// testMsisdn                  = "2250000000000"
// correspondent               = "VODAFONE-WALLET"
// operator                    = "VODAFONE"
// description                 = "test description"
// depositPreAuthorisationCode = "0"
// )

// var (
// amount        = strconv.FormatFloat(float64(amountCents)/100.0, 'f', 2, 64)
// testTimestamp = time.Now()
// )

// func TestPawapayClient_GetDeposit(t *testing.T) {
// t.Run("DepositFound", func(tt *testing.T) {
// mockTransport := func(req *http.Request) *http.Response {
// // Test request parameters
// assert.Equal(tt, req.URL.String(), mopay.DepositEndpoint+"/"+depositId)

// // fake deposit info
// deposit := `{"depositId": "03cb753f-5e03-4c97-8e47-625115476c72", "status": "ACCEPTED", "amount": "200"}`

// // Send mock response
// return &http.Response{
// StatusCode: 200,
// Status:     "200 OK",

// // Send response to be tested
// Body: io.NopCloser(bytes.NewBufferString(deposit)),

// // Must be set to non-nil value or it panics
// Header: make(http.Header),
// }
// }

// testPawapayClient := mopay.NewPawapayClient("", "")

// // We replace the transport mechanism by mocking the http request
// // so that the test stays a unit test e.g no server/network call.
// testPawapayClient.R.C.SetTransport(RoundTripFunc(mockTransport))

// d, err := testPawapayClient.GetDeposit(depositId)

// assert.NoError(tt, err)
// assert.Equal(tt, depositId, d.DepositId)
// })

// t.Run("DepositNotFound", func(tt *testing.T) {
// mockTransport := func(req *http.Request) *http.Response {
// assert.Equal(tt, req.URL.String(), mopay.DepositEndpoint+"/"+depositId)

// return &http.Response{
// StatusCode: 404,
// Status:     "404 NOT FOUND",
// Header:     make(http.Header),
// }
// }

// testPawapayClient := mopay.NewPawapayClient("", "")

// testPawapayClient.R.C.SetTransport(RoundTripFunc(mockTransport))

// d, err := testPawapayClient.GetDeposit(depositId)

// assert.Error(tt, err)
// assert.Nil(tt, d)
// })

// t.Run("InvalidResponsePayload", func(tt *testing.T) {
// mockTransport := func(req *http.Request) *http.Response {
// assert.Equal(tt, req.URL.String(), mopay.DepositEndpoint+"/"+depositId)

// return &http.Response{
// StatusCode: 200,
// Status:     "200 OK",
// Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
// Header:     make(http.Header),
// }
// }

// testPawapayClient := mopay.NewPawapayClient("", "")

// testPawapayClient.R.C.SetTransport(RoundTripFunc(mockTransport))

// d, err := testPawapayClient.GetDeposit(depositId)

// assert.Error(tt, err)
// assert.Nil(tt, d)
// })

// t.Run("RequestFailure", func(tt *testing.T) {
// mockTransport := func(req *http.Request) *http.Response {
// assert.Equal(tt, req.URL.String(), mopay.DepositEndpoint+"/"+depositId)

// return nil
// }

// testPawapayClient := mopay.NewPawapayClient("", "")

// testPawapayClient.R.C.SetTransport(RoundTripFunc(mockTransport))

// d, err := testPawapayClient.GetDeposit(depositId)

// assert.Error(tt, err)
// assert.Nil(tt, d)
// })
// }

// func TestPawapayClient_AddDeposit(t *testing.T) {
// t.Run("DepositAdded", func(tt *testing.T) {
// mockTransport := func(req *http.Request) *http.Response {
// // Test request parameters
// assert.Equal(tt, req.URL.String(), mopay.DepositEndpoint)

// // fake deposit info
// deposit := `{"depositId": "03cb753f-5e03-4c97-8e47-625115476c72", "status": "ACCEPTED", "amount_cents": 200}`

// // Send mock response
// return &http.Response{
// StatusCode: 200,
// Status:     "200 OK",

// // Send response to be tested
// Body: io.NopCloser(bytes.NewBufferString(deposit)),

// // Must be set to non-nil value or it panics
// Header: make(http.Header),
// }
// }

// testPawapayClient := mopay.NewPawapayClient("", "")

// // We replace the transport mechanism by mocking the http request
// // so that the test stays a unit test e.g no server/network call.
// testPawapayClient.R.C.SetTransport(RoundTripFunc(mockTransport))

// d, err := testPawapayClient.AddDeposit(
// mopay.AddDepositRequest{
// DepositId: depositId,
// Amount:    amount,
// Currency:  currency,
// })

// assert.NoError(tt, err)
// assert.Equal(tt, depositId, d.DepositId)
// })

// t.Run("InvalidResponseHeader", func(tt *testing.T) {
// mockTransport := func(req *http.Request) *http.Response {
// assert.Equal(tt, req.URL.String(), mopay.DepositEndpoint)

// return &http.Response{
// StatusCode: 500,
// Status:     "500 INTERNAL SERVER ERROR",
// Body:       io.NopCloser(bytes.NewBufferString(`INTERNAL SERVER ERROR`)),
// Header:     make(http.Header),
// }
// }

// testPawapayClient := mopay.NewPawapayClient("", "")

// testPawapayClient.R.C.SetTransport(RoundTripFunc(mockTransport))

// d, err := testPawapayClient.AddDeposit(
// mopay.AddDepositRequest{
// DepositId: depositId,
// Amount:    amount,
// Currency:  currency,
// })

// assert.Error(tt, err)
// assert.Nil(tt, d)
// })

// t.Run("InvalidResponsePayload", func(tt *testing.T) {
// mockTransport := func(req *http.Request) *http.Response {
// assert.Equal(tt, req.URL.String(), mopay.DepositEndpoint)

// return &http.Response{
// StatusCode: 201,
// Status:     "201 CREATED",
// Body:       io.NopCloser(bytes.NewBufferString(`CREATED`)),
// Header:     make(http.Header),
// }
// }

// testPawapayClient := mopay.NewPawapayClient("", "")

// testPawapayClient.R.C.SetTransport(RoundTripFunc(mockTransport))

// d, err := testPawapayClient.AddDeposit(
// mopay.AddDepositRequest{
// DepositId: depositId,
// Amount:    amount,
// Currency:  currency,
// })

// assert.Error(tt, err)
// assert.Nil(tt, d)
// })

// t.Run("RequestFailure", func(tt *testing.T) {
// mockTransport := func(req *http.Request) *http.Response {
// assert.Equal(tt, req.URL.String(), mopay.DepositEndpoint)

// return nil
// }

// testPawapayClient := mopay.NewPawapayClient("", "")

// testPawapayClient.R.C.SetTransport(RoundTripFunc(mockTransport))

// d, err := testPawapayClient.AddDeposit(
// mopay.AddDepositRequest{
// DepositId: depositId,
// Amount:    amount,
// Currency:  currency,
// })

// assert.Error(tt, err)
// assert.Nil(tt, d)
// })
// }

// func TestPawapayClient_ResendDepositCallback(t *testing.T) {
// t.Run("CallbackSent", func(tt *testing.T) {
// mockTransport := func(req *http.Request) *http.Response {
// // Test request parameters
// assert.Equal(tt, req.URL.String(), mopay.DepositEndpoint+"/resend-callback")

// // fake deposit info
// deposit := `{"depositId": "03cb753f-5e03-4c97-8e47-625115476c72", "status": "ACCEPTED", "amount": "200"}`

// // Send mock response
// return &http.Response{
// StatusCode: 201,
// Status:     "201 CREATED",

// // Send response to be tested
// Body: io.NopCloser(bytes.NewBufferString(deposit)),

// // Must be set to non-nil value or it panics
// Header: make(http.Header),
// }
// }

// testPawapayClient := mopay.NewPawapayClient("", "")

// // We replace the transport mechanism by mocking the http request
// // so that the test stays a unit test e.g no server/network call.
// testPawapayClient.R.C.SetTransport(RoundTripFunc(mockTransport))

// d, err := testPawapayClient.ResendDepositCallback(
// mopay.CallbackRequest{
// DepositId: depositId,
// })

// assert.NoError(tt, err)
// assert.Equal(tt, depositId, d.DepositId)
// })

// t.Run("InvalidResponseHeader", func(tt *testing.T) {
// mockTransport := func(req *http.Request) *http.Response {
// assert.Equal(tt, req.URL.String(), mopay.DepositEndpoint+"/resend-callback")

// return &http.Response{
// StatusCode: 500,
// Status:     "500 INTERNAL SERVER ERROR",
// Body:       io.NopCloser(bytes.NewBufferString(`INTERNAL SERVER ERROR`)),
// Header:     make(http.Header),
// }
// }

// testPawapayClient := mopay.NewPawapayClient("", "")

// testPawapayClient.R.C.SetTransport(RoundTripFunc(mockTransport))

// d, err := testPawapayClient.ResendDepositCallback(
// mopay.CallbackRequest{
// DepositId: depositId,
// })

// assert.Error(tt, err)
// assert.Nil(tt, d)
// })

// t.Run("InvalidResponsePayload", func(tt *testing.T) {
// mockTransport := func(req *http.Request) *http.Response {
// assert.Equal(tt, req.URL.String(), mopay.DepositEndpoint+"/resend-callback")

// return &http.Response{
// StatusCode: 201,
// Status:     "201 CREATED",
// Body:       io.NopCloser(bytes.NewBufferString(`CREATED`)),
// Header:     make(http.Header),
// }
// }

// testPawapayClient := mopay.NewPawapayClient("", "")

// testPawapayClient.R.C.SetTransport(RoundTripFunc(mockTransport))

// d, err := testPawapayClient.ResendDepositCallback(
// mopay.CallbackRequest{
// DepositId: depositId,
// })

// assert.Error(tt, err)
// assert.Nil(tt, d)
// })

// t.Run("RequestFailure", func(tt *testing.T) {
// mockTransport := func(req *http.Request) *http.Response {
// assert.Equal(tt, req.URL.String(), mopay.DepositEndpoint+"/resend-callback")

// return nil
// }

// testPawapayClient := mopay.NewPawapayClient("", "")

// testPawapayClient.R.C.SetTransport(RoundTripFunc(mockTransport))

// d, err := testPawapayClient.ResendDepositCallback(
// mopay.CallbackRequest{
// DepositId: depositId,
// })

// assert.Error(tt, err)
// assert.Nil(tt, d)
// })
// }

// func TestPawapayClient_PredictMno(t *testing.T) {
// t.Run("MnoFound", func(tt *testing.T) {
// mockTransport := func(req *http.Request) *http.Response {
// // Test request parameters
// assert.Equal(tt, req.URL.String(), "/v1/predict-correspondent")

// // fake deposit info
// operator := `{"country": "USA", "operator": "VODAFONE", "msisdn": "2250000000000"}`

// // Send mock response
// return &http.Response{
// StatusCode: 201,
// Status:     "201 CREATED",

// // Send response to be tested
// Body: io.NopCloser(bytes.NewBufferString(operator)),

// // Must be set to non-nil value or it panics
// Header: make(http.Header),
// }
// }

// testPawapayClient := mopay.NewPawapayClient("", "")

// // We replace the transport mechanism by mocking the http request
// // so that the test stays a unit test e.g no server/network call.
// testPawapayClient.R.C.SetTransport(RoundTripFunc(mockTransport))

// o, err := testPawapayClient.PredictMno(
// mopay.MsisdnRequest{
// Msisdn: testMsisdn,
// })

// assert.NoError(tt, err)
// assert.Equal(tt, testMsisdn, o.Msisdn)
// })

// t.Run("InvalidResponseHeader", func(tt *testing.T) {
// mockTransport := func(req *http.Request) *http.Response {
// assert.Equal(tt, req.URL.String(), "/v1/predict-correspondent")

// return &http.Response{
// StatusCode: 500,
// Status:     "500 INTERNAL SERVER ERROR",
// Body:       io.NopCloser(bytes.NewBufferString(`INTERNAL SERVER ERROR`)),
// Header:     make(http.Header),
// }
// }

// testPawapayClient := mopay.NewPawapayClient("", "")

// testPawapayClient.R.C.SetTransport(RoundTripFunc(mockTransport))

// o, err := testPawapayClient.PredictMno(
// mopay.MsisdnRequest{
// Msisdn: testMsisdn,
// })

// assert.Error(tt, err)
// assert.Nil(tt, o)
// })

// t.Run("InvalidResponsePayload", func(tt *testing.T) {
// mockTransport := func(req *http.Request) *http.Response {
// assert.Equal(tt, req.URL.String(), "/v1/predict-correspondent")

// return &http.Response{
// StatusCode: 201,
// Status:     "201 CREATED",
// Body:       io.NopCloser(bytes.NewBufferString(`CREATED`)),
// Header:     make(http.Header),
// }
// }

// testPawapayClient := mopay.NewPawapayClient("", "")

// testPawapayClient.R.C.SetTransport(RoundTripFunc(mockTransport))

// o, err := testPawapayClient.PredictMno(
// mopay.MsisdnRequest{
// Msisdn: testMsisdn,
// })

// assert.Error(tt, err)
// assert.Nil(tt, o)
// })

// t.Run("RequestFailure", func(tt *testing.T) {
// mockTransport := func(req *http.Request) *http.Response {
// assert.Equal(tt, req.URL.String(), "/v1/predict-correspondent")

// return nil
// }

// testPawapayClient := mopay.NewPawapayClient("", "")

// testPawapayClient.R.C.SetTransport(RoundTripFunc(mockTransport))

// o, err := testPawapayClient.PredictMno(
// mopay.MsisdnRequest{
// Msisdn: testMsisdn,
// })

// assert.Error(tt, err)
// assert.Nil(tt, o)
// })
// }

// func TestPawapayClient_GetMnosAvailability(t *testing.T) {
// t.Run("CountriesFound", func(tt *testing.T) {
// mockTransport := func(req *http.Request) *http.Response {
// // Test request parameters
// assert.Equal(tt, req.URL.String(), "/availability")

// // fake deposit info
// countries := `[{"country": "USA"}]`

// // Send mock response
// return &http.Response{
// StatusCode: 200,
// Status:     "200 OK",

// // Send response to be tested
// Body: io.NopCloser(bytes.NewBufferString(countries)),

// // Must be set to non-nil value or it panics
// Header: make(http.Header),
// }
// }

// testPawapayClient := mopay.NewPawapayClient("", "")

// // We replace the transport mechanism by mocking the http request
// // so that the test stays a unit test e.g no server/network call.
// testPawapayClient.R.C.SetTransport(RoundTripFunc(mockTransport))

// c, err := testPawapayClient.GetMnosAvailability()

// assert.NoError(tt, err)
// assert.Equal(tt, country, c[0].Country)
// })

// t.Run("CountriesNotFound", func(tt *testing.T) {
// mockTransport := func(req *http.Request) *http.Response {
// assert.Equal(tt, req.URL.String(), "/availability")

// return &http.Response{
// StatusCode: 404,
// Status:     "404 NOT FOUND",
// Header:     make(http.Header),
// }
// }

// testPawapayClient := mopay.NewPawapayClient("", "")

// testPawapayClient.R.C.SetTransport(RoundTripFunc(mockTransport))

// c, err := testPawapayClient.GetMnosAvailability()

// assert.Error(tt, err)
// assert.Nil(tt, c)
// })

// t.Run("InvalidResponsePayload", func(tt *testing.T) {
// mockTransport := func(req *http.Request) *http.Response {
// assert.Equal(tt, req.URL.String(), "/availability")

// return &http.Response{
// StatusCode: 200,
// Status:     "200 OK",
// Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
// Header:     make(http.Header),
// }
// }

// testPawapayClient := mopay.NewPawapayClient("", "")

// testPawapayClient.R.C.SetTransport(RoundTripFunc(mockTransport))

// c, err := testPawapayClient.GetMnosAvailability()

// assert.Error(tt, err)
// assert.Nil(tt, c)
// })

// t.Run("RequestFailure", func(tt *testing.T) {
// mockTransport := func(req *http.Request) *http.Response {
// assert.Equal(tt, req.URL.String(), "/availability")

// return nil
// }

// testPawapayClient := mopay.NewPawapayClient("", "")

// testPawapayClient.R.C.SetTransport(RoundTripFunc(mockTransport))

// c, err := testPawapayClient.GetMnosAvailability()

// assert.Error(tt, err)
// assert.Nil(tt, c)
// })
// }

// type RoundTripFunc func(req *http.Request) *http.Response

// func (r RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
// return r(req), nil
// }
