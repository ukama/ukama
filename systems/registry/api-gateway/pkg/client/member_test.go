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

	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	pb "github.com/ukama/ukama/systems/registry/member/pb/gen"
	memmocks "github.com/ukama/ukama/systems/registry/member/pb/gen/mocks"
)

func TestNewMemberRegistry(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// This test is limited since NewMemberRegistry creates a real gRPC connection
		// In a real scenario, you might want to use a test server or mock the connection
		memberHost := "localhost:9090"
		timeout := 5 * time.Second

		// Note: This will fail if there's no actual member service running
		// In practice, you might want to use a test server or skip this test
		registry := NewMemberRegistry(memberHost, timeout)

		assert.NotNil(t, registry)
		assert.Equal(t, memberHost, registry.host)
		assert.Equal(t, timeout, registry.timeout)
		assert.NotNil(t, registry.client)
		assert.NotNil(t, registry.conn)

		// Clean up
		registry.Close()
	})
}

func TestNewRegistryFromClient(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &memmocks.MemberServiceClient{}
		registry := NewRegistryFromClient(mockClient)

		assert.NotNil(t, registry)
		assert.Equal(t, "localhost", registry.host)
		assert.Equal(t, 1*time.Second, registry.timeout)
		assert.Nil(t, registry.conn)
		assert.Equal(t, mockClient, registry.client)
	})
}

func TestMemberRegistry_Close(t *testing.T) {
	t.Run("WithoutConnection", func(t *testing.T) {
		registry := &MemberRegistry{
			conn: nil,
		}

		// Should not panic
		assert.NotPanics(t, func() {
			registry.Close()
		})
	})
}

func TestMemberRegistry_GetMember(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &memmocks.MemberServiceClient{}
		registry := NewRegistryFromClient(mockClient)

		expectedResponse := &pb.MemberResponse{
			Member: &pb.Member{
				MemberId: "test-member-id",
				UserId:   "test-user-uuid",
				Role:     upb.RoleType_ROLE_OWNER,
			},
		}

		mockClient.On("GetMember", mock.Anything, &pb.MemberRequest{MemberId: "test-member-id"}).
			Return(expectedResponse, nil)

		response, err := registry.GetMember("test-member-id")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &memmocks.MemberServiceClient{}
		registry := NewRegistryFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "member not found")
		mockClient.On("GetMember", mock.Anything, &pb.MemberRequest{MemberId: "non-existent-id"}).
			Return(nil, expectedError)

		response, err := registry.GetMember("non-existent-id")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestMemberRegistry_GetMemberByUserId(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &memmocks.MemberServiceClient{}
		registry := NewRegistryFromClient(mockClient)

		expectedResponse := &pb.GetMemberByUserIdResponse{
			Member: &pb.Member{
				MemberId: "test-member-id",
				UserId:   "test-user-uuid",
				Role:     upb.RoleType_ROLE_OWNER,
			},
		}

		mockClient.On("GetMemberByUserId", mock.Anything, &pb.GetMemberByUserIdRequest{MemberId: "test-user-uuid"}).
			Return(expectedResponse, nil)

		response, err := registry.GetMemberByUserId("test-user-uuid")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &memmocks.MemberServiceClient{}
		registry := NewRegistryFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "member not found for user")
		mockClient.On("GetMemberByUserId", mock.Anything, &pb.GetMemberByUserIdRequest{MemberId: "non-existent-user"}).
			Return(nil, expectedError)

		response, err := registry.GetMemberByUserId("non-existent-user")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestMemberRegistry_GetMembers(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &memmocks.MemberServiceClient{}
		registry := NewRegistryFromClient(mockClient)

		expectedResponse := &pb.GetMembersResponse{
			Members: []*pb.Member{
				{
					MemberId: "member-1",
					UserId:   "user-1",
					Role:     upb.RoleType_ROLE_OWNER,
				},
				{
					MemberId: "member-2",
					UserId:   "user-2",
					Role:     upb.RoleType_ROLE_USER,
				},
			},
		}

		mockClient.On("GetMembers", mock.Anything, &pb.GetMembersRequest{}).
			Return(expectedResponse, nil)

		response, err := registry.GetMembers()

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &memmocks.MemberServiceClient{}
		registry := NewRegistryFromClient(mockClient)

		expectedError := status.Error(codes.Internal, "internal server error")
		mockClient.On("GetMembers", mock.Anything, &pb.GetMembersRequest{}).
			Return(nil, expectedError)

		response, err := registry.GetMembers()

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestMemberRegistry_AddMember(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &memmocks.MemberServiceClient{}
		registry := NewRegistryFromClient(mockClient)

		expectedResponse := &pb.MemberResponse{
			Member: &pb.Member{
				MemberId: "new-member-id",
				UserId:   "test-user-uuid",
				Role:     upb.RoleType_ROLE_USER,
			},
		}

		mockClient.On("AddMember", mock.Anything, mock.MatchedBy(func(req *pb.AddMemberRequest) bool {
			return req.UserUuid == "test-user-uuid" &&
				req.Role == upb.RoleType_ROLE_USER
		})).Return(expectedResponse, nil)

		response, err := registry.AddMember("test-user-uuid", "ROLE_USER")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &memmocks.MemberServiceClient{}
		registry := NewRegistryFromClient(mockClient)

		expectedError := status.Error(codes.AlreadyExists, "member already exists")
		mockClient.On("AddMember", mock.Anything, mock.Anything).Return(nil, expectedError)

		response, err := registry.AddMember("existing-user-uuid", "ROLE_USER")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestMemberRegistry_UpdateMember(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &memmocks.MemberServiceClient{}
		registry := NewRegistryFromClient(mockClient)

		// UpdateMember returns only an error, no response

		mockClient.On("UpdateMember", mock.Anything, mock.MatchedBy(func(req *pb.UpdateMemberRequest) bool {
			return req.MemberId == "test-member-id" &&
				req.IsDeactivated == true
		})).Return(nil, nil)

		err := registry.UpdateMember("test-member-id", true, "ROLE_ADMIN")

		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &memmocks.MemberServiceClient{}
		registry := NewRegistryFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "member not found")
		mockClient.On("UpdateMember", mock.Anything, mock.Anything).Return(nil, expectedError)

		err := registry.UpdateMember("non-existent-id", false, "ROLE_USER")

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestMemberRegistry_RemoveMember(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &memmocks.MemberServiceClient{}
		registry := NewRegistryFromClient(mockClient)

		expectedResponse := &pb.MemberResponse{
			Member: &pb.Member{
				MemberId: "test-member-id",
				UserId:   "test-user-uuid",
				Role:     upb.RoleType_ROLE_OWNER,
			},
		}

		mockClient.On("RemoveMember", mock.Anything, &pb.MemberRequest{MemberId: "test-member-id"}).
			Return(expectedResponse, nil)

		err := registry.RemoveMember("test-member-id")

		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &memmocks.MemberServiceClient{}
		registry := NewRegistryFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "member not found")
		mockClient.On("RemoveMember", mock.Anything, &pb.MemberRequest{MemberId: "non-existent-id"}).
			Return(nil, expectedError)

		err := registry.RemoveMember("non-existent-id")

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}
