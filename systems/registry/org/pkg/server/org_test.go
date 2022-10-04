package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	bstmock "github.com/ukama/ukama/services/bootstrap/client/mocks"
	"github.com/ukama/ukama/services/cloud/org/mocks"
	pb "github.com/ukama/ukama/services/cloud/org/pb/gen"
	"github.com/ukama/ukama/services/cloud/org/pkg/db"
	"github.com/ukama/ukama/services/common/msgbus/stub"
	"testing"
)

var testDeviceGatewayHost = "1.1.1.1"

const testOrgName = "test-org"

func TestOrgServer_AddOrg(t *testing.T) {
	// Arrange
	orgName := "org-1"
	ownerId := uuid.NewString()
	orgRepo := &mocks.OrgRepo{}
	pub := &stub.QPubStub{}

	// trick to call nested bootstrap call
	orgRepo.On("Add", mock.Anything, mock.Anything).
		Return(func(o *db.Org, f ...func() error) error {
			return f[0]()
		}).Once()

	bootstrapClient := &bstmock.Client{}
	bootstrapClient.On("AddOrUpdateOrg", orgName, mock.Anything, testDeviceGatewayHost).Return(nil)

	s := NewOrgServer(orgRepo, bootstrapClient, testDeviceGatewayHost, pub)

	// Act
	res, err := s.Add(context.TODO(), &pb.AddRequest{Org: &pb.Organization{
		Name: orgName, Owner: ownerId,
	}})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, orgName, res.Org.Name)
	assert.Equal(t, ownerId, res.Org.Owner)
	orgRepo.AssertExpectations(t)
	bootstrapClient.AssertExpectations(t)
}

func TestNetworkServer_GetOrg(t *testing.T) {
	orgName := "org-1"

	orgRepo := &mocks.OrgRepo{}
	pub := &stub.QPubStub{}
	orgRepo.On("GetByName", mock.Anything).Return(&db.Org{Name: orgName}, nil).Once()

	s := NewOrgServer(orgRepo, &bstmock.Client{}, testDeviceGatewayHost, pub)
	org, err := s.Get(context.TODO(), &pb.GetRequest{Name: orgName})
	assert.NoError(t, err)
	assert.Equal(t, orgName, org.GetOrg().GetName())
	orgRepo.AssertExpectations(t)
}

func TestNetworkServer_AddOrg_fails_without_owner_id(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	pub := &stub.QPubStub{}
	s := NewOrgServer(orgRepo, &bstmock.Client{}, testDeviceGatewayHost, pub)
	_, err := s.Add(context.TODO(), &pb.AddRequest{
		Org: &pb.Organization{Name: testOrgName},
	})
	assert.Error(t, err)
}

func TestNetworkServer_AddOrg_fails_with_bad_owner_id(t *testing.T) {
	orgName := "org-1"
	orgRepo := &mocks.OrgRepo{}
	pub := &stub.QPubStub{}
	s := NewOrgServer(orgRepo, &bstmock.Client{}, testDeviceGatewayHost, pub)
	_, err := s.Add(context.TODO(), &pb.AddRequest{
		Org: &pb.Organization{Name: orgName},
	})
	assert.Error(t, err)
}
