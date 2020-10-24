package location

import (
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_NewDepthOnlyDiscovery(t *testing.T) {
	actual := NewDepthOnlyDiscovery()

	assert.ToBeSame(t, depthOnlyDiscoveryV, actual)
}

func Test_depthOnlyDiscovery_DiscoverLocation(t *testing.T) {
	instance := NewDepthOnlyDiscovery()

	actual := instance.DiscoverLocation(nil, 666)

	assert.ToBeEqual(t, DepthOnly(666), actual)
}

func Test_DepthOnly_String(t *testing.T) {
	instance := DepthOnly(666)

	actual := instance.String()

	assert.ToBeEqual(t, "666", actual)
}
