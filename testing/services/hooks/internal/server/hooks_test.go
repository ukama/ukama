/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server_test

// import (
// "context"
// "strconv"
// "testing"
// "time"

// "github.com/stretchr/testify/mock"
// "github.com/tj/assert"

// "github.com/ukama/ukama/systems/common/uuid"
// "github.com/ukama/ukama/testing/services/hooks/internal/server"
// "github.com/ukama/ukama/testing/services/hooks/mocks"

// mbmocks "github.com/ukama/ukama/systems/common/mocks"
// pb "github.com/ukama/ukama/testing/services/hooks/pb/gen"
// )

// const (
// OrgName               = "testorg"
// pawapayAcceptedStatus = "accepted"

// itemTypePackge      = "package"
// itemTypeInvoice     = "invoice"
// paymentMethodMopay  = "mopay"
// paymentMethodStripe = "stripe"
// statusPending       = "pending"
// strAmount           = "2"
// currency            = "USD"
// description         = "Some description"
// payerName           = "John Doe"
// payerEmail          = "johndoe@example.com"
// payerPhone          = "260763456789"
// correspondent       = "MTN_MOMO_ZMB"
// country             = "ZMB"
// token               = "fake token"
// )

// var (
// paymentId = uuid.NewV4()
// itemId    = uuid.NewV4()
// amount, _ = strconv.ParseInt(strAmount, 10, 64)
// paidAt    = time.Now().UTC().Truncate(time.Second)
// )

// func TestProcessorServer_Add(t *testing.T) {
// msgbusClient := &mbmocks.MsgBusServiceClient{}

// msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

// repo := &mocks.PaymentRepo{}
// s := server.NewProcessorServer(OrgName, repo, nil, nil, msgbusClient)

// t.Run("PaymentIsValid", func(tt *testing.T) {
// paymtReq := &pb.AddRequest{
// ItemId:        itemId.String(),
// ItemType:      itemTypePackge,
// Amount:        strAmount,
// Currency:      currency,
// PaymentMethod: paymentMethodMopay,
// PayerPhone:    payerPhone,
// Country:       country,
// Description:   description,
// }

// repo.On("Add", mock.Anything).Return(nil)

// resp, err := s.Add(context.TODO(), paymtReq)

// assert.NoError(t, err)
// assert.NotNil(t, resp)
// repo.AssertExpectations(t)
// })

// t.Run("ItemIdNotValid", func(tt *testing.T) {
// paymtReq := &pb.AddRequest{
// ItemId:        "lol",
// ItemType:      itemTypePackge,
// Amount:        strAmount,
// Currency:      currency,
// PaymentMethod: paymentMethodMopay,
// PayerPhone:    payerPhone,
// Country:       country,
// Description:   description,
// }
// resp, err := s.Add(context.TODO(), paymtReq)

// assert.Error(t, err)
// assert.Nil(t, resp)
// repo.AssertExpectations(t)
// })

// t.Run("ItemTypeNotValid", func(tt *testing.T) {
// paymtReq := &pb.AddRequest{
// ItemId:        itemId.String(),
// ItemType:      "lol",
// Amount:        strAmount,
// Currency:      currency,
// PaymentMethod: paymentMethodMopay,
// PayerPhone:    payerPhone,
// Country:       country,
// Description:   description,
// }

// resp, err := s.Add(context.TODO(), paymtReq)

// assert.Error(t, err)
// assert.Nil(t, resp)
// repo.AssertExpectations(t)
// })

// t.Run("AmountNotValid", func(tt *testing.T) {
// paymtReq := &pb.AddRequest{
// ItemId:        itemId.String(),
// ItemType:      itemTypePackge,
// Amount:        "0.3e",
// Currency:      currency,
// PaymentMethod: paymentMethodMopay,
// PayerPhone:    payerPhone,
// Country:       country,
// Description:   description,
// }

// resp, err := s.Add(context.TODO(), paymtReq)

// assert.Error(t, err)
// assert.Nil(t, resp)
// repo.AssertExpectations(t)
// })

// t.Run("AmountLessThanMinimum", func(tt *testing.T) {
// paymtReq := &pb.AddRequest{
// ItemId:        itemId.String(),
// ItemType:      itemTypePackge,
// Amount:        "0.3",
// Currency:      currency,
// PaymentMethod: paymentMethodMopay,
// PayerPhone:    payerPhone,
// Country:       country,
// Description:   description,
// }

// resp, err := s.Add(context.TODO(), paymtReq)

// assert.Error(t, err)
// assert.Nil(t, resp)
// repo.AssertExpectations(t)
// })
// }
