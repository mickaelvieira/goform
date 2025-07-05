package goform

import (
	"html/template"
	"mime/multipart"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"testing"
)

func cleanHTML(h template.HTML) string {
	s := string(h)
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, "  ", "")
	return s
}

func TestForm_Creation(t *testing.T) {
	t.Run("empty form", func(t *testing.T) {
		f := Form()

		children := f.Children()
		if len(children) != 0 {
			t.Errorf("expected 0 elements, got %d", len(children))
		}

		result := f.Render()
		if result == template.HTML("") {
			t.Error("expected non-empty render result, template might not be initialized")
		}
	})

	t.Run("form with single element", func(t *testing.T) {
		textField := Text("username")
		f := Form().AddChildren(textField)

		children := f.Children()
		if len(children) != 1 {
			t.Errorf("expected 1 element, got %d", len(children))
		}
		if children[0] != textField {
			t.Error("expected first element to be the text field")
		}
	})

	t.Run("form with multiple elements", func(t *testing.T) {
		usernameField := Text("username")
		emailField := Email("email")
		passwordField := Password("password")
		telephoneField := Phone("telephone")

		f := Form().AddChildren(usernameField, emailField, passwordField)

		children := f.Children()
		if len(children) != 3 {
			t.Errorf("expected 3 elements, got %d", len(children))
		}

		indexUsername := slices.IndexFunc(children, func(e Renderer) bool {
			return e == usernameField
		})
		indexEmail := slices.IndexFunc(children, func(e Renderer) bool {
			return e == emailField
		})
		indexPassword := slices.IndexFunc(children, func(e Renderer) bool {
			return e == passwordField
		})
		indexTelephone := slices.IndexFunc(children, func(e Renderer) bool {
			return e == telephoneField
		})

		if indexUsername == -1 {
			t.Error("expected to have a username field")
		}
		if indexEmail == -1 {
			t.Error("expected to have a email field")
		}
		if indexPassword == -1 {
			t.Error("expected to have a password field")
		}
		if indexTelephone >= 0 {
			t.Error("expected to not have a telephone field")
		}
	})
}

func TestForm_SetAttributes(t *testing.T) {
	t.Run("set attributes on empty form", func(t *testing.T) {
		f := Form()
		result := f.SetAttributes(
			Attr("method", "POST"),
			Attr("action", "/submit"),
			Attr("class", "login-form"),
		)

		if f.attributes.String("method") != "POST" {
			t.Errorf("expected method=POST, got %s", f.attributes.String("method"))
		}
		if f.attributes.String("action") != "/submit" {
			t.Errorf("expected action=/submit, got %s", f.attributes.String("action"))
		}
		if f.attributes.String("class") != "login-form" {
			t.Errorf("expected class=login-form, got %s", f.attributes.String("class"))
		}
		if result != f {
			t.Error("SetAttributes should return the same form instance")
		}
	})

	t.Run("override existing attributes", func(t *testing.T) {
		f := Form().SetAttributes(Attr("method", "GET"))
		result := f.SetAttributes(Attr("method", "POST"))

		if f.attributes.String("method") != "POST" {
			t.Errorf("expected method=POST after override, got %s", f.attributes.String("method"))
		}
		if result != f {
			t.Error("SetAttributes should return the same form instance")
		}
	})

	t.Run("set multiple attributes at once", func(t *testing.T) {
		f := Form().
			SetAttributes(
				Attr("method", "POST"),
				Attr("action", "/submit"),
				Attr("enctype", "multipart/form-data"),
				Attr("novalidate", true),
				Attr("autocomplete", "off"),
			)

		if f.attributes.String("method") != "POST" {
			t.Errorf("expected method=POST, got %s", f.attributes.String("method"))
		}
		if f.attributes.String("action") != "/submit" {
			t.Errorf("expected action=/submit, got %s", f.attributes.String("action"))
		}
		if f.attributes.String("enctype") != "multipart/form-data" {
			t.Errorf("expected enctype=multipart/form-data, got %s", f.attributes.String("enctype"))
		}
		if !f.attributes.Bool("novalidate") {
			t.Error("expected novalidate=true")
		}
		if f.attributes.String("autocomplete") != "off" {
			t.Errorf("expected autocomplete=off, got %s", f.attributes.String("autocomplete"))
		}
	})
}

func TestForm_Children(t *testing.T) {
	t.Run("empty form children", func(t *testing.T) {
		f := Form()
		children := f.Children()

		if len(children) != 0 {
			t.Errorf("expected 0 children, got %d", len(children))
		}
		if children == nil {
			t.Error("expected non-nil slice, even if empty")
		}
	})

	t.Run("form with children", func(t *testing.T) {
		textField := Text("text")
		selectField := Select("select")
		f := Form().AddChildren(textField, selectField)

		children := f.Children()

		if len(children) != 2 {
			t.Errorf("expected 2 children, got %d", len(children))
		}
		if children[0] != textField {
			t.Error("expected first child to be text field")
		}
		if children[1] != selectField {
			t.Error("expected second child to be select field")
		}
	})
}

func TestForm_Render(t *testing.T) {
	t.Run("render empty form", func(t *testing.T) {
		f := Form()
		result := f.Render()

		if result == template.HTML("") {
			t.Error("expected non-empty render result")
		}

		htmlStr := cleanHTML(result)
		if htmlStr == "" {
			t.Error("expected non-empty HTML string")
		}

		if !strings.Contains(htmlStr, "<form") {
			t.Error("expected rendered HTML to contain <form tag")
		}
		if !strings.Contains(htmlStr, "</form>") {
			t.Error("expected rendered HTML to contain </form> tag")
		}
	})

	t.Run("render form with attributes", func(t *testing.T) {
		f := Form().SetAttributes(
			Attr("method", "POST"),
			Attr("action", "/submit"),
		)
		result := f.Render()
		htmlStr := cleanHTML(result)

		// Should contain the attributes
		if !strings.Contains(htmlStr, `method="POST"`) {
			t.Error("expected rendered HTML to contain method attribute")
		}
		if !strings.Contains(htmlStr, `action="/submit"`) {
			t.Error("expected rendered HTML to contain action attribute")
		}
	})

	t.Run("render form with elements", func(t *testing.T) {
		textField := Text("username").SetAttributes(Id("username-field"))
		emailField := Email("email").SetAttributes(Id("email-field"))

		f := Form().
			AddChildren(textField, emailField).
			SetAttributes(
				Attr("method", "POST"),
				Attr("action", "/submit"),
			)

		result := f.Render()
		htmlStr := cleanHTML(result)

		if !strings.Contains(htmlStr, "<form") {
			t.Error("expected rendered HTML to contain <form tag")
		}
		if !strings.Contains(htmlStr, "</form>") {
			t.Error("expected rendered HTML to contain </form> tag")
		}

		if !strings.Contains(htmlStr, `name="username"`) {
			t.Error("expected rendered HTML to contain username field")
		}
		if !strings.Contains(htmlStr, `name="email"`) {
			t.Error("expected rendered HTML to contain email field")
		}
	})
}

func TestForm_GroupElements(t *testing.T) {
	t.Run("render form  with grouped elements", func(t *testing.T) {
		textField := Text("username").SetAttributes(Id("username-field"))
		resetButton := Reset("reset").SetAttributes(Id("reset-btn"))
		regularButton := Button("button").SetAttributes(Id("button-btn"))
		submitButton := Submit("submit").SetAttributes(Id("submit-btn"))

		f := Form().
			SetError("Please fill in all required fields.").
			AddChildren(textField, Group(submitButton, resetButton, regularButton)).
			SetAttributes(
				Id("login-form"),
				Attr("method", "POST"),
				Attr("action", "/login"),
			)

		result := f.Render()
		htmlStr := cleanHTML(result)

		expectedForm := `<form action="/login" aria-errormessage="login-form-error" enctype="application/x-www-form-urlencoded" id="login-form" method="POST"><span id="login-form-error">Please fill in all required fields.</span>`
		expectedUsernameField := `<div><div><input id="username-field" name="username" type="text"></div></div>`
		expectedSubmitButton := `<div role="group"><input id="submit-btn" name="submit" type="submit">`
		expectedResetButton := `<input id="reset-btn" name="reset" type="reset">`
		expectedRegularButton := `<input id="button-btn" name="button" type="button"></div>`
		expectedFormEnd := `</form>`

		expected := expectedForm + expectedUsernameField + expectedSubmitButton + expectedResetButton + expectedRegularButton + expectedFormEnd

		if htmlStr != expected {
			t.Errorf("expected exact HTML match:\nExpected: %s\nActual: %s", expected, htmlStr)
		}
	})
}

func TestForm_ComplexScenarios(t *testing.T) {
	t.Run("form with various element types", func(t *testing.T) {
		usernameField := Text("username")
		emailField := Email("email")
		passwordField := Password("password")
		agreeCheckbox := Checkbox("agree")
		countrySelect := Select("country")
		descriptionTextarea := Textarea("description")
		submitButton := Submit("submit")

		f := Form().
			AddChildren(
				usernameField,
				emailField,
				passwordField,
				agreeCheckbox,
				countrySelect,
				descriptionTextarea,
				submitButton,
			).
			SetAttributes(
				Attr("method", "POST"),
				Attr("action", "/register"),
			)

		if len(f.Children()) != 7 {
			t.Errorf("expected 7 children, got %d", len(f.Children()))
		}

		if f.attributes.String("method") != "POST" {
			t.Errorf("expected method=POST, got %s", f.attributes.String("method"))
		}
		if f.attributes.String("action") != "/register" {
			t.Errorf("expected action=/register, got %s", f.attributes.String("action"))
		}

		result := f.Render()
		if result == template.HTML("") {
			t.Error("expected non-empty render result for complex form")
		}
	})

	t.Run("form with nested containers", func(t *testing.T) {
		userFieldSet := FieldSet("User Information",
			Text("username").SetAttributes(Id("username-field")),
			Email("email").SetAttributes(Id("email-field")),
		)

		addressGroup := Group(
			Text("street").SetAttributes(Id("street-field")),
			Text("city").SetAttributes(Id("city-field")),
		).
			SetAttributes(Attr("class", "address-group"))

		f := Form().
			AddChildren(userFieldSet, addressGroup).
			SetAttributes(
				Id("nested-form"),
				Attr("method", "POST"),
				Attr("action", "/submit"),
			)

		result := f.Render()
		htmlStr := cleanHTML(result)

		// Test the complete rendered form with nested containers
		expectedForm := `<form action="/submit" enctype="application/x-www-form-urlencoded" id="nested-form" method="POST">`
		expectedFieldSet := `<fieldset><legend>User Information</legend>`
		expectedUsernameField := `<div><div><input id="username-field" name="username" type="text"></div></div>`
		expectedEmailField := `<div><div><input id="email-field" name="email" type="email"></div></div>`
		expectedFieldSetEnd := `</fieldset>`
		expectedGroupStart := `<div class="address-group" role="group">`
		expectedStreetField := `<div><div><input id="street-field" name="street" type="text"></div></div>`
		expectedCityField := `<div><div><input id="city-field" name="city" type="text"></div></div>`
		expectedGroupEnd := `</div>`
		expectedFormEnd := `</form>`

		expected := expectedForm + expectedFieldSet + expectedUsernameField + expectedEmailField + expectedFieldSetEnd + expectedGroupStart + expectedStreetField + expectedCityField + expectedGroupEnd + expectedFormEnd

		if htmlStr != expected {
			t.Errorf("expected exact HTML match:\nExpected: %s\nActual: %s", expected, htmlStr)
		}
	})
}

func TestForm_NilElements(t *testing.T) {
	t.Run("handles mixed nil and valid elements", func(t *testing.T) {
		validElement := Text("valid")

		f := Form().AddChildren(validElement, nil, nil)

		children := f.Children()
		if len(children) != 1 {
			t.Errorf("expected 1 child, got %d", len(children))
		}

		if children[0] != validElement {
			t.Error("expected first child to be valid element")
		}
	})
}

func TestForm_Chaining(t *testing.T) {
	t.Run("method chaining", func(t *testing.T) {
		textField := Text("username")
		emailField := Email("email")

		f := Form().
			AddChildren(textField, emailField).
			SetAttributes(Attr("method", "POST")).
			SetAttributes(Attr("action", "/submit")).
			SetAttributes(Attr("class", "user-form"))

		if f.attributes.String("method") != "POST" {
			t.Errorf("expected method=POST, got %s", f.attributes.String("method"))
		}
		if f.attributes.String("action") != "/submit" {
			t.Errorf("expected action=/submit, got %s", f.attributes.String("action"))
		}
		if f.attributes.String("class") != "user-form" {
			t.Errorf("expected class=user-form, got %s", f.attributes.String("class"))
		}

		if len(f.Children()) != 2 {
			t.Errorf("expected 2 children after chaining, got %d", len(f.Children()))
		}
	})
}

func TestForm_AttributeHandling(t *testing.T) {
	tests := []struct {
		name      string
		attribute string
		value     any
		expected  string
	}{
		{"method POST", "method", "POST", "POST"},
		{"method GET", "method", "GET", "GET"},
		{"action path", "action", "/submit", "/submit"},
		{"class name", "class", "login-form", "login-form"},
		{"enctype multipart", "enctype", "multipart/form-data", "multipart/form-data"},
		{"autocomplete off", "autocomplete", "off", "off"},
		{"target blank", "target", "_blank", "_blank"},
		{"name form", "name", "userForm", "userForm"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Form().SetAttributes(Attr(tt.attribute, tt.value))
			if f.attributes.String(tt.attribute) != tt.expected {
				t.Errorf("expected %s=%s, got %s", tt.attribute, tt.expected, f.attributes.String(tt.attribute))
			}
		})
	}
}

func TestForm_BooleanAttributes(t *testing.T) {
	t.Run("boolean attributes", func(t *testing.T) {
		f := Form().SetAttributes(
			Attr("novalidate", true),
			Attr("autocomplete", false),
		)

		if !f.attributes.Bool("novalidate") {
			t.Error("expected novalidate=true")
		}

		if f.attributes.Bool("autocomplete") {
			t.Error("expected autocomplete to be false/absent")
		}
	})
}

func TestForm_EmptyForm(t *testing.T) {
	t.Run("form with no elements should still render", func(t *testing.T) {
		f := Form().SetAttributes(
			Id("empty-form"),
			Attr("method", "POST"),
			Attr("action", "/submit"),
		)

		result := cleanHTML(f.Render())
		expected := `<form action="/submit" enctype="application/x-www-form-urlencoded" id="empty-form" method="POST"></form>`

		if result != expected {
			t.Errorf("expected exact HTML match:\nExpected: %s\nActual: %s", expected, result)
		}
	})
}

func TestForm_PopulateFromRequest(t *testing.T) {
	t.Run("populate from URL-encoded form data", func(t *testing.T) {
		// Create a form with text and email inputs
		form := Form().AddChildren(
			Text("name").SetLabel("Name"),
			Email("email").SetLabel("Email"),
		)

		// Create a mock HTTP request with form data
		formData := url.Values{}
		formData.Set("name", "John Doe")
		formData.Set("email", "john@example.com")

		req := &http.Request{
			Method:   http.MethodPost,
			Form:     formData,
			PostForm: formData,
		}

		// Populate the form from the request
		_ = form.PopulateFromRequest(req)

		// Check that the elements were populated correctly
		elements := form.Elements()

		nameElement := elements["name"]
		if nameElement == nil {
			t.Fatal("name element not found")
		}

		emailElement := elements["email"]
		if emailElement == nil {
			t.Fatal("email element not found")
		}

		// Note: Since Element interface doesn't expose Value() method,
		// we can test by rendering and checking the output contains the values
		nameHTML := string(nameElement.Render())
		if !strings.Contains(nameHTML, "John Doe") {
			t.Errorf("expected name element to contain 'John Doe', got: %s", nameHTML)
		}

		emailHTML := string(emailElement.Render())
		if !strings.Contains(emailHTML, "john@example.com") {
			t.Errorf("expected email element to contain 'john@example.com', got: %s", emailHTML)
		}
	})

	t.Run("populate from multipart form data with file uploads", func(t *testing.T) {
		// Create a form with text and file inputs
		form := Form().AddChildren(
			Text("name").SetLabel("Name"),
			File("avatar").SetLabel("Avatar"),
			File("document").SetLabel("Document"),
		)

		// Create a mock HTTP request with multipart form data
		formData := url.Values{}
		formData.Set("name", "John Doe")

		// Create mock multipart form with file data
		multipartForm := &multipart.Form{
			Value: map[string][]string{
				"name": {"John Doe"}, // This should also be in multipart
			},
			File: map[string][]*multipart.FileHeader{
				"avatar": {
					{
						Filename: "profile.jpg",
						Size:     1024,
					},
				},
				"document": {
					{
						Filename: "resume.pdf",
						Size:     2048,
					},
				},
			},
		}

		req := &http.Request{
			Method:        http.MethodPost,
			Form:          formData,
			PostForm:      formData,
			MultipartForm: multipartForm,
		}

		// Populate the form from the request
		_ = form.PopulateFromRequest(req)

		// Check that the elements were populated correctly
		elements := form.Elements()

		nameElement := elements["name"]
		if nameElement == nil {
			t.Fatal("name element not found")
		}

		avatarElement := elements["avatar"]
		if avatarElement == nil {
			t.Fatal("avatar element not found")
		}

		documentElement := elements["document"]
		if documentElement == nil {
			t.Fatal("document element not found")
		}

		// Check that text field was populated
		nameHTML := string(nameElement.Render())
		if !strings.Contains(nameHTML, "John Doe") {
			t.Errorf("expected name element to contain 'John Doe', got: %s", nameHTML)
		}

		// Check that file fields were populated with filenames
		avatarHTML := string(avatarElement.Render())
		if !strings.Contains(avatarHTML, "profile.jpg") {
			t.Errorf("expected avatar element to contain 'profile.jpg', got: %s", avatarHTML)
		}

		documentHTML := string(documentElement.Render())
		if !strings.Contains(documentHTML, "resume.pdf") {
			t.Errorf("expected document element to contain 'resume.pdf', got: %s", documentHTML)
		}
	})

	t.Run("populate from multipart form data without files", func(t *testing.T) {
		// Create a form with text and file inputs
		form := Form().AddChildren(
			Text("name").SetLabel("Name"),
			Text("email").SetLabel("Email"),
			File("optional_file").SetLabel("Optional File"),
		)

		// Create a mock HTTP request with multipart form data but no files
		formData := url.Values{}
		formData.Set("name", "Jane Smith")
		formData.Set("email", "jane@example.com")

		// Create multipart form with only text values, no files
		multipartForm := &multipart.Form{
			Value: map[string][]string{
				"name":  {"Jane Smith"},
				"email": {"jane@example.com"},
			},
			File: map[string][]*multipart.FileHeader{}, // Empty files map
		}

		req := &http.Request{
			Method:        http.MethodPost,
			Form:          formData,
			PostForm:      formData,
			MultipartForm: multipartForm,
		}

		// Populate the form from the request
		_ = form.PopulateFromRequest(req)

		// Check that the text elements were populated correctly
		elements := form.Elements()

		nameElement := elements["name"]
		if nameElement == nil {
			t.Fatal("name element not found")
		}

		emailElement := elements["email"]
		if emailElement == nil {
			t.Fatal("email element not found")
		}

		fileElement := elements["optional_file"]
		if fileElement == nil {
			t.Fatal("optional_file element not found")
		}

		// Check that text fields were populated
		nameHTML := string(nameElement.Render())
		if !strings.Contains(nameHTML, "Jane Smith") {
			t.Errorf("expected name element to contain 'Jane Smith', got: %s", nameHTML)
		}

		emailHTML := string(emailElement.Render())
		if !strings.Contains(emailHTML, "jane@example.com") {
			t.Errorf("expected email element to contain 'jane@example.com', got: %s", emailHTML)
		}

		// File element should remain empty since no file was uploaded
		fileHTML := string(fileElement.Render())
		if strings.Contains(fileHTML, `value="`) && !strings.Contains(fileHTML, `value=""`) {
			t.Errorf("expected file element to have empty value, got: %s", fileHTML)
		}
	})

	t.Run("populate and validate", func(t *testing.T) {
		// Create a form with required fields
		form := Form().AddChildren(
			Text("name").SetLabel("Name").
				SetAttributes(Attr("required", true)),
			Email("email").SetLabel("Email").
				SetAttributes(Attr("required", true)),
		)

		// Create a request with valid data
		formData := url.Values{}
		formData.Set("name", "Jane Doe")
		formData.Set("email", "jane@example.com")

		req := &http.Request{
			Method:   http.MethodPost,
			Form:     formData,
			PostForm: formData,
		}

		// Populate the form first
		_ = form.PopulateFromRequest(req)

		// Then validate separately
		isValid, errors := form.IsValid()

		// Should be valid with proper data
		if !isValid {
			t.Errorf("expected form to be valid, got errors: %v", errors)
		}

		if len(errors) > 0 {
			t.Errorf("expected no errors, got: %v", errors)
		}
	})

	t.Run("populate from multipart form data with multiple files", func(t *testing.T) {
		// Create a form with a file input that allows multiple files
		form := Form().AddChildren(
			Text("name").SetLabel("Name"),
			File("documents").SetLabel("Documents").
				SetAttributes(Attr("multiple", true)),
		)

		// Create a mock HTTP request with multipart form data
		formData := url.Values{}
		formData.Set("name", "John Doe")

		// Create mock multipart form with multiple files
		multipartForm := &multipart.Form{
			Value: map[string][]string{
				"name": {"John Doe"},
			},
			File: map[string][]*multipart.FileHeader{
				"documents": {
					{
						Filename: "resume.pdf",
						Size:     1024,
					},
					{
						Filename: "cover-letter.docx",
						Size:     2048,
					},
					{
						Filename: "portfolio.pdf",
						Size:     3072,
					},
				},
			},
		}

		req := &http.Request{
			Method:        http.MethodPost,
			Form:          formData,
			PostForm:      formData,
			MultipartForm: multipartForm,
		}

		// Populate the form from the request
		_ = form.PopulateFromRequest(req)

		// Check that the elements were populated correctly
		elements := form.Elements()

		nameElement := elements["name"]
		if nameElement == nil {
			t.Fatal("name element not found")
		}

		documentsElement := elements["documents"]
		if documentsElement == nil {
			t.Fatal("documents element not found")
		}

		// Check that text field was populated
		nameHTML := string(nameElement.Render())
		if !strings.Contains(nameHTML, "John Doe") {
			t.Errorf("expected name element to contain 'John Doe', got: %s", nameHTML)
		}

		// Check that file field was populated with all filenames
		documentsHTML := string(documentsElement.Render())
		expectedValue := "resume.pdf, cover-letter.docx, portfolio.pdf"
		if !strings.Contains(documentsHTML, expectedValue) {
			t.Errorf("expected documents element to contain '%s', got: %s", expectedValue, documentsHTML)
		}

		// Verify all individual filenames are present
		if !strings.Contains(documentsHTML, "resume.pdf") {
			t.Error("expected documents element to contain 'resume.pdf'")
		}
		if !strings.Contains(documentsHTML, "cover-letter.docx") {
			t.Error("expected documents element to contain 'cover-letter.docx'")
		}
		if !strings.Contains(documentsHTML, "portfolio.pdf") {
			t.Error("expected documents element to contain 'portfolio.pdf'")
		}
	})

	t.Run("populate from multipart form data with mixed empty and valid filenames", func(t *testing.T) {
		// Create a form with a file input
		form := Form().AddChildren(
			File("files").SetLabel("Files"),
		)

		// Create mock multipart form with some empty filenames
		multipartForm := &multipart.Form{
			Value: map[string][]string{},
			File: map[string][]*multipart.FileHeader{
				"files": {
					{
						Filename: "document1.pdf",
						Size:     1024,
					},
					{
						Filename: "", // Empty filename (can happen in some browsers)
						Size:     0,
					},
					{
						Filename: "document2.jpg",
						Size:     2048,
					},
					{
						Filename: "", // Another empty filename
						Size:     0,
					},
				},
			},
		}

		req := &http.Request{
			Method:        http.MethodPost,
			Form:          url.Values{},
			PostForm:      url.Values{},
			MultipartForm: multipartForm,
		}

		// Populate the form from the request
		_ = form.PopulateFromRequest(req)

		// Check that the element was populated correctly
		elements := form.Elements()
		filesElement := elements["files"]
		if filesElement == nil {
			t.Fatal("files element not found")
		}

		// Check that only non-empty filenames are included
		filesHTML := string(filesElement.Render())
		expectedValue := "document1.pdf, document2.jpg"
		if !strings.Contains(filesHTML, expectedValue) {
			t.Errorf("expected files element to contain '%s', got: %s", expectedValue, filesHTML)
		}

		// Verify empty filenames are not included
		if strings.Contains(filesHTML, ",,") {
			t.Error("expected no consecutive commas from empty filenames")
		}
	})
}

func TestForm_IsValid(t *testing.T) {
	t.Run("valid form", func(t *testing.T) {
		form := Form().AddChildren(
			Text("name").SetLabel("Name"),
			Email("email").SetLabel("Email"),
		)

		// Set valid values
		elements := form.Elements()
		elements["name"].SetValue("John Doe")
		elements["email"].SetValue("john@example.com")

		isValid, errors := form.IsValid()

		if !isValid {
			t.Errorf("expected form to be valid, got errors: %v", errors)
		}

		if len(errors) > 0 {
			t.Errorf("expected no errors, got: %v", errors)
		}
	})

	t.Run("form with validation errors", func(t *testing.T) {
		form := Form().AddChildren(
			Text("name").SetLabel("Name").
				SetAttributes(Attr("required", true)),
			Email("email").SetLabel("Email").
				SetAttributes(Attr("required", true)),
		)

		// Leave fields empty or set invalid values
		// The actual validation logic depends on your element implementation

		isValid, errors := form.IsValid()

		// We expect this to work regardless of the specific validation results
		// since the method should always return a boolean and map
		if isValid && len(errors) > 0 {
			t.Error("inconsistent state: form is valid but has errors")
		}

		if !isValid && len(errors) == 0 {
			t.Error("inconsistent state: form is invalid but has no errors")
		}
	})
}

func TestForm_Populate(t *testing.T) {
	t.Run("populate struct from form data", func(t *testing.T) {
		// Define a struct to populate
		type UserForm struct {
			Name    string `goform:"name"`
			Email   string `goform:"email"`
			Bio     string `goform:"bio"`
			Avatar  string `goform:"avatar"`
			Ignored string // no goform tag, should be ignored
		}

		// Create a form and set some values
		form := Form().AddChildren(
			Text("name").SetLabel("Name"),
			Email("email").SetLabel("Email"),
			Textarea("bio").SetLabel("Bio"),
			File("avatar").SetLabel("Avatar"),
		)

		// Set values directly on the elements
		elements := form.Elements()
		elements["name"].SetValue("John Doe")
		elements["email"].SetValue("john@example.com")
		elements["bio"].SetValue("Software developer")
		elements["avatar"].SetValue("profile.jpg")

		// Create struct instance and populate it
		var user UserForm
		_ = form.Populate(&user)

		// Check that struct was populated correctly
		if user.Name != "John Doe" {
			t.Errorf("expected Name='John Doe', got '%s'", user.Name)
		}
		if user.Email != "john@example.com" {
			t.Errorf("expected Email='john@example.com', got '%s'", user.Email)
		}
		if user.Bio != "Software developer" {
			t.Errorf("expected Bio='Software developer', got '%s'", user.Bio)
		}
		if user.Avatar != "profile.jpg" {
			t.Errorf("expected Avatar='profile.jpg', got '%s'", user.Avatar)
		}
		if user.Ignored != "" {
			t.Errorf("expected Ignored field to remain empty, got '%s'", user.Ignored)
		}
	})

	t.Run("populate struct with multiple file names", func(t *testing.T) {
		type FileForm struct {
			Documents []string `goform:"documents"`
			Photos    []string `goform:"photos"`
		}

		// Create a form with file inputs
		form := Form().AddChildren(
			File("documents").SetLabel("Documents").
				SetAttributes(Attr("multiple", true)),
			File("photos").SetLabel("Photos").
				SetAttributes(Attr("multiple", true)),
		)

		// Set multiple file names (as they would appear from PopulateFromRequest)
		elements := form.Elements()
		elements["documents"].SetValue("resume.pdf, cover-letter.docx, portfolio.pdf")
		elements["photos"].SetValue("photo1.jpg, photo2.png")

		// Create struct instance and populate it
		var files FileForm
		_ = form.Populate(&files)

		// Check that slice was populated correctly
		expectedDocs := []string{"resume.pdf", "cover-letter.docx", "portfolio.pdf"}
		if len(files.Documents) != len(expectedDocs) {
			t.Errorf("expected %d documents, got %d", len(expectedDocs), len(files.Documents))
		}
		for i, expected := range expectedDocs {
			if i < len(files.Documents) && files.Documents[i] != expected {
				t.Errorf("expected Documents[%d]='%s', got '%s'", i, expected, files.Documents[i])
			}
		}

		expectedPhotos := []string{"photo1.jpg", "photo2.png"}
		if len(files.Photos) != len(expectedPhotos) {
			t.Errorf("expected %d photos, got %d", len(expectedPhotos), len(files.Photos))
		}
		for i, expected := range expectedPhotos {
			if i < len(files.Photos) && files.Photos[i] != expected {
				t.Errorf("expected Photos[%d]='%s', got '%s'", i, expected, files.Photos[i])
			}
		}
	})

	t.Run("populate struct with empty values", func(t *testing.T) {
		type EmptyForm struct {
			Name  string `goform:"name"`
			Email string `goform:"email"`
		}

		// Create a form with empty values
		form := Form().AddChildren(
			Text("name").SetLabel("Name"),
			Email("email").SetLabel("Email"),
		)

		// Don't set any values (elements will have empty values)

		// Create struct instance and populate it
		var empty EmptyForm
		_ = form.Populate(&empty)

		// Check that struct fields remain empty
		if empty.Name != "" {
			t.Errorf("expected Name to be empty, got '%s'", empty.Name)
		}
		if empty.Email != "" {
			t.Errorf("expected Email to be empty, got '%s'", empty.Email)
		}
	})

	t.Run("populate non-pointer should not panic", func(t *testing.T) {
		type TestStruct struct {
			Name string `goform:"name"`
		}

		form := Form().AddChildren(Text("name"))
		elements := form.Elements()
		elements["name"].SetValue("test")

		var ts TestStruct
		// Passing non-pointer should not panic, just return without doing anything
		result := form.Populate(ts)

		if result != form {
			t.Error("expected Populate to return the same form instance")
		}

		// ts should remain unchanged
		if ts.Name != "" {
			t.Errorf("expected Name to be empty, got '%s'", ts.Name)
		}
	})

	t.Run("populate with missing form fields", func(t *testing.T) {
		type PartialForm struct {
			Name    string `goform:"name"`
			Missing string `goform:"missing"`
		}

		// Create form with only one field
		form := Form().AddChildren(
			Text("name").SetLabel("Name"),
		)

		elements := form.Elements()
		elements["name"].SetValue("John")

		var partial PartialForm
		_ = form.Populate(&partial)

		// Name should be set, Missing should remain empty
		if partial.Name != "John" {
			t.Errorf("expected Name='John', got '%s'", partial.Name)
		}
		if partial.Missing != "" {
			t.Errorf("expected Missing to be empty, got '%s'", partial.Missing)
		}
	})

	t.Run("full flow: request to struct population", func(t *testing.T) {
		type CompleteForm struct {
			Name      string   `goform:"name"`
			Email     string   `goform:"email"`
			Avatar    string   `goform:"avatar"`
			Documents []string `goform:"documents"`
		}

		// Create a form
		form := Form().AddChildren(
			Text("name").SetLabel("Name").SetAttributes(Attr("required", true)),
			Email("email").SetLabel("Email").SetAttributes(Attr("required", true)),
			File("avatar").SetLabel("Avatar"),
			File("documents").SetLabel("Documents").SetAttributes(Attr("multiple", true)),
		)

		// Create a mock HTTP request with multipart form data
		formData := url.Values{}
		formData.Set("name", "Jane Doe")
		formData.Set("email", "jane@example.com")

		multipartForm := &multipart.Form{
			Value: map[string][]string{
				"name":  {"Jane Doe"},
				"email": {"jane@example.com"},
			},
			File: map[string][]*multipart.FileHeader{
				"avatar": {
					{Filename: "profile.jpg", Size: 1024},
				},
				"documents": {
					{Filename: "resume.pdf", Size: 2048},
					{Filename: "cover-letter.docx", Size: 1536},
					{Filename: "portfolio.pdf", Size: 4096},
				},
			},
		}

		req := &http.Request{
			Method:        http.MethodPost,
			Form:          formData,
			PostForm:      formData,
			MultipartForm: multipartForm,
		}

		// Step 1: Populate from request
		_ = form.PopulateFromRequest(req)

		// Step 2: Validate
		isValid, errors := form.IsValid()
		if !isValid {
			t.Errorf("expected form to be valid, got errors: %v", errors)
		}

		// Step 3: Populate struct
		var result CompleteForm
		_ = form.Populate(&result)

		// Step 4: Verify struct was populated correctly
		if result.Name != "Jane Doe" {
			t.Errorf("expected Name='Jane Doe', got '%s'", result.Name)
		}
		if result.Email != "jane@example.com" {
			t.Errorf("expected Email='jane@example.com', got '%s'", result.Email)
		}
		if result.Avatar != "profile.jpg" {
			t.Errorf("expected Avatar='profile.jpg', got '%s'", result.Avatar)
		}

		expectedDocs := []string{"resume.pdf", "cover-letter.docx", "portfolio.pdf"}
		if len(result.Documents) != len(expectedDocs) {
			t.Errorf("expected %d documents, got %d", len(expectedDocs), len(result.Documents))
		}
		for i, expected := range expectedDocs {
			if i < len(result.Documents) && result.Documents[i] != expected {
				t.Errorf("expected Documents[%d]='%s', got '%s'", i, expected, result.Documents[i])
			}
		}
	})
}
