package types

import (
	fmt "fmt"
)

// Key computes the execution key
func (e *Execution) Key() string {
	sa := e.StartedAt.AsTime()
	return fmt.Sprintf("%d-%s", sa.UnixNano(), e.NodeName)
}
