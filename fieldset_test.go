package goform

import (
	"html/template"
	"strings"
	"testing"
)

func TestFieldSet_Creation(t *testing.T) {
	t.Run("empty fieldset", func(t *testing.T) {
		fs := FieldSet("Test Legend")

		if fs.legend != "Test Legend" {
			t.Errorf("expected legend=Test Legend, got %s", fs.legend)
		}
		if len(fs.children) != 0 {
			t.Errorf("expected 0 elements, got %d", len(fs.children))
		}
		if fs.renderer == nil {
			t.Error("expected template to be initialized")
		}
	})

	t.Run("fieldset with single element", func(t *testing.T) {
		textField := Text("username")
		fs := FieldSet("User Details", textField)

		if fs.legend != "User Details" {
			t.Errorf("expected legend=User Details, got %s", fs.legend)
		}
		if len(fs.children) != 1 {
			t.Errorf("expected 1 element, got %d", len(fs.children))
		}
		if fs.children[0] != textField {
			t.Error("expected first element to be the text field")
		}
	})

	t.Run("fieldset with multiple elements", func(t *testing.T) {
		usernameField := Text("username")
		emailField := Email("email")
		passwordField := Password("password")

		fs := FieldSet("Registration Form", usernameField, emailField, passwordField)

		if fs.legend != "Registration Form" {
			t.Errorf("expected legend=Registration Form, got %s", fs.legend)
		}
		if len(fs.children) != 3 {
			t.Errorf("expected 3 elements, got %d", len(fs.children))
		}

		// Check elements are in correct order
		if fs.children[0] != usernameField {
			t.Error("expected first element to be username field")
		}
		if fs.children[1] != emailField {
			t.Error("expected second element to be email field")
		}
		if fs.children[2] != passwordField {
			t.Error("expected third element to be password field")
		}
	})
}

func TestFieldSet_Children(t *testing.T) {
	t.Run("empty fieldset children", func(t *testing.T) {
		fs := FieldSet("Empty")
		children := fs.Children()

		if len(children) != 0 {
			t.Errorf("expected 0 children, got %d", len(children))
		}
		if children == nil {
			t.Error("expected non-nil slice, even if empty")
		}
	})

	t.Run("fieldset with children", func(t *testing.T) {
		textField := Text("text")
		selectField := Select("select")
		fs := FieldSet("Form Fields", textField, selectField)

		children := fs.Children()

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

	t.Run("children returns same slice reference", func(t *testing.T) {
		textField := Text("text")
		fs := FieldSet("Test", textField)

		children1 := fs.Children()
		children2 := fs.Children()

		// Should return the same underlying slice
		if &children1[0] != &children2[0] {
			t.Error("expected Children() to return reference to same slice")
		}
	})
}

func TestFieldSet_Render(t *testing.T) {
	t.Run("render empty fieldset", func(t *testing.T) {
		fs := FieldSet("Empty Legend")
		result := fs.Render()

		if result == template.HTML("") {
			t.Error("expected non-empty render result")
		}

		htmlStr := string(result)
		if htmlStr == "" {
			t.Error("expected non-empty HTML string")
		}

		// Should contain the legend
		if !strings.Contains(htmlStr, "Empty Legend") {
			t.Error("expected rendered HTML to contain legend text")
		}
	})

	t.Run("render fieldset with elements", func(t *testing.T) {
		textField := Text("username").SetAttributes(Attr("placeholder", "Enter username"))
		emailField := Email("email").SetAttributes(Required(true))

		fs := FieldSet("User Information", textField, emailField)
		result := fs.Render()

		if result == template.HTML("") {
			t.Error("expected non-empty render result")
		}

		htmlStr := string(result)
		if htmlStr == "" {
			t.Error("expected non-empty HTML string")
		}

		// Should contain the legend
		if !strings.Contains(htmlStr, "User Information") {
			t.Error("expected rendered HTML to contain legend text")
		}
	})

	t.Run("render returns template.HTML type", func(t *testing.T) {
		fs := FieldSet("Test")
		result := fs.Render()

		// Verify it's the correct type
		var _ template.HTML = result
	})
}

func TestFieldSet_ComplexScenarios(t *testing.T) {
	t.Run("fieldset with various element types", func(t *testing.T) {
		textField := Text("name")
		emailField := Email("email")
		passwordField := Password("password")
		checkboxField := Checkbox("agree")
		selectField := Select("country")
		textareaField := Textarea("description")

		fs := FieldSet("Complete Form",
			textField,
			emailField,
			passwordField,
			checkboxField,
			selectField,
			textareaField,
		)

		if len(fs.Children()) != 6 {
			t.Errorf("expected 6 children, got %d", len(fs.Children()))
		}

		// Verify each element type
		children := fs.Children()
		if children[0].(*element).template != "input" {
			t.Error("expected first element to be input template")
		}
		if children[1].(*element).template != "input" {
			t.Error("expected second element to be input template")
		}
		if children[2].(*element).template != "input" {
			t.Error("expected third element to be input template")
		}
		if children[3].(*element).template != "checkbox" {
			t.Error("expected fourth element to be checkbox template")
		}
		if children[4].(*element).template != "select" {
			t.Error("expected fifth element to be select template")
		}
		if children[5].(*element).template != "textarea" {
			t.Error("expected sixth element to be textarea template")
		}
	})

	t.Run("fieldset with configured elements", func(t *testing.T) {
		usernameField := Text("username").
			SetAttributes(
				Required(true),
				Attr("placeholder", "Enter username"),
				Attr("class", "form-control"),
			).
			SetLabel("Username").
			SetHint("Must be unique")

		emailField := Email("email").
			SetAttributes(
				Required(true),
				Attr("placeholder", "Enter email"),
			).
			SetLabel("Email Address").
			SetError("Invalid email format")

		fs := FieldSet("Account Details", usernameField, emailField)

		children := fs.Children()
		if len(children) != 2 {
			t.Errorf("expected 2 children, got %d", len(children))
		}

		// Check first element configuration
		firstElem := children[0].(*element)
		if firstElem.label != "Username" {
			t.Errorf("expected first element label=Username, got %s", firstElem.label)
		}
		if firstElem.hint != "Must be unique" {
			t.Errorf("expected first element hint=Must be unique, got %s", firstElem.hint)
		}
		if !firstElem.IsRequired() {
			t.Error("expected first element to be required")
		}

		// Check second element configuration
		secondElem := children[1].(*element)
		if secondElem.label != "Email Address" {
			t.Errorf("expected second element label=Email Address, got %s", secondElem.label)
		}
		if secondElem.error != "Invalid email format" {
			t.Errorf("expected second element error=Invalid email format, got %s", secondElem.error)
		}
		if !secondElem.IsRequired() {
			t.Error("expected second element to be required")
		}
	})
}

func TestFieldSet_ContainerInterface(t *testing.T) {
	t.Run("implements Container interface", func(t *testing.T) {
		fs := FieldSet("Test")

		// This should compile without issues if fieldSet implements Container
		var container Container = fs

		children := container.Children()
		if children == nil {
			t.Error("expected non-nil children from Container interface")
		}
	})
}

func TestFieldSet_EmptyLegend(t *testing.T) {
	t.Run("empty legend string", func(t *testing.T) {
		fs := FieldSet("")

		if fs.legend != "" {
			t.Errorf("expected empty legend, got %s", fs.legend)
		}

		result := fs.Render()
		if result == template.HTML("") {
			t.Error("expected non-empty render result even with empty legend")
		}
	})

	t.Run("whitespace legend", func(t *testing.T) {
		fs := FieldSet("   ")

		if fs.legend != "   " {
			t.Errorf("expected whitespace legend to be preserved, got %s", fs.legend)
		}
	})
}

func TestFieldSet_LegendHandling(t *testing.T) {
	tests := []struct {
		name     string
		legend   string
		expected string
	}{
		{"simple legend", "User Details", "User Details"},
		{"legend with spaces", "User Profile Information", "User Profile Information"},
		{"legend with special characters", "User's Profile & Settings", "User's Profile & Settings"},
		{"empty legend", "", ""},
		{"legend with unicode", "用户详情", "用户详情"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := FieldSet(tt.legend)
			if fs.legend != tt.expected {
				t.Errorf("expected legend=%s, got %s", tt.expected, fs.legend)
			}
		})
	}
}

func TestFieldSet_ElementOrder(t *testing.T) {
	t.Run("maintains element order", func(t *testing.T) {
		elements := []Renderer{
			Text("field1"),
			Email("field2"),
			Password("field3"),
			Checkbox("field4"),
			Select("field5"),
		}

		fs := FieldSet("Test", elements...)
		children := fs.Children()

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

func TestFieldSet_NilElements(t *testing.T) {
	t.Run("handles mixed nil and valid elements", func(t *testing.T) {
		validElement := Text("valid")

		fs := FieldSet("Test", validElement, nil, nil)

		children := fs.Children()
		if len(children) != 3 {
			t.Errorf("expected 3 children (including nils), got %d", len(children))
		}

		if children[0] != validElement {
			t.Error("expected first child to be valid element")
		}
		if children[1] != nil {
			t.Error("expected second child to be nil")
		}
		if children[2] != nil {
			t.Error("expected third child to be nil")
		}
	})
}
