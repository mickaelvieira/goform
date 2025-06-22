package goform

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"time"
)

const (
	AriaHintTemplate   = "%s-hint"
	AriaHintAttribute  = "aria-describedby"
	AriaErrorTemplate  = "%s-error"
	AriaErrorAttribute = "aria-errormessage"
)

var attributes = []string{
	"accept",
	"accept-charset",
	"accesskey",
	"action",
	"alt",
	"autocomplete",
	"autofocus",
	"checked",
	"class",
	"cols",
	"content",
	"contenteditable",
	"dir",
	"disabled",
	"draggable",
	"enctype",
	"hidden",
	"href",
	"id",
	"label",
	"lang",
	"list",
	"max",
	"maxlength",
	"method",
	"min",
	"minlength",
	"multiple",
	"name",
	"novalidate",
	"pattern",
	"placeholder",
	"readonly",
	"required",
	"role",
	"rows",
	"selected",
	"size",
	"spellcheck",
	"step",
	"style",
	"tabindex",
	"target",
	"title",
	"type",
	"value",
}

func isAria(name string) bool {
	return strings.HasPrefix(name, "aria-")
}

func isData(name string) bool {
	return strings.HasPrefix(name, "data-")
}

type attrModifier func(attrs Attrs)

func newModifier(name string, value any) attrModifier {
	n := strings.ToLower(strings.TrimSpace(name))

	if !slices.Contains(attributes, n) && !isAria(n) && !isData(n) {
		panic(fmt.Sprintf("unsupported attribute %s", n))
	}

	switch n {
	case "id":
		return Id(value.(string))
	default:
		return func(attrs Attrs) {
			switch value := value.(type) {
			case bool:
				attrs[n] = value
			case string:
				attrs[n] = strings.TrimSpace(value)
			default:
				panic(
					fmt.Sprintf("unsupported attribute %s type %T, only booleans & strings are supported", n, value),
				)
			}
		}
	}
}

type Attrs map[string]any

func (a Attrs) Set(name string, value any) Attrs {
	modifier := newModifier(name, value)
	modifier(a)
	return a
}

func (a Attrs) Unset(name string) Attrs {
	delete(a, name)
	return a
}

func (a Attrs) Get(name string) any {
	value, ok := a[name]
	if !ok {
		return ""
	}
	return value
}

func (a Attrs) String(name string) string {
	value, ok := a[name].(string)
	if !ok {
		return ""
	}
	return value
}

func (a Attrs) Bool(name string) bool {
	value, ok := a[name].(bool)
	if !ok {
		return false
	}
	return value
}

func Attributes(modifiers ...attrModifier) Attrs {
	a := make(Attrs)
	for _, mod := range modifiers {
		mod(a)
	}
	return a
}

func Id(id string) attrModifier {
	return func(attrs Attrs) {
		attrs["id"] = id

		// make sure ID-depend attributes are kept up to date
		_, hasError := attrs[AriaErrorAttribute]
		if hasError {
			attrs[AriaErrorAttribute] = fmt.Sprintf(AriaErrorTemplate, id)
		}

		_, hasHint := attrs[AriaHintAttribute]
		if hasHint {
			attrs[AriaHintAttribute] = fmt.Sprintf(AriaHintTemplate, id)
		}
	}
}

func Attr(name string, value any) attrModifier {
	return newModifier(name, value)
}

func GenId() string {
	rand.NewSource(time.Now().UnixNano())

	l := 10
	c := "abcdefghijklmnopqrstuvwxyz0123456789"

	b := make([]byte, l)
	for i := range b {
		b[i] = c[rand.Intn(len(c))] //nolint:gosec // G404: Use of weak random number generator
	}
	return string(b)
}
