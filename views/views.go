package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var (
	LayoutDir   string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

type View struct {
	Template *template.Template
	Layout   string
}

func NewView(layout string, files ...string) *View {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
	for i, f := range files {
		files[i] = f + TemplateExt
	}

	layoutFiles, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}

	files = append(files, layoutFiles...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

// Render builds a template using data
func (v *View) RenderJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := v.Template.ExecuteTemplate(w, v.Layout, data)
	if err != nil {
		fmt.Println(err)
	}
}

// Render is used to render the view with the predefined layout.
func (v *View) RenderHTML(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	if err := v.Template.ExecuteTemplate(w, v.Layout, data); err != nil {
		log.Println(err)
	}
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.RenderHTML(w, r, nil)
}
