package goforms

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/mickaelvieira/goforms/attr"
)

const (
	InputTypeTel           = "tel"
	InputTypeNumber        = "number"
	InputTypeSearch        = "search"
	InputTypeUrl           = "url"
	InputTypeColor         = "color"
	InputTypeRange         = "range"
	InputTypeDate          = "date"
	InputTypeDateTimeLocal = "datetime-local"
	InputTypeFile          = "file"
	InputTypeCheckbox      = "checkbox"
	InputTypeRadio         = "radio"
	InputTypeHidden        = "hidden"
	InputTypeSubmit        = "submit"
	InputTypeButton        = "button"
	InputTypeReset         = "reset"
	InputTypeImage         = "image"
	InputTypeTime          = "time"
	InputTypeMonth         = "month"
	InputTypeWeek          = "week"
	InputTypeDatetime      = "datetime"
	InputTypeText          = "text"
	InputTypeEmail         = "email"
	InputTypePassword      = "password"
)

const (
	SelectElement   = "select"
	TextareaElement = "textarea"
)

var defaultModifiers = []attr.Modifier{
	attr.Attr("aria-invalid", "false"),
	attr.Attr("aria-required", "false"),
}

func isInputType(t string) bool {
	return t != SelectElement && t != TextareaElement
}

func Phone(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeTel, modifiers...)
}

func Number(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeNumber, modifiers...)
}

func Search(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeSearch, modifiers...)
}

func Url(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeUrl, modifiers...)
}

func Color(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeColor, modifiers...)
}

func Range(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeRange, modifiers...)
}

func Date(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeDate, modifiers...)
}

func DateTimeLocal(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeDateTimeLocal, modifiers...)
}

func File(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeFile, modifiers...)
}

func Checkbox(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeCheckbox, modifiers...)
}

func Radio(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeRadio, modifiers...)
}

func Hidden(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeHidden, modifiers...)
}

func Submit(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeSubmit, modifiers...)
}

func Button(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeButton, modifiers...)
}

func Reset(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeReset, modifiers...)
}

func Image(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeImage, modifiers...)
}

func Time(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeTime, modifiers...)
}

func Month(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeMonth, modifiers...)
}

func Week(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeWeek, modifiers...)
}

func Datetime(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeDatetime, modifiers...)
}

func Text(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeText, modifiers...)
}

func Email(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypeEmail, modifiers...)
}

func Password(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, InputTypePassword, modifiers...)
}

func Textarea(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, TextareaElement, modifiers...)
}

func Select(name string, modifiers ...attr.Modifier) *element {
	return newElement(name, SelectElement, modifiers...)
}

func Option(label, value string) option {
	return option{
		Label: strings.TrimSpace(label),
		Value: strings.TrimSpace(value),
	}
}

type option struct {
	Label string
	Value string
}

type element struct {
	renderer   TemplateRenderer
	template   string
	Label      string
	Error      string
	Hint       string
	Attributes attr.Attrs
	Options    []option
}

func newElement(name, kind string, modifiers ...attr.Modifier) *element {
	t := kind
	a := attr.Attributes()

	if isInputType(kind) {
		if kind == InputTypeCheckbox {
			t = "checkbox"
		} else if kind == InputTypeRadio {
			t = "radio"
		} else if kind == InputTypeSubmit {
			t = "button"
		} else {
			t = "input"
		}
		a.Set("type", kind)
	}

	a.Set("name", name).
		Set("id", attr.GenId())

	i := &element{
		template:   t,
		renderer:   parseTemplates(),
		Attributes: a,
	}

	for _, mod := range defaultModifiers {
		mod(i.Attributes)
	}

	for _, mod := range modifiers {
		mod(i.Attributes)
	}

	return i
}

func (e *element) Render() template.HTML {
	return e.renderer.Render(fmt.Sprintf("%s.html", e.template), e)
}

func (e *element) RenderError() template.HTML {
	return e.renderer.Render("error.html", struct {
		ID      string
		Message string
	}{
		ID:      e.Id(),
		Message: e.Error,
	})
}

func (e *element) Id() string {
	return e.Attributes.String("id")
}

func (e *element) Name() string {
	return e.Attributes.String("name")
}

func (e *element) Attribute(name string) any {
	return e.Attributes.Get(name)
}

func (e *element) Value() string {
	return e.Attributes.String("value")
}

func (e *element) IsRequired() bool {
	return e.Attributes.Bool("required")
}

func (e *element) IsValid() bool {
	if !e.IsRequired() {
		return true
	}
	// @TODO check against a pattern if provided
	if e.Value() != "" {
		return true
	}
	return false
}

func (e *element) SetError(value string) *element {
	e.Error = strings.TrimSpace(value)
	return e
}

func (e *element) SetLabel(value string) *element {
	e.Label = strings.TrimSpace(value)
	return e
}

func (e *element) SetValue(value string) {
	e.Attributes.Set("value", value)
}

func (e *element) SetHint(value string) *element {
	e.Hint = strings.TrimSpace(value)
	return e
}

func (e *element) SetOptions(options ...option) *element {
	e.Options = make([]option, len(options))
	for i, opt := range options {
		e.Options[i] = option{
			Label: opt.Label,
			Value: opt.Value,
		}
	}
	return e
}

func (e *element) SetAttributes(modifiers ...attr.Modifier) *element {
	for _, mod := range modifiers {
		mod(e.Attributes)
	}
	return e
}

func (e *element) MarkAsInvalid() {
	attr.Invalid(e.Attributes)
}

var _ Element = (*element)(nil)
