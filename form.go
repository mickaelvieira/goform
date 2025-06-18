package goforms

import (
	"html/template"
	"maps"
	"net/http"
	"reflect"
	"slices"

	"github.com/mickaelvieira/goforms/attr"
)

const (
	MultipartData  = "multipart/form-data"
	URLEncodedData = "application/x-www-form-urlencoded"
)

type Element interface {
	Renderer
	ErrorRenderer
	Name() string
	IsValid() bool
	SetValue(string)
	MarkAsInvalid()
}

type Container interface {
	Children() []Element
}

// FormChild represents anything that can be added to a form
type FormChild interface {
	Renderer
}

type FormOptions struct {
	method  string
	enctype string
}

type modifiers func(*FormOptions)

func Multipart() modifiers {
	return func(o *FormOptions) {
		o.enctype = MultipartData
	}
}

func WithGetMethod() modifiers {
	return func(o *FormOptions) {
		o.method = http.MethodGet
	}
}

type form struct {
	template   TemplateRenderer
	elements   []Renderer
	attributes attr.Attrs
}

func Form(elements ...Renderer) *form {
	f := &form{
		elements: make([]Renderer, 0),
		template: parseTemplates(),
		attributes: attr.Attributes(
			attr.Attr("id", attr.GenId()),
			attr.Attr("method", http.MethodPost),
			attr.Attr("enctype", URLEncodedData),
		),
	}

	if elements == nil {
		elements = make([]Renderer, 0)
	}

	for _, el := range elements {
		if el == nil {
			continue
		}
		f.elements = append(f.elements, el)
	}

	return f
}

func (f *form) SetAttributes(modifiers ...attr.Modifier) *form {
	for _, mod := range modifiers {
		mod(f.attributes)
	}
	return f
}

func (f *form) AddElements(elements ...Renderer) *form {
	f.elements = append(f.elements, elements...)
	return f
}

func (f *form) Render() template.HTML {
	return f.template.Render("form.html", struct {
		Attributes attr.Attrs
		Elements   []Renderer
	}{
		Attributes: f.attributes,
		Elements:   f.elements,
	})
}

func (f *form) Populate(obj any) *form {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	elements := f.children()

	for i := range t.NumField() {
		field := t.Field(i)
		value := v.Field(i).String()
		name := field.Tag.Get("goforms")

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

func (f *form) Children() []Element {
	c := f.children()
	if len(c) == 0 {
		return make([]Element, 0)
	}
	return slices.Collect(maps.Values(c))
}

func (f *form) children() map[string]Element {
	elements := make(map[string]Element)

	for _, el := range f.elements {
		switch el := el.(type) {
		case Container:
			for _, ch := range el.Children() {
				elements[ch.Name()] = ch
			}
		case Element:
			elements[el.Name()] = el
		}
	}

	return elements
}
