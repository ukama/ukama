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
	"github.com/ukama/ukama/systems/inventory/component/mocks"
	pb "github.com/ukama/ukama/systems/inventory/component/pb/gen"
	"github.com/ukama/ukama/systems/inventory/component/pkg/db"
	"gorm.io/gorm"
)

const OrgName = "testorg"

func TestComponentServer_Get(t *testing.T) {
	t.Run("Component record found", func(t *testing.T) {
		var cId = uuid.NewV4()

		compRepo := &mocks.ComponentRepo{}

		compRepo.On("Get", cId).Return(
			&db.Component{Id: cId}, nil).Once()

		s := NewComponentServer(OrgName, compRepo, nil, "", nil, "", "", "")
		compResp, err := s.Get(context.TODO(), &pb.GetRequest{
			Id: cId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, compResp)
		assert.Equal(t, cId.String(), compResp.Component.Id)

		compRepo.AssertExpectations(t)
	})

	t.Run("Component record not found", func(t *testing.T) {
		var cId = uuid.NewV4()

		compRepo := &mocks.ComponentRepo{}

		compRepo.On("Get", cId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewComponentServer(OrgName, compRepo, nil, "", nil, "", "", "")
		compResp, err := s.Get(context.TODO(), &pb.GetRequest{
			Id: cId.String()})

		assert.Error(t, err)
		assert.Nil(t, compResp)
		compRepo.AssertExpectations(t)
	})
}

func TestComponentServer_GetByUser(t *testing.T) {
	t.Run("Component records by user found", func(t *testing.T) {
		var uId = uuid.NewV4()
		var cId = uuid.NewV4()

		compRepo := &mocks.ComponentRepo{}

		compRepo.On("GetByUser", uId.String(), int32(1)).Return(
			[]*db.Component{
				{
					Id:            cId,
					UserId:        uId,
					Inventory:     "2",
					Category:      1,
					Type:          "tower node",
					Description:   "best tower node",
					DatasheetURL:  "https://datasheet.com",
					ImagesURL:     "https://images.com",
					PartNumber:    "1234",
					Manufacturer:  "ukama",
					Managed:       "ukama",
					Warranty:      1,
					Specification: "spec",
				}}, nil).Once()

		s := NewComponentServer(OrgName, compRepo, nil, "", nil, "", "", "")

		compResp, err := s.GetByUser(context.TODO(),
			&pb.GetByUserRequest{
				UserId:   uId.String(),
				Category: pb.ComponentCategory_ACCESS})

		assert.NoError(t, err)
		assert.NotNil(t, compResp)
		assert.Equal(t, uId.String(), compResp.GetComponents()[0].GetUserId())
		assert.Equal(t, cId.String(), compResp.GetComponents()[0].GetId())
		compRepo.AssertExpectations(t)
	})
}
