package integration

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestTopLevelTestForSuite(t *testing.T) {
	// Run all tests in suite
	suite.Run(t, &IntegrationTestSuite{})
}
