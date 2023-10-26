package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	// "strings"
	// "unicode/utf8"

	"github.com/dheepika6/LetsGoWebProgram/internal/models"
	"github.com/dheepika6/LetsGoWebProgram/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, r, err)
	}

	data := templateData{
		Snippets: snippets,
	}
	app.serveTemplate(w, r, http.StatusOK, "home.tmpl", data)

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// id, err := strconv.Atoi(r.URL.Query().Get("id"))

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
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

	data := templateData{
		Snippet: snippet,
	}

	app.serveTemplate(w, r, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.serveTemplate(w, r, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:   title,
		Content: content,
		Expires: expires,
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "The field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This filed cannot be 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "The field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData()
		data.Form = form
		app.serveTemplate(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)

		return
	}

	latestId, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.logger.Error("Error creating snippet", "title", title, "content", content, "expires", expires)
		app.serverError(w, r, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", latestId), http.StatusSeeOther)
}
