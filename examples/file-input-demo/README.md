# File Input Demo

This example demonstrates the file upload capabilities of the goform library using dedicated template files and modern Pico CSS styling, showcasing both single and multiple file inputs with various configurations.

## Features Demonstrated

### Modern Template Architecture

- **Dedicated template files** - Clean separation with `form.tmpl` and `results.tmpl`
- **Pico CSS integration** - Beautiful, accessible styling out of the box
- **Progressive enhancement** - JavaScript validation with graceful fallbacks
- **Responsive design** - Optimized for all device sizes

### File Input Types

- **Single file input** - Upload individual files with type restrictions
- **Multiple file inputs** - Upload multiple files simultaneously
- **File type restrictions** - Control accepted file types using the `accept` attribute
- **Custom styling** - Enhanced presentation with CSS

### Form Configuration

- **Multipart encoding** - Proper form encoding for file uploads (`enctype="multipart/form-data"`)
- **POST method** - Required HTTP method for file uploads
- **File processing** - Server-side handling of uploaded files

### File Input Examples

```go
// Single file input with image restriction
goform.File("avatar").
    SetLabel("Profile Picture").
    SetAttributes(goform.Attr("accept", "image/*"))

// Multiple file input with document types
goform.File("documents").
    SetLabel("Documents (Multiple)").
    SetAttributes(
        goform.Attr("multiple", true),
        goform.Attr("accept", ".pdf,.doc,.docx"),
    )

// Multiple photo uploads
goform.File("photos").
    SetLabel("Photos (Multiple)").
    SetAttributes(
        goform.Attr("multiple", true),
        goform.Attr("accept", "image/*"),
    )
```

## Key Components

### Form Setup
- **File inputs** - Using `goform.File()` for file upload fields
- **Multiple attribute** - Enabling multiple file selection
- **Accept attribute** - Restricting file types (MIME types or extensions)
- **Proper encoding** - Setting `enctype="multipart/form-data"`

### Server Handling
- **Request processing** - Using `PopulateFromRequest()` to handle uploaded files
- **File data extraction** - Accessing uploaded file information
- **Response generation** - Displaying upload results to users

## File Type Restrictions

The example demonstrates different ways to restrict file types:

1. **MIME type patterns**: `image/*` (all image types)
2. **Specific extensions**: `.pdf,.doc,.docx` (document types)
3. **Combined restrictions**: Multiple MIME types or extensions

## Running the Example

```bash
cd examples/file-input-demo
go run main.go
```

Then visit `http://localhost:9000` in your browser to:

1. **Upload files** - Select single or multiple files using the form
2. **View results** - See the processed file information on the results page
3. **Test restrictions** - Try uploading different file types to see validation

## Implementation Details

### Form Structure
- Three different file input fields with various configurations
- Submit button to process the uploads
- Clean, accessible HTML structure

### Upload Processing
- Server creates a matching form structure for processing
- Uses `PopulateFromRequest()` to extract file data
- Displays results showing uploaded file information
- Provides navigation back to the form

### Security Considerations
- File type validation through `accept` attributes
- Server-side processing of multipart form data
- Clean separation between form display and upload handling

## Use Cases

This example is perfect for understanding:

1. **File upload forms** - Basic file upload functionality
2. **Multiple file handling** - Processing several files at once
3. **File type validation** - Restricting uploads to specific types
4. **Form processing** - Server-side handling of file uploads
5. **User feedback** - Displaying upload results and status

The demo provides a foundation for building robust file upload functionality in web applications using goform.
