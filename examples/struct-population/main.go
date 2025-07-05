package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/mickaelvieira/goform"
)

// UserRegistration represents the form data structure
type UserRegistration struct {
	Name      string   `goform:"name"`
	Email     string   `goform:"email"`
	Bio       string   `goform:"bio"`
	Avatar    string   `goform:"avatar"`
	Documents []string `goform:"documents"`
}

func main() {
	http.HandleFunc("/", handleForm)
	fmt.Println("Server starting on http://localhost:9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	// Parse templates
	formTemplate := template.Must(
		template.
			New("templates").
			Funcs(
				template.FuncMap{
					"form": goform.FormRenderer(),
				},
			).
			ParseFiles("form.tmpl"))

	// Create the form structure
	form := goform.Form().
		SetAttributes(
			goform.Attr("action", "/"),
			goform.Attr("method", "POST"),
			goform.Attr("enctype", "multipart/form-data"),
		).
		AddChildren(
			goform.Text("name").
				SetLabel("Full Name").
				SetAttributes(goform.Attr("required", true)),
			goform.Email("email").
				SetLabel("Email Address").
				SetAttributes(goform.Attr("required", true)),
			goform.Textarea("bio").
				SetLabel("Biography"),
			goform.File("avatar").
				SetLabel("Profile Picture"),
			goform.File("documents").
				SetLabel("Documents").
				SetAttributes(goform.Attr("multiple", true)),
			goform.Submit("submit").
				SetAttributes(goform.Attr("value", "Register")),
		)

	if r.Method == "POST" {
		// Populate form from HTTP request
		form.PopulateFromRequest(r)

		// Validate the form
		isValid, errors := form.IsValid()

		if isValid {
			// Populate struct from form data
			var user UserRegistration
			form.Populate(&user)

			// Parse success template
			successTemplate := template.Must(template.ParseFiles("success.tmpl"))

			// Display the populated struct
			successTemplate.ExecuteTemplate(w, "success.tmpl", struct {
				User      UserRegistration
				RawStruct string
			}{
				User:      user,
				RawStruct: fmt.Sprintf("%+v", user),
			})
			return
		} else {
			// Handle validation errors - for simplicity, we'll just show them
			form.SetError("Please correct the errors below")
			for field, err := range errors {
				fmt.Printf("Field %s has error: %s\n", field, err)
			}
		}
	}

	// Display the form (GET request or after validation errors)
	formTemplate.ExecuteTemplate(w, "form.tmpl", struct {
		Form goform.Renderer
	}{
		Form: form,
	})
}
