package goform

import (
	"html/template"
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
		f := Form()
		f.SetAttributes(
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
			AddChildren(textField, Group(submitButton, resetButton, regularButton)).
			SetAttributes(
				Id("login-form"),
				Attr("method", "POST"),
				Attr("action", "/login"),
			)

		result := f.Render()
		htmlStr := cleanHTML(result)

		expectedForm := `<form action="/login" aria-errormessage="login-form-error" enctype="application/x-www-form-urlencoded" id="login-form" method="POST">`
		expectedUsernameField := `<div><div><input aria-errormessage="username-field-error" aria-invalid="false" aria-required="false" id="username-field" name="username" type="text"></div></div>`
		expectedSubmitButton := `<div><input aria-errormessage="submit-btn-error" aria-invalid="false" aria-required="false" id="submit-btn" name="submit" type="submit">`
		expectedResetButton := `<input aria-errormessage="reset-btn-error" aria-invalid="false" aria-required="false" id="reset-btn" name="reset" type="reset">`
		expectedRegularButton := `<input aria-errormessage="button-btn-error" aria-invalid="false" aria-required="false" id="button-btn" name="button" type="button"></div>`
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
		expectedForm := `<form action="/submit" aria-errormessage="nested-form-error" enctype="application/x-www-form-urlencoded" id="nested-form" method="POST">`
		expectedFieldSet := `<fieldset><legend>User Information</legend>`
		expectedUsernameField := `<div><div><input aria-errormessage="username-field-error" aria-invalid="false" aria-required="false" id="username-field" name="username" type="text"></div></div>`
		expectedEmailField := `<div><div><input aria-errormessage="email-field-error" aria-invalid="false" aria-required="false" id="email-field" name="email" type="email"></div></div>`
		expectedFieldSetEnd := `</fieldset>`
		expectedGroupStart := `<div class="address-group">`
		expectedStreetField := `<div><div><input aria-errormessage="street-field-error" aria-invalid="false" aria-required="false" id="street-field" name="street" type="text"></div></div>`
		expectedCityField := `<div><div><input aria-errormessage="city-field-error" aria-invalid="false" aria-required="false" id="city-field" name="city" type="text"></div></div>`
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
		expected := `<form action="/submit" aria-errormessage="empty-form-error" enctype="application/x-www-form-urlencoded" id="empty-form" method="POST"></form>`

		if result != expected {
			t.Errorf("expected exact HTML match:\nExpected: %s\nActual: %s", expected, result)
		}
	})
}
