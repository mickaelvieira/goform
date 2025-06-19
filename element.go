package goform

import (
	"fmt"
	"html/template"
	"strings"
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

var defaultModifiers = []attrModifier{
	Attr("aria-invalid", "false"),
	Attr("aria-required", "false"),
}

func isInputType(t string) bool {
	return t != SelectElement && t != TextareaElement
}

func Phone(name string) *element {
	return newElement(name, InputTypeTel)
}

func Number(name string) *element {
	return newElement(name, InputTypeNumber)
}

func Search(name string) *element {
	return newElement(name, InputTypeSearch)
}

func Url(name string) *element {
	return newElement(name, InputTypeUrl)
}

func Color(name string) *element {
	return newElement(name, InputTypeColor)
}

func Range(name string) *element {
	return newElement(name, InputTypeRange)
}

func Date(name string) *element {
	return newElement(name, InputTypeDate)
}

func DateTimeLocal(name string) *element {
	return newElement(name, InputTypeDateTimeLocal)
}

func File(name string) *element {
	return newElement(name, InputTypeFile)
}

func Checkbox(name string) *element {
	return newElement(name, InputTypeCheckbox)
}

func Radio(name string) *element {
	return newElement(name, InputTypeRadio)
}

func Hidden(name string) *element {
	return newElement(name, InputTypeHidden)
}

func Submit(name string) *element {
	return newElement(name, InputTypeSubmit)
}

func Button(name string) *element {
	return newElement(name, InputTypeButton)
}

func Reset(name string) *element {
	return newElement(name, InputTypeReset)
}

func Image(name string) *element {
	return newElement(name, InputTypeImage)
}

func Time(name string) *element {
	return newElement(name, InputTypeTime)
}

func Month(name string) *element {
	return newElement(name, InputTypeMonth)
}

func Week(name string) *element {
	return newElement(name, InputTypeWeek)
}

func Datetime(name string) *element {
	return newElement(name, InputTypeDatetime)
}

func Text(name string) *element {
	return newElement(name, InputTypeText)
}

func Email(name string) *element {
	return newElement(name, InputTypeEmail)
}

func Password(name string) *element {
	return newElement(name, InputTypePassword)
}

func Textarea(name string) *element {
	return newElement(name, TextareaElement)
}

func Select(name string) *element {
	return newElement(name, SelectElement)
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

type Element interface {
	Renderer
	HintRenderer
	ErrorRenderer
	Name() string
	IsValid() bool
	SetValue(string)
	MarkAsInvalid()
}

type element struct {
	hint       string
	label      string
	error      string
	template   string
	options    []option
	attributes Attrs
	renderer   TemplateRenderer
}

func newElement(name, kind string) *element {
	t := kind
	a := Attributes()

	if isInputType(kind) {
		switch kind {
		case InputTypeCheckbox:
			t = "checkbox"
		case InputTypeRadio:
			t = "radio"
		default:
			t = "input"
		}
		a.Set("type", kind)
	}

	a.Set("name", name).
		Set("id", GenId())

	i := &element{
		template:   t,
		renderer:   getTemplateRenderer(),
		attributes: a,
	}

	for _, mod := range defaultModifiers {
		mod(i.attributes)
	}

	return i
}

func (e *element) Render() template.HTML {
	return e.renderer.Render(fmt.Sprintf("%s.tmpl", e.template), e)
}

func (e *element) RenderError() template.HTML {
	return e.renderer.Render("error.tmpl", struct {
		Id    string
		Error string
	}{
		Id:    e.Id(),
		Error: e.error,
	})
}

func (e *element) RenderHint() template.HTML {
	return e.renderer.Render("hint.tmpl", struct {
		Id   string
		Hint string
	}{
		Id:   e.Id(),
		Hint: e.hint,
	})
}

func (e *element) Id() string {
	return e.attributes.String("id")
}

func (e *element) Name() string {
	return e.attributes.String("name")
}

func (e *element) Attribute(name string) any {
	return e.attributes.Get(name)
}

func (e *element) SetValue(value string) {
	e.attributes.Set("value", value)
}

func (e *element) Value() string {
	return e.attributes.String("value")
}

func (e *element) IsRequired() bool {
	return e.attributes.Bool("required")
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
	e.error = strings.TrimSpace(value)
	return e
}

func (e *element) Error() string {
	return e.error
}

func (e *element) SetLabel(value string) *element {
	e.label = strings.TrimSpace(value)
	return e
}

func (e *element) Label() string {
	return e.label
}

func (e *element) SetHint(value string) *element {
	e.hint = strings.TrimSpace(value)
	if e.hint == "" {
		e.attributes.Unset("aria-describedby")
	} else {
		e.attributes.Set("aria-describedby", fmt.Sprintf("%s-hint", e.Id()))
	}
	return e
}

func (e *element) Hint() string {
	return e.hint
}

func (e *element) SetOptions(options ...option) *element {
	e.options = make([]option, len(options))
	for i, opt := range options {
		e.options[i] = option{
			Label: opt.Label,
			Value: opt.Value,
		}
	}
	return e
}

func (e *element) Options() []option {
	return e.options
}

func (e *element) SetAttributes(modifiers ...attrModifier) *element {
	for _, mod := range modifiers {
		mod(e.attributes)
	}
	return e
}

func (e *element) Attributes() Attrs {
	return e.attributes
}

func (e *element) MarkAsInvalid() {
	Invalid(e.attributes)
}

var _ Element = (*element)(nil)
