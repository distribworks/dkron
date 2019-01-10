// +build dev

package static

import "net/http"

//go:generate vfsgendev -source="github.com/victorcoder/dkron/static".Assets
// Templates contains project templates.
var Assets http.FileSystem = http.Dir(".")
