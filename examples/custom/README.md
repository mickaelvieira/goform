# Custom Templates Example

This example demonstrates how to customize the rendering of goform components by overriding the default templates with your own custom templates.

## Features Demonstrated

### Custom Template Overrides

- **Template embedding** - Using Go's `embed` package to include custom templates
- **Template overriding** - Replacing default goform templates with custom ones
- **Custom rendering logic** - Creating personalized form layouts and styling

### Advanced Form Features

- **Form-level error messages** - Setting error messages at the form level
- **Custom template structure** - Organizing form components with custom HTML
- **Enhanced styling** - Custom CSS for improved visual presentation
- **Client-side validation** - JavaScript for interactive form validation

### Template Customization

The example shows how to override the default `form.tmpl` template:

```go
//go:embed form.tmpl
var formFS embed.FS

// Set custom templates to override defaults
goform.SetOverridingTemplates(formFS, "form.tmpl")
```

### Custom Form Template

The custom `form.tmpl` demonstrates:
- Accessing form children directly (`{{ $fieldset := index .Children 0 }}`)
- Custom HTML structure with CSS classes
- Manual placement of error messages
- Flexible component rendering

```gotmpl
{{ $fieldset := index .Children 0 }}
{{ $group := index .Children 1 }}

<div class="form-example">
  <form{{ if gt (len .Attributes) 0 }} {{ form_attributes .Attributes }}{{ end }}>
    {{ form_component $fieldset }}
    {{ .RenderError }}
    {{ form_component $group }}
  </form>
</div>
```

## Key Components

- **Embedded Templates** - Custom templates packaged with the binary
- **Template Override System** - Mechanism to replace default goform templates
- **Custom Styling** - Enhanced CSS with PicoCSS framework
- **Interactive Validation** - JavaScript-powered client-side validation
- **Form-level Errors** - Global error messages for the entire form

## Template Override Benefits

1. **Complete Control** - Full control over HTML structure and styling
2. **Brand Consistency** - Match your application's design system
3. **Enhanced UX** - Add custom interactions and animations
4. **Flexible Layout** - Organize form elements exactly as needed
5. **Progressive Enhancement** - Add JavaScript functionality while maintaining accessibility

## Running the Example

```bash
cd examples/custom
go run main.go
```

Then visit `http://localhost:9000` in your browser to see the custom-styled form with enhanced validation.

## Files Structure

- `main.go` - Main application with custom template setup
- `form.tmpl` - Custom form template that overrides the default
- `login.tmpl` - Complete HTML page with styling and JavaScript
- `go.mod` - Module dependencies

## Customization Tips

1. **Template Names** - Custom templates must match the original template names
2. **Template Functions** - Use goform's template functions like `form_attributes` and `form_component`
3. **Component Access** - Access form children using Go template syntax
4. **Styling** - Combine with CSS frameworks or custom styles for enhanced appearance
5. **JavaScript** - Add client-side functionality for better user experience

This example serves as a foundation for creating sophisticated, branded forms that integrate seamlessly with your application's design and functionality requirements.
