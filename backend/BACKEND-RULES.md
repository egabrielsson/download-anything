# Golang Backend Best Practices (Basics for Starters)

These are foundational rules for a clean, maintainable Go web API (e.g., using net/http, Gin, Echo, or Chi).

## Project Structure
Use a clean, layered structure to separate concerns:

project/
├── cmd/
│   └── server/
│       └── main.go          # Entry point, wires everything
├── internal/
│   ├── handlers/            # HTTP handlers / controllers
│   ├── services/            # Business logic
│   ├── repositories/        # Data access (DB, files, etc.)
│   ├── models/              # Domain structs / entities
│   └── config/              # Config loading
├── pkg/                     # Reusable public packages (if any)
├── migrations/              # DB migrations (optional)
└── go.mod

- Keep `internal/` private to your app.
- Group by feature (e.g., `internal/scraper/`) as the project grows.

## Error Handling
- Always check errors: `if err != nil { ... }`
- Use custom error types or wrap errors (`fmt.Errorf("failed to fetch: %w", err)`) for context.
- Return proper HTTP status codes (e.g., 400 Bad Request, 500 Internal Server Error).
- Use structured logging (e.g., zap or log/slog) with context.

## HTTP API Basics
- Use standard `net/http` or a lightweight router (Gin/Echo/Chi recommended for simplicity).
- Keep handlers thin: parse request → call service → respond.
- Validate inputs early (e.g., with ozzo-validation or simple checks).
- Use JSON for requests/responses: `json.NewDecoder(r.Body).Decode(&req)`

## Dependency Management & Config
- Use `go mod tidy` regularly.
- Load config from env vars + files (e.g., viper or built-in env).
- Prefer interfaces for dependencies (dependency injection via constructors).

## Testing
- Write unit tests for services/repositories (table-driven tests).
- Use `httptest` for handler tests.
- Aim for good coverage on core logic.

## Other Basics
- Graceful shutdown: catch signals and close servers cleanly.
- Use context.Context everywhere (requests, DB calls).
- Keep functions short and single-responsibility.
- Avoid global variables.
- Format code with `go fmt`, lint with golangci-lint.

Start simple, refactor as needed. Follow Go's "Effective Go" and "Go Proverbs" for more philosophy.

