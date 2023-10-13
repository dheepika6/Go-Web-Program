package main

import "github.com/dheepika6/LetsGoWebProgram/internal/models"

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
	Form     any
}
