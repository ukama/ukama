/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"
	"testing"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/inventory/accounting/mocks"
	pb "github.com/ukama/ukama/systems/inventory/accounting/pb/gen"
	"github.com/ukama/ukama/systems/inventory/accounting/pkg/db"
	"gorm.io/gorm"
)

const OrgName = "testorg"

func TestAccountingServer_Get(t *testing.T) {
	t.Run("Account record found", func(t *testing.T) {
		var aId = uuid.NewV4()

		accRepo := &mocks.AccountingRepo{}

		accRepo.On("Get", aId).Return(
			&db.Accounting{Id: aId}, nil).Once()

		s := NewAccountingServer(OrgName, accRepo, nil, "", nil, "")
		accResp, err := s.Get(context.TODO(), &pb.GetRequest{
			Id: aId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, accResp)
		assert.Equal(t, aId.String(), accResp.Accounting.Id)

		accRepo.AssertExpectations(t)
	})

	t.Run("Account record not found", func(t *testing.T) {
		var aId = uuid.NewV4()

		accRepo := &mocks.AccountingRepo{}

		accRepo.On("Get", aId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewAccountingServer(OrgName, accRepo, nil, "", nil, "")
		accResp, err := s.Get(context.TODO(), &pb.GetRequest{
			Id: aId.String()})

		assert.Error(t, err)
		assert.Nil(t, accResp)
		accRepo.AssertExpectations(t)
	})
}

func TestAccountingServer_GetByUser(t *testing.T) {
	t.Run("Account records by user found", func(t *testing.T) {
		var uId = uuid.NewV4()
		var aId = uuid.NewV4()

		accRepo := &mocks.AccountingRepo{}

		accRepo.On("GetByUser", uId.String()).Return(
			[]*db.Accounting{
				{
					Id:            aId,
					Vat:           "10",
					Item:          "Product-1",
					UserId:        uId.String(),
					Inventory:     "1",
					OpexFee:       "100",
					EffectiveDate: "2023-01-01",
					Description:   "Product-1 description",
				}}, nil).Once()

		s := NewAccountingServer(OrgName, accRepo, nil, "", nil, "")

		accResp, err := s.GetByUser(context.TODO(),
			&pb.GetByUserRequest{
				UserId: uId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, accResp)
		assert.Equal(t, uId.String(), accResp.GetAccounting()[0].GetUserId())
		accRepo.AssertExpectations(t)
	})
}
