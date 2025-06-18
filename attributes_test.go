package goform

import (
	"strings"
	"testing"
)

func TestNewModifier_ValidAttributes(t *testing.T) {
	tests := []struct {
		name     string
		attr     string
		value    any
		expected map[string]any
	}{
		{
			name:  "string attribute",
			attr:  "class",
			value: "form-control",
			expected: map[string]any{
				"class": "form-control",
			},
		},
		{
			name:  "boolean attribute true",
			attr:  "disabled",
			value: true,
			expected: map[string]any{
				"disabled": true,
			},
		},
		{
			name:  "boolean attribute false",
			attr:  "disabled",
			value: false,
			expected: map[string]any{
				"disabled": false,
			},
		},
		{
			name:  "string with whitespace trimmed",
			attr:  "placeholder",
			value: "  Enter text  ",
			expected: map[string]any{
				"placeholder": "Enter text",
			},
		},
		{
			name:  "attribute name with whitespace trimmed",
			attr:  "  class  ",
			value: "test",
			expected: map[string]any{
				"class": "test",
			},
		},
		{
			name:  "uppercase attribute name",
			attr:  "CLASS",
			value: "test",
			expected: map[string]any{
				"class": "test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := make(Attrs)
			modifier := newModifier(tt.attr, tt.value)
			modifier(attrs)

			for key, expectedValue := range tt.expected {
				if attrs[key] != expectedValue {
					t.Errorf("expected %s=%v, got %v", key, expectedValue, attrs[key])
				}
			}
		})
	}
}

func TestNewModifier_AriaAttributes(t *testing.T) {
	tests := []struct {
		name     string
		attr     string
		value    any
		expected map[string]any
	}{
		{
			name:  "aria string attribute",
			attr:  "aria-label",
			value: "Close button",
			expected: map[string]any{
				"aria-label": "Close button",
			},
		},
		{
			name:  "aria boolean attribute",
			attr:  "aria-hidden",
			value: true,
			expected: map[string]any{
				"aria-hidden": true,
			},
		},
		{
			name:  "aria attribute case insensitive",
			attr:  "ARIA-LABEL",
			value: "test",
			expected: map[string]any{
				"aria-label": "test",
			},
		},
		{
			name:  "boolean remove attribute",
			attr:  "aria-hidden",
			value: false,
			expected: map[string]any{
				"aria-hidden": false,
			},
		},
		{
			name:  "aria attribute case insensitive",
			attr:  "aria-label",
			value: "",
			expected: map[string]any{
				"aria-label": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := make(Attrs)
			modifier := newModifier(tt.attr, tt.value)
			modifier(attrs)

			for key, expectedValue := range tt.expected {
				if attrs[key] != expectedValue {
					t.Errorf("expected %s=%v, got %v", key, expectedValue, attrs[key])
				}
			}
		})
	}
}

func TestNewModifier_DataAttributes(t *testing.T) {
	tests := []struct {
		name     string
		attr     string
		value    any
		expected map[string]any
	}{
		{
			name:  "data string attribute",
			attr:  "data-testid",
			value: "submit-button",
			expected: map[string]any{
				"data-testid": "submit-button",
			},
		},
		{
			name:  "data boolean attribute",
			attr:  "data-toggle",
			value: true,
			expected: map[string]any{
				"data-toggle": true,
			},
		},
		{
			name:  "data attribute case insensitive",
			attr:  "DATA-ID",
			value: "test",
			expected: map[string]any{
				"data-id": "test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := make(Attrs)
			modifier := newModifier(tt.attr, tt.value)
			modifier(attrs)

			for key, expectedValue := range tt.expected {
				if attrs[key] != expectedValue {
					t.Errorf("expected %s=%v, got %v", key, expectedValue, attrs[key])
				}
			}
		})
	}
}

func TestNewModifier_SpecialCases(t *testing.T) {
	t.Run("id attribute", func(t *testing.T) {
		attrs := make(Attrs)
		modifier := newModifier("id", "test-id")
		modifier(attrs)

		if attrs["id"] != "test-id" {
			t.Errorf("expected id=test-id, got %v", attrs["id"])
		}
		if attrs["aria-errormessage"] != "test-id-error" {
			t.Errorf("expected aria-errormessage=test-id-error, got %v", attrs["aria-errormessage"])
		}
	})

	t.Run("required true", func(t *testing.T) {
		attrs := make(Attrs)
		modifier := newModifier("required", true)
		modifier(attrs)

		if attrs["required"] != true {
			t.Errorf("expected required=true, got %v", attrs["required"])
		}
		if attrs["aria-required"] != "true" {
			t.Errorf("expected aria-required=true, got %v", attrs["aria-required"])
		}
	})

	t.Run("required false", func(t *testing.T) {
		attrs := make(Attrs)
		modifier := newModifier("required", false)
		modifier(attrs)

		if attrs["required"] != false {
			t.Errorf("expected required=false, got %v", attrs["required"])
		}
		if attrs["aria-required"] != "false" {
			t.Errorf("expected aria-required=false, got %v", attrs["aria-required"])
		}
	})
}

func TestNewModifier_UnsupportedAttribute(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for unsupported attribute")
		} else {
			// Check that panic message contains expected content
			panicMsg, ok := r.(string)
			if !ok {
				t.Error("expected panic to be a string")
			} else if !strings.Contains(panicMsg, "unsupported attribute") {
				t.Errorf("expected panic message to contain 'unsupported attribute', got: %s", panicMsg)
			}
		}
	}()

	newModifier("unsupported-attr", "value")
}

func TestNewModifier_UnsupportedType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for unsupported type")
		} else {
			// Check that panic message contains expected content
			panicMsg, ok := r.(string)
			if !ok {
				t.Error("expected panic to be a string")
			} else if !strings.Contains(panicMsg, "unsupported attribute class type") {
				t.Errorf("expected panic message to contain 'unsupported attribute class type', got: %s", panicMsg)
			}
		}
	}()

	attrs := make(Attrs)
	modifier := newModifier("class", 123) // int is unsupported
	modifier(attrs)
}

func TestAttrs_Set(t *testing.T) {
	attrs := make(Attrs)
	result := attrs.Set("class", "test-class")

	if result["class"] != "test-class" {
		t.Errorf("expected class=test-class, got %v", result["class"])
	}
	// Check that it's the same instance (returns itself)
	if len(attrs) == 0 {
		t.Error("Set should modify the original attrs instance")
	}
}

func TestAttrs_Get(t *testing.T) {
	attrs := Attrs{
		"class":    "test-class",
		"disabled": true,
	}

	if attrs.Get("class") != "test-class" {
		t.Errorf("expected class=test-class, got %v", attrs.Get("class"))
	}
	if attrs.Get("nonexistent") != "" {
		t.Errorf("expected empty string for nonexistent key, got %v", attrs.Get("nonexistent"))
	}
}

func TestAttrs_String(t *testing.T) {
	attrs := Attrs{
		"class":    "test-class",
		"disabled": true,
	}

	if attrs.String("class") != "test-class" {
		t.Errorf("expected class=test-class, got %v", attrs.String("class"))
	}
	if attrs.String("disabled") != "" {
		t.Errorf("expected empty string for non-string value, got %v", attrs.String("disabled"))
	}
	if attrs.String("nonexistent") != "" {
		t.Errorf("expected empty string for nonexistent key, got %v", attrs.String("nonexistent"))
	}
}

func TestAttrs_Bool(t *testing.T) {
	attrs := Attrs{
		"class":    "test-class",
		"disabled": true,
		"hidden":   false,
	}

	if attrs.Bool("disabled") != true {
		t.Errorf("expected disabled=true, got %v", attrs.Bool("disabled"))
	}
	if attrs.Bool("hidden") != false {
		t.Errorf("expected hidden=false, got %v", attrs.Bool("hidden"))
	}
	if attrs.Bool("class") != false {
		t.Errorf("expected false for non-bool value, got %v", attrs.Bool("class"))
	}
	if attrs.Bool("nonexistent") != false {
		t.Errorf("expected false for nonexistent key, got %v", attrs.Bool("nonexistent"))
	}
}

func TestAttributes(t *testing.T) {
	attrs := Attributes(
		Id("test-id"),
		Required(true),
		Attr("class", "form-control"),
	)

	expected := map[string]any{
		"id":                "test-id",
		"aria-errormessage": "test-id-error",
		"required":          true,
		"aria-required":     "true",
		"class":             "form-control",
	}

	for key, expectedValue := range expected {
		if attrs[key] != expectedValue {
			t.Errorf("expected %s=%v, got %v", key, expectedValue, attrs[key])
		}
	}
}

func TestId(t *testing.T) {
	attrs := make(Attrs)
	modifier := Id("test-id")
	modifier(attrs)

	if attrs["id"] != "test-id" {
		t.Errorf("expected id=test-id, got %v", attrs["id"])
	}
	if attrs["aria-errormessage"] != "test-id-error" {
		t.Errorf("expected aria-errormessage=test-id-error, got %v", attrs["aria-errormessage"])
	}
}

func TestRequired(t *testing.T) {
	t.Run("required true", func(t *testing.T) {
		attrs := make(Attrs)
		modifier := Required(true)
		modifier(attrs)

		if attrs["required"] != true {
			t.Errorf("expected required=true, got %v", attrs["required"])
		}
		if attrs["aria-required"] != "true" {
			t.Errorf("expected aria-required=true, got %v", attrs["aria-required"])
		}
	})

	t.Run("required false", func(t *testing.T) {
		attrs := make(Attrs)
		modifier := Required(false)
		modifier(attrs)

		if attrs["required"] != false {
			t.Errorf("expected required=false, got %v", attrs["required"])
		}
		if attrs["aria-required"] != "false" {
			t.Errorf("expected aria-required=false, got %v", attrs["aria-required"])
		}
	})
}

func TestInvalid(t *testing.T) {
	attrs := make(Attrs)
	Invalid(attrs)

	if attrs["aria-invalid"] != "true" {
		t.Errorf("expected aria-invalid=true, got %v", attrs["aria-invalid"])
	}
}

func TestAttr(t *testing.T) {
	attrs := make(Attrs)
	modifier := Attr("class", "test-class")
	modifier(attrs)

	if attrs["class"] != "test-class" {
		t.Errorf("expected class=test-class, got %v", attrs["class"])
	}
}

func TestGenId(t *testing.T) {
	id1 := GenId()
	id2 := GenId()

	if len(id1) != 10 {
		t.Errorf("expected id length 10, got %d", len(id1))
	}
	if len(id2) != 10 {
		t.Errorf("expected id length 10, got %d", len(id2))
	}
	// Note: There's a small chance they could be the same, but very unlikely
	if id1 == id2 {
		t.Error("expected different ids, got the same (this could rarely happen by chance)")
	}

	// Check that all characters are valid
	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for _, char := range id1 {
		if !strings.ContainsRune(validChars, char) {
			t.Errorf("id contains invalid character: %c", char)
		}
	}
}

func TestIsAria(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"aria prefix", "aria-label", true},
		{"ARIA prefix uppercase", "ARIA-label", false}, // function is case sensitive
		{"no aria prefix", "class", false},
		{"partial match", "not-aria-label", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAria(tt.input)
			if result != tt.expected {
				t.Errorf("isAria(%s) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsData(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"data prefix", "data-testid", true},
		{"DATA prefix uppercase", "DATA-testid", false}, // function is case sensitive
		{"no data prefix", "class", false},
		{"partial match", "not-data-testid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isData(tt.input)
			if result != tt.expected {
				t.Errorf("isData(%s) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
