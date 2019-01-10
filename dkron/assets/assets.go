// +build dev

package assets

import "net/http"

// Templates contains project templates.
var Assets http.FileSystem = http.Dir("../../static")
