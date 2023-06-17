package views

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/robbridges/webapp_v2/context"
	"github.com/robbridges/webapp_v2/models"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"
)

type Template struct {
	HtmlTpl *template.Template
}

type public interface {
	Public() string
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, pattern ...string) (Template, error) {
	tpl := template.New(path.Base(pattern[0]))

	tpl.Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", fmt.Errorf("csrf function not implmented")
		},
		"currentUser": func() (template.HTML, error) {
			return "", fmt.Errorf("current user not implemented")
		},
		"errors": func() []string {
			return nil
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

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	// clone the template instead of using the same one every time to prevent data races
	tpl, err := t.HtmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error rendering the page", http.StatusInternalServerError)
		return
	}
	errMsgs := errMessages(errs...)
	tpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
		"currentUser": func() *models.User {
			return context.User(r.Context())
		},
		"errors": func() []string {

			return errMsgs
		},
	})

	w.Header().Set("Content-Type", "text/html")
	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		log.Printf("Executing template %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func errMessages(errs ...error) []string {
	var messages []string
	//temp this just needs to work for the time being
	for _, err := range errs {
		var pubErr public
		if errors.As(err, &pubErr) {
			messages = append(messages, pubErr.Public())
		} else {
			fmt.Println(err)
			messages = append(messages, "Something went wrong")
		}
	}
	return messages
}
