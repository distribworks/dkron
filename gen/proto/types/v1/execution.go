package typesv1

import (
	"fmt"
)

// Key computes the execution key
func (e *Execution) Key() string {
	sa := e.StartedAt.AsTime()
	return fmt.Sprintf("%d-%s", sa.UnixNano(), e.NodeName)
}
