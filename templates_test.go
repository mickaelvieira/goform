package goform

import (
	"html/template"
	"strings"
	"testing"
)

// Mock renderer for testing
type mockRenderer struct {
	html string
}

func (m *mockRenderer) Render() template.HTML {
	return template.HTML(m.html) //nolint:gosec // G203
}

func TestTemplateRenderer_render(t *testing.T) {
	t.Run("executes template successfully", func(t *testing.T) {
		tmpl := template.Must(template.New("test").Parse("<div>{{.}}</div>"))
		tr := &templateRenderer{}

		result := tr.render(tmpl, "hello")
		expected := template.HTML("<div>hello</div>")

		if result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})

	t.Run("returns error message on execution failure", func(t *testing.T) {
		// Template that will cause an error (accessing non-existent field)
		tmpl := template.Must(template.New("test").Parse("<div>{{.NonExistent.Field}}</div>"))
		tr := &templateRenderer{}

		result := tr.render(tmpl, "hello")
		resultStr := string(result)

		// Should contain error message
		if !strings.Contains(resultStr, "can't evaluate field NonExistent") {
			t.Errorf("expected error message about NonExistent field, got %s", resultStr)
		}
	})
}

func TestElementRenderer(t *testing.T) {
	t.Run("renders element correctly", func(t *testing.T) {
		renderer := &mockRenderer{html: "<input type=\"text\" name=\"test\">"}
		elementRendererFunc := componentRenderer()

		result := elementRendererFunc(renderer)
		expected := template.HTML("<input type=\"text\" name=\"test\">")

		if result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})

	t.Run("renders empty element", func(t *testing.T) {
		renderer := &mockRenderer{html: ""}
		elementRendererFunc := componentRenderer()

		result := elementRendererFunc(renderer)
		expected := template.HTML("")

		if result != expected {
			t.Errorf("expected empty string, got %s", result)
		}
	})
}

func TestFormRenderer(t *testing.T) {
	t.Run("returns element renderer function", func(t *testing.T) {
		formRendererFunc := FormRenderer()
		elementRendererFunc := componentRenderer()

		// Both should be the same type of function
		renderer := &mockRenderer{html: "<div>test</div>"}

		formResult := formRendererFunc(renderer)
		elementResult := elementRendererFunc(renderer)

		if formResult != elementResult {
			t.Errorf("FormRenderer should return same result as elementRenderer")
		}
	})
}

func TestAttributesRenderer(t *testing.T) {
	attributesRendererFunc := attributesRenderer()

	t.Run("renders empty attributes", func(t *testing.T) {
		attrs := map[string]any{}
		result, err := attributesRendererFunc(attrs)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != template.HTMLAttr("") {
			t.Errorf("expected empty string, got %s", result)
		}
	})

	t.Run("renders string attributes", func(t *testing.T) {
		attrs := map[string]any{
			"name":  "username",
			"class": "form-control",
			"id":    "user-field",
		}
		result, err := attributesRendererFunc(attrs)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Should be sorted alphabetically
		expected := template.HTMLAttr(`class="form-control" id="user-field" name="username"`)
		if result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})

	t.Run("renders boolean attributes", func(t *testing.T) {
		attrs := map[string]any{
			"required": true,
			"disabled": false,
			"readonly": true,
		}
		result, err := attributesRendererFunc(attrs)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Only true boolean attributes should appear, sorted alphabetically
		expected := template.HTMLAttr(`readonly required`)
		if result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})

	t.Run("skips empty string attributes", func(t *testing.T) {
		attrs := map[string]any{
			"name":        "username",
			"placeholder": "",
			"class":       "form-control",
			"title":       "",
		}
		result, err := attributesRendererFunc(attrs)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Empty strings should not appear
		expected := template.HTMLAttr(`class="form-control" name="username"`)
		if result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})

	t.Run("handles mixed attribute types", func(t *testing.T) {
		attrs := map[string]any{
			"name":     "username",
			"required": true,
			"disabled": false,
			"class":    "form-control",
			"readonly": true,
		}
		result, err := attributesRendererFunc(attrs)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Should include string and true boolean attributes, sorted
		expected := template.HTMLAttr(`class="form-control" name="username" readonly required`)
		if result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})

	t.Run("escapes HTML in attribute names and values", func(t *testing.T) {
		attrs := map[string]any{
			"data-test": "<script>alert('xss')</script>",
			"class":     "user&admin",
			"<onclick>": "malicious",
		}
		result, err := attributesRendererFunc(attrs)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		resultStr := string(result)

		// Should escape HTML entities
		if !strings.Contains(resultStr, "&lt;script&gt;") {
			t.Error("expected script tags to be escaped")
		}
		if !strings.Contains(resultStr, "user&amp;admin") {
			t.Error("expected ampersand to be escaped")
		}
		if !strings.Contains(resultStr, "&lt;onclick&gt;") {
			t.Error("expected attribute name to be escaped")
		}
	})

	t.Run("handles unsupported data types", func(t *testing.T) {
		attrs := map[string]any{
			"name":   "username",
			"number": 42,
			"slice":  []string{"a", "b"},
			"struct": struct{ X int }{X: 1},
		}
		result, err := attributesRendererFunc(attrs)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		resultStr := string(result)

		// Should contain error messages for unsupported types
		if !strings.Contains(resultStr, "attribute number has an unsupported data type int") {
			t.Error("expected error message for int type")
		}
		if !strings.Contains(resultStr, "attribute slice has an unsupported data type []string") {
			t.Error("expected error message for slice type")
		}
		if !strings.Contains(resultStr, "attribute struct has an unsupported data type struct") {
			t.Error("expected error message for struct type")
		}

		// Should still include valid attributes
		if !strings.Contains(resultStr, `name="username"`) {
			t.Error("expected valid attributes to still be included")
		}
	})

	t.Run("maintains consistent ordering", func(t *testing.T) {
		attrs := map[string]any{
			"z-index":    "5",
			"class":      "test",
			"id":         "element",
			"name":       "field",
			"aria-label": "Label",
		}

		// Run multiple times to ensure consistent ordering
		for i := 0; i < 5; i++ {
			result, err := attributesRendererFunc(attrs)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			expected := template.HTMLAttr(`aria-label="Label" class="test" id="element" name="field" z-index="5"`)
			if result != expected {
				t.Errorf("iteration %d: expected %s, got %s", i, expected, result)
			}
		}
	})
}
