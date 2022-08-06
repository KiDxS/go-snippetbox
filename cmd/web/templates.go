package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/KiDxS/snippetbox/pkg/forms"
	"github.com/KiDxS/snippetbox/pkg/models"
	"github.com/justinas/nosurf"
)

type templateData struct {
	CSRFToken         string
	CurrentYear       int
	AuthenticatedUser int
	Snippet           *models.Snippet
	Snippets          []*models.Snippet
	Form              *forms.Form
	FormData          url.Values
	FormErrors        map[string]string
	Flash             string
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}
	// Use the filepath.Glob function to get a slice of all filepaths with
	// the extension '.page.tmpl'. This essentially gives us a slice of all the
	// 'page' templates for the application.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}
	// Loop through the pages one-by-one.
	for _, page := range pages {
		// Extract the file name (like 'home.page.tmpl') from the full file pat
		// and assign it to the name variable.
		name := filepath.Base(page)
		// Parse the page template file in to a template set.
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// Use the ParseGlob method to add any 'layout' templates to the
		// template set (in our case, it's just the 'base' layout at the
		// moment).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}
		// Use the ParseGlob method to add any 'partial' templates to the
		// template set (in our case, it's just the 'footer' partial at the
		// moment).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {

	if td == nil {
		td = &templateData{}
	}
	td.CSRFToken = nosurf.Token(r)
	td.AuthenticatedUser = app.authenticatedUser(r)
	td.Flash = app.session.PopString(r, "flash")
	td.CurrentYear = time.Now().Year()
	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {

	ts, ok := app.templateCache[name]

	if !ok {
		app.serveError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serveError(w, err)
		return
	}

	buf.WriteTo(w)
}
