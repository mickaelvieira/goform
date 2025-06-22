package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/mickaelvieira/goform"
)

//go:embed form.tmpl
var formFS embed.FS

func main() {
	// Set the embedded filesystem for overriding templates
	// this allows you to use custom templates for form rendering.
	// You can override any template present in the templates directory
	// you just need to make sure that the file name matches that of the original template.
	goform.SetOverridingTemplates(formFS, "form.tmpl")

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
			SetError("Please fill in all required fields.").
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
