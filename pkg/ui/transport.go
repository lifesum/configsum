package ui

import (
	"html/template"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
)

const tmplIndex = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <base href="{{ .Base }}">
    <link href="https://fonts.googleapis.com/css?family=Roboto+Mono|Roboto:300,400,500,700,900" rel="stylesheet">
    <link href="styles/nucleo.css" rel="stylesheet">
    <link href="styles/normalize.css" rel="stylesheet">
    <link href="styles/console.css" rel="stylesheet">
    <script src="scripts/console.js" type="text/javascript"></script>
  </head>
  <body>
    <script type="text/javascript">
		var app = Elm.Main.fullscreen({"now": (new Date()).getTime()});
    </script>
 </body>
</html>`

// MakeHandler returns am http.Handler for the UI.
func MakeHandler(logger log.Logger, base string, local bool) http.Handler {
	r := mux.NewRouter()

	r.Methods("GET").PathPrefix("/fonts").Name("fonts").Handler(
		http.FileServer(_escFS(local)),
	)

	r.Methods("GET").PathPrefix("/scripts").Name("scripts").Handler(
		http.FileServer(_escFS(local)),
	)

	r.Methods("GET").PathPrefix("/styles").Name("styles").Handler(
		http.FileServer(_escFS(local)),
	)

	tplRoot := template.Must(template.New("root").Parse(tmplIndex))

	r.Methods("GET").PathPrefix("/").Name("root").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_ = tplRoot.Execute(w, struct {
				Base string
			}{
				Base: base,
			})
		},
	)

	return r
}
