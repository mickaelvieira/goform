package attr

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"time"
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

type Modifier func(attrs Attrs)
type Attrs map[string]any

func newModifier(name string, value any) Modifier {
	n := strings.ToLower(strings.TrimSpace(name))

	if !slices.Contains(attributes, n) && !isAria(n) && !isData(n) {
		panic(fmt.Sprintf("unsupported attribute %s", n))
	}

	switch n {
	case "id":
		return Id(value.(string))
	case "required":
		return Required(value.(bool))
	default:
		return func(attrs Attrs) {
			switch value := value.(type) {
			case bool:
				attrs[n] = value
			case string:
				attrs[n] = strings.TrimSpace(value)
			default:
				panic(fmt.Sprintf("unsupported type %T, only booleans & strings are supported", value))
			}
		}
	}
}

func (a Attrs) Set(name string, value any) Attrs {
	modifier := newModifier(name, value)
	modifier(a)
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

func Attributes(modifiers ...Modifier) Attrs {
	a := make(Attrs)
	for _, mod := range modifiers {
		mod(a)
	}
	return a
}

func Id(value string) Modifier {
	return func(attrs Attrs) {
		attrs["id"] = value
		attrs["aria-errormessage"] = value + "-error"
	}
}

func Required(value bool) Modifier {
	return func(attrs Attrs) {
		if value {
			attrs["aria-required"] = "true"
		} else {
			attrs["aria-required"] = "false"
		}
		attrs["required"] = value
	}
}

func Invalid(attrs Attrs) {
	attrs["aria-invalid"] = "true"
}

func Attr(name string, value any) Modifier {
	return newModifier(name, value)
}

func GenId() string {
	rand.NewSource(time.Now().UnixNano())

	l := 8
	c := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, l)
	for i := range b {
		b[i] = c[rand.Intn(len(c))]
	}
	return string(b)
}
