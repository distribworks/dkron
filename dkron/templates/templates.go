// +build dev

package templates

import "net/http"

// Templates contains project templates.
var Templates http.FileSystem = http.Dir("../../templates")
