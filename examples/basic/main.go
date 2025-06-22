package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/mickaelvieira/goform"
)

func main() {
	t := template.Must(
		template.
			New("templates").
			Funcs(
				template.FuncMap{
					"form": goform.FormRenderer(),
				},
			).
			ParseFiles("login.tmpl"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		form := goform.Form().
			SetAttributes(
				goform.Attr("action", "/login"),
				goform.Attr("id", "login-form"),
			).
			AddChildren(
				goform.FieldSet(
					"Login",
					goform.Text("username").
						SetLabel("Username").
						SetError("The username is required").
						SetAttributes(
							goform.Attr("id", "username"),
							goform.Attr("placeholder", "Enter your username"),
							goform.Attr("required", true),
							goform.Attr("autofocus", true),
						),
					goform.Password("password").
						SetLabel("Password").
						SetHint("Password must be at least 12 characters long").
						SetError("The password is required").
						SetAttributes(
							goform.Attr("id", "password"),
							goform.Attr("required", true),
						),
				),
				goform.Group(
					goform.Reset("reset").
						SetAttributes(
							goform.Attr("value", "Reset"),
						),
					goform.Submit("submit").
						SetAttributes(
							goform.Attr("value", "Login"),
						),
				),
			)

		t.ExecuteTemplate(w, "login.tmpl", struct {
			Form goform.Renderer
		}{
			Form: form,
		})
	})

	fmt.Println("Starting server on http://localhost:9000")
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal(err)
	}
}
