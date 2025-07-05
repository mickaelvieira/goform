package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/mickaelvieira/goform"
)

func main() {
	http.HandleFunc("/", showForm)
	http.HandleFunc("/upload", handleUpload)

	fmt.Println("Server starting on http://localhost:9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func showForm(w http.ResponseWriter, r *http.Request) {
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

	// Create a form with different file input types
	form := goform.Form().
		SetAttributes(
			goform.Attr("action", "/upload"),
			goform.Attr("method", "POST"),
			goform.Attr("enctype", "multipart/form-data"),
		).
		AddChildren(
			// Single file input
			goform.File("avatar").
				SetLabel("Profile Picture").
				SetAttributes(goform.Attr("accept", "image/*")),

			// Multiple file input using SetAttributes
			goform.File("documents").
				SetLabel("Documents (Multiple)").
				SetAttributes(
					goform.Attr("multiple", true),
					goform.Attr("accept", ".pdf,.doc,.docx"),
				),

			// Another way to create multiple file input using SetAttributes
			goform.File("photos").
				SetLabel("Photos (Multiple)").
				SetAttributes(
					goform.Attr("multiple", true),
					goform.Attr("accept", "image/*"),
				),

			// Submit button
			goform.Submit("submit").SetAttributes(goform.Attr("value", "Upload Files")),
		)

	t.ExecuteTemplate(w, "form.tmpl", struct {
		Form goform.Renderer
	}{
		Form: form,
	})
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Parse templates
	t := template.Must(template.ParseFiles("results.tmpl"))

	// Create and populate the form from the request
	form := goform.Form().
		AddChildren(
			goform.File("avatar"),
			goform.File("documents").SetAttributes(goform.Attr("multiple", true)),
			goform.File("photos").SetAttributes(goform.Attr("multiple", true)),
		)

	form.PopulateFromRequest(r)

	// Get the elements to see what files were uploaded
	elements := form.Elements()

	t.ExecuteTemplate(w, "results.tmpl", struct {
		Avatar    string
		Documents string
		Photos    string
	}{
		Avatar:    getFileValue(elements["avatar"]),
		Documents: getFileValue(elements["documents"]),
		Photos:    getFileValue(elements["photos"]),
	})
}

func getFileValue(element goform.Element) string {
	if element == nil {
		return "No element found"
	}

	value := element.Value()
	if value == "" {
		return "No files uploaded"
	}
	return value
}
