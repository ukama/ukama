package server

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ukama/ukama/systems/registry/org/mocks"
	pb "github.com/ukama/ukama/systems/registry/org/pb/gen"
	"github.com/ukama/ukama/systems/registry/org/pkg/db"
)

const testOrgName = "test-org"

func TestOrgServer_AddOrg(t *testing.T) {
	// Arrange
	orgName := "org-1"
	ownerUUID := uuid.New()
	certificate := "ukama_certs"

	orgRepo := &mocks.OrgRepo{}
	userRepo := &mocks.UserRepo{}

	org := &db.Org{
		ID:          0,
		Owner:       ownerUUID,
		Certificate: certificate,
		Name:        orgName,
	}

	orgRepo.On("Add", org, mock.Anything).Return(nil).Once()

	userRepo.On("Get", ownerUUID).Return(&db.User{
		ID:   1,
		Uuid: ownerUUID,
	}, nil).Once()

	// orgRepo.On("AddMember", &db.OrgUser{
	// OrgID:  org.ID,
	// UserID: 1,
	// Uuid:   ownerUUID,
	// }).Return(nil).Once()

	s := NewOrgServer(orgRepo, userRepo, nil)

	// Act
	res, err := s.Add(context.TODO(), &pb.AddRequest{Org: &pb.Organization{
		Name: orgName, Owner: ownerUUID.String(), Certificate: certificate}})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, orgName, res.Org.Name)
	assert.Equal(t, ownerUUID.String(), res.Org.Owner)
	orgRepo.AssertExpectations(t)
}

func TestOrgServer_GetOrg(t *testing.T) {
	orgID := uint64(0)

	orgRepo := &mocks.OrgRepo{}

	orgRepo.On("Get", mock.Anything).Return(&db.Org{ID: uint(orgID)}, nil).Once()

	s := NewOrgServer(orgRepo, nil, nil)
	orgResp, err := s.Get(context.TODO(), &pb.GetRequest{Id: orgID})

	assert.NoError(t, err)
	assert.Equal(t, orgID, orgResp.GetOrg().GetId())
	orgRepo.AssertExpectations(t)
}

func TestOrgServer_AddOrg_fails_without_owner_uuid(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}

	s := NewOrgServer(orgRepo, nil, nil)
	_, err := s.Add(context.TODO(), &pb.AddRequest{
		Org: &pb.Organization{Name: testOrgName},
	})

	assert.Error(t, err)
}

func TestOrgServer_AddOrg_fails_with_bad_owner_id(t *testing.T) {
	owner := "org-1"
	orgRepo := &mocks.OrgRepo{}

	s := NewOrgServer(orgRepo, nil, nil)
	_, err := s.Add(context.TODO(), &pb.AddRequest{
		Org: &pb.Organization{Owner: owner},
	})

	assert.Error(t, err)
}
