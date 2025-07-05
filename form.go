package goform

import (
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"strings"
)

const (
	MultipartData  = "multipart/form-data"
	URLEncodedData = "application/x-www-form-urlencoded"
)

type Container interface {
	Children() []Renderer
}

type form struct {
	error      string
	children   []Renderer
	renderer   TemplateRenderer
	attributes Attrs
}

func Form() *form {
	f := &form{
		children: make([]Renderer, 0),
		renderer: getTemplateRenderer(),
		attributes: Attributes(
			Attr("id", GenId()),
			Attr("method", http.MethodPost),
			Attr("enctype", URLEncodedData),
		),
	}

	return f
}

func (f *form) Id() string {
	return f.attributes.String("id")
}

func (f *form) SetError(value string) *form {
	f.error = strings.TrimSpace(value)
	if f.error == "" {
		f.attributes.Unset(AriaErrorAttribute)
	} else {
		f.attributes.Set(AriaErrorAttribute, fmt.Sprintf(AriaErrorTemplate, f.Id()))
	}
	return f
}

func (f *form) Error() string {
	return f.error
}

func (f *form) SetAttributes(modifiers ...attrModifier) *form {
	for _, mod := range modifiers {
		mod(f.attributes)
	}
	return f
}

func (f *form) Attributes() Attrs {
	return f.attributes
}

func (f *form) AddChildren(children ...Renderer) *form {
	for _, c := range children {
		if c != nil {
			f.children = append(f.children, c)
		}
	}
	return f
}

func (f *form) Children() []Renderer {
	return f.children
}

func (f *form) Render() template.HTML {
	return f.renderer.Render("form.tmpl", f)
}

func (f *form) RenderError() template.HTML {
	return f.renderer.Render("error.tmpl", struct {
		Id    string
		Error string
	}{
		Id:    f.Id(),
		Error: f.error,
	})
}

func (f *form) PopulateFromStruct(obj any) *form {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	elements := f.Elements()

	for i := range t.NumField() {
		field := t.Field(i)
		value := v.Field(i).String()
		name := field.Tag.Get("goform")

		element, ok := elements[name]
		if ok {
			element.SetValue(value)
		}
	}

	// Mark invalid elements
	for _, element := range elements {
		if !element.IsValid() {
			element.MarkAsInvalid()
		}
	}

	return f
}

func (f *form) PopulateFromRequest(r *http.Request) *form {
	// Parse form data - this handles both URL-encoded and multipart forms
	if err := r.ParseForm(); err != nil {
		return f
	}

	// Also parse multipart form if present (for file uploads)
	if r.MultipartForm == nil {
		_ = r.ParseMultipartForm(32 << 20) // 32 MB max memory, ignore errors as it's optional
	}

	elements := f.Elements()

	// Iterate through all form values in the request
	for name, values := range r.Form {
		element, ok := elements[name]
		if !ok {
			continue
		}

		// For elements that can have multiple values (like checkboxes with same name),
		// we'll take the first value for now. This could be extended later.
		if len(values) > 0 {
			element.SetValue(values[0])
		}
	}

	// Handle file uploads separately if multipart form exists
	if r.MultipartForm != nil {
		for name, files := range r.MultipartForm.File {
			element, ok := elements[name]
			if !ok {
				continue
			}

			// For file inputs, handle multiple files
			if len(files) > 0 {
				if len(files) == 1 {
					// Single file: just set the filename
					element.SetValue(files[0].Filename)
				} else {
					// Multiple files: concatenate filenames with comma separation
					var filenames []string
					for _, file := range files {
						if file.Filename != "" {
							filenames = append(filenames, file.Filename)
						}
					}
					element.SetValue(strings.Join(filenames, ", "))
				}
			}
		}
	}

	// Validate all elements after population
	for _, element := range elements {
		if !element.IsValid() {
			element.MarkAsInvalid()
		}
	}

	return f
}

func (f *form) IsValid() (bool, map[string]string) {
	elements := f.Elements()
	errors := make(map[string]string)
	isValid := true

	// Collect validation errors from all elements
	for name, element := range elements {
		if !element.IsValid() {
			errors[name] = "Invalid value"
			isValid = false
		}
	}

	return isValid, errors
}

func (f *form) Elements() map[string]Element {
	elements := make(map[string]Element)

	for _, el := range f.children {
		switch el := el.(type) {
		case Container:
			// @TODO we should handle nested containers
			for _, c := range el.Children() {
				e, ok := c.(Element)
				if !ok {
					continue
				}
				elements[e.Name()] = e
			}
		case Element:
			elements[el.Name()] = el
		}
	}

	return elements
}

func (f *form) Populate(obj any) *form {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return f // obj must be a pointer to a struct
	}

	structValue := v.Elem()
	structType := structValue.Type()
	elements := f.Elements()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		// Skip unexported fields
		if !fieldValue.CanSet() {
			continue
		}

		// Get the goform tag
		tag := field.Tag.Get("goform")
		if tag == "" {
			continue
		}

		// Find the corresponding form element
		element, ok := elements[tag]
		if !ok {
			continue
		}

		// Get the value from the form element
		formValue := element.Value()
		if formValue == "" {
			continue
		}

		// Set the field value based on its type
		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(formValue)
		case reflect.Slice:
			// Handle file uploads or multiple values
			if fieldValue.Type().Elem().Kind() == reflect.String {
				// For []string, split comma-separated values (like multiple file names)
				if formValue != "" {
					values := strings.Split(formValue, ", ")
					// Clean up any empty strings
					cleanValues := make([]string, 0, len(values))
					for _, val := range values {
						val = strings.TrimSpace(val)
						if val != "" {
							cleanValues = append(cleanValues, val)
						}
					}
					fieldValue.Set(reflect.ValueOf(cleanValues))
				}
			}
		// Add more type handling as needed
		default:
			// For other types, try to set as string if the field type supports it
			if fieldValue.Type().ConvertibleTo(reflect.TypeOf("")) {
				fieldValue.Set(reflect.ValueOf(formValue).Convert(fieldValue.Type()))
			}
		}
	}

	return f
}
