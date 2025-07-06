# Agent Instructions for GoForm

## Project Overview
GoForm is a Go library for building HTML forms with server-side validation and population from HTTP requests. It provides a fluent API for creating form elements, handling multipart file uploads, and populating Go structs from form data.

## Coding Standards

### Go Style
- Never write comments
- Follow standard Go conventions and idioms
- Use `gofmt` and `golangci-lint` for code formatting and linting
- All linter issues must be resolved (0 issues tolerance)
- Prefer explicit error handling over ignoring return values
- Use meaningful variable and function names

### Design Patterns
- **Functional Options Pattern**: Use for configurable constructors (e.g., `Form(options ...FormOption)`)
- **Builder Pattern**: Continue the fluent API style for method chaining
- **Interface Segregation**: Keep interfaces small and focused
- **Composition over Inheritance**: Use struct embedding and interfaces

### Error Handling
- Always handle errors explicitly - never use `_ = someFunction()`
- Return wrapped errors with context using `fmt.Errorf("description: %w", err)`
- Prefer returning errors over panicking
- Use error returns over internal error states when possible

### Naming Conventions
- Use PascalCase for exported types, functions, and methods
- Use camelCase for unexported items
- Use meaningful names that describe purpose
- Prefix interface names with capital letter (avoid "I" prefix)
- Use noun phrases for types, verb phrases for functions

## Git Workflow

### Commit Messages
Use Conventional Commits format:
```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

Types:
- `feat`: New features
- `fix`: Bug fixes
- `refactor`: Code refactoring without functionality changes
- `test`: Adding or updating tests
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `perf`: Performance improvements

### Breaking Changes
- Use `BREAKING CHANGE:` in commit footer for breaking changes
- Update version numbers appropriately
- Document migration path in commit message

## Testing Requirements

### Test Coverage
- Write comprehensive unit tests for all public APIs
- Test both happy paths and error conditions
- Use table-driven tests for multiple scenarios
- Test names should be descriptive: `TestFunctionName_Scenario`

### Test Organization
- Group related tests using `t.Run()` with descriptive names
- Use setup and teardown functions when needed
- Keep tests focused and independent
- Use meaningful assertions with clear error messages

### Running Tests
- Always run `go test -v ./...` before committing
- Ensure all tests pass before making changes
- Add new tests for new functionality
- Update existing tests when changing behavior

## Architecture Guidelines

### Package Organization
- Keep related functionality together in logical files
- Use descriptive file names (e.g., `form.go`, `element.go`, `attributes.go`)
- Avoid circular dependencies
- Export only what needs to be public

### Interface Design
- Define interfaces at the point of use
- Keep interfaces small and focused (Interface Segregation Principle)
- Use composition to build complex behaviors
- Prefer accepting interfaces, returning concrete types

### Configuration
- Use functional options pattern for optional configuration
- Provide sensible defaults
- Make configuration immutable after creation
- Document configuration options clearly

## HTML Form Specifics

### Form Elements
- Support all standard HTML input types
- Provide fluent API for setting attributes
- Handle validation consistently across element types
- Support accessibility attributes (aria-*, role, etc.)

### Data Handling
- Handle both URL-encoded and multipart form data
- Support file uploads with configurable memory limits
- Provide struct population from form data
- Validate data at appropriate points

### Template Rendering
- Use Go's `html/template` package for safety
- Escape user input appropriately
- Support custom templates
- Maintain clean separation between data and presentation

## Specific Instructions

### File Editing
- When making changes, read the full context first
- Use appropriate editing tools (`replace_string_in_file`, `insert_edit_into_file`)
- Include sufficient context in replacements (3-5 lines before/after)
- Avoid repeating existing code in edits

### Terminal Usage
- Use descriptive explanations for terminal commands
- Check command output for errors
- Run tests and linter after making changes
- Commit with appropriate conventional commit messages

### Documentation
- Comment exported functions and types
- Use GoDoc conventions for documentation
- Include usage examples in comments when helpful
- Keep comments concise but informative

### Dependencies
- Minimize external dependencies
- Use only standard library when possible
- Document any external dependencies and their purpose
- Consider alternatives before adding new dependencies

## Preferences

### Code Organization
- Group related functionality logically
- Keep functions focused and single-purpose
- Use early returns to reduce nesting
- Prefer explicit code over clever code

### API Design
- Maintain backward compatibility when possible
- Use consistent naming across the API
- Provide both simple and advanced usage patterns
- Make the common case easy, advanced cases possible

### Performance
- Consider memory allocations in hot paths
- Use appropriate data structures for the use case
- Profile when performance matters
- Don't optimize prematurely

## Notes
- This project emphasizes clean, idiomatic Go code
- Accessibility and web standards compliance are important
- The API should be intuitive for Go developers familiar with HTML forms
- Maintain high code quality standards throughout
