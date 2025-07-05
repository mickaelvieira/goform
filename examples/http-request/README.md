# HTTP Request Population Example

This example demonstrates how to use the HTTP request population feature in the goform library with dedicated template files and modern Pico CSS styling.

## Features Demonstrated

### Template-Based Architecture

- **Dedicated template files** - Clean separation with `form.tmpl` and `success.tmpl`
- **Pico CSS integration** - Modern, accessible styling framework
- **Progressive enhancement** - Client-side validation with JavaScript
- **Responsive design** - Mobile-first, accessible forms

### `PopulateFromRequest(r *http.Request)`

This method automatically populates form elements from HTTP request data:

```go
form := goform.Form().AddChildren(
    goform.Text("name").SetLabel("Name"),
    goform.Email("email").SetLabel("Email"),
)

// Populate from HTTP request
form.PopulateFromRequest(r)
```

### `IsValid() (bool, map[string]string)`

This method validates the form and returns validation status with errors:

```go
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

## Features Supported

- **URL-encoded form data** (`application/x-www-form-urlencoded`)
- **Multipart form data** (`multipart/form-data`) including file uploads
- **Multiple file uploads** - When multiple files are uploaded to the same input, filenames are concatenated with commas
- **Automatic validation** after population
- **Error collection** with field-specific error messages

## Running the Example

```bash
cd examples/http-request
go run main.go
```

Then visit `http://localhost:9000` in your browser.

## Template Integration

The example uses dedicated templates for:

- **Form display** - `form.tmpl` with Pico CSS styling and JavaScript validation
- **Success page** - `success.tmpl` with consistent styling and navigation
- **Error handling** - Integrated error display within the form template

## How It Works

1. The form is created with various input types
2. On GET requests, an empty form is displayed
3. On POST requests:
   - `PopulateFromRequest()` extracts form data from the request
   - `IsValid()` performs validation and collects errors
   - Form elements are populated with submitted values
   - If valid, the user is redirected to a success page
   - If invalid, errors are displayed and the form is re-rendered with user input

## Use Cases

This feature is particularly useful for:

- **Server-side form processing** - Handle form submissions with clean separation of concerns
- **Form validation** - Validate forms independently of how they were populated
- **File uploads** - Handle multipart forms with file inputs
- **Repopulating forms** - Keep user input when validation fails
- **RESTful APIs** - Process form data in HTTP handlers
