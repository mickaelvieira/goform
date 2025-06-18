package goform

import (
	"html/template"
)

func Group(children ...Renderer) *group {
	if children == nil {
		children = []Renderer{}
	}

	return &group{
		children:   children,
		attributes: Attributes(),
		renderer:   getTemplateRenderer(),
	}
}

type group struct {
	class      string
	children   []Renderer
	attributes Attrs
	renderer   TemplateRenderer
}

func (g *group) SetAttributes(modifiers ...attrModifier) *group {
	for _, mod := range modifiers {
		mod(g.attributes)
	}
	return g
}

func (g *group) Attributes() Attrs {
	return g.attributes
}

func (g *group) Children() []Renderer {
	return g.children
}

func (g *group) Render() template.HTML {
	return g.renderer.Render("group.tmpl", g)
}

var _ Container = (*group)(nil)
