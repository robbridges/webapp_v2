package views

import (
	"fmt"
	"github.com/gorilla/csrf"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type Template struct {
	HtmlTpl *template.Template
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, pattern ...string) (Template, error) {
	tpl := template.New(pattern[0])

	tpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return `<input type="hidden" />`
		},
	})
	tpl, err := tpl.ParseFS(fs, pattern...)
	if err != nil {
		return Template{}, fmt.Errorf("parse FS: %w", err)
	}

	return Template{
		HtmlTpl: tpl,
	}, nil

}

//func Parse(filepath string) (Template, error) {
//	tpl, err := template.ParseFiles(filepath)
//	if err != nil {
//		return Template{}, fmt.Errorf("parsing template: %w", err)
//	}
//	return Template{
//		HtmlTpl: tpl,
//	}, nil
//}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {
	tpl, err := t.HtmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error rendering the page", http.StatusInternalServerError)
		return
	}
	tpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
	})

	w.Header().Set("Content-Type", "text/html")
	err = tpl.Execute(w, data)
	if err != nil {
		log.Printf("Executing template %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
}
