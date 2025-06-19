package goform

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"sort"
	"strings"
	"sync"
)

//go:embed templates
var templateFS embed.FS
var templatesOptions *templatesOverrideOptions

type templatesOverrideOptions struct {
	filesystem fs.FS
	patterns   []string
}

func SetOverridingTemplates(filesystem fs.FS, patterns ...string) {
	templatesOptions = &templatesOverrideOptions{
		filesystem: filesystem,
		patterns:   patterns,
	}
}

type templateRenderer struct {
	base      *template.Template
	overwrite *template.Template
}

func (tr *templateRenderer) Render(name string, data any) template.HTML {
	if tr.overwrite != nil {
		if t := tr.overwrite.Lookup(name); t != nil {
			return tr.render(t, data)
		}
	}

	if tr.base != nil {
		if t := tr.base.Lookup(name); t != nil {
			return tr.render(t, data)
		}
	}

	//nolint:gosec // G203
	return template.HTML(fmt.Sprintf("template %s was not found", name))
}

func (tr *templateRenderer) render(t *template.Template, data any) template.HTML {
	var buf strings.Builder
	if err := t.Execute(&buf, data); err != nil {
		return template.HTML(fmt.Sprintf("%s", err)) //nolint:gosec // G203
	}
	return template.HTML(buf.String()) //nolint:gosec // G203
}

var getTemplateRenderer = sync.OnceValue(
	func() *templateRenderer {
		fn := template.FuncMap{
			"form_attributes": attributesRenderer(),
			"form_component":  componentRenderer(),
		}

		t, err := template.New("base").
			Funcs(fn).
			ParseFS(templateFS, "templates/*.tmpl")

		if err != nil {
			panic(err)
		}

		r := &templateRenderer{
			base: t,
		}

		if templatesOptions != nil {
			t, err := template.New("override").
				Funcs(fn).
				ParseFS(
					templatesOptions.filesystem,
					templatesOptions.patterns...,
				)

			if err != nil {
				panic(err)
			}

			r.overwrite = t
		}

		return r
	})

type Renderer interface {
	Render() template.HTML
}

type ErrorRenderer interface {
	RenderError() template.HTML
}

type HintRenderer interface {
	RenderHint() template.HTML
}

type TemplateRenderer interface {
	Render(name string, data any) template.HTML
}

func FormRenderer() func(f Renderer) template.HTML {
	return componentRenderer()
}

func componentRenderer() func(f Renderer) template.HTML {
	return func(f Renderer) template.HTML {
		return f.Render()
	}
}

func attributesRenderer() func(map[string]any) (template.HTMLAttr, error) {
	return func(attributes map[string]any) (template.HTMLAttr, error) {
		if len(attributes) == 0 {
			return template.HTMLAttr(""), nil
		}

		s := strings.Builder{}

		// Sort keys for consistent output
		keys := make([]string, 0, len(attributes))
		for name := range attributes {
			keys = append(keys, name)
		}
		sort.Strings(keys)

		for _, name := range keys {
			value := attributes[name]
			switch value := value.(type) {
			case bool:
				if value {
					s.WriteString(fmt.Sprintf("%s ", template.HTMLEscapeString(name)))
				}
			case string:
				if value != "" {
					s.WriteString(fmt.Sprintf(`%s="%s" `, template.HTMLEscapeString(name), template.HTMLEscapeString(value)))
				}
			default:
				s.WriteString(fmt.Sprintf(`attribute %s has an unsupported data type %T, only boolean & string are allowed`, name, value))
			}
		}
		//nolint:gosec // G203
		return template.HTMLAttr(strings.TrimSpace(s.String())), nil
	}
}
