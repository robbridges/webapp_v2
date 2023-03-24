package controllers

import (
	"github.com/robbridges/webapp_v2/views"
	"html/template"
	"net/http"
)

func StaticHandler(tpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}

func FAQ(tpl views.Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "Is there a free trial?",
			Answer:   "Yes we offer a 30 day free trial for all sign-ups.",
		},
		{
			Question: "What are your support hours?",
			Answer: "We have support staff answering emails 24/7, response times mayb e a bit slower on weekends and " +
				"holidays",
		},
		{
			Question: "How do I contact support?",
			Answer:   `Email us - <a href="mailto:support@lenslocked.com">support@lenslocked.com</a>`,
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, questions)
	}
}
