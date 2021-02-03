// +build dev

package assets_ui

import (
	"net/http"
)

// Assets contains project assets.
var Assets http.FileSystem = http.Dir("../../ui-dist")
