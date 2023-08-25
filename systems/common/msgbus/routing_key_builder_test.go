package msgbus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {

	t.Run("basic_usage", func(t *testing.T) {
		rk, err := NewRoutingKeyBuilder().SetEventType().SetCloudSource().SetObject("some-obj").
			SetAction("create").SetCloudSource().SetGlobalScope().SetOrgName("org").SetSystem("system").SetService("some_container").Build()
		assert.NoError(t, err)
		assert.Equal(t, "event.cloud.global.org.system.some_container.some-obj.create", rk)
	})

	t.Run("use_star_segment", func(t *testing.T) {
		rk, err := NewRoutingKeyBuilder().SetEventType().SetCloudSource().SetObject("some-obj").
			SetAction("*").SetCloudSource().SetOrgName("org").SetSystem("system").SetService("some_container").Build()
		assert.NoError(t, err)
		assert.Equal(t, "event.cloud.local.org.system.some_container.some-obj.*", rk)
	})

	t.Run("error_missing_segment", func(t *testing.T) {
		_, err := NewRoutingKeyBuilder().SetEventType().SetCloudSource().
			SetAction("*").SetGlobalScope().SetOrgName("org").SetSystem("system").SetService("some_container").Build()
		assert.Error(t, err, "")
		assert.EqualErrorf(t, err, "object segment is not set", "")

	})

	t.Run("make_sure_new_instace_is_created", func(t *testing.T) {
		rk := NewRoutingKeyBuilder()
		rk1 := rk.SetEventType().SetCloudSource().
			SetAction("*").SetCloudSource().SetGlobalScope().SetOrgName("org").SetSystem("system").SetService("container1")

		rk2 := rk.SetEventType().SetCloudSource().
			SetAction("*").SetCloudSource().SetGlobalScope().SetOrgName("org").SetSystem("system").SetService("container2")

		assert.NotEqual(t, rk1, rk2)
	})
}
