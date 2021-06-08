package dkron

import (
	"strconv"
	"strings"

	"github.com/hashicorp/serf/serf"
	"github.com/sirupsen/logrus"
)

// cleanTags takes the tag spec and returns strictly key:value pairs
// along with the lowest cardinality specified
func cleanTags(tags map[string]string, logger *logrus.Entry) (map[string]string, int) {
	cardinality := int(^uint(0) >> 1) // MaxInt

	cleanTags := make(map[string]string, len(tags))

	for k, v := range tags {
		vparts := strings.Split(v, ":")

		cleanTags[k] = vparts[0]

		// If a cardinality is specified (i.e. "value:3") and it is lower than our
		// max cardinality, lower the max
		if len(vparts) == 2 {
			tagCard, err := strconv.Atoi(vparts[1])
			if err != nil {
				// Tag value is malformed
				tagCard = 0
				logger.Errorf("improper cardinality specified for tag %s: %v", k, vparts[1])
			}

			if tagCard < cardinality {
				cardinality = tagCard
			}
		}
	}

	return cleanTags, cardinality
}

// nodeMatchesTags tests if a node matches all of the provided tags
func nodeMatchesTags(node serf.Member, tags map[string]string) bool {
	for k, v := range tags {
		nodeVal, present := node.Tags[k]
		if !present {
			return false
		}
		if nodeVal != v {
			return false
		}
	}
	// If we matched all key:value pairs, the node matches the tags
	return true
}
