# Basic Form Example

This example demonstrates the fundamental features of the goform library with a simple login form using dedicated template files and Pico CSS for modern styling.

## Features Demonstrated

### Template-Based Rendering

- **Dedicated template files** - Clean separation between Go code and HTML templates
- **Pico CSS integration** - Modern, accessible styling with minimal CSS framework
- **Form renderer integration** - Using `goform.FormRenderer()` for seamless template integration
- **Responsive design** - Mobile-friendly forms with Pico's responsive features

### Basic Form Creation

- **Form structure** - Creating a form with fieldsets and groups
- **Input types** - Text inputs, password inputs, and buttons
- **Attributes** - Setting HTML attributes like `id`, `placeholder`, `required`, etc.
- **Labels and hints** - Adding descriptive text and help information
- **Error messages** - Displaying validation error messages

### Template Integration

- **Custom template rendering** - Using Go's `html/template` with goform
- **Form renderer** - Using the `goform.FormRenderer()` function for template rendering
- **Template data** - Passing form data to templates

### Template Structure

The example uses a dedicated `login.tmpl` file with:

- **Pico CSS** - Beautiful, semantic styling without custom CSS
- **Accessibility features** - ARIA attributes and semantic HTML
- **Client-side validation** - Progressive enhancement with JavaScript
- **Responsive layout** - Works perfectly on all device sizes

### Form Structure

```go
goform.Form().
    AddChildren(
        goform.FieldSet("Login",
            goform.Text("username").SetLabel("Username"),
            goform.Password("password").SetLabel("Password"),
        ),
        goform.Group(
            goform.Reset("reset"),
            goform.Submit("submit"),
        ),
    )
```

## Key Components

- **FieldSet** - Groups related form fields with a legend
- **Text Input** - Standard text input with validation
- **Password Input** - Secure password input
- **Reset Button** - Clears form data
- **Submit Button** - Submits the form
- **Group** - Organizes form elements without a legend

## Running the Example

```bash
cd examples/basic
go run main.go
```

Then visit `http://localhost:9000` in your browser.

## Template File

The example uses a custom template file (`login.tmpl`) that demonstrates how to:
- Render the complete form structure
- Apply custom styling
- Handle form display

## Use Cases

This basic example is perfect for:

- **Learning goform fundamentals** - Understanding the core concepts
- **Simple forms** - Login, contact, or registration forms
- **Template integration** - Seeing how goform works with Go templates
- **Getting started** - Your first goform implementation

## Next Steps

After understanding this basic example, explore:
- [Custom Example](../custom/) - Custom templates and styling
- [HTTP Request Example](../http-request/) - Form processing and validation
- [Struct Population Example](../struct-population/) - Data binding with structs
