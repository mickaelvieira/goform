package goforms

import (
	"html/template"
)

func FieldSet(legend string, elements ...Element) *fieldSet {
	// Ensure elements is never nil, even when no elements are provided
	if elements == nil {
		elements = []Element{}
	}

	return &fieldSet{
		legend:   legend,
		elements: elements,
		template: parseTemplates(),
	}
}

type fieldSet struct {
	legend   string
	template TemplateRenderer
	elements []Element
}

func (f *fieldSet) Children() []Element {
	return f.elements
}

func (f *fieldSet) Render() template.HTML {
	return f.template.Render("fieldset.html", struct {
		Legend   string
		Elements []Element
	}{
		Legend:   f.legend,
		Elements: f.elements,
	})
}

var _ Container = (*fieldSet)(nil)
