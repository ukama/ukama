/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	uconf "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/notification"
	"github.com/ukama/ukama/systems/common/uuid"
	enpb "github.com/ukama/ukama/systems/notification/event-notify/pb/gen"
	"google.golang.org/grpc"
)

// MockEventNotifyClientProvider is a mock implementation of the EventNotifyClientProvider interface
type MockEventNotifyClientProvider struct {
	mock.Mock
}

func (m *MockEventNotifyClientProvider) GetClient() (enpb.EventToNotifyServiceClient, error) {
	args := m.Called()
	return args.Get(0).(enpb.EventToNotifyServiceClient), args.Error(1)
}

// MockEventToNotifyServiceClient is a mock implementation of the EventToNotifyServiceClient interface
type MockEventToNotifyServiceClient struct {
	mock.Mock
}

func (m *MockEventToNotifyServiceClient) Get(ctx context.Context, in *enpb.GetRequest, opts ...grpc.CallOption) (*enpb.GetResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*enpb.GetResponse), args.Error(1)
}

func (m *MockEventToNotifyServiceClient) GetAll(ctx context.Context, in *enpb.GetAllRequest, opts ...grpc.CallOption) (*enpb.GetAllResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*enpb.GetAllResponse), args.Error(1)
}

func (m *MockEventToNotifyServiceClient) UpdateStatus(ctx context.Context, in *enpb.UpdateStatusRequest, opts ...grpc.CallOption) (*enpb.UpdateStatusResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*enpb.UpdateStatusResponse), args.Error(1)
}

func TestNewNotifyHandler_Success(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	// Act
	handler := NewNotifyHandler(db, mockProvider)

	// Assert
	assert.NotNil(t, handler)
	assert.Equal(t, db, handler.Db)
	assert.Equal(t, mockClient, handler.c)
	assert.NotNil(t, handler.done)
	assert.NotNil(t, handler.subs)
	assert.Equal(t, 0, len(handler.subs))

	mockProvider.AssertExpectations(t)
}

func TestNewNotifyHandler_GetClientError(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	expectedError := errors.New("connection failed")
	mockProvider.On("GetClient").Return(nil, expectedError).Once()

	// Act & Assert
	// This should panic due to log.Fatalf in the original code
	// We'll test this by checking if it panics
	assert.Panics(t, func() {
		NewNotifyHandler(db, mockProvider)
	})

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Register_Success(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	orgId := "test-org-id"
	networkId := "test-network-id"
	subscriberId := "test-subscriber-id"
	userId := "test-user-id"
	scopes := []notification.NotificationScope{
		notification.SCOPE_ORG,
		notification.SCOPE_NETWORK,
	}

	// Act
	id, sub := handler.Register(orgId, networkId, subscriberId, userId, scopes)

	// Assert
	assert.NotEmpty(t, id)
	assert.NotNil(t, sub)
	assert.Equal(t, orgId, sub.OrgId)
	assert.Equal(t, networkId, sub.NetworkId)
	assert.Equal(t, subscriberId, sub.SubscriberId)
	assert.Equal(t, userId, sub.UserId)
	assert.Equal(t, scopes, sub.Scopes)
	assert.NotNil(t, sub.DataChan)
	assert.NotNil(t, sub.QuitChan)
	assert.Equal(t, BufferCapacity, cap(sub.DataChan))

	// Verify the subscription is stored in the handler
	assert.Equal(t, 1, len(handler.subs))
	assert.Contains(t, handler.subs, id)
	assert.Equal(t, *sub, handler.subs[id])

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Register_EmptyValues(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	orgId := ""
	networkId := ""
	subscriberId := ""
	userId := ""
	scopes := []notification.NotificationScope{}

	// Act
	id, sub := handler.Register(orgId, networkId, subscriberId, userId, scopes)

	// Assert
	assert.NotEmpty(t, id)
	assert.NotNil(t, sub)
	assert.Equal(t, orgId, sub.OrgId)
	assert.Equal(t, networkId, sub.NetworkId)
	assert.Equal(t, subscriberId, sub.SubscriberId)
	assert.Equal(t, userId, sub.UserId)
	assert.Equal(t, scopes, sub.Scopes)
	assert.NotNil(t, sub.DataChan)
	assert.NotNil(t, sub.QuitChan)

	// Verify the subscription is stored in the handler
	assert.Equal(t, 1, len(handler.subs))
	assert.Contains(t, handler.subs, id)

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Register_MultipleSubscriptions(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	// Act - Register multiple subscriptions
	id1, sub1 := handler.Register("org1", "network1", "sub1", "user1", []notification.NotificationScope{notification.SCOPE_ORG})
	id2, sub2 := handler.Register("org2", "network2", "sub2", "user2", []notification.NotificationScope{notification.SCOPE_NETWORK})
	id3, sub3 := handler.Register("org3", "network3", "sub3", "user3", []notification.NotificationScope{notification.SCOPE_SUBSCRIBER})

	// Assert
	assert.NotEqual(t, id1, id2)
	assert.NotEqual(t, id2, id3)
	assert.NotEqual(t, id1, id3)

	assert.NotNil(t, sub1)
	assert.NotNil(t, sub2)
	assert.NotNil(t, sub3)

	// Verify all subscriptions are stored
	assert.Equal(t, 3, len(handler.subs))
	assert.Contains(t, handler.subs, id1)
	assert.Contains(t, handler.subs, id2)
	assert.Contains(t, handler.subs, id3)

	// Verify subscription details
	assert.Equal(t, "org1", sub1.OrgId)
	assert.Equal(t, "org2", sub2.OrgId)
	assert.Equal(t, "org3", sub3.OrgId)

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Register_UniqueIds(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	// Act - Register multiple subscriptions with same parameters
	id1, _ := handler.Register("org1", "network1", "sub1", "user1", []notification.NotificationScope{notification.SCOPE_ORG})
	id2, _ := handler.Register("org1", "network1", "sub1", "user1", []notification.NotificationScope{notification.SCOPE_ORG})
	id3, _ := handler.Register("org1", "network1", "sub1", "user1", []notification.NotificationScope{notification.SCOPE_ORG})

	// Assert - Each registration should have a unique ID
	assert.NotEqual(t, id1, id2)
	assert.NotEqual(t, id2, id3)
	assert.NotEqual(t, id1, id3)

	// Verify all subscriptions are stored
	assert.Equal(t, 3, len(handler.subs))

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Register_ChannelCapacity(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	// Act
	_, sub := handler.Register("org1", "network1", "sub1", "user1", []notification.NotificationScope{notification.SCOPE_ORG})

	// Assert - Check channel capacity
	assert.Equal(t, BufferCapacity, cap(sub.DataChan))
	assert.Equal(t, 0, cap(sub.QuitChan)) // QuitChan should be unbuffered

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Register_AllScopes(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	// Test all available scopes
	allScopes := []notification.NotificationScope{
		notification.SCOPE_INVALID,
		notification.SCOPE_OWNER,
		notification.SCOPE_ORG,
		notification.SCOPE_NETWORKS,
		notification.SCOPE_NETWORK,
		notification.SCOPE_SITES,
		notification.SCOPE_SITE,
		notification.SCOPE_SUBSCRIBERS,
		notification.SCOPE_SUBSCRIBER,
		notification.SCOPE_USERS,
		notification.SCOPE_USER,
		notification.SCOPE_NODE,
	}

	// Act
	_, sub := handler.Register("org1", "network1", "sub1", "user1", allScopes)

	// Assert
	assert.Equal(t, allScopes, sub.Scopes)
	assert.Equal(t, len(allScopes), len(sub.Scopes))

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Register_WithUUID(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	// Act
	id, sub := handler.Register("org1", "network1", "sub1", "user1", []notification.NotificationScope{notification.SCOPE_ORG})

	// Assert - Verify UUID format
	_, err := uuid.FromString(id)
	assert.NoError(t, err)
	assert.Equal(t, id, sub.Id.String())

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Register_ConcurrentAccess(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	// Act - Register multiple subscriptions concurrently
	done := make(chan bool, 10)
	ids := make([]string, 10)

	for i := 0; i < 10; i++ {
		go func(index int) {
			id, _ := handler.Register(
				"org"+string(rune(index)),
				"network"+string(rune(index)),
				"sub"+string(rune(index)),
				"user"+string(rune(index)),
				[]notification.NotificationScope{notification.SCOPE_ORG},
			)
			ids[index] = id
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Assert
	assert.Equal(t, 10, len(handler.subs))

	// Verify all IDs are unique
	idSet := make(map[string]bool)
	for _, id := range ids {
		assert.NotEmpty(t, id)
		assert.False(t, idSet[id], "Duplicate ID found: %s", id)
		idSet[id] = true
	}

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Deregister_Success(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	// Register a subscription first
	id, _ := handler.Register("org1", "network1", "sub1", "user1", []notification.NotificationScope{notification.SCOPE_ORG})

	// Verify subscription exists
	assert.Equal(t, 1, len(handler.subs))
	assert.Contains(t, handler.subs, id)

	// Act
	err := handler.Deregister(id)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, len(handler.subs))
	assert.NotContains(t, handler.subs, id)

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Deregister_NotFound(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	// Act
	err := handler.Deregister("non-existent-id")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sub with id non-existent-id not found")
	assert.Equal(t, 0, len(handler.subs))

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Deregister_EmptyId(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	// Act
	err := handler.Deregister("")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sub with id  not found")
	assert.Equal(t, 0, len(handler.subs))

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Deregister_MultipleSubscriptions(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	// Register multiple subscriptions
	id1, _ := handler.Register("org1", "network1", "sub1", "user1", []notification.NotificationScope{notification.SCOPE_ORG})
	id2, _ := handler.Register("org2", "network2", "sub2", "user2", []notification.NotificationScope{notification.SCOPE_NETWORK})
	id3, _ := handler.Register("org3", "network3", "sub3", "user3", []notification.NotificationScope{notification.SCOPE_SUBSCRIBER})

	// Verify all subscriptions exist
	assert.Equal(t, 3, len(handler.subs))

	// Act - Deregister one subscription
	err := handler.Deregister(id2)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(handler.subs))
	assert.Contains(t, handler.subs, id1)
	assert.NotContains(t, handler.subs, id2)
	assert.Contains(t, handler.subs, id3)

	// Deregister another subscription
	err = handler.Deregister(id1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(handler.subs))
	assert.NotContains(t, handler.subs, id1)
	assert.NotContains(t, handler.subs, id2)
	assert.Contains(t, handler.subs, id3)

	// Deregister the last subscription
	err = handler.Deregister(id3)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(handler.subs))
	assert.NotContains(t, handler.subs, id1)
	assert.NotContains(t, handler.subs, id2)
	assert.NotContains(t, handler.subs, id3)

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Stop_NoSubscriptions(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	// Verify no subscriptions exist
	assert.Equal(t, 0, len(handler.subs))

	// Act
	handler.Stop()

	// Assert

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Stop_WithSingleSubscription(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	// Register a subscription
	handler.Register("org1", "network1", "sub1", "user1", []notification.NotificationScope{notification.SCOPE_ORG})

	// Verify subscription exists
	assert.Equal(t, 1, len(handler.subs))

	// Act
	handler.Stop()

	// Assert

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Stop_WithMultipleSubscriptions(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	// Register multiple subscriptions
	handler.Register("org1", "network1", "sub1", "user1", []notification.NotificationScope{notification.SCOPE_ORG})
	handler.Register("org2", "network2", "sub2", "user2", []notification.NotificationScope{notification.SCOPE_NETWORK})
	handler.Register("org3", "network3", "sub3", "user3", []notification.NotificationScope{notification.SCOPE_SUBSCRIBER})

	// Verify all subscriptions exist
	assert.Equal(t, 3, len(handler.subs))

	// Act
	handler.Stop()

	// Assert

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Stop_AfterDeregister(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	// Register multiple subscriptions
	handler.Register("org1", "network1", "sub1", "user1", []notification.NotificationScope{notification.SCOPE_ORG})
	id2, _ := handler.Register("org2", "network2", "sub2", "user2", []notification.NotificationScope{notification.SCOPE_NETWORK})
	handler.Register("org3", "network3", "sub3", "user3", []notification.NotificationScope{notification.SCOPE_SUBSCRIBER})

	// Verify all subscriptions exist
	assert.Equal(t, 3, len(handler.subs))

	// Deregister one subscription
	err := handler.Deregister(id2)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(handler.subs))

	// Act
	handler.Stop()

	// Assert

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Start_Success(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	// Act
	handler.Start()

	// Assert

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Start_MultipleCalls(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	// Act - Call Start multiple times
	handler.Start()
	handler.Start()
	handler.Start()

	// Assert

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Start_WithSubscriptions(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	// Register subscriptions before starting
	handler.Register("org1", "network1", "sub1", "user1", []notification.NotificationScope{notification.SCOPE_ORG})
	handler.Register("org2", "network2", "sub2", "user2", []notification.NotificationScope{notification.SCOPE_NETWORK})

	// Verify subscriptions exist
	assert.Equal(t, 2, len(handler.subs))

	// Act
	handler.Start()

	// Assert
	// Start should work correctly even when subscriptions already exist

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Start_EmptyHandler(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	// Verify no subscriptions exist
	assert.Equal(t, 0, len(handler.subs))

	// Act
	handler.Start()

	// Assert
	// Start should work correctly even with no subscriptions

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Start_ConcurrentAccess(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	// Act - Start multiple goroutines concurrently
	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func() {
			handler.Start()
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 5; i++ {
		<-done
	}

	// Assert

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Start_Stop_Integration(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	// Register a subscription
	handler.Register("org1", "network1", "sub1", "user1", []notification.NotificationScope{notification.SCOPE_ORG})

	// Act - Start the handler
	handler.Start()

	// Give the goroutine a moment to start (though it will fail to connect to DB in test environment)
	time.Sleep(10 * time.Millisecond)

	// Stop the handler
	handler.Stop()

	// Assert

	mockProvider.AssertExpectations(t)
}

func TestNotifyHandler_Start_Stop_WithoutStart(t *testing.T) {
	// Arrange
	mockProvider := &MockEventNotifyClientProvider{}
	mockClient := &MockEventToNotifyServiceClient{}
	db := &uconf.Database{
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
		Username: "testuser",
		Password: "testpass",
	}

	mockProvider.On("GetClient").Return(mockClient, nil).Once()

	handler := NewNotifyHandler(db, mockProvider)

	// Register a subscription
	handler.Register("org1", "network1", "sub1", "user1", []notification.NotificationScope{notification.SCOPE_ORG})

	// Act - Stop without starting
	handler.Stop()

	// Assert

	mockProvider.AssertExpectations(t)
}
