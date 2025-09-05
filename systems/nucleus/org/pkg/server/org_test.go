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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/nucleus/org/mocks"
	pb "github.com/ukama/ukama/systems/nucleus/org/pb/gen"
	"github.com/ukama/ukama/systems/nucleus/org/pkg/db"
	"github.com/ukama/ukama/systems/nucleus/org/pkg/providers"
	"gorm.io/gorm"
)

const OrgName = "testorg"

func TestOrgServer_Add(t *testing.T) {
	// Arrange
	orgName := "org-1"
	ownerUuid := uuid.NewV4()
	certificate := "ukama_certs"
	msgclientRepo := &cmocks.MsgBusServiceClient{}

	orgRepo := &mocks.OrgRepo{}
	userRepo := &mocks.UserRepo{}
	orchSystem := &mocks.OrchestratorProvider{}
	registryClient := &cmocks.MemberClient{}
	org := &db.Org{
		Owner:       ownerUuid,
		Certificate: certificate,
		Name:        orgName,
	}

	pOrg := providers.DeployOrgRequest{
		OrgId:   uuid.NewV4().String(),
		OrgName: org.Name,
		OwnerId: org.Owner.String(),
	}

	orgRepo.On("Add", org, mock.Anything).Return(nil).Once()

	orchSystem.On("DeployOrg", pOrg).Return(&providers.DeployOrgResponse{}, nil)

	msgclientRepo.On("PublishRequest", mock.Anything, mock.MatchedBy(func(e *epb.EventOrgCreate) bool {
		return e.Name == org.Name
	})).Return(nil).Once()

	orgRepo.On("GetOrgCount").Return(int64(1), int64(0), nil).Once()

	s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registryClient, msgclientRepo, "", true)

	t.Run("AddValidOrg", func(tt *testing.T) {
		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Org: &pb.Organization{
				Name:        orgName,
				Owner:       ownerUuid.String(),
				Certificate: certificate,
			}})

		// Assert
		msgclientRepo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, orgName, res.Org.Name)
		assert.Equal(t, ownerUuid.String(), res.Org.Owner)
		orgRepo.AssertExpectations(t)
	})

	t.Run("AddOrgWithoutOwner", func(tt *testing.T) {
		// Act
		orgResp, err := s.Add(context.TODO(), &pb.AddRequest{
			Org: &pb.Organization{Name: OrgName},
		})

		assert.Error(t, err)
		assert.Nil(t, orgResp)
	})

	t.Run("AddOrgWithInvalidOwner", func(tt *testing.T) {
		owner := "org-1"

		// Act
		orgResp, err := s.Add(context.TODO(), &pb.AddRequest{
			Org: &pb.Organization{Owner: owner},
		})

		assert.Error(t, err)
		assert.Nil(t, orgResp)
	})

	t.Run("AddOrgWithAdditionalFields", func(tt *testing.T) {
		// Arrange
		orgName := "org-with-fields"
		ownerUuid := uuid.NewV4()
		certificate := "ukama_certs"
		country := "US"
		currency := "USD"
		msgclientRepo := &cmocks.MsgBusServiceClient{}

		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registryClient := &cmocks.MemberClient{}
		org := &db.Org{
			Owner:       ownerUuid,
			Certificate: certificate,
			Name:        orgName,
			Country:     country,
			Currency:    currency,
		}

		orgRepo.On("Add", org, mock.Anything).Return(nil).Once()
		orchSystem.On("DeployOrg", mock.Anything).Return(&providers.DeployOrgResponse{}, nil)
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()
		orgRepo.On("GetOrgCount").Return(int64(1), int64(0), nil).Once()

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registryClient, msgclientRepo, "", true)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Org: &pb.Organization{
				Name:        orgName,
				Owner:       ownerUuid.String(),
				Certificate: certificate,
				Country:     country,
				Currency:    currency,
			}})

		// Assert
		assert.NoError(tt, err)
		assert.NotNil(tt, res)
		assert.Equal(tt, orgName, res.Org.Name)
		assert.Equal(tt, ownerUuid.String(), res.Org.Owner)
		assert.Equal(tt, country, res.Org.Country)
		assert.Equal(tt, currency, res.Org.Currency)
		orgRepo.AssertExpectations(tt)
	})

	t.Run("AddOrgWithEmptyName", func(tt *testing.T) {
		// Arrange
		ownerUuid := uuid.NewV4()
		certificate := "ukama_certs"
		msgclientRepo := &cmocks.MsgBusServiceClient{}

		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registryClient := &cmocks.MemberClient{}
		org := &db.Org{
			Owner:       ownerUuid,
			Certificate: certificate,
			Name:        "", // Empty name
		}

		orgRepo.On("Add", org, mock.Anything).Return(nil).Once()
		orchSystem.On("DeployOrg", mock.Anything).Return(&providers.DeployOrgResponse{}, nil)
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()
		orgRepo.On("GetOrgCount").Return(int64(1), int64(0), nil).Once()

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registryClient, msgclientRepo, "", true)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Org: &pb.Organization{
				Name:        "",
				Owner:       ownerUuid.String(),
				Certificate: certificate,
			}})

		assert.NoError(tt, err)
		assert.NotNil(tt, res)
		assert.Equal(tt, "", res.Org.Name)
		assert.Equal(tt, ownerUuid.String(), res.Org.Owner)
		orgRepo.AssertExpectations(tt)
	})

	t.Run("AddOrgWithDatabaseRecordNotFoundError", func(tt *testing.T) {
		// Arrange
		orgName := "org-db-notfound"
		ownerUuid := uuid.NewV4()
		certificate := "ukama_certs"
		msgclientRepo := &cmocks.MsgBusServiceClient{}

		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registryClient := &cmocks.MemberClient{}
		org := &db.Org{
			Owner:       ownerUuid,
			Certificate: certificate,
			Name:        orgName,
		}

		orgRepo.On("Add", org, mock.Anything).Return(gorm.ErrRecordNotFound).Once()

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registryClient, msgclientRepo, "", true)

		// Act
		orgResp, err := s.Add(context.TODO(), &pb.AddRequest{
			Org: &pb.Organization{
				Name:        orgName,
				Owner:       ownerUuid.String(),
				Certificate: certificate,
			}})

		// Assert
		assert.Error(tt, err)
		assert.Nil(tt, orgResp)
		orgRepo.AssertExpectations(tt)
	})

	t.Run("AddOrgWithDatabaseError", func(tt *testing.T) {
		// Arrange
		orgName := "org-db-error"
		ownerUuid := uuid.NewV4()
		certificate := "ukama_certs"
		msgclientRepo := &cmocks.MsgBusServiceClient{}

		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registryClient := &cmocks.MemberClient{}
		org := &db.Org{
			Owner:       ownerUuid,
			Certificate: certificate,
			Name:        orgName,
		}

		orgRepo.On("Add", org, mock.Anything).Return(errors.New("database connection error")).Once()

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registryClient, msgclientRepo, "", true)

		// Act
		orgResp, err := s.Add(context.TODO(), &pb.AddRequest{
			Org: &pb.Organization{
				Name:        orgName,
				Owner:       ownerUuid.String(),
				Certificate: certificate,
			}})

		// Assert
		assert.Error(tt, err)
		assert.Nil(tt, orgResp)
		orgRepo.AssertExpectations(tt)
	})

	t.Run("AddOrgWithMessageBusPublishFailure", func(tt *testing.T) {
		// Arrange
		orgName := "org-msg-fail"
		ownerUuid := uuid.NewV4()
		certificate := "ukama_certs"
		msgclientRepo := &cmocks.MsgBusServiceClient{}

		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registryClient := &cmocks.MemberClient{}
		org := &db.Org{
			Owner:       ownerUuid,
			Certificate: certificate,
			Name:        orgName,
		}

		orgRepo.On("Add", org, mock.Anything).Return(nil).Once()
		orchSystem.On("DeployOrg", mock.Anything).Return(&providers.DeployOrgResponse{}, nil)
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(errors.New("message bus error")).Once()
		orgRepo.On("GetOrgCount").Return(int64(1), int64(0), nil).Once()

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registryClient, msgclientRepo, "", true)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Org: &pb.Organization{
				Name:        orgName,
				Owner:       ownerUuid.String(),
				Certificate: certificate,
			}})

		// Assert - Should still succeed even if message bus fails
		assert.NoError(tt, err)
		assert.NotNil(tt, res)
		assert.Equal(tt, orgName, res.Org.Name)
		assert.Equal(tt, ownerUuid.String(), res.Org.Owner)
		orgRepo.AssertExpectations(tt)
		msgclientRepo.AssertExpectations(tt)
	})

	t.Run("AddOrgWithMetricPushFailure", func(tt *testing.T) {
		// Arrange
		orgName := "org-metric-fail"
		ownerUuid := uuid.NewV4()
		certificate := "ukama_certs"
		msgclientRepo := &cmocks.MsgBusServiceClient{}

		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registryClient := &cmocks.MemberClient{}
		org := &db.Org{
			Owner:       ownerUuid,
			Certificate: certificate,
			Name:        orgName,
		}

		orgRepo.On("Add", org, mock.Anything).Return(nil).Once()
		orchSystem.On("DeployOrg", mock.Anything).Return(&providers.DeployOrgResponse{}, nil)
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()
		orgRepo.On("GetOrgCount").Return(int64(0), int64(0), errors.New("metric error")).Once()

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registryClient, msgclientRepo, "", true)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Org: &pb.Organization{
				Name:        orgName,
				Owner:       ownerUuid.String(),
				Certificate: certificate,
			}})

		// Assert - Should still succeed even if metric push fails
		assert.NoError(tt, err)
		assert.NotNil(tt, res)
		assert.Equal(tt, orgName, res.Org.Name)
		assert.Equal(tt, ownerUuid.String(), res.Org.Owner)
		orgRepo.AssertExpectations(tt)
		msgclientRepo.AssertExpectations(tt)
	})

	t.Run("AddOrgWithOrchestratorError", func(tt *testing.T) {
		// Arrange
		orgName := "org-orch-error"
		ownerUuid := uuid.NewV4()
		certificate := "ukama_certs"
		msgclientRepo := &cmocks.MsgBusServiceClient{}

		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registryClient := &cmocks.MemberClient{}
		org := &db.Org{
			Owner:       ownerUuid,
			Certificate: certificate,
			Name:        orgName,
		}

		orgRepo.On("Add", org, mock.Anything).Return(errors.New("orchestrator error")).Once()

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registryClient, msgclientRepo, "", true)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Org: &pb.Organization{
				Name:        orgName,
				Owner:       ownerUuid.String(),
				Certificate: certificate,
			}})

		// Assert
		assert.Error(tt, err)
		assert.Nil(tt, res)
		assert.Contains(tt, err.Error(), "orchestrator error")
		orgRepo.AssertExpectations(tt)
	})

	t.Run("AddOrgWithInvalidUUID", func(tt *testing.T) {
		// Arrange
		orgName := "org-invalid-uuid"
		msgclientRepo := &cmocks.MsgBusServiceClient{}

		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registryClient := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registryClient, msgclientRepo, "", true)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Org: &pb.Organization{
				Name:        orgName,
				Owner:       "invalid-uuid",
				Certificate: "ukama_certs",
			}})

		// Assert
		assert.Error(tt, err)
		assert.Nil(tt, res)
		assert.Contains(tt, err.Error(), "invalid format of owner uuid")
	})

}

func TestOrgServer_Get(t *testing.T) {
	orgId := uuid.NewV4()
	msgclientRepo := &cmocks.MsgBusServiceClient{}

	orgRepo := &mocks.OrgRepo{}
	userRepo := &mocks.UserRepo{}
	orchSystem := &mocks.OrchestratorProvider{}
	registry := &cmocks.MemberClient{}

	s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

	t.Run("OrgFound", func(tt *testing.T) {
		orgRepo.On("Get", mock.Anything).Return(&db.Org{Id: orgId}, nil).Once()

		// Act
		orgResp, err := s.Get(context.TODO(), &pb.GetRequest{Id: orgId.String()})

		assert.NoError(t, err)
		assert.Equal(t, orgId.String(), orgResp.GetOrg().GetId())
		orgRepo.AssertExpectations(t)
	})

	t.Run("OrgNotFound", func(tt *testing.T) {
		orgRepo.On("Get", mock.Anything).Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		orgResp, err := s.Get(context.TODO(), &pb.GetRequest{Id: orgId.String()})

		assert.Error(t, err)
		assert.Nil(t, orgResp)
		orgRepo.AssertExpectations(t)
	})

	t.Run("InvalidUUID", func(tt *testing.T) {
		// Act
		orgResp, err := s.Get(context.TODO(), &pb.GetRequest{Id: "invalid-uuid"})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, orgResp)
		assert.Contains(t, err.Error(), "invalid format of org uuid")
		orgRepo.AssertExpectations(t)
	})

	t.Run("DatabaseError", func(tt *testing.T) {
		orgRepo.On("Get", mock.Anything).Return(nil, errors.New("database error")).Once()

		// Act
		orgResp, err := s.Get(context.TODO(), &pb.GetRequest{Id: orgId.String()})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, orgResp)
		assert.Contains(t, err.Error(), "database error")
		orgRepo.AssertExpectations(t)
	})

}

func TestOrgServer_GetByName(t *testing.T) {
	orgName := "test-org"
	msgclientRepo := &cmocks.MsgBusServiceClient{}

	orgRepo := &mocks.OrgRepo{}
	userRepo := &mocks.UserRepo{}
	orchSystem := &mocks.OrchestratorProvider{}
	registry := &cmocks.MemberClient{}

	s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

	t.Run("OrgFound", func(tt *testing.T) {
		orgRepo.On("GetByName", mock.Anything).Return(&db.Org{Name: orgName}, nil).Once()

		// Act
		orgResp, err := s.GetByName(context.TODO(), &pb.GetByNameRequest{Name: orgName})

		assert.NoError(t, err)
		assert.Equal(t, orgName, orgResp.GetOrg().GetName())
		orgRepo.AssertExpectations(t)
	})

	t.Run("OrgNotFound", func(tt *testing.T) {
		orgRepo.On("GetByName", mock.Anything).Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		orgResp, err := s.GetByName(context.TODO(), &pb.GetByNameRequest{Name: orgName})

		assert.Error(t, err)
		assert.Nil(t, orgResp)
		orgRepo.AssertExpectations(t)
	})
}

func TestOrgServer_GetByOwner(t *testing.T) {
	ownerId := uuid.NewV4()
	msgclientRepo := &cmocks.MsgBusServiceClient{}

	orgRepo := &mocks.OrgRepo{}
	userRepo := &mocks.UserRepo{}
	orchSystem := &mocks.OrchestratorProvider{}
	registry := &cmocks.MemberClient{}

	s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

	t.Run("OwnerFound", func(tt *testing.T) {
		orgRepo.On("GetByOwner", mock.Anything).
			Return([]db.Org{{Id: ownerId}}, nil).Once()

		// Act
		orgResp, err := s.GetByOwner(context.TODO(), &pb.GetByOwnerRequest{
			UserUuid: ownerId.String()})

		assert.NoError(t, err)
		assert.Equal(t, ownerId.String(), orgResp.GetOwner())
		orgRepo.AssertExpectations(t)
	})

	t.Run("OwnerNotFound", func(tt *testing.T) {
		orgRepo.On("GetByOwner", mock.Anything).Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		orgResp, err := s.GetByOwner(context.TODO(), &pb.GetByOwnerRequest{
			UserUuid: ownerId.String()})

		assert.Error(t, err)
		assert.Nil(t, orgResp)
		orgRepo.AssertExpectations(t)
	})

	t.Run("InvalidUUID", func(tt *testing.T) {
		// Act
		orgResp, err := s.GetByOwner(context.TODO(), &pb.GetByOwnerRequest{
			UserUuid: "invalid-uuid"})

		assert.Error(t, err)
		assert.Nil(t, orgResp)
		assert.Contains(t, err.Error(), "invalid format of owner uuid")
		orgRepo.AssertExpectations(t)
	})

	t.Run("DatabaseError", func(tt *testing.T) {
		orgRepo.On("GetByOwner", mock.Anything).Return(nil, errors.New("database error")).Once()

		// Act
		orgResp, err := s.GetByOwner(context.TODO(), &pb.GetByOwnerRequest{
			UserUuid: ownerId.String()})

		assert.Error(t, err)
		assert.Nil(t, orgResp)
		assert.Contains(t, err.Error(), "database error")
		orgRepo.AssertExpectations(t)
	})
}

func TestOrgServer_GetByUser(t *testing.T) {
	ownerId := uuid.NewV4()
	orgId := uuid.NewV4()
	userId := uuid.NewV4()
	var id uint = 1

	msgclientRepo := &cmocks.MsgBusServiceClient{}

	orgRepo := &mocks.OrgRepo{}
	userRepo := &mocks.UserRepo{}
	orchSystem := &mocks.OrchestratorProvider{}
	registry := &cmocks.MemberClient{}

	s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

	t.Run("UserFoundOnOwnersAndMembers", func(tt *testing.T) {
		userRepo.On("Get", userId).Return(&db.User{Id: 1, Uuid: userId}, nil).Once()

		orgRepo.On("GetByOwner", userId).
			Return([]db.Org{{Id: orgId, Owner: ownerId}}, nil).Once()

		orgRepo.On("GetByMember", id).
			Return([]db.Org{{Id: orgId}}, nil).Once()

		// Act
		orgResp, err := s.GetByUser(context.TODO(), &pb.GetByOwnerRequest{
			UserUuid: userId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, orgResp)
		assert.Equal(t, 1, len(orgResp.OwnerOf))
		assert.Equal(t, 1, len(orgResp.MemberOf))
		assert.Equal(t, userId.String(), orgResp.GetUser())
		orgRepo.AssertExpectations(t)
	})

	t.Run("UserNotFoundOnOwners", func(tt *testing.T) {
		userRepo.On("Get", userId).Return(&db.User{Id: 1, Uuid: userId}, nil).Once()

		orgRepo.On("GetByOwner", userId).
			Return(nil, gorm.ErrRecordNotFound).Once()

		orgRepo.On("GetByMember", id).
			Return([]db.Org{{Id: orgId}}, nil).Once()
		// Act
		orgResp, err := s.GetByUser(context.TODO(), &pb.GetByOwnerRequest{
			UserUuid: userId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, orgResp)
		assert.Equal(t, 0, len(orgResp.OwnerOf))
		assert.Equal(t, 1, len(orgResp.MemberOf))
		orgRepo.AssertExpectations(t)
	})

	t.Run("UserNotFoundMembers", func(tt *testing.T) {
		userRepo.On("Get", userId).Return(&db.User{Id: 1, Uuid: userId}, nil).Once()

		orgRepo.On("GetByOwner", userId).
			Return([]db.Org{{Id: orgId, Owner: ownerId}}, nil).Once()

		orgRepo.On("GetByMember", id).
			Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		orgResp, err := s.GetByUser(context.TODO(), &pb.GetByOwnerRequest{
			UserUuid: userId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, orgResp)
		assert.Equal(t, 1, len(orgResp.OwnerOf))
		assert.Equal(t, 0, len(orgResp.MemberOf))
		orgRepo.AssertExpectations(t)
	})

	t.Run("InvalidUUID", func(tt *testing.T) {
		// Act
		orgResp, err := s.GetByUser(context.TODO(), &pb.GetByOwnerRequest{
			UserUuid: "invalid-uuid"})

		assert.Error(t, err)
		assert.Nil(t, orgResp)
		assert.Contains(t, err.Error(), "invalid format of user uuid")
		userRepo.AssertExpectations(t)
		orgRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(tt *testing.T) {
		userRepo.On("Get", userId).Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		orgResp, err := s.GetByUser(context.TODO(), &pb.GetByOwnerRequest{
			UserUuid: userId.String()})

		assert.Error(t, err)
		assert.Nil(t, orgResp)
		assert.Contains(t, err.Error(), "user doesn't exist")
		userRepo.AssertExpectations(t)
		orgRepo.AssertExpectations(t)
	})

	t.Run("UserRepoError", func(tt *testing.T) {
		userRepo.On("Get", userId).Return(nil, errors.New("database error")).Once()

		// Act
		orgResp, err := s.GetByUser(context.TODO(), &pb.GetByOwnerRequest{
			UserUuid: userId.String()})

		assert.Error(t, err)
		assert.Nil(t, orgResp)
		assert.Contains(t, err.Error(), "invalid format of user uuid")
		userRepo.AssertExpectations(t)
		orgRepo.AssertExpectations(t)
	})

	t.Run("OwnedOrgsError", func(tt *testing.T) {
		userRepo.On("Get", userId).Return(&db.User{Id: 1, Uuid: userId}, nil).Once()

		orgRepo.On("GetByOwner", userId).
			Return(nil, errors.New("database error")).Once()

		// Act
		orgResp, err := s.GetByUser(context.TODO(), &pb.GetByOwnerRequest{
			UserUuid: userId.String()})

		assert.Error(t, err)
		assert.Nil(t, orgResp)
		assert.Contains(t, err.Error(), "database error")
		userRepo.AssertExpectations(t)
		orgRepo.AssertExpectations(t)
	})

	t.Run("MemberOrgsError", func(tt *testing.T) {
		userRepo.On("Get", userId).Return(&db.User{Id: 1, Uuid: userId}, nil).Once()

		orgRepo.On("GetByOwner", userId).
			Return([]db.Org{{Id: orgId, Owner: ownerId}}, nil).Once()

		orgRepo.On("GetByMember", id).
			Return(nil, errors.New("database error")).Once()

		// Act
		orgResp, err := s.GetByUser(context.TODO(), &pb.GetByOwnerRequest{
			UserUuid: userId.String()})

		assert.Error(t, err)
		assert.Nil(t, orgResp)
		assert.Contains(t, err.Error(), "database error")
		userRepo.AssertExpectations(t)
		orgRepo.AssertExpectations(t)
	})

	t.Run("BothOrgsNotFound", func(tt *testing.T) {
		userRepo.On("Get", userId).Return(&db.User{Id: 1, Uuid: userId}, nil).Once()

		orgRepo.On("GetByOwner", userId).
			Return(nil, gorm.ErrRecordNotFound).Once()

		orgRepo.On("GetByMember", id).
			Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		orgResp, err := s.GetByUser(context.TODO(), &pb.GetByOwnerRequest{
			UserUuid: userId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, orgResp)
		assert.Equal(t, 0, len(orgResp.OwnerOf))
		assert.Equal(t, 0, len(orgResp.MemberOf))
		userRepo.AssertExpectations(t)
		orgRepo.AssertExpectations(t)
	})

}

func TestOrgServer_UpdateUser(t *testing.T) {
	userUUID := uuid.NewV4()
	msgclientRepo := &cmocks.MsgBusServiceClient{}

	orgRepo := &mocks.OrgRepo{}
	userRepo := &mocks.UserRepo{}
	orchSystem := &mocks.OrchestratorProvider{}
	registry := &cmocks.MemberClient{}

	s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

	t.Run("UpdateUserSuccess", func(tt *testing.T) {
		// Arrange
		updatedUser := &db.User{
			Uuid:        userUUID,
			Deactivated: true,
		}

		userRepo.On("Update", mock.MatchedBy(func(user *db.User) bool {
			return user.Uuid == userUUID && user.Deactivated == true
		})).Return(updatedUser, nil).Once()

		// Act
		resp, err := s.UpdateUser(context.TODO(), &pb.UpdateUserRequest{
			UserUuid: userUUID.String(),
			Attributes: &pb.UserAttributes{
				IsDeactivated: true,
			},
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.User)
		assert.Equal(t, userUUID.String(), resp.User.Uuid)
		assert.Equal(t, true, resp.User.IsDeactivated)
		userRepo.AssertExpectations(t)
	})

	t.Run("UpdateUserInvalidUUID", func(tt *testing.T) {
		// Act
		resp, err := s.UpdateUser(context.TODO(), &pb.UpdateUserRequest{
			UserUuid: "invalid-uuid",
			Attributes: &pb.UserAttributes{
				IsDeactivated: false,
			},
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid format of owner uuid")
	})

	t.Run("UpdateUserDatabaseError", func(tt *testing.T) {
		// Arrange
		userRepo.On("Update", mock.Anything).Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		resp, err := s.UpdateUser(context.TODO(), &pb.UpdateUserRequest{
			UserUuid: userUUID.String(),
			Attributes: &pb.UserAttributes{
				IsDeactivated: false,
			},
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		userRepo.AssertExpectations(t)
	})
}

func TestOrgServer_RegisterUser(t *testing.T) {
	userUUID := uuid.NewV4()
	orgUUID := uuid.NewV4()

	t.Run("RegisterUserSuccess", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		org := &db.Org{
			Id:   orgUUID,
			Name: OrgName,
		}

		orgRepo.On("GetByName", OrgName).Return(org, nil).Once()
		userRepo.On("Get", userUUID).Return(nil, gorm.ErrRecordNotFound).Once()
		userRepo.On("Add", mock.MatchedBy(func(user *db.User) bool {
			return user.Uuid == userUUID
		}), mock.Anything).Return(nil).Once()
		msgclientRepo.On("PublishRequest", mock.Anything, mock.MatchedBy(func(e *epb.EventOrgRegisterUser) bool {
			return e.OrgId == orgUUID.String() && e.UserId == userUUID.String()
		})).Return(nil).Once()
		userRepo.On("GetUserCount").Return(int64(1), int64(0), nil).Once()

		// Act
		resp, err := s.RegisterUser(context.TODO(), &pb.RegisterUserRequest{
			UserUuid: userUUID.String(),
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		orgRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("RegisterUserInvalidUUID", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		org := &db.Org{
			Id:   orgUUID,
			Name: OrgName,
		}

		orgRepo.On("GetByName", OrgName).Return(org, nil).Once()

		// Act
		resp, err := s.RegisterUser(context.TODO(), &pb.RegisterUserRequest{
			UserUuid: "invalid-uuid",
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid format of user uuid")
		orgRepo.AssertExpectations(t)
	})

	t.Run("RegisterUserOrgNotFound", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		orgRepo.On("GetByName", OrgName).Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		resp, err := s.RegisterUser(context.TODO(), &pb.RegisterUserRequest{
			UserUuid: userUUID.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		orgRepo.AssertExpectations(t)
	})

	t.Run("RegisterUserAlreadyExists", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		org := &db.Org{
			Id:   orgUUID,
			Name: OrgName,
		}
		existingUser := &db.User{
			Uuid: userUUID,
		}

		orgRepo.On("GetByName", OrgName).Return(org, nil).Once()
		userRepo.On("Get", userUUID).Return(existingUser, nil).Once()

		// Act
		resp, err := s.RegisterUser(context.TODO(), &pb.RegisterUserRequest{
			UserUuid: userUUID.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "user already exist")
		orgRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	t.Run("RegisterUserRegistryError", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		org := &db.Org{
			Id:   orgUUID,
			Name: OrgName,
		}

		orgRepo.On("GetByName", OrgName).Return(org, nil).Once()
		userRepo.On("Get", userUUID).Return(nil, gorm.ErrRecordNotFound).Once()
		userRepo.On("Add", mock.Anything, mock.Anything).Return(errors.New("failed to add user to member registry")).Once()

		// Act
		resp, err := s.RegisterUser(context.TODO(), &pb.RegisterUserRequest{
			UserUuid: userUUID.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to add user to member registry")
		orgRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	t.Run("RegisterUserMessageBusError", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		org := &db.Org{
			Id:   orgUUID,
			Name: OrgName,
		}

		orgRepo.On("GetByName", OrgName).Return(org, nil).Once()
		userRepo.On("Get", userUUID).Return(nil, gorm.ErrRecordNotFound).Once()
		userRepo.On("Add", mock.Anything, mock.Anything).Return(nil).Once()
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(errors.New("message bus error")).Once()
		userRepo.On("GetUserCount").Return(int64(1), int64(0), nil).Once()

		// Act
		resp, err := s.RegisterUser(context.TODO(), &pb.RegisterUserRequest{
			UserUuid: userUUID.String(),
		})

		// Assert - Should still succeed even if message bus fails
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		orgRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("RegisterUserUserRepoError", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		org := &db.Org{
			Id:   orgUUID,
			Name: OrgName,
		}

		orgRepo.On("GetByName", OrgName).Return(org, nil).Once()
		userRepo.On("Get", userUUID).Return(nil, errors.New("database error")).Once()

		// Act
		resp, err := s.RegisterUser(context.TODO(), &pb.RegisterUserRequest{
			UserUuid: userUUID.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to check user")
		orgRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	t.Run("RegisterUserMetricError", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		org := &db.Org{
			Id:   orgUUID,
			Name: OrgName,
		}

		orgRepo.On("GetByName", OrgName).Return(org, nil).Once()
		userRepo.On("Get", userUUID).Return(nil, gorm.ErrRecordNotFound).Once()
		userRepo.On("Add", mock.MatchedBy(func(user *db.User) bool {
			return user.Uuid == userUUID
		}), mock.Anything).Return(nil).Once()
		msgclientRepo.On("PublishRequest", mock.Anything, mock.MatchedBy(func(e *epb.EventOrgRegisterUser) bool {
			return e.OrgId == orgUUID.String() && e.UserId == userUUID.String()
		})).Return(nil).Once()
		userRepo.On("GetUserCount").Return(int64(0), int64(0), errors.New("metric error")).Once()

		// Act
		resp, err := s.RegisterUser(context.TODO(), &pb.RegisterUserRequest{
			UserUuid: userUUID.String(),
		})

		// Assert - Should still succeed even if metric push fails
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		orgRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})
}

func TestOrgServer_UpdateOrgForUser(t *testing.T) {
	userUUID := uuid.NewV4()
	orgUUID := uuid.NewV4()

	t.Run("UpdateOrgForUserSuccess", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		org := &db.Org{
			Id:   orgUUID,
			Name: "test-org",
		}
		user := &db.User{
			Uuid: userUUID,
		}

		orgRepo.On("Get", orgUUID).Return(org, nil).Once()
		userRepo.On("Get", userUUID).Return(user, nil).Once()
		orgRepo.On("AddUser", org, user).Return(nil).Once()
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()
		userRepo.On("GetUserCount").Return(int64(1), int64(0), nil).Once()

		// Act
		resp, err := s.UpdateOrgForUser(context.TODO(), &pb.UpdateOrgForUserRequest{
			UserId: userUUID.String(),
			OrgId:  orgUUID.String(),
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		orgRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("UpdateOrgForUserInvalidOrgUUID", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		// Act
		resp, err := s.UpdateOrgForUser(context.TODO(), &pb.UpdateOrgForUserRequest{
			UserId: userUUID.String(),
			OrgId:  "invalid-uuid",
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid format of user uuid")
	})

	t.Run("UpdateOrgForUserInvalidUserUUID", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		org := &db.Org{
			Id:   orgUUID,
			Name: "test-org",
		}

		orgRepo.On("Get", orgUUID).Return(org, nil).Once()

		// Act
		resp, err := s.UpdateOrgForUser(context.TODO(), &pb.UpdateOrgForUserRequest{
			UserId: "invalid-uuid",
			OrgId:  orgUUID.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid format of user uuid")
		orgRepo.AssertExpectations(t)
	})

	t.Run("UpdateOrgForUserOrgNotFound", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		orgRepo.On("Get", orgUUID).Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		resp, err := s.UpdateOrgForUser(context.TODO(), &pb.UpdateOrgForUserRequest{
			UserId: userUUID.String(),
			OrgId:  orgUUID.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		orgRepo.AssertExpectations(t)
	})

	t.Run("UpdateOrgForUserUserNotFound", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		org := &db.Org{
			Id:   orgUUID,
			Name: "test-org",
		}

		orgRepo.On("Get", orgUUID).Return(org, nil).Once()
		userRepo.On("Get", userUUID).Return(nil, gorm.ErrRecordNotFound).Once()
		orgRepo.On("AddUser", org, (*db.User)(nil)).Return(errors.New("user not found")).Once()

		// Act
		resp, err := s.UpdateOrgForUser(context.TODO(), &pb.UpdateOrgForUserRequest{
			UserId: userUUID.String(),
			OrgId:  orgUUID.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "user not found")
		orgRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	t.Run("UpdateOrgForUserDatabaseError", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		org := &db.Org{
			Id:   orgUUID,
			Name: "test-org",
		}
		user := &db.User{
			Uuid: userUUID,
		}

		orgRepo.On("Get", orgUUID).Return(org, nil).Once()
		userRepo.On("Get", userUUID).Return(user, nil).Once()
		orgRepo.On("AddUser", org, user).Return(errors.New("database error")).Once()

		// Act
		resp, err := s.UpdateOrgForUser(context.TODO(), &pb.UpdateOrgForUserRequest{
			UserId: userUUID.String(),
			OrgId:  orgUUID.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "database error")
		orgRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	t.Run("UpdateOrgForUserMessageBusError", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		org := &db.Org{
			Id:   orgUUID,
			Name: "test-org",
		}
		user := &db.User{
			Uuid: userUUID,
		}

		orgRepo.On("Get", orgUUID).Return(org, nil).Once()
		userRepo.On("Get", userUUID).Return(user, nil).Once()
		orgRepo.On("AddUser", org, user).Return(nil).Once()
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(errors.New("message bus error")).Once()
		userRepo.On("GetUserCount").Return(int64(1), int64(0), nil).Once()

		// Act
		resp, err := s.UpdateOrgForUser(context.TODO(), &pb.UpdateOrgForUserRequest{
			UserId: userUUID.String(),
			OrgId:  orgUUID.String(),
		})

		// Assert - Should still succeed even if message bus fails
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		orgRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("UpdateOrgForUserUserRepoError", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		org := &db.Org{
			Id:   orgUUID,
			Name: "test-org",
		}

		orgRepo.On("Get", orgUUID).Return(org, nil).Once()
		userRepo.On("Get", userUUID).Return(nil, errors.New("database error")).Once()

		// Act
		resp, err := s.UpdateOrgForUser(context.TODO(), &pb.UpdateOrgForUserRequest{
			UserId: userUUID.String(),
			OrgId:  orgUUID.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to check user")
		orgRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	t.Run("UpdateOrgForUserMetricError", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		org := &db.Org{
			Id:   orgUUID,
			Name: "test-org",
		}
		user := &db.User{
			Uuid: userUUID,
		}

		orgRepo.On("Get", orgUUID).Return(org, nil).Once()
		userRepo.On("Get", userUUID).Return(user, nil).Once()
		orgRepo.On("AddUser", org, user).Return(nil).Once()
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()
		userRepo.On("GetUserCount").Return(int64(0), int64(0), errors.New("metric error")).Once()

		// Act
		resp, err := s.UpdateOrgForUser(context.TODO(), &pb.UpdateOrgForUserRequest{
			UserId: userUUID.String(),
			OrgId:  orgUUID.String(),
		})

		// Assert - Should still succeed even if metric push fails
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		orgRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})
}

func TestOrgServer_RemoveOrgForUser(t *testing.T) {
	userUUID := uuid.NewV4()
	orgUUID := uuid.NewV4()

	t.Run("RemoveOrgForUserSuccess", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		org := &db.Org{
			Id:   orgUUID,
			Name: "test-org",
		}
		user := &db.User{
			Uuid: userUUID,
		}

		orgRepo.On("Get", orgUUID).Return(org, nil).Once()
		userRepo.On("Get", userUUID).Return(user, nil).Once()
		userRepo.On("RemoveOrgFromUser", user, org).Return(nil).Once()
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()
		userRepo.On("GetUserCount").Return(int64(1), int64(0), nil).Once()

		// Act
		resp, err := s.RemoveOrgForUser(context.TODO(), &pb.RemoveOrgForUserRequest{
			UserId: userUUID.String(),
			OrgId:  orgUUID.String(),
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		orgRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("RemoveOrgForUserInvalidOrgUUID", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		// Act
		resp, err := s.RemoveOrgForUser(context.TODO(), &pb.RemoveOrgForUserRequest{
			UserId: userUUID.String(),
			OrgId:  "invalid-uuid",
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid format of user uuid")
	})

	t.Run("RemoveOrgForUserInvalidUserUUID", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		org := &db.Org{
			Id:   orgUUID,
			Name: "test-org",
		}

		orgRepo.On("Get", orgUUID).Return(org, nil).Once()

		// Act
		resp, err := s.RemoveOrgForUser(context.TODO(), &pb.RemoveOrgForUserRequest{
			UserId: "invalid-uuid",
			OrgId:  orgUUID.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid format of user uuid")
		orgRepo.AssertExpectations(t)
	})

	t.Run("RemoveOrgForUserOrgNotFound", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		orgRepo.On("Get", orgUUID).Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		resp, err := s.RemoveOrgForUser(context.TODO(), &pb.RemoveOrgForUserRequest{
			UserId: userUUID.String(),
			OrgId:  orgUUID.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		orgRepo.AssertExpectations(t)
	})

	t.Run("RemoveOrgForUserUserNotFound", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		org := &db.Org{
			Id:   orgUUID,
			Name: "test-org",
		}

		orgRepo.On("Get", orgUUID).Return(org, nil).Once()
		userRepo.On("Get", userUUID).Return(nil, gorm.ErrRecordNotFound).Once()
		userRepo.On("RemoveOrgFromUser", (*db.User)(nil), org).Return(errors.New("user not found")).Once()

		// Act
		resp, err := s.RemoveOrgForUser(context.TODO(), &pb.RemoveOrgForUserRequest{
			UserId: userUUID.String(),
			OrgId:  orgUUID.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "user not found")
		orgRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	t.Run("RemoveOrgForUserDatabaseError", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		org := &db.Org{
			Id:   orgUUID,
			Name: "test-org",
		}
		user := &db.User{
			Uuid: userUUID,
		}

		orgRepo.On("Get", orgUUID).Return(org, nil).Once()
		userRepo.On("Get", userUUID).Return(user, nil).Once()
		userRepo.On("RemoveOrgFromUser", user, org).Return(errors.New("database error")).Once()

		// Act
		resp, err := s.RemoveOrgForUser(context.TODO(), &pb.RemoveOrgForUserRequest{
			UserId: userUUID.String(),
			OrgId:  orgUUID.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "database error")
		orgRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	t.Run("RemoveOrgForUserMessageBusError", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		org := &db.Org{
			Id:   orgUUID,
			Name: "test-org",
		}
		user := &db.User{
			Uuid: userUUID,
		}

		orgRepo.On("Get", orgUUID).Return(org, nil).Once()
		userRepo.On("Get", userUUID).Return(user, nil).Once()
		userRepo.On("RemoveOrgFromUser", user, org).Return(nil).Once()
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(errors.New("message bus error")).Once()
		userRepo.On("GetUserCount").Return(int64(1), int64(0), nil).Once()

		// Act
		resp, err := s.RemoveOrgForUser(context.TODO(), &pb.RemoveOrgForUserRequest{
			UserId: userUUID.String(),
			OrgId:  orgUUID.String(),
		})

		// Assert - Should still succeed even if message bus fails
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		orgRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("RemoveOrgForUserUserRepoError", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		org := &db.Org{
			Id:   orgUUID,
			Name: "test-org",
		}

		orgRepo.On("Get", orgUUID).Return(org, nil).Once()
		userRepo.On("Get", userUUID).Return(nil, errors.New("database error")).Once()

		// Act
		resp, err := s.RemoveOrgForUser(context.TODO(), &pb.RemoveOrgForUserRequest{
			UserId: userUUID.String(),
			OrgId:  orgUUID.String(),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to check user")
		orgRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	t.Run("RemoveOrgForUserMetricError", func(tt *testing.T) {
		// Arrange
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		orgRepo := &mocks.OrgRepo{}
		userRepo := &mocks.UserRepo{}
		orchSystem := &mocks.OrchestratorProvider{}
		registry := &cmocks.MemberClient{}

		s := NewOrgServer(OrgName, orgRepo, userRepo, orchSystem, registry, msgclientRepo, "", true)

		org := &db.Org{
			Id:   orgUUID,
			Name: "test-org",
		}
		user := &db.User{
			Uuid: userUUID,
		}

		orgRepo.On("Get", orgUUID).Return(org, nil).Once()
		userRepo.On("Get", userUUID).Return(user, nil).Once()
		userRepo.On("RemoveOrgFromUser", user, org).Return(nil).Once()
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()
		userRepo.On("GetUserCount").Return(int64(0), int64(0), errors.New("metric error")).Once()

		// Act
		resp, err := s.RemoveOrgForUser(context.TODO(), &pb.RemoveOrgForUserRequest{
			UserId: userUUID.String(),
			OrgId:  orgUUID.String(),
		})

		// Assert - Should still succeed even if metric push fails
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		orgRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})
}
