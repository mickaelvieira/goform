package goform

import (
	"html/template"
)

func FieldSet(legend string, elements ...Renderer) *fieldSet {
	if elements == nil {
		elements = []Renderer{}
	}

	return &fieldSet{
		legend:     legend,
		children:   elements,
		attributes: Attributes(),
		renderer:   getTemplateRenderer(),
	}
}

type fieldSet struct {
	legend     string
	attributes Attrs
	renderer   TemplateRenderer
	children   []Renderer
}

func (f *fieldSet) SetAttributes(modifiers ...attrModifier) *fieldSet {
	for _, mod := range modifiers {
		mod(f.attributes)
	}
	return f
}

func (f *fieldSet) Attributes() Attrs {
	return f.attributes
}

func (f *fieldSet) Legend() string {
	return f.legend
}

func (f *fieldSet) Children() []Renderer {
	return f.children
}

func (f *fieldSet) Render() template.HTML {
	return f.renderer.Render("fieldset.tmpl", f)
}

var _ Container = (*fieldSet)(nil)
