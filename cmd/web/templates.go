package main

import (
	"html/template"
	"path/filepath"
)

func newTemplateCache() (map[string]*template.Template, error) {
	templateMap := map[string]*template.Template{}

	pagesPath, err := filepath.Glob("./ui/html/pages/*.tmpl")

	if err != nil {
		return nil, err
	}

	for _, page := range pagesPath {
		name := filepath.Base(page)

		ts, err := template.ParseFiles("./ui/html/base.tmpl")

		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")

		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)

		if err != nil {
			return nil, err
		}

		templateMap[name] = ts
	}

	return templateMap, nil
}
