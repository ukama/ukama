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

}
