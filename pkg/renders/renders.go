package renders

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/ptamarov/go-cards/pkg/config"
	"github.com/ptamarov/go-cards/pkg/models"
)

var app *config.AppConfig

func NewRenders(a *config.AppConfig) { // Set the app variable.
	app = a
}

func AddDefaultTemplateData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders a template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		var err error
		tc, err = CreateTemplateCache()
		if err != nil {
			log.Fatal("renders; could not create template cache")
		}
	}
	// create a template cache

	// get requested template from cache
	t, ok := tc[tmpl]

	if !ok {
		log.Fatal("missing template", tmpl)
	}

	buf := new(bytes.Buffer)

	td = AddDefaultTemplateData(td, r)

	err := t.Execute(buf, td) // need to pass TemplateData to the Execute function!
	if err != nil {
		log.Println("renders; I got the following error in buffer:", err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("renders; I got the following error when writing:", err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{}
	// get all of the files named *.page.tmpl from the ./templates folder
	pages, err := filepath.Glob("./cmd/templates/*.page.tmpl") // Glob returns all files matching expression
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page) // Base gets the last element of the path
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		// now get all layouts
		matches, err := filepath.Glob("./cmd/templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./cmd/templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}

	return myCache, nil
}
