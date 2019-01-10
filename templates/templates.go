// +build dev

package templates

import "net/http"

//go:generate vfsgendev -source="github.com/victorcoder/dkron/templates".Templates
// Templates contains project templates.
var Templates http.FileSystem = http.Dir(".")
