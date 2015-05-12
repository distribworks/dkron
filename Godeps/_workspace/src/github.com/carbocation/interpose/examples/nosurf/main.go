package main

import (
	"html/template"
	"net/http"
	"time"

	"github.com/carbocation/interpose"
	"github.com/carbocation/interpose/middleware"
	"github.com/gorilla/mux"
	"github.com/justinas/nosurf"
	"github.com/stretchr/graceful"
)

var templ = template.Must(template.New("t1").Parse(`
<!doctype html>
<html>
<body>
{{ if .name }}
<p>Your name: {{ .name }}</p>
{{ end }}
<form action="/" method="POST">
<input type="text" name="name">

<!-- Try removing this or changing its value
     and see what happens -->
<input type="hidden" name="csrf_token" value="{{ .token }}">
<input type="submit" value="Send">
</form>
</body>
</html>
`))

func main() {
	mw := interpose.New()

	// Invoke NoSurf (it modifies headers so must be called before your router)
	mw.Use(middleware.Nosurf())

	// Create and apply the router
	router := mux.NewRouter()
	mw.UseHandler(router)

	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		context := make(map[string]string)
		context["token"] = nosurf.Token(req)
		if req.Method == "POST" {
			context["name"] = req.FormValue("name")
		}

		templ.Execute(w, context)
	})

	// Launch and permit graceful shutdown, allowing up to 10 seconds for existing
	// connections to end
	graceful.Run(":3001", 10*time.Second, mw)
}
