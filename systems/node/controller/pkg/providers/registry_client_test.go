package providers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	//import the mock
	"github.com/ukama/ukama/systems/node/controller/mocks"
)

func TestValidateSite(t *testing.T) {
	mockRegistry := mocks.NewRegistryProvider(t)

	mockRegistry.On("ValidateSite", "network_id", "site_name", "org_name").Return(nil)
	mockRegistry.On("ValidateSite", "", "", "").Return(errors.New("invalid arguments"))

	err := mockRegistry.ValidateSite("network_id", "site_name", "org_name")

	assert.NoError(t, err)

	err = mockRegistry.ValidateSite("", "", "")

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid arguments")

	mockRegistry.AssertCalled(t, "ValidateSite", "network_id", "site_name", "org_name")
}

func TestValidateNetwork(t *testing.T) {
	mockRegistry := mocks.NewRegistryProvider(t)

	mockRegistry.On("ValidateNetwork", "network_id", "org_name").Return(nil)
	mockRegistry.On("ValidateNetwork", "", "").Return(errors.New("invalid arguments"))

	err := mockRegistry.ValidateNetwork("network_id", "org_name")

	assert.NoError(t, err)

	err = mockRegistry.ValidateNetwork("", "")

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid arguments")

	mockRegistry.AssertCalled(t, "ValidateNetwork", "network_id", "org_name")
}

func TestGetNodesBySite(t *testing.T) {
	mockRegistry := mocks.NewRegistryProvider(t)

	mockRegistry.On("GetNodesBySite", "site_id").Return([]string{"node_id"}, nil)
	mockRegistry.On("GetNodesBySite", "").Return(nil, errors.New("invalid arguments"))

	nodes, err := mockRegistry.GetNodesBySite("site_id")

	assert.NoError(t, err)
	assert.Equal(t, []string{"node_id"}, nodes)

	nodes, err = mockRegistry.GetNodesBySite("")

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid arguments")
	assert.Nil(t, nodes)

	mockRegistry.AssertCalled(t, "GetNodesBySite", "site_id")
}