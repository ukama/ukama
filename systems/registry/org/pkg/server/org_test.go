package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/org/mocks"
	pb "github.com/ukama/ukama/systems/registry/org/pb/gen"
	"github.com/ukama/ukama/systems/registry/org/pkg/db"
)

const testOrgName = "test-org"

func TestOrgServer_AddOrg(t *testing.T) {
	// Arrange
	orgName := "org-1"
	ownerUUID := uuid.NewV4()
	certificate := "ukama_certs"
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	orgRepo := &mocks.OrgRepo{}
	userRepo := &mocks.UserRepo{}

	org := &db.Org{
		Owner:       ownerUUID,
		Certificate: certificate,
		Name:        orgName,
	}

	orgRepo.On("Add", org, mock.Anything).Return(nil).Once()

	userRepo.On("Get", ownerUUID).Return(&db.User{
		Id:   1,
		Uuid: ownerUUID,
	}, nil).Once()

	msgclientRepo.On("PublishRequest", mock.Anything, &pb.AddRequest{Org: &pb.Organization{
		Name:        orgName,
		Owner:       ownerUUID.String(),
		Certificate: certificate,
	}}).Return(nil).Once()
	orgRepo.On("GetOrgCount").Return(int64(1), int64(0), nil).Once()
	orgRepo.On("GetMemberCount", mock.Anything).Return(int64(1), int64(0), nil).Once()
	userRepo.On("GetUserCount").Return(int64(1), int64(0), nil).Once()

	s := NewOrgServer(orgRepo, userRepo, "", msgclientRepo, "")

	// Act
	res, err := s.Add(context.TODO(), &pb.AddRequest{Org: &pb.Organization{
		Name:        orgName,
		Owner:       ownerUUID.String(),
		Certificate: certificate,
	}})

	// Assert
	msgclientRepo.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, orgName, res.Org.Name)
	assert.Equal(t, ownerUUID.String(), res.Org.Owner)
	orgRepo.AssertExpectations(t)
}

func TestOrgServer_GetOrg(t *testing.T) {
	orgID := uuid.NewV4()
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	orgRepo := &mocks.OrgRepo{}

	orgRepo.On("Get", mock.Anything).Return(&db.Org{Id: orgID}, nil).Once()

	s := NewOrgServer(orgRepo, nil, "", msgclientRepo, "")
	orgResp, err := s.Get(context.TODO(), &pb.GetRequest{Id: orgID.String()})

	assert.NoError(t, err)
	assert.Equal(t, orgID.String(), orgResp.GetOrg().GetId())
	orgRepo.AssertExpectations(t)
}

func TestOrgServer_AddOrg_fails_without_owner_uuid(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	s := NewOrgServer(orgRepo, nil, "", msgclientRepo, "")
	_, err := s.Add(context.TODO(), &pb.AddRequest{
		Org: &pb.Organization{Name: testOrgName},
	})

	assert.Error(t, err)
}

func TestOrgServer_AddOrg_fails_with_bad_owner_id(t *testing.T) {
	owner := "org-1"
	orgRepo := &mocks.OrgRepo{}
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	s := NewOrgServer(orgRepo, nil, "", msgclientRepo, "")
	_, err := s.Add(context.TODO(), &pb.AddRequest{
		Org: &pb.Organization{Owner: owner},
	})

	assert.Error(t, err)
}
