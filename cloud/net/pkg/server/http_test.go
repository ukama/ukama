package server

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_marshallTargets(t *testing.T) {
	output, err := marshallTargets(map[string]string{"a": "1", "b": "2"}, 10250)
	assert.NoError(t, err)

	actual := strings.ReplaceAll(string(output), " ", "")
	actual = strings.Replace(actual, `{"targets":["1:10250"],"labels":{"nodeid":"a"}}`, "", 1)
	actual = strings.Replace(actual, `{"targets":["2:10250"],"labels":{"nodeid":"b"}}`, "", 1)

	assert.Equal(t, "[,]", actual)
}
