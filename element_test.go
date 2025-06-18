package goforms

import (
	"testing"

	"github.com/mickaelvieira/goforms/attr"
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
		createFunc   func(string, ...attr.Modifier) *element
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
		{"Submit", Submit, InputTypeSubmit, "button"},
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
				if elem.Attributes.String("type") != tt.expectedType {
					t.Errorf("expected type=%s, got %s", tt.expectedType, elem.Attributes.String("type"))
				}
			}

			// Check ID is generated
			if elem.Id() == "" {
				t.Error("expected ID to be generated")
			}

			// Check default attributes
			if elem.Attributes.String("aria-invalid") != "false" {
				t.Errorf("expected aria-invalid=false, got %s", elem.Attributes.String("aria-invalid"))
			}
			if elem.Attributes.String("aria-required") != "false" {
				t.Errorf("expected aria-required=false, got %s", elem.Attributes.String("aria-required"))
			}
		})
	}
}

func TestNewElementWithModifiers(t *testing.T) {
	elem := Text("test-field",
		attr.Required(true),
		attr.Attr("class", "form-control"),
		attr.Attr("placeholder", "Enter text"),
	)

	if elem.Name() != "test-field" {
		t.Errorf("expected name=test-field, got %s", elem.Name())
	}
	if elem.Attributes.String("type") != InputTypeText {
		t.Errorf("expected type=text, got %s", elem.Attributes.String("type"))
	}
	if !elem.Attributes.Bool("required") {
		t.Error("expected required=true")
	}
	if elem.Attributes.String("class") != "form-control" {
		t.Errorf("expected class=form-control, got %s", elem.Attributes.String("class"))
	}
	if elem.Attributes.String("placeholder") != "Enter text" {
		t.Errorf("expected placeholder=Enter text, got %s", elem.Attributes.String("placeholder"))
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
	if len(id) != 8 {
		t.Errorf("expected ID length 8, got %d", len(id))
	}
}

func TestElement_Name(t *testing.T) {
	elem := Text("test-name")
	if elem.Name() != "test-name" {
		t.Errorf("expected name=test-name, got %s", elem.Name())
	}
}

func TestElement_Attribute(t *testing.T) {
	elem := Text("test", attr.Attr("class", "form-control"))

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
	elem.SetAttributes(attr.Required(true))
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
		elem := Text("test", attr.Required(true))
		elem.SetValue("some-value")
		if !elem.IsValid() {
			t.Error("expected required field with value to be valid")
		}
	})

	t.Run("required field without value is invalid", func(t *testing.T) {
		elem := Text("test", attr.Required(true))
		if elem.IsValid() {
			t.Error("expected required field without value to be invalid")
		}
	})
}

func TestElement_SetError(t *testing.T) {
	elem := Text("test")
	result := elem.SetError("  This field is required  ")

	if elem.Error != "This field is required" {
		t.Errorf("expected error=This field is required, got %s", elem.Error)
	}
	if result != elem {
		t.Error("SetError should return the same element instance")
	}
}

func TestElement_SetLabel(t *testing.T) {
	elem := Text("test")
	result := elem.SetLabel("  Test Label  ")

	if elem.Label != "Test Label" {
		t.Errorf("expected label=Test Label, got %s", elem.Label)
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

	if elem.Hint != "This is a hint" {
		t.Errorf("expected hint=This is a hint, got %s", elem.Hint)
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

	if len(elem.Options) != 2 {
		t.Errorf("expected 2 options, got %d", len(elem.Options))
	}
	if elem.Options[0].Label != "Label 1" {
		t.Errorf("expected first option label=Label 1, got %s", elem.Options[0].Label)
	}
	if elem.Options[0].Value != "value1" {
		t.Errorf("expected first option value=value1, got %s", elem.Options[0].Value)
	}
	if elem.Options[1].Label != "Label 2" {
		t.Errorf("expected second option label=Label 2, got %s", elem.Options[1].Label)
	}
	if elem.Options[1].Value != "value2" {
		t.Errorf("expected second option value=value2, got %s", elem.Options[1].Value)
	}
	if result != elem {
		t.Error("SetOptions should return the same element instance")
	}
}

func TestElement_SetAttributes(t *testing.T) {
	elem := Text("test")
	result := elem.SetAttributes(
		attr.Attr("class", "form-control"),
		attr.Required(true),
	)

	if elem.Attributes.String("class") != "form-control" {
		t.Errorf("expected class=form-control, got %s", elem.Attributes.String("class"))
	}
	if !elem.Attributes.Bool("required") {
		t.Error("expected required=true")
	}
	if result != elem {
		t.Error("SetAttributes should return the same element instance")
	}
}

func TestElement_MarkAsInvalid(t *testing.T) {
	elem := Text("test")
	elem.MarkAsInvalid()

	if elem.Attributes.String("aria-invalid") != "true" {
		t.Errorf("expected aria-invalid=true, got %s", elem.Attributes.String("aria-invalid"))
	}
}

func TestElement_Render(t *testing.T) {
	elem := Text("test", attr.Id("test-id"))
	result := elem.Render()
	htmlStr := string(result)

	// Attributes in alphabetical order: aria-errormessage, aria-invalid, aria-required, id, name, type
	expected := `<input aria-errormessage="test-id-error" aria-invalid="false" aria-required="false" id="test-id" name="test" type="text">`

	if htmlStr != expected {
		t.Errorf("expected exact HTML match:\nExpected: %s\nActual:   %s", expected, htmlStr)
	}
}

func TestElement_RenderWithAttributes(t *testing.T) {
	elem := Text("username",
		attr.Id("username-field"),
		attr.Required(true),
		attr.Attr("class", "form-control"),
		attr.Attr("placeholder", "Enter username"),
	)
	result := elem.Render()
	htmlStr := string(result)

	// Alphabetical order: aria-errormessage, aria-invalid, aria-required, class, id, name, placeholder, required, type
	expected := `<input aria-errormessage="username-field-error" aria-invalid="false" aria-required="true" class="form-control" id="username-field" name="username" placeholder="Enter username" required type="text">`

	if htmlStr != expected {
		t.Errorf("expected exact HTML match:\nExpected: %s\nActual:   %s", expected, htmlStr)
	}
}

func TestElement_RenderCheckbox(t *testing.T) {
	elem := Checkbox("agree",
		attr.Id("agree-checkbox"),
		attr.Attr("value", "yes"),
	)
	result := elem.Render()
	htmlStr := string(result)

	// Alphabetical order: aria-errormessage, aria-invalid, aria-required, id, name, type, value
	expected := `<input aria-errormessage="agree-checkbox-error" aria-invalid="false" aria-required="false" id="agree-checkbox" name="agree" type="checkbox" value="yes">`

	if htmlStr != expected {
		t.Errorf("expected exact HTML match:\nExpected: %s\nActual:   %s", expected, htmlStr)
	}
}

func TestElement_RenderSelect(t *testing.T) {
	elem := Select("country", attr.Id("country-select")).SetOptions(
		Option("United States", "us"),
		Option("Canada", "ca"),
	)
	result := elem.Render()
	htmlStr := string(result)

	// Alphabetical order: aria-errormessage, aria-invalid, aria-required, id, name
	expected := `<select aria-errormessage="country-select-error" aria-invalid="false" aria-required="false" id="country-select" name="country"><option value="us">United States</option><option value="ca">Canada</option></select>`

	if htmlStr != expected {
		t.Errorf("expected exact HTML match:\nExpected: %s\nActual:   %s", expected, htmlStr)
	}
}

func TestElement_RenderError(t *testing.T) {
	elem := Text("test", attr.Id("test-field"))
	elem.SetError("This field is required")
	result := elem.RenderError()
	htmlStr := string(result)

	expected := `<div class="error" id="test-field">This field is required</div>`

	if htmlStr != expected {
		t.Errorf("expected exact HTML match:\nExpected: %s\nActual:   %s", expected, htmlStr)
	}
}

func TestElement_RenderTextarea(t *testing.T) {
	elem := Textarea("description",
		attr.Id("desc-textarea"),
		attr.Attr("rows", "5"),
		attr.Attr("cols", "40"),
	)
	result := elem.Render()
	htmlStr := string(result)

	// Alphabetical order: aria-errormessage, aria-invalid, aria-required, cols, id, name, rows
	expected := `<textarea aria-errormessage="desc-textarea-error" aria-invalid="false" aria-required="false" cols="40" id="desc-textarea" name="description" rows="5"></textarea>`

	if htmlStr != expected {
		t.Errorf("expected exact HTML match:\nExpected: %s\nActual:   %s", expected, htmlStr)
	}
}

func TestElement_RenderButton(t *testing.T) {
	elem := Submit("submit-btn",
		attr.Id("submit-button"),
		attr.Attr("class", "btn-primary"),
	)
	result := elem.Render()
	htmlStr := string(result)

	// Alphabetical order: aria-errormessage, aria-invalid, aria-required, class, id, name, type
	expected := `<button aria-errormessage="submit-button-error" aria-invalid="false" aria-required="false" class="btn-primary" id="submit-button" name="submit-btn" type="submit">Submit</button>`

	if htmlStr != expected {
		t.Errorf("expected exact HTML match:\nExpected: %s\nActual:   %s", expected, htmlStr)
	}
}

func TestElement_RenderRadio(t *testing.T) {
	elem := Radio("gender",
		attr.Id("gender-male"),
		attr.Attr("value", "male"),
	)
	result := elem.Render()
	htmlStr := string(result)

	// Alphabetical order: aria-errormessage, aria-invalid, aria-required, id, name, type, value
	expected := `<input aria-errormessage="gender-male-error" aria-invalid="false" aria-required="false" id="gender-male" name="gender" type="radio" value="male">`

	if htmlStr != expected {
		t.Errorf("expected exact HTML match:\nExpected: %s\nActual:   %s", expected, htmlStr)
	}
}

func TestElement_RenderWithBooleanAttributes(t *testing.T) {
	elem := Text("test",
		attr.Id("test-bool"),
		attr.Required(true),
		attr.Attr("disabled", true),
		attr.Attr("readonly", false), // Should not appear in output
	)
	result := elem.Render()
	htmlStr := string(result)

	// Alphabetical order: aria-errormessage, aria-required, disabled, id, name, required, type
	// Note: aria-invalid="false" is set by default, readonly=false is omitted
	expected := `<input aria-errormessage="test-bool-error" aria-invalid="false" aria-required="true" disabled id="test-bool" name="test" required type="text">`

	if htmlStr != expected {
		t.Errorf("expected exact HTML match:\nExpected: %s\nActual:   %s", expected, htmlStr)
	}
}

func TestElement_RenderEmptyAttributes(t *testing.T) {
	elem := Text("test",
		attr.Id("test-empty"),
		attr.Attr("placeholder", ""),       // Empty string should not appear
		attr.Attr("class", "form-control"), // Non-empty should appear
	)
	result := elem.Render()
	htmlStr := string(result)

	// Alphabetical order: aria-errormessage, aria-invalid, aria-required, class, id, name, type
	// Note: placeholder="" is omitted because it's empty
	expected := `<input aria-errormessage="test-empty-error" aria-invalid="false" aria-required="false" class="form-control" id="test-empty" name="test" type="text">`

	if htmlStr != expected {
		t.Errorf("expected exact HTML match:\nExpected: %s\nActual:   %s", expected, htmlStr)
	}
}

func TestElementTemplateMapping(t *testing.T) {
	tests := []struct {
		name         string
		elementType  string
		expectedTmpl string
	}{
		{"checkbox input", InputTypeCheckbox, "checkbox"},
		{"radio input", InputTypeRadio, "radio"},
		{"submit input", InputTypeSubmit, "button"},
		{"text input", InputTypeText, "input"},
		{"email input", InputTypeEmail, "input"},
		{"password input", InputTypePassword, "input"},
		{"select element", SelectElement, SelectElement},
		{"textarea element", TextareaElement, TextareaElement},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elem := newElement("test", tt.elementType)
			if elem.template != tt.expectedTmpl {
				t.Errorf("expected template=%s, got %s", tt.expectedTmpl, elem.template)
			}
		})
	}
}

func TestElementChaining(t *testing.T) {
	elem := Text("test").
		SetLabel("Test Field").
		SetHint("Enter some text").
		SetError("This field is required").
		SetAttributes(attr.Required(true), attr.Attr("class", "form-control"))

	if elem.Label != "Test Field" {
		t.Errorf("expected label=Test Field, got %s", elem.Label)
	}
	if elem.Hint != "Enter some text" {
		t.Errorf("expected hint=Enter some text, got %s", elem.Hint)
	}
	if elem.Error != "This field is required" {
		t.Errorf("expected error=This field is required, got %s", elem.Error)
	}
	if !elem.IsRequired() {
		t.Error("expected required=true")
	}
	if elem.Attributes.String("class") != "form-control" {
		t.Errorf("expected class=form-control, got %s", elem.Attributes.String("class"))
	}
}

func TestDefaultModifiers(t *testing.T) {
	elem := Text("test")

	// Check that default modifiers are applied
	if elem.Attributes.String("aria-invalid") != "false" {
		t.Errorf("expected aria-invalid=false by default, got %s", elem.Attributes.String("aria-invalid"))
	}
	if elem.Attributes.String("aria-required") != "false" {
		t.Errorf("expected aria-required=false by default, got %s", elem.Attributes.String("aria-required"))
	}
}

func TestModifiersOverrideDefaults(t *testing.T) {
	elem := Text("test", attr.Required(true))

	// Required modifier should override the default aria-required=false
	if elem.Attributes.String("aria-required") != "true" {
		t.Errorf("expected aria-required=true after Required modifier, got %s", elem.Attributes.String("aria-required"))
	}
}
