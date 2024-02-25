package types

import (
	fmt "fmt"

	"github.com/golang/protobuf/ptypes"
)

// Key computes the execution key
func (e *Execution) Key() string {
	sa, _ := ptypes.Timestamp(e.StartedAt)
	return fmt.Sprintf("%d-%s", sa.UnixNano(), e.NodeName)
}
