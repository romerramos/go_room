# Bills Manager

A modern web application for managing bills built with Go, Gin, and HTMX.

## Features

- Clean architecture
- Server-side rendering with Go templates
- Modern UI with Tailwind CSS
- Interactive UI updates with HTMX
- Environment configuration support

## Prerequisites

- Go 1.21 or higher
- Git

## Project Structure

```
bills/
├── main.go           # Application entry point
├── templates/        # HTML templates
│   └── index.html    # Main template
├── static/          # Static assets
├── internal/        # Internal packages
│   ├── handlers/    # HTTP handlers
│   ├── models/      # Domain models
│   └── services/    # Business logic
└── .env            # Environment variables
```

## Getting Started

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Create a `.env` file:
   ```
   PORT=8080
   GIN_MODE=debug
   ```
4. Run the application:
   ```bash
   go run main.go
   ```
5. Visit `http://localhost:8080` in your browser

## Development

The application uses:
- Gin for routing and middleware
- HTMX for dynamic UI updates without JavaScript
- Tailwind CSS for styling
- Clean Architecture principles for maintainability

## Best Practices

- Environment configuration through `.env`
- Structured logging
- Clean separation of concerns
- Responsive design
- Progressive enhancement with HTMX 