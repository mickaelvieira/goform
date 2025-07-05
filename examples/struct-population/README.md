# Struct Population Example

This example demonstrates how to use the `Populate` method to populate Go structs from form data using `goform` tags, with dedicated template files and modern Pico CSS styling.

## Features Demonstrated

### Template-Based Architecture

- **Dedicated template files** - Clean separation with `form.tmpl` and `success.tmpl`
- **Pico CSS integration** - Beautiful, accessible styling framework
- **Progressive enhancement** - JavaScript validation with graceful fallbacks
- **Responsive design** - Optimized for all screen sizes

### `Populate(obj any)`

This method populates a Go struct from the current form values:

```go
type UserRegistration struct {
    Name      string   `goform:"name"`
    Email     string   `goform:"email"`
    Bio       string   `goform:"bio"`
    Avatar    string   `goform:"avatar"`
    Documents []string `goform:"documents"`
}

var user UserRegistration
form.Populate(&user)
```

### Template Structure

The example showcases:

- **Form template** - `form.tmpl` with Pico CSS styling and real-time validation
- **Success template** - `success.tmpl` displaying populated struct data beautifully
- **Consistent styling** - Modern, accessible design across all pages

## How It Works

1. **Define a struct** with `goform` tags that match your form field names
2. **Populate the form** from HTTP request using `PopulateFromRequest(r)`
3. **Populate the struct** from form data using `Populate(&struct)`
4. **Access the structured data** in your Go code

## File Handling

The `Populate` method handles files intelligently:

- **Single files**: The filename is stored as a string
- **Multiple files**: Filenames are parsed from comma-separated values into a `[]string` slice
- **Empty values**: Fields with no data are left at their zero value

## Running the Example

```bash
cd examples/struct-population
go run main.go
```

Then visit `http://localhost:9000` in your browser.

## Use Cases

This feature is particularly useful for:

- **Data binding** - Convert form data directly to Go structs
- **API endpoints** - Clean separation between form processing and business logic
- **File upload processing** - Handle single and multiple file uploads with proper typing
- **Validation preparation** - Get structured data ready for validation libraries
- **Database operations** - Convert form data to structs that match your database models
