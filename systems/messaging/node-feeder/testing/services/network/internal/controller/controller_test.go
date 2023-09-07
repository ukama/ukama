package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/testing/services/network/internal"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func fakeNewController() *Controller {
	client := fake.NewSimpleClientset()
	return &Controller{
		cs: client,
		ns: "default",
	}
}

func Test_PowerOnNode(t *testing.T) {
	internal.ServiceConfig = &internal.Config{
		NodeImage: "",
		NodeCmd:   []string{},
	}

	c := fakeNewController()
	node := ukama.NewVirtualHomeNodeId()
	err := c.PowerOnNode(node.String(), "ukama")
	assert.Nil(t, err)
}

func Test_PowerOnNodeFail(t *testing.T) {
	internal.ServiceConfig = &internal.Config{
		NodeImage: "",
		NodeCmd:   []string{},
	}

	c := fakeNewController()
	node := "abcd"
	err := c.PowerOnNode(node, "ukama")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "not expected nodeid format")
	}
}

func Test_PowerOffNode(t *testing.T) {
	c := fakeNewController()

	internal.ServiceConfig = &internal.Config{
		NodeImage: "",
		NodeCmd:   []string{},
	}

	node := ukama.NewVirtualHomeNodeId()
	err := c.PowerOnNode(node.String(), "ukama")
	assert.Nil(t, err)

	err = c.PowerOffNode(node.String())
	assert.Nil(t, err)
}

func Test_PowerOffNodeDoesNotExist(t *testing.T) {
	c := fakeNewController()
	node := ukama.NewVirtualHomeNodeId()
	err := c.PowerOffNode(node.String())
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "not found")
	}

}

func TestGetNodeRuntimeStatus(t *testing.T) {

	testVector := map[v1.PodPhase]string{
		v1.PodPending:   VNodeBooting,
		v1.PodRunning:   VNodeActive,
		v1.PodSucceeded: VNodeHalted,
		v1.PodFailed:    VNodeFaulty,
		v1.PodUnknown:   VNodeUnkown,
	}

	for k, v := range testVector {
		str := getNodeRuntimeStatus(k)
		assert.Equal(t, str, v)
	}
}
