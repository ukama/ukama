package multipl_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/services/cloud/device-feeder/mocks"
	"github.com/ukama/ukama/services/cloud/device-feeder/pkg"
	"github.com/ukama/ukama/services/cloud/device-feeder/pkg/multipl"
	pb "github.com/ukama/ukama/services/cloud/registry/pb/gen"
)

func Test_requestMultiplier_Process(t *testing.T) {
	registry := mocks.RegistryClient{}
	registry.On("GetNodesList", "test-org").Return([]*pb.Node{
		{
			NodeId: "node-1",
		},
		{
			NodeId: "node-2",
		},
	}, nil)

	t.Run("PublishAllMessages", func(t *testing.T) {
		pub := mocks.QueuePublisher{}
		pub.On("Publish", mock.Anything).Return(nil).Twice()

		m := multipl.NewRequestMultiplier(&registry, &pub)

		err := m.Process(&pkg.DevicesUpdateRequest{
			HttpMethod: "POST",
			Target:     "test-org.*",
			Body:       `{ "body": "test" }`,
			Path:       "/devices/update",
		})

		assert.NoError(t, err)
		pub.AssertExpectations(t)
	})

	t.Run("FailsToPublishMessage", func(t *testing.T) {
		pub := mocks.QueuePublisher{}
		pub.On("Publish", mock.MatchedBy(func(m pkg.DevicesUpdateRequest) bool {
			return strings.HasSuffix(m.Target, "node-1")
		})).Return(nil).Once()

		pub.On("Publish", mock.MatchedBy(func(m pkg.DevicesUpdateRequest) bool {
			return strings.HasSuffix(m.Target, "node-2")
		})).Return(fmt.Errorf("error publishing the message")).Once()

		m := multipl.NewRequestMultiplier(&registry, &pub)

		err := m.Process(&pkg.DevicesUpdateRequest{
			HttpMethod: "POST",
			Target:     "test-org.*",
			Body:       `{ "body": "test" }`,
			Path:       "/devices/update",
		})

		assert.Error(t, err)
		pub.AssertExpectations(t)
	})

}
