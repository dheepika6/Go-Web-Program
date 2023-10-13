package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)
	app.logger.Error(err.Error(), "method", method, "URI", uri, "trace", trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}

func (app *application) clientError(w http.ResponseWriter, status int) {

	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) serveTemplate(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	parsedTemplate, ok := app.templates[page]

	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	err := parsedTemplate.ExecuteTemplate(buf, "base", data)

	if err != nil {
		app.serverError(w, r, err)
	}

	w.WriteHeader(status)

	buf.WriteTo(w)

}

func (app *application) newTemplateData() templateData {
	return templateData{}
}
