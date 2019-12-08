// +build dev

package assets

import "net/http"

// Assets contains project assets.
var Assets http.FileSystem = http.Dir("../../static")
