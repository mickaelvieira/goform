package goforms

import (
	"html/template"
	"testing"

	"github.com/mickaelvieira/goforms/attr"
)

func TestGroup_Creation(t *testing.T) {
	t.Run("empty group", func(t *testing.T) {
		g := Group()

		if len(g.elements) != 0 {
			t.Errorf("expected 0 elements, got %d", len(g.elements))
		}
		if g.class != "" {
			t.Errorf("expected empty class, got %s", g.class)
		}
		if g.template == nil {
			t.Error("expected template to be initialized")
		}
	})

	t.Run("group with single element", func(t *testing.T) {
		textField := Text("username")
		g := Group(textField)

		if len(g.elements) != 1 {
			t.Errorf("expected 1 element, got %d", len(g.elements))
		}
		if g.elements[0] != textField {
			t.Error("expected first element to be the text field")
		}
	})

	t.Run("group with multiple elements", func(t *testing.T) {
		usernameField := Text("username")
		emailField := Email("email")
		passwordField := Password("password")

		g := Group(usernameField, emailField, passwordField)

		if len(g.elements) != 3 {
			t.Errorf("expected 3 elements, got %d", len(g.elements))
		}

		// Check elements are in correct order
		if g.elements[0] != usernameField {
			t.Error("expected first element to be username field")
		}
		if g.elements[1] != emailField {
			t.Error("expected second element to be email field")
		}
		if g.elements[2] != passwordField {
			t.Error("expected third element to be password field")
		}
	})
}

func TestGroup_SetClass(t *testing.T) {
	t.Run("set class on empty group", func(t *testing.T) {
		g := Group()
		result := g.SetClass("form-group")

		if g.class != "form-group" {
			t.Errorf("expected class=form-group, got %s", g.class)
		}
		if result != g {
			t.Error("SetClass should return the same group instance")
		}
	})

	t.Run("set class on group with elements", func(t *testing.T) {
		textField := Text("test")
		g := Group(textField)
		result := g.SetClass("input-group")

		if g.class != "input-group" {
			t.Errorf("expected class=input-group, got %s", g.class)
		}
		if result != g {
			t.Error("SetClass should return the same group instance")
		}
	})

	t.Run("override existing class", func(t *testing.T) {
		g := Group().SetClass("old-class")
		result := g.SetClass("new-class")

		if g.class != "new-class" {
			t.Errorf("expected class=new-class, got %s", g.class)
		}
		if result != g {
			t.Error("SetClass should return the same group instance")
		}
	})

	t.Run("set empty class", func(t *testing.T) {
		g := Group().SetClass("some-class")
		result := g.SetClass("")

		if g.class != "" {
			t.Errorf("expected empty class, got %s", g.class)
		}
		if result != g {
			t.Error("SetClass should return the same group instance")
		}
	})

	t.Run("set class with whitespace", func(t *testing.T) {
		g := Group()
		result := g.SetClass("  form-group  ")

		// Class should preserve whitespace as provided
		if g.class != "  form-group  " {
			t.Errorf("expected class with whitespace preserved, got %s", g.class)
		}
		if result != g {
			t.Error("SetClass should return the same group instance")
		}
	})
}

func TestGroup_Children(t *testing.T) {
	t.Run("empty group children", func(t *testing.T) {
		g := Group()
		children := g.Children()

		if len(children) != 0 {
			t.Errorf("expected 0 children, got %d", len(children))
		}
		// Unlike fieldSet, group doesn't guarantee non-nil for empty case
		// but let's check the actual behavior
	})

	t.Run("group with children", func(t *testing.T) {
		textField := Text("text")
		selectField := Select("select")
		g := Group(textField, selectField)

		children := g.Children()

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
		g := Group(textField)

		children1 := g.Children()
		children2 := g.Children()

		// Should return the same underlying slice
		if &children1[0] != &children2[0] {
			t.Error("expected Children() to return reference to same slice")
		}
	})
}

func TestGroup_Render(t *testing.T) {
	t.Run("render empty group", func(t *testing.T) {
		g := Group()
		result := g.Render()

		if result == template.HTML("") {
			t.Error("expected non-empty render result")
		}

		htmlStr := string(result)
		if htmlStr == "" {
			t.Error("expected non-empty HTML string")
		}
	})

	t.Run("render group with class", func(t *testing.T) {
		g := Group().SetClass("form-group")
		result := g.Render()

		if result == template.HTML("") {
			t.Error("expected non-empty render result")
		}

		htmlStr := string(result)
		if htmlStr == "" {
			t.Error("expected non-empty HTML string")
		}

		// Should contain the class (depending on template implementation)
		// Note: This test might need adjustment based on actual template content
	})

	t.Run("render group with elements", func(t *testing.T) {
		textField := Text("username", attr.Attr("placeholder", "Enter username"))
		emailField := Email("email", attr.Required(true))

		g := Group(textField, emailField).SetClass("user-fields")
		result := g.Render()

		if result == template.HTML("") {
			t.Error("expected non-empty render result")
		}

		htmlStr := string(result)
		if htmlStr == "" {
			t.Error("expected non-empty HTML string")
		}
	})

	t.Run("render returns template.HTML type", func(t *testing.T) {
		g := Group()
		result := g.Render()

		// Verify it's the correct type
		var _ template.HTML = result
	})
}

func TestGroup_ComplexScenarios(t *testing.T) {
	t.Run("group with various element types", func(t *testing.T) {
		textField := Text("name")
		emailField := Email("email")
		passwordField := Password("password")
		checkboxField := Checkbox("agree")
		selectField := Select("country")
		textareaField := Textarea("description")

		g := Group(
			textField,
			emailField,
			passwordField,
			checkboxField,
			selectField,
			textareaField,
		).SetClass("complete-form")

		if len(g.Children()) != 6 {
			t.Errorf("expected 6 children, got %d", len(g.Children()))
		}

		if g.class != "complete-form" {
			t.Errorf("expected class=complete-form, got %s", g.class)
		}

		// Verify each element type
		children := g.Children()
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

	t.Run("group with configured elements", func(t *testing.T) {
		usernameField := Text("username",
			attr.Required(true),
			attr.Attr("placeholder", "Enter username"),
			attr.Attr("class", "form-control"),
		).SetLabel("Username").SetHint("Must be unique")

		emailField := Email("email",
			attr.Required(true),
			attr.Attr("placeholder", "Enter email"),
		).SetLabel("Email Address").SetError("Invalid email format")

		g := Group(usernameField, emailField).SetClass("user-info")

		children := g.Children()
		if len(children) != 2 {
			t.Errorf("expected 2 children, got %d", len(children))
		}

		if g.class != "user-info" {
			t.Errorf("expected class=user-info, got %s", g.class)
		}

		// Check first element configuration
		firstElem := children[0].(*element)
		if firstElem.Label != "Username" {
			t.Errorf("expected first element label=Username, got %s", firstElem.Label)
		}
		if firstElem.Hint != "Must be unique" {
			t.Errorf("expected first element hint=Must be unique, got %s", firstElem.Hint)
		}
		if !firstElem.IsRequired() {
			t.Error("expected first element to be required")
		}

		// Check second element configuration
		secondElem := children[1].(*element)
		if secondElem.Label != "Email Address" {
			t.Errorf("expected second element label=Email Address, got %s", secondElem.Label)
		}
		if secondElem.Error != "Invalid email format" {
			t.Errorf("expected second element error=Invalid email format, got %s", secondElem.Error)
		}
		if !secondElem.IsRequired() {
			t.Error("expected second element to be required")
		}
	})
}

func TestGroup_ContainerInterface(t *testing.T) {
	t.Run("implements Container interface", func(t *testing.T) {
		g := Group()

		// This should compile without issues if group implements Container
		var container Container = g

		children := container.Children()
		if len(children) != 0 {
			t.Errorf("expected 0 children, got %d", len(children))
		}
	})

	t.Run("container interface with elements", func(t *testing.T) {
		textField := Text("test")
		g := Group(textField)

		var container Container = g
		children := container.Children()

		if len(children) != 1 {
			t.Errorf("expected 1 child, got %d", len(children))
		}
		if len(children) > 0 && children[0] != textField {
			t.Error("expected first child to be the text field")
		}
	})
}

func TestGroup_ElementOrder(t *testing.T) {
	t.Run("maintains element order", func(t *testing.T) {
		elements := []Element{
			Text("field1"),
			Email("field2"),
			Password("field3"),
			Checkbox("field4"),
			Select("field5"),
		}

		g := Group(elements...)
		children := g.Children()

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

func TestGroup_NilElements(t *testing.T) {
	t.Run("handles mixed nil and valid elements", func(t *testing.T) {
		validElement := Text("valid")

		// Create group with valid and nil elements
		g := Group(validElement, nil, nil)

		children := g.Children()
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

func TestGroup_Chaining(t *testing.T) {
	t.Run("method chaining", func(t *testing.T) {
		textField := Text("username")
		emailField := Email("email")

		g := Group(textField, emailField).
			SetClass("user-form").
			SetClass("updated-class") // Test overriding

		if g.class != "updated-class" {
			t.Errorf("expected class=updated-class, got %s", g.class)
		}

		if len(g.Children()) != 2 {
			t.Errorf("expected 2 children after chaining, got %d", len(g.Children()))
		}
	})
}

func TestGroup_ClassHandling(t *testing.T) {
	tests := []struct {
		name     string
		class    string
		expected string
	}{
		{"simple class", "form-group", "form-group"},
		{"class with spaces", "form group container", "form group container"},
		{"class with special characters", "form-group_v2", "form-group_v2"},
		{"empty class", "", ""},
		{"class with unicode", "表单组", "表单组"},
		{"class with leading/trailing spaces", "  form-group  ", "  form-group  "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Group().SetClass(tt.class)
			if g.class != tt.expected {
				t.Errorf("expected class=%s, got %s", tt.expected, g.class)
			}
		})
	}
}

func TestGroup_NilSafety(t *testing.T) {
	t.Run("group with nil elements slice", func(t *testing.T) {
		// This tests the behavior when elements variadic param is nil
		g := Group()

		// Should not panic when accessing children
		children := g.Children()

		// The exact behavior depends on how Go handles nil variadic params
		// but it should not panic
		if children != nil && len(children) != 0 {
			t.Errorf("expected empty or nil children, got %d", len(children))
		}
	})
}
