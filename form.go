package goform

import (
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

func (f *form) Populate(obj any) *form {
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

	for _, element := range elements {
		if !element.IsValid() {
			element.MarkAsInvalid()
		}
	}

	return f
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
