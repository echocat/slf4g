package location

import (
	"strconv"

	log "github.com/echocat/slf4g"
)

// DepthOnly is a location which expresses only the depth where a log event was
// captured.
type DepthOnly uint16

// String implements fmt.Stringer
func (instance DepthOnly) String() string {
	return strconv.Itoa(int(instance))
}

// NewDepthOnlyDiscovery provides a Discovery to create DepthOnly instances.
func NewDepthOnlyDiscovery() Discovery {
	return depthOnlyDiscoveryV
}

var depthOnlyDiscoveryV = &depthOnlyDiscovery{}

type depthOnlyDiscovery struct{}

func (instance *depthOnlyDiscovery) DiscoverLocation(_ log.Event, skipFrames uint16) Location {
	return DepthOnly(skipFrames)
}
