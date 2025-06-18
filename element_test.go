package goform

import (
	"testing"
)

func TestIsInputType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"select element", SelectElement, false},
		{"textarea element", TextareaElement, false},
		{"input type text", InputTypeText, true},
		{"input type email", InputTypeEmail, true},
		{"input type password", InputTypePassword, true},
		{"custom input type", "custom", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isInputType(tt.input)
			if result != tt.expected {
				t.Errorf("isInputType(%s) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestElementCreationFunctions(t *testing.T) {
	tests := []struct {
		name         string
		createFunc   func(string) *element
		expectedType string
		expectedTmpl string
	}{
		{"Phone", Phone, InputTypeTel, "input"},
		{"Number", Number, InputTypeNumber, "input"},
		{"Search", Search, InputTypeSearch, "input"},
		{"Url", Url, InputTypeUrl, "input"},
		{"Color", Color, InputTypeColor, "input"},
		{"Range", Range, InputTypeRange, "input"},
		{"Date", Date, InputTypeDate, "input"},
		{"DateTimeLocal", DateTimeLocal, InputTypeDateTimeLocal, "input"},
		{"File", File, InputTypeFile, "input"},
		{"Checkbox", Checkbox, InputTypeCheckbox, "checkbox"},
		{"Radio", Radio, InputTypeRadio, "radio"},
		{"Hidden", Hidden, InputTypeHidden, "input"},
		{"Submit", Submit, InputTypeSubmit, "input"},
		{"Button", Button, InputTypeButton, "input"},
		{"Reset", Reset, InputTypeReset, "input"},
		{"Image", Image, InputTypeImage, "input"},
		{"Time", Time, InputTypeTime, "input"},
		{"Month", Month, InputTypeMonth, "input"},
		{"Week", Week, InputTypeWeek, "input"},
		{"Datetime", Datetime, InputTypeDatetime, "input"},
		{"Text", Text, InputTypeText, "input"},
		{"Email", Email, InputTypeEmail, "input"},
		{"Password", Password, InputTypePassword, "input"},
		{"Textarea", Textarea, TextareaElement, TextareaElement},
		{"Select", Select, SelectElement, SelectElement},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elem := tt.createFunc("test-field")

			// Check name
			if elem.Name() != "test-field" {
				t.Errorf("expected name=test-field, got %s", elem.Name())
			}

			// Check template
			if elem.template != tt.expectedTmpl {
				t.Errorf("expected template=%s, got %s", tt.expectedTmpl, elem.template)
			}

			// Check type attribute for input elements
			if isInputType(tt.expectedType) {
				if elem.attributes.String("type") != tt.expectedType {
					t.Errorf("expected type=%s, got %s", tt.expectedType, elem.attributes.String("type"))
				}
			}

			// Check ID is generated
			if elem.Id() == "" {
				t.Error("expected ID to be generated")
			}

			// Check default attributes
			if elem.attributes.String("aria-invalid") != "false" {
				t.Errorf("expected aria-invalid=false, got %s", elem.attributes.String("aria-invalid"))
			}
			if elem.attributes.String("aria-required") != "false" {
				t.Errorf("expected aria-required=false, got %s", elem.attributes.String("aria-required"))
			}
		})
	}
}

func TestNewElementWithModifiers(t *testing.T) {
	elem := Text("test-field").SetAttributes(
		Required(true),
		Attr("class", "form-control"),
		Attr("placeholder", "Enter text"),
	)

	if elem.Name() != "test-field" {
		t.Errorf("expected name=test-field, got %s", elem.Name())
	}
	if elem.attributes.String("type") != InputTypeText {
		t.Errorf("expected type=text, got %s", elem.attributes.String("type"))
	}
	if !elem.attributes.Bool("required") {
		t.Error("expected required=true")
	}
	if elem.attributes.String("class") != "form-control" {
		t.Errorf("expected class=form-control, got %s", elem.attributes.String("class"))
	}
	if elem.attributes.String("placeholder") != "Enter text" {
		t.Errorf("expected placeholder=Enter text, got %s", elem.attributes.String("placeholder"))
	}
}

func TestOption(t *testing.T) {
	opt := Option("  Test Label  ", "  test-value  ")

	if opt.Label != "Test Label" {
		t.Errorf("expected label=Test Label, got %s", opt.Label)
	}
	if opt.Value != "test-value" {
		t.Errorf("expected value=test-value, got %s", opt.Value)
	}
}

func TestElement_Id(t *testing.T) {
	elem := Text("test")
	id := elem.Id()

	if id == "" {
		t.Error("expected ID to be set")
	}
	if len(id) != 10 {
		t.Errorf("expected ID length 10, got %d", len(id))
	}
}

func TestElement_Name(t *testing.T) {
	elem := Text("test-name")
	if elem.Name() != "test-name" {
		t.Errorf("expected name=test-name, got %s", elem.Name())
	}
}

func TestElement_Attribute(t *testing.T) {
	elem := Text("test").SetAttributes(Attr("class", "form-control"))

	if elem.Attribute("class") != "form-control" {
		t.Errorf("expected class=form-control, got %v", elem.Attribute("class"))
	}
	if elem.Attribute("nonexistent") != "" {
		t.Errorf("expected empty for nonexistent attribute, got %v", elem.Attribute("nonexistent"))
	}
}

func TestElement_Value(t *testing.T) {
	elem := Text("test")

	// Initially empty
	if elem.Value() != "" {
		t.Errorf("expected empty value, got %s", elem.Value())
	}

	// Set value
	elem.SetValue("test-value")
	if elem.Value() != "test-value" {
		t.Errorf("expected value=test-value, got %s", elem.Value())
	}
}

func TestElement_IsRequired(t *testing.T) {
	elem := Text("test")

	// Initially not required
	if elem.IsRequired() {
		t.Error("expected required=false initially")
	}

	// Set required
	elem.SetAttributes(Required(true))
	if !elem.IsRequired() {
		t.Error("expected required=true after setting")
	}
}

func TestElement_IsValid(t *testing.T) {
	t.Run("not required field is always valid", func(t *testing.T) {
		elem := Text("test")
		if !elem.IsValid() {
			t.Error("expected non-required field to be valid")
		}
	})

	t.Run("required field with value is valid", func(t *testing.T) {
		elem := Text("test").SetAttributes(Required(true))
		elem.SetValue("some-value")
		if !elem.IsValid() {
			t.Error("expected required field with value to be valid")
		}
	})

	t.Run("required field without value is invalid", func(t *testing.T) {
		elem := Text("test").SetAttributes(Required(true))
		if elem.IsValid() {
			t.Error("expected required field without value to be invalid")
		}
	})
}

func TestElement_SetError(t *testing.T) {
	elem := Text("test")
	result := elem.SetError("  This field is required  ")

	if elem.error != "This field is required" {
		t.Errorf("expected error=This field is required, got %s", elem.error)
	}
	if result != elem {
		t.Error("SetError should return the same element instance")
	}
}

func TestElement_SetLabel(t *testing.T) {
	elem := Text("test")
	result := elem.SetLabel("  Test Label  ")

	if elem.label != "Test Label" {
		t.Errorf("expected label=Test Label, got %s", elem.label)
	}
	if result != elem {
		t.Error("SetLabel should return the same element instance")
	}
}

func TestElement_SetValue(t *testing.T) {
	elem := Text("test")
	elem.SetValue("test-value")

	if elem.Value() != "test-value" {
		t.Errorf("expected value=test-value, got %s", elem.Value())
	}
}

func TestElement_SetHint(t *testing.T) {
	elem := Text("test")
	result := elem.SetHint("  This is a hint  ")

	if elem.hint != "This is a hint" {
		t.Errorf("expected hint=This is a hint, got %s", elem.hint)
	}
	if result != elem {
		t.Error("SetHint should return the same element instance")
	}
}

func TestElement_SetOptions(t *testing.T) {
	elem := Select("test")
	opt1 := Option("Label 1", "value1")
	opt2 := Option("Label 2", "value2")

	result := elem.SetOptions(opt1, opt2)

	if len(elem.options) != 2 {
		t.Errorf("expected 2 options, got %d", len(elem.options))
	}
	if elem.options[0].Label != "Label 1" {
		t.Errorf("expected first option label=Label 1, got %s", elem.options[0].Label)
	}
	if elem.options[0].Value != "value1" {
		t.Errorf("expected first option value=value1, got %s", elem.options[0].Value)
	}
	if elem.options[1].Label != "Label 2" {
		t.Errorf("expected second option label=Label 2, got %s", elem.options[1].Label)
	}
	if elem.options[1].Value != "value2" {
		t.Errorf("expected second option value=value2, got %s", elem.options[1].Value)
	}
	if result != elem {
		t.Error("SetOptions should return the same element instance")
	}
}

func TestElement_SetAttributes(t *testing.T) {
	elem := Text("test")
	result := elem.SetAttributes(
		Attr("class", "form-control"),
		Required(true),
	)

	if elem.attributes.String("class") != "form-control" {
		t.Errorf("expected class=form-control, got %s", elem.attributes.String("class"))
	}
	if !elem.attributes.Bool("required") {
		t.Error("expected required=true")
	}
	if result != elem {
		t.Error("SetAttributes should return the same element instance")
	}
}

func TestElement_MarkAsInvalid(t *testing.T) {
	elem := Text("test")
	elem.MarkAsInvalid()

	if elem.attributes.String("aria-invalid") != "true" {
		t.Errorf("expected aria-invalid=true, got %s", elem.attributes.String("aria-invalid"))
	}
}

func TestElement_Render(t *testing.T) {
	elem := Text("test").SetAttributes(Id("test-id"))
	result := elem.Render()
	htmlStr := cleanHTML(result)

	expected := `<div><div><input aria-errormessage="test-id-error" aria-invalid="false" aria-required="false" id="test-id" name="test" type="text"></div></div>`

	if htmlStr != expected {
		t.Errorf("expected exact HTML match:\nExpected: %s\nActual: %s", expected, htmlStr)
	}
}

func TestElement_RenderWithAttributes(t *testing.T) {
	elem := Text("username").
		SetAttributes(
			Id("username-field"),
			Required(true),
			Attr("class", "form-control"),
			Attr("placeholder", "Enter username"),
		).
		SetLabel("Username").
		SetError("This field is required").
		SetHint("Must be unique")

	result := elem.Render()
	htmlStr := cleanHTML(result)

	expected := `<div><label for="username-field">Username <span>*</span></label><div><input aria-describedby="username-field-hint" aria-errormessage="username-field-error" aria-invalid="false" aria-required="true" class="form-control" id="username-field" name="username" placeholder="Enter username" required type="text"><span id="username-field-error">This field is required</span><i id="username-field-hint">Must be unique</i></div></div>`

	if htmlStr != expected {
		t.Errorf("expected exact HTML match:\nExpected: %s\nActual:   %s", expected, htmlStr)
	}
}

func TestElement_RenderCheckbox(t *testing.T) {
	elem := Checkbox("agree").SetAttributes(
		Id("agree-checkbox"),
		Attr("value", "yes"),
	)
	result := elem.Render()
	htmlStr := cleanHTML(result)

	expected := `<div><label for="agree-checkbox"><input aria-errormessage="agree-checkbox-error" aria-invalid="false" aria-required="false" id="agree-checkbox" name="agree" type="checkbox" value="yes"></label></div>`

	if htmlStr != expected {
		t.Errorf("expected exact HTML match:\nExpected: %s\nActual: %s", expected, htmlStr)
	}
}

func TestElement_RenderSelect(t *testing.T) {
	elem := Select("country").SetAttributes(
		Id("country-select"),
	).SetOptions(
		Option("United States", "us"),
		Option("Canada", "ca"),
	)
	result := elem.Render()
	htmlStr := cleanHTML(result)

	expected := `<div><div><select aria-errormessage="country-select-error" aria-invalid="false" aria-required="false" id="country-select" name="country"><option value="us" >United States</option><option value="ca" >Canada</option></select></div></div>`

	if htmlStr != expected {
		t.Errorf("expected exact HTML match:\nExpected: %s\nActual: %s", expected, htmlStr)
	}
}

func TestElement_RenderError(t *testing.T) {
	elem := Text("test").SetAttributes(Id("test-field"))
	elem.SetError("This field is required")
	result := elem.RenderError()
	htmlStr := cleanHTML(result)

	expected := `<span id="test-field-error">This field is required</span>`

	if htmlStr != expected {
		t.Errorf("expected exact HTML match:\nExpected: %s\nActual:   %s", expected, htmlStr)
	}
}

func TestElement_RenderTextarea(t *testing.T) {
	elem := Textarea("description").SetAttributes(
		Id("desc-textarea"),
		Attr("rows", "5"),
		Attr("cols", "40"),
	)
	result := elem.Render()
	htmlStr := cleanHTML(result)

	expected := `<div><div><textarea aria-errormessage="desc-textarea-error" aria-invalid="false" aria-required="false" cols="40" id="desc-textarea" name="description" rows="5"></textarea></div></div>`

	if htmlStr != expected {
		t.Errorf("expected exact HTML match:\nExpected: %s\nActual: %s", expected, htmlStr)
	}
}

func TestElement_RenderButton(t *testing.T) {
	elem := Submit("submit-btn").SetAttributes(
		Id("submit-button"),
		Attr("class", "btn-primary"),
	)
	result := elem.Render()
	htmlStr := cleanHTML(result)

	expected := `<input aria-errormessage="submit-button-error" aria-invalid="false" aria-required="false" class="btn-primary" id="submit-button" name="submit-btn" type="submit">`

	if htmlStr != expected {
		t.Errorf("expected exact HTML match:\nExpected: %s\nActual: %s", expected, htmlStr)
	}
}

func TestElement_RenderRadio(t *testing.T) {
	elem := Radio("gender").SetAttributes(
		Id("gender-male"),
		Attr("value", "male"),
	)
	result := elem.Render()
	htmlStr := cleanHTML(result)

	expected := `<div><label for="gender-male"><input aria-errormessage="gender-male-error" aria-invalid="false" aria-required="false" id="gender-male" name="gender" type="radio" value="male"></label></div>`

	if htmlStr != expected {
		t.Errorf("expected exact HTML match:\nExpected: %s\nActual: %s", expected, htmlStr)
	}
}

func TestElement_RenderWithBooleanAttributes(t *testing.T) {
	elem := Text("test").SetAttributes(
		Id("test-bool"),
		Required(true),
		Attr("disabled", true),
		Attr("readonly", false), // Should not appear in output
	)
	result := elem.Render()
	htmlStr := cleanHTML(result)

	expected := `<div><div><input aria-errormessage="test-bool-error" aria-invalid="false" aria-required="true" disabled id="test-bool" name="test" required type="text"></div></div>`

	if htmlStr != expected {
		t.Errorf("expected exact HTML match:\nExpected: %s\nActual: %s", expected, htmlStr)
	}
}

func TestElement_RenderEmptyAttributes(t *testing.T) {
	elem := Text("test").SetAttributes(
		Id("test-empty"),
		Attr("placeholder", ""),       // Empty string should not appear
		Attr("class", "form-control"), // Non-empty should appear
	)
	result := elem.Render()
	htmlStr := cleanHTML(result)

	expected := `<div><div><input aria-errormessage="test-empty-error" aria-invalid="false" aria-required="false" class="form-control" id="test-empty" name="test" type="text"></div></div>`

	if htmlStr != expected {
		t.Errorf("expected exact HTML match:\nExpected: %s\nActual: %s", expected, htmlStr)
	}
}

func TestElementChaining(t *testing.T) {
	elem := Text("test").
		SetLabel("Test Field").
		SetHint("Enter some text").
		SetError("This field is required").
		SetAttributes(Required(true), Attr("class", "form-control"))

	if elem.label != "Test Field" {
		t.Errorf("expected label=Test Field, got %s", elem.label)
	}
	if elem.hint != "Enter some text" {
		t.Errorf("expected hint=Enter some text, got %s", elem.hint)
	}
	if elem.error != "This field is required" {
		t.Errorf("expected error=This field is required, got %s", elem.error)
	}
	if !elem.IsRequired() {
		t.Error("expected required=true")
	}
	if elem.attributes.String("class") != "form-control" {
		t.Errorf("expected class=form-control, got %s", elem.attributes.String("class"))
	}
}

func TestDefaultModifiers(t *testing.T) {
	elem := Text("test")

	if elem.attributes.String("aria-invalid") != "false" {
		t.Errorf("expected aria-invalid=false by default, got %s", elem.attributes.String("aria-invalid"))
	}
	if elem.attributes.String("aria-required") != "false" {
		t.Errorf("expected aria-required=false by default, got %s", elem.attributes.String("aria-required"))
	}
}

func TestModifiersOverrideDefaults(t *testing.T) {
	elem := Text("test").SetAttributes(Required(true))

	if elem.attributes.String("aria-required") != "true" {
		t.Errorf("expected aria-required=true after Required modifier, got %s", elem.attributes.String("aria-required"))
	}
}
