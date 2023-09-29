package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/software-manager/mocks"
	pb "github.com/ukama/ukama/systems/node/software-manager/pb/gen"
	"github.com/ukama/ukama/systems/node/software-manager/pkg/db"
)

const testOrgName = "test-org"

var orgId = uuid.NewV4()

func TestMemberServer_AddMember(t *testing.T) {
	// Arrange
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	mRepo := &mocks.SoftwareManagerRepo{}
	nOrg := &mocks.NucleusClientProvider{}

	sw := db.Software{
		Id: uuid.NewV4(),
		NodeId: uuid.NewV4(),
		Tag: "0.12",

		
	}

	mRepo.On("AddMember", mock.Anything, orgId.String(), mock.Anything).Return(nil).Once()
	// mOrg.On("GetUserById", member.UserId.String()).Return(&providers.UserInfo{
	// 	Id: member.UserId.String(),
	// }, nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, &pb.AddMemberRequest{
		UserUuid: member.UserId.String(),
		Role:     pb.RoleType(db.Users),
	}).Return(nil).Once()
	mRepo.On("GetMemberCount").Return(int64(1), int64(1), nil).Once()
	s := NewMemberServer(testOrgName, mRepo, nOrg, msgclientRepo, "", orgId)

	// Act
	_, err := s.AddMember(context.TODO(), &pb.AddMemberRequest{
		UserUuid: member.UserId.String(),
		Role:     pb.RoleType(db.Users),
	})

	// Assert
	msgclientRepo.AssertExpectations(t)
	assert.NoError(t, err)

	mRepo.AssertExpectations(t)
}