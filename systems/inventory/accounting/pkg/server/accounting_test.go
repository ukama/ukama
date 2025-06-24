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
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
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

	t.Run("Invalid UUID format", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}

		s := NewAccountingServer(OrgName, accRepo, nil, "", nil, "")
		accResp, err := s.Get(context.TODO(), &pb.GetRequest{
			Id: "invalid-uuid-format"})

		assert.Error(t, err)
		assert.Nil(t, accResp)
		assert.Contains(t, err.Error(), "invalid format of accounting uuid")
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
					UserId:        uId,
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

	t.Run("Database error when getting accountings by user", func(t *testing.T) {
		var uId = uuid.NewV4()

		accRepo := &mocks.AccountingRepo{}

		accRepo.On("GetByUser", uId.String()).Return(nil, gorm.ErrInvalidDB).Once()

		s := NewAccountingServer(OrgName, accRepo, nil, "", nil, "")

		accResp, err := s.GetByUser(context.TODO(),
			&pb.GetByUserRequest{
				UserId: uId.String()})

		assert.Error(t, err)
		assert.Nil(t, accResp)
		assert.Contains(t, err.Error(), "invalid db")
		accRepo.AssertExpectations(t)
	})
}

func TestAccountingServer_SyncAccounting(t *testing.T) {
	t.Run("Successfully sync accounting", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		userId := uuid.NewV4()

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(`{"test":[{"company":"TestCompany","git_branch_name":"test-branch","email":"test@example.com","user_id":"`+userId.String()+`"}]}`), nil)
		gitClient.On("BranchCheckout", "test-branch").Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(`{"ukama":[{"item":"TestItem","description":"Test Description","inventory":"10","opex_fee":"100.00","vat":"20.00","effective_date":"2023-01-01"}],"backhaul":[{"item":"BackhaulItem","description":"Backhaul Description","inventory":"5","opex_fee":"50.00","vat":"10.00","effective_date":"2023-01-01"}]}`), nil)

		accRepo.On("Delete").Return(nil)
		accRepo.On("Add", mock.AnythingOfType("[]*db.Accounting")).Return(nil)
		accRepo.On("GetByUser", userId.String()).Return([]*db.Accounting{
			{
				Id:            uuid.NewV4(),
				UserId:        userId,
				Item:          "TestItem",
				Description:   "Test Description",
				Inventory:     "10",
				OpexFee:       "100.00",
				Vat:           "20.00",
				EffectiveDate: "2023-01-01",
			},
		}, nil)

		msgBus.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)

		s := NewAccountingServer(OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		gitClient.AssertExpectations(t)
		accRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Failed to clone git repo", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(errors.New("git clone failed"))

		s := NewAccountingServer(OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to clone git repo")

		gitClient.AssertExpectations(t)
	})
}
