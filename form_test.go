package goforms

import (
	"html/template"
	"log"
	"os"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/mickaelvieira/goforms/attr"
)

func TestMain(m *testing.M) {
	testFS := fstest.MapFS{
		"form.html": &fstest.MapFile{
			Data: []byte(`<form{{attributes .Attributes}}>{{range .Elements}}{{element .}}{{end}}</form>`),
		},
		"input.html": &fstest.MapFile{
			Data: []byte(`<input{{attributes .Attributes}}>`),
		},
		"checkbox.html": &fstest.MapFile{
			Data: []byte(`<input{{attributes .Attributes}}>`),
		},
		"select.html": &fstest.MapFile{
			Data: []byte(`<select{{attributes .Attributes}}>{{range .Options}}<option value="{{.Value}}">{{.Label}}</option>{{end}}</select>`),
		},
		"textarea.html": &fstest.MapFile{
			Data: []byte(`<textarea{{attributes .Attributes}}></textarea>`),
		},
		"button.html": &fstest.MapFile{
			Data: []byte(`<button{{attributes .Attributes}}>Submit</button>`),
		},
		"radio.html": &fstest.MapFile{
			Data: []byte(`<input{{attributes .Attributes}}>`),
		},
		"error.html": &fstest.MapFile{
			Data: []byte(`<div class="error" id="{{.ID}}">{{.Message}}</div>`),
		},
		"fieldset.html": &fstest.MapFile{
			Data: []byte(`<fieldset><legend>{{.Legend}}</legend>{{range .Elements}}{{element .}}{{end}}</fieldset>`),
		},
		"group.html": &fstest.MapFile{
			Data: []byte(`<div class="{{.Class}}">{{range .Elements}}{{element .}}{{end}}</div>`),
		},
	}

	SetOverridingTemplates(testFS, "*.html")

	code := m.Run()

	log.Printf("TestMain completed with exit code %d", code)

	os.Exit(code)
}

func TestForm_Creation(t *testing.T) {
	t.Run("empty form", func(t *testing.T) {
		f := Form()

		// Use public methods instead of accessing private fields
		children := f.Children()
		if len(children) != 0 {
			t.Errorf("expected 0 elements, got %d", len(children))
		}

		// Test that render works (indicates template is initialized)
		result := f.Render()
		if result == template.HTML("") {
			t.Error("expected non-empty render result, template might not be initialized")
		}
	})

	t.Run("form with single element", func(t *testing.T) {
		textField := Text("username")
		f := Form(textField)

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

		f := Form(usernameField, emailField, passwordField)

		children := f.Children()
		if len(children) != 3 {
			t.Errorf("expected 3 elements, got %d", len(children))
		}

		indexUsername := slices.IndexFunc(children, func(e Element) bool {
			return e == usernameField
		})
		indexEmail := slices.IndexFunc(children, func(e Element) bool {
			return e == emailField
		})
		indexPassword := slices.IndexFunc(children, func(e Element) bool {
			return e == passwordField
		})
		indexTelephone := slices.IndexFunc(children, func(e Element) bool {
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
			attr.Attr("method", "POST"),
			attr.Attr("action", "/submit"),
			attr.Attr("class", "login-form"),
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
		f := Form().SetAttributes(attr.Attr("method", "GET"))
		result := f.SetAttributes(attr.Attr("method", "POST"))

		if f.attributes.String("method") != "POST" {
			t.Errorf("expected method=POST after override, got %s", f.attributes.String("method"))
		}
		if result != f {
			t.Error("SetAttributes should return the same form instance")
		}
	})

	t.Run("set multiple attributes at once", func(t *testing.T) {
		f := Form()
		f.SetAttributes(
			attr.Attr("method", "POST"),
			attr.Attr("action", "/submit"),
			attr.Attr("enctype", "multipart/form-data"),
			attr.Attr("novalidate", true),
			attr.Attr("autocomplete", "off"),
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
		f := Form(textField, selectField)

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

		htmlStr := string(result)
		if htmlStr == "" {
			t.Error("expected non-empty HTML string")
		}

		// Should contain form tags
		if !strings.Contains(htmlStr, "<form") {
			t.Error("expected rendered HTML to contain <form tag")
		}
		if !strings.Contains(htmlStr, "</form>") {
			t.Error("expected rendered HTML to contain </form> tag")
		}
	})

	t.Run("render form with attributes", func(t *testing.T) {
		f := Form().SetAttributes(
			attr.Attr("method", "POST"),
			attr.Attr("action", "/submit"),
		)
		result := f.Render()
		htmlStr := string(result)

		// Should contain the attributes
		if !strings.Contains(htmlStr, `method="POST"`) {
			t.Error("expected rendered HTML to contain method attribute")
		}
		if !strings.Contains(htmlStr, `action="/submit"`) {
			t.Error("expected rendered HTML to contain action attribute")
		}
	})

	t.Run("render form with elements", func(t *testing.T) {
		textField := Text("username", attr.Id("username-field"))
		emailField := Email("email", attr.Id("email-field"))

		f := Form(textField, emailField).SetAttributes(
			attr.Attr("method", "POST"),
			attr.Attr("action", "/submit"),
		)
		result := f.Render()
		htmlStr := string(result)

		// Should contain form structure
		if !strings.Contains(htmlStr, "<form") {
			t.Error("expected rendered HTML to contain <form tag")
		}
		if !strings.Contains(htmlStr, "</form>") {
			t.Error("expected rendered HTML to contain </form> tag")
		}

		// Should contain elements (though exact content depends on element rendering)
		if !strings.Contains(htmlStr, `name="username"`) {
			t.Error("expected rendered HTML to contain username field")
		}
		if !strings.Contains(htmlStr, `name="email"`) {
			t.Error("expected rendered HTML to contain email field")
		}
	})

	t.Run("render returns template.HTML type", func(t *testing.T) {
		f := Form()
		result := f.Render()

		// Verify it's the correct type
		var _ template.HTML = result
	})
}

func TestForm_ExactRender(t *testing.T) {
	t.Run("exact form render with fixed attributes", func(t *testing.T) {
		f := Form().SetAttributes(
			attr.Id("test-form"), // Use controlled ID
			attr.Attr("method", "POST"),
			attr.Attr("action", "/submit"),
			attr.Attr("class", "login-form"),
		)
		result := f.Render()
		htmlStr := string(result)

		// Alphabetical order: action, aria-errormessage, class, enctype, id, method
		expected := `<form action="/submit" aria-errormessage="test-form-error" class="login-form" enctype="application/x-www-form-urlencoded" id="test-form" method="POST"></form>`

		if htmlStr != expected {
			t.Errorf("expected exact HTML match:\nExpected: %s\nActual:   %s", expected, htmlStr)
		}
	})

	t.Run("exact form render with elements", func(t *testing.T) {
		textField := Text("username", attr.Id("username-field"))
		submitButton := Submit("submit", attr.Id("submit-btn"))

		f := Form(textField, submitButton).SetAttributes(
			attr.Id("login-form"), // Use controlled ID
			attr.Attr("method", "POST"),
			attr.Attr("action", "/login"),
		)
		result := f.Render()
		htmlStr := string(result)

		expectedForm := `<form action="/login" aria-errormessage="login-form-error" enctype="application/x-www-form-urlencoded" id="login-form" method="POST">`
		expectedUsernameField := `<input aria-errormessage="username-field-error" aria-invalid="false" aria-required="false" id="username-field" name="username" type="text">`
		expectedSubmitButton := `<button aria-errormessage="submit-btn-error" aria-invalid="false" aria-required="false" id="submit-btn" name="submit" type="submit">Submit</button>`
		expectedFormEnd := `</form>`

		expected := expectedForm + expectedUsernameField + expectedSubmitButton + expectedFormEnd

		if htmlStr != expected {
			t.Errorf("expected exact HTML match:\nExpected: %s\nActual:   %s", expected, htmlStr)
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

		f := Form(
			usernameField,
			emailField,
			passwordField,
			agreeCheckbox,
			countrySelect,
			descriptionTextarea,
			submitButton,
		).SetAttributes(
			attr.Attr("method", "POST"),
			attr.Attr("action", "/register"),
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

		// Test that we can render the form without errors
		result := f.Render()
		if result == template.HTML("") {
			t.Error("expected non-empty render result for complex form")
		}
	})

	t.Run("form with nested containers", func(t *testing.T) {
		userFieldset := FieldSet("User Information",
			Text("username", attr.Id("username-field")),
			Email("email", attr.Id("email-field")),
		)

		addressGroup := Group(
			Text("street", attr.Id("street-field")),
			Text("city", attr.Id("city-field")),
		).SetClass("address-group")

		f := Form(userFieldset, addressGroup).SetAttributes(
			attr.Id("nested-form"),
			attr.Attr("method", "POST"),
			attr.Attr("action", "/submit"),
		)

		result := f.Render()
		htmlStr := string(result)

		// Test the complete rendered form with nested containers
		expectedForm := `<form action="/submit" aria-errormessage="nested-form-error" enctype="application/x-www-form-urlencoded" id="nested-form" method="POST">`
		expectedFieldset := `<fieldset><legend>User Information</legend>`
		expectedUsernameField := `<input aria-errormessage="username-field-error" aria-invalid="false" aria-required="false" id="username-field" name="username" type="text">`
		expectedEmailField := `<input aria-errormessage="email-field-error" aria-invalid="false" aria-required="false" id="email-field" name="email" type="email">`
		expectedFieldsetEnd := `</fieldset>`
		expectedGroupStart := `<div class="address-group">`
		expectedStreetField := `<input aria-errormessage="street-field-error" aria-invalid="false" aria-required="false" id="street-field" name="street" type="text">`
		expectedCityField := `<input aria-errormessage="city-field-error" aria-invalid="false" aria-required="false" id="city-field" name="city" type="text">`
		expectedGroupEnd := `</div>`
		expectedFormEnd := `</form>`

		expected := expectedForm + expectedFieldset + expectedUsernameField + expectedEmailField + expectedFieldsetEnd + expectedGroupStart + expectedStreetField + expectedCityField + expectedGroupEnd + expectedFormEnd

		if htmlStr != expected {
			t.Errorf("expected exact HTML match:\nExpected: %s\nActual:   %s", expected, htmlStr)
		}
	})
}

func TestForm_ElementOrder(t *testing.T) {
	t.Run("maintains element order", func(t *testing.T) {
		elements := []Element{
			Text("field1"),
			Email("field2"),
			Password("field3"),
			Checkbox("field4"),
			Select("field5"),
		}

		// Convert to Renderer slice for Form constructor
		renderers := make([]Renderer, len(elements))
		for i, elem := range elements {
			renderers[i] = elem
		}

		f := Form(renderers...)
		children := f.Children()

		if len(children) != len(elements) {
			t.Errorf("expected %d children, got %d", len(elements), len(children))
		}

		for i, expectedElem := range elements {
			if children[i] != expectedElem {
				t.Errorf("element at position %d doesn't match expected element", i)
			}
		}
	})
}

func TestForm_NilElements(t *testing.T) {
	t.Run("handles mixed nil and valid elements", func(t *testing.T) {
		validElement := Text("valid")

		// Create form with valid and nil elements
		f := Form(validElement, nil, nil)

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

		f := Form(textField, emailField).
			SetAttributes(attr.Attr("method", "POST")).
			SetAttributes(attr.Attr("action", "/submit")).
			SetAttributes(attr.Attr("class", "user-form"))

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
			f := Form().SetAttributes(attr.Attr(tt.attribute, tt.value))
			if f.attributes.String(tt.attribute) != tt.expected {
				t.Errorf("expected %s=%s, got %s", tt.attribute, tt.expected, f.attributes.String(tt.attribute))
			}
		})
	}
}

func TestForm_BooleanAttributes(t *testing.T) {
	t.Run("boolean attributes", func(t *testing.T) {
		f := Form().SetAttributes(
			attr.Attr("novalidate", true),
			attr.Attr("autocomplete", false), // Should not appear
		)

		if !f.attributes.Bool("novalidate") {
			t.Error("expected novalidate=true")
		}

		// Boolean false attributes should not be present
		if f.attributes.Bool("autocomplete") {
			t.Error("expected autocomplete to be false/absent")
		}
	})
}

func TestForm_EmptyForm(t *testing.T) {
	t.Run("form with no elements should still render", func(t *testing.T) {
		f := Form().SetAttributes(
			attr.Id("empty-form"), // Use controlled ID
			attr.Attr("method", "POST"),
			attr.Attr("action", "/submit"),
		)

		result := f.Render()
		htmlStr := string(result)

		expected := `<form action="/submit" aria-errormessage="empty-form-error" enctype="application/x-www-form-urlencoded" id="empty-form" method="POST"></form>`

		if htmlStr != expected {
			t.Errorf("expected exact HTML match:\nExpected: %s\nActual:   %s", expected, htmlStr)
		}
	})
}
