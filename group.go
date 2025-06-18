package goforms

import (
	"html/template"
)

func Group(elements ...Element) *group {
	return &group{
		elements: elements,
		template: parseTemplates(),
	}
}

type group struct {
	class    string
	template TemplateRenderer
	elements []Element
}

func (f *group) SetClass(class string) *group {
	f.class = class
	return f
}

func (f *group) Children() []Element {
	return f.elements
}

func (f *group) Render() template.HTML {
	return f.template.Render("group.html", struct {
		Class    string
		Elements []Element
	}{
		Class:    f.class,
		Elements: f.elements,
	})
}

var _ Container = (*group)(nil)
