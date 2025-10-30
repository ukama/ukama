/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	uType "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	pb "github.com/ukama/ukama/systems/registry/invitation/pb/gen"
	invmocks "github.com/ukama/ukama/systems/registry/invitation/pb/gen/mocks"
)

func TestNewInvitationRegistry(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// This test is limited since NewInvitationRegistry creates a real gRPC connection
		// In a real scenario, you might want to use a test server or mock the connection
		invitationHost := "localhost:9090"
		timeout := 5 * time.Second

		// Note: This will fail if there's no actual invitation service running
		// In practice, you might want to use a test server or skip this test
		registry := NewInvitationRegistry(invitationHost, timeout)

		assert.NotNil(t, registry)
		assert.Equal(t, invitationHost, registry.host)
		assert.Equal(t, timeout, registry.timeout)
		assert.NotNil(t, registry.client)
		assert.NotNil(t, registry.conn)

		// Clean up
		registry.Close()
	})
}

func TestNewInvitationRegistryFromClient(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &invmocks.InvitationServiceClient{}
		registry := NewInvitationRegistryFromClient(mockClient)

		assert.NotNil(t, registry)
		assert.Equal(t, "localhost", registry.host)
		assert.Equal(t, 1*time.Second, registry.timeout)
		assert.Nil(t, registry.conn)
		assert.Equal(t, mockClient, registry.client)
	})
}

func TestInvitationRegistry_Close(t *testing.T) {
	t.Run("WithoutConnection", func(t *testing.T) {
		registry := &InvitationRegistry{
			conn: nil,
		}

		// Should not panic
		assert.NotPanics(t, func() {
			registry.Close()
		})
	})
}

func TestInvitationRegistry_RemoveInvitation(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &invmocks.InvitationServiceClient{}
		registry := NewInvitationRegistryFromClient(mockClient)

		expectedResponse := &pb.DeleteResponse{
			Id: "test-invitation-id",
		}

		mockClient.On("Delete", mock.Anything, &pb.DeleteRequest{Id: "test-invitation-id"}).
			Return(expectedResponse, nil)

		response, err := registry.RemoveInvitation("test-invitation-id")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &invmocks.InvitationServiceClient{}
		registry := NewInvitationRegistryFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "invitation not found")
		mockClient.On("Delete", mock.Anything, &pb.DeleteRequest{Id: "non-existent-id"}).
			Return(nil, expectedError)

		response, err := registry.RemoveInvitation("non-existent-id")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestInvitationRegistry_GetInvitationById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &invmocks.InvitationServiceClient{}
		registry := NewInvitationRegistryFromClient(mockClient)

		expectedResponse := &pb.GetResponse{
			Invitation: &pb.Invitation{
				Id:    "test-invitation-id",
				Name:  "John Doe",
				Email: "john@example.com",
				Role:  uType.RoleType_ROLE_ADMIN,
			},
		}

		mockClient.On("Get", mock.Anything, &pb.GetRequest{Id: "test-invitation-id"}).
			Return(expectedResponse, nil)

		response, err := registry.GetInvitationById("test-invitation-id")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &invmocks.InvitationServiceClient{}
		registry := NewInvitationRegistryFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "invitation not found")
		mockClient.On("Get", mock.Anything, &pb.GetRequest{Id: "non-existent-id"}).
			Return(nil, expectedError)

		response, err := registry.GetInvitationById("non-existent-id")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestInvitationRegistry_AddInvitation(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &invmocks.InvitationServiceClient{}
		registry := NewInvitationRegistryFromClient(mockClient)

		expectedResponse := &pb.AddResponse{
			Invitation: &pb.Invitation{
				Id:    "new-invitation-id",
				Name:  "Jane Doe",
				Email: "jane@example.com",
				Role:  uType.RoleType_ROLE_USER,
			},
		}

		mockClient.On("Add", mock.Anything, mock.MatchedBy(func(req *pb.AddRequest) bool {
			return req.Name == "Jane Doe" &&
				req.Email == "jane@example.com" &&
				req.Role == uType.RoleType_ROLE_USER
		})).Return(expectedResponse, nil)

		response, err := registry.AddInvitation("Jane Doe", "jane@example.com", "ROLE_USER")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &invmocks.InvitationServiceClient{}
		registry := NewInvitationRegistryFromClient(mockClient)

		expectedError := status.Error(codes.InvalidArgument, "invalid role")
		mockClient.On("Add", mock.Anything, mock.Anything).Return(nil, expectedError)

		response, err := registry.AddInvitation("Jane Doe", "jane@example.com", "INVALID_ROLE")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestInvitationRegistry_GetAllInvitations(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &invmocks.InvitationServiceClient{}
		registry := NewInvitationRegistryFromClient(mockClient)

		expectedResponse := &pb.GetAllResponse{
			Invitations: []*pb.Invitation{
				{
					Id:    "invitation-1",
					Name:  "John Doe",
					Email: "john@example.com",
					Role:  uType.RoleType_ROLE_OWNER,
				},
				{
					Id:    "invitation-2",
					Name:  "Jane Doe",
					Email: "jane@example.com",
					Role:  uType.RoleType_ROLE_USER,
				},
			},
		}

		mockClient.On("GetAll", mock.Anything, &pb.GetAllRequest{}).
			Return(expectedResponse, nil)

		response, err := registry.GetAllInvitations()

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &invmocks.InvitationServiceClient{}
		registry := NewInvitationRegistryFromClient(mockClient)

		expectedError := status.Error(codes.Internal, "internal server error")
		mockClient.On("GetAll", mock.Anything, &pb.GetAllRequest{}).
			Return(nil, expectedError)

		response, err := registry.GetAllInvitations()

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestInvitationRegistry_UpdateInvitation(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &invmocks.InvitationServiceClient{}
		registry := NewInvitationRegistryFromClient(mockClient)

		expectedResponse := &pb.UpdateStatusResponse{
			Id:     "test-invitation-id",
			Status: uType.InvitationStatus_INVITE_ACCEPTED,
		}

		mockClient.On("UpdateStatus", mock.Anything, mock.MatchedBy(func(req *pb.UpdateStatusRequest) bool {
			return req.Id == "test-invitation-id" &&
				req.Email == "john@example.com" &&
				req.Status == uType.InvitationStatus_INVITE_ACCEPTED
		})).Return(expectedResponse, nil)

		response, err := registry.UpdateInvitation("test-invitation-id", "INVITE_ACCEPTED", "john@example.com")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &invmocks.InvitationServiceClient{}
		registry := NewInvitationRegistryFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "invitation not found")
		mockClient.On("UpdateStatus", mock.Anything, mock.Anything).Return(nil, expectedError)

		response, err := registry.UpdateInvitation("non-existent-id", "INVITE_ACCEPTED", "john@example.com")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestInvitationRegistry_GetInvitationsByEmail(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &invmocks.InvitationServiceClient{}
		registry := NewInvitationRegistryFromClient(mockClient)

		expectedResponse := &pb.GetByEmailResponse{
			Invitation: &pb.Invitation{
				Id:    "invitation-1",
				Name:  "John Doe",
				Email: "john@example.com",
				Role:  uType.RoleType_ROLE_OWNER,
			},
		}

		mockClient.On("GetByEmail", mock.Anything, &pb.GetByEmailRequest{Email: "john@example.com"}).
			Return(expectedResponse, nil)

		response, err := registry.GetInvitationsByEmail("john@example.com")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &invmocks.InvitationServiceClient{}
		registry := NewInvitationRegistryFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "no invitations found for email")
		mockClient.On("GetByEmail", mock.Anything, &pb.GetByEmailRequest{Email: "nonexistent@example.com"}).
			Return(nil, expectedError)

		response, err := registry.GetInvitationsByEmail("nonexistent@example.com")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}
