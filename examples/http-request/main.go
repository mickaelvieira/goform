package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/mickaelvieira/goform"
)

func main() {
	http.HandleFunc("/", handleForm)
	http.HandleFunc("/success", handleSuccess)

	fmt.Println("Server starting on :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	// Parse templates
	t := template.Must(
		template.
			New("templates").
			Funcs(
				template.FuncMap{
					"form": goform.FormRenderer(),
				},
			).
			ParseFiles("form.tmpl"))

	// Create a form
	form := goform.Form().
		AddChildren(
			goform.Text("name").SetLabel("Name").
				SetAttributes(goform.Attr("required", true)),
			goform.Email("email").SetLabel("Email").
				SetAttributes(goform.Attr("required", true)),
			goform.Number("age").SetLabel("Age"),
			goform.Select("country").SetLabel("Country").
				SetOptions(
					goform.Option("United States", "us"),
					goform.Option("Canada", "ca"),
					goform.Option("United Kingdom", "uk"),
				),
			goform.Submit("submit").
				SetAttributes(goform.Attr("value", "Submit")),
		)

	// If this is a POST request, populate from the request
	if r.Method == http.MethodPost {
		// Populate form with request data
		form.PopulateFromRequest(r)

		// Validate the form
		isValid, errors := form.IsValid()

		if isValid {
			// Form is valid, redirect to success page
			http.Redirect(w, r, "/success", http.StatusSeeOther)
			return
		} else {
			// Form has errors, display them
			form.SetError("Please correct the errors below")
			for field, err := range errors {
				fmt.Printf("Field %s has error: %s\n", field, err)
			}
		}
	}

	t.ExecuteTemplate(w, "form.tmpl", struct {
		Form goform.Renderer
	}{
		Form: form,
	})
}

func handleSuccess(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("success.tmpl"))
	t.ExecuteTemplate(w, "success.tmpl", nil)
}
