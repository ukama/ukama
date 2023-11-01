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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/software/mocks"
	"github.com/ukama/ukama/systems/node/software/pb/gen"
	"github.com/ukama/ukama/systems/node/software/pkg/db"
)

const testOrgName = "test-org"

var orgId = uuid.NewV4()

func Test_CreateSoftwareUpdate(t *testing.T) {
	ctx := context.Background()

	// Create a mock for the SoftwareManagerRepo interface
	softwareManager := &mocks.SoftwareRepo{}
	healthsvc := &mocks.HealthClientProvider{}

	// Configure the mock to expect a call to CreateSoftwareUpdate with specific arguments
	softwareManager.On("CreateSoftwareUpdate", mock.Anything, mock.Anything).
		Return(nil) // You can specify the expected return value here

	// Create a mock for the MsgBusServiceClient interface
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	// Configure the mock to expect a call to PublishRequest with specific arguments

	// Create an instance of the SoftwareManagerServer with the mocks
	s := NewSoftwareServer(testOrgName, softwareManager, msgclientRepo, false, healthsvc)

	// Test
	r, err := s.CreateSoftwareUpdate(ctx, &gen.CreateSoftwareUpdateRequest{
		Name: "test",
		Tag:  "test",
	})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, r)

	// Verify that the expected methods on the mocks were called
	softwareManager.AssertExpectations(t)
	msgclientRepo.AssertExpectations(t)
}
func Test_GetLatestSoftwareUpdate(t *testing.T) {
	ctx := context.Background()

	// Create a mock for the SoftwareManagerRepo interface
	softwareManager := &mocks.SoftwareRepo{}
	healthsvc := &mocks.HealthClientProvider{}

	// Configure the mock to expect a call to GetLatestSoftwareUpdate and return a *db.Software
	softwareManager.On("GetLatestSoftwareUpdate").
		Return(&db.Software{
			Id:   uuid.NewV4(),
			Name: "test",
			Tag:  "test",
			// Add other fields as needed
		}, nil)

	// Create a mock for the MsgBusServiceClient interface
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	// Create an instance of the SoftwareManagerServer with the mocks
	s := NewSoftwareServer(testOrgName, softwareManager, msgclientRepo, false, healthsvc)

	// Test
	r, err := s.GetLatestSoftwareUpdate(ctx, &gen.GetLatestSoftwareUpdateRequest{})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, r)

	// Verify that the expected methods on the mocks were called
	softwareManager.AssertExpectations(t)
	msgclientRepo.AssertExpectations(t)
}
