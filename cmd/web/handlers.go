package main

import (
	"errors"
	"fmt"
	// "html/template"
	"net/http"
	"strconv"

	"github.com/dheepika6/LetsGoWebProgram/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, r, err)
	}

	for _, s := range snippets {
		fmt.Fprintf(w, "%+v\n", s)
	}
	// files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/pages/home.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// }

	// ts, err := template.ParseFiles(files...)

	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	// err = ts.ExecuteTemplate(w, "base", nil)

	// if err != nil {
	// 	app.serverError(w, r, err)
	// }
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		app.notFound(w)
	}

	snippet, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {

			app.serverError(w, r, err)
		}
		return
	}

	fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "A new title"
	content := "O Snail\nClimb Mount\n Work hard\n ypu can do it"
	expires := 7

	latestId, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.logger.Error("Error creating snippet", "title", title, "content", content, "expires", expires)
		app.serverError(w, r, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", latestId), http.StatusSeeOther)
}
