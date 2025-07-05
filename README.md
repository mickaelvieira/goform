# Go Form

A type-safe, template-driven Go package for building HTML forms with comprehensive validation, flexible styling, and clean API design.

## Features

- **HTTP Request Population** - Automatically populate forms from HTTP request data
- **Struct Population** - Populate Go structs from form data using struct tags
- **Type-safe form building** - Fluent API for creating forms with various input types
- **Template-driven rendering** - Customizable HTML templates for all form elements
- **Built-in validation** - Form and field-level validation with error handling
- **Accessibility support** - ARIA attributes and semantic HTML structure
- **File upload support** - Handle multipart forms and file uploads

## Quick Start

Examples can be found in this [directory](./examples/)

### Basic Form Creation

```go
form := goform.Form().
    AddChildren(
        goform.Text("name").SetLabel("Name"),
        goform.Email("email").SetLabel("Email"),
        goform.Select("country").SetLabel("Country").
            SetOptions(
                goform.Option("United States", "us"),
                goform.Option("Canada", "ca"),
            ),
    )
```

### HTTP Request Population

```go
// Populate form from HTTP request data
form.PopulateFromRequest(r)

// Validate the form
isValid, errors := form.IsValid()
if isValid {
    // Process successful form submission
} else {
    // Handle validation errors
    for field, err := range errors {
        fmt.Printf("Field %s: %s\n", field, err)
    }
}
```

### Struct Population

```go
// Define a struct with goform tags
type UserForm struct {
    Name      string   `goform:"name"`
    Email     string   `goform:"email"`
    Documents []string `goform:"documents"` // Handles multiple files
}

// Populate struct from form data
var user UserForm
form.Populate(&user)

// Access structured data
fmt.Printf("User: %+v\n", user)
```
