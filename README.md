# go-dbx

A lightweight Go database utility library that provides enhanced database operations with asynchronous support and automatic struct mapping.

## Features

- **Asynchronous Operations**: Execute database queries and transactions asynchronously
- **Generic Type Support**: Type-safe query results with Go generics
- **Automatic Struct Mapping**: Map query results to structs using reflection and `db` tags
- **Transaction Management**: Simplified transaction handling with automatic commit/rollback
- **Lightweight**: Minimal abstraction over Go's standard `database/sql` package
- **Context Support**: Full context support for cancellation and timeouts

## Installation

```bash
go get github.com/uoul/go-dbx
```

## Requirements

- Go 1.25.4 or higher 
- Dependencies: `github.com/uoul/go-async v1.0.0` 

## Core Interfaces

### IDbSession
Basic database session interface for executing queries:
```go
type IDbSession interface {
    QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}
```

### IDbConnection
Extended interface that includes transaction support:
```go
type IDbConnection interface {
    IDbSession
    BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}
```

## Usage Examples

### Basic Queries

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    
    "github.com/uoul/go-dbx/db"
    _ "github.com/go-sql-driver/mysql"
)

type User struct {
    ID    int    `db:"id"`
    Name  string `db:"name"`
    Email string `db:"email"`
}

func main() {
    // Open database connection
    database, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/dbname")
    if err != nil {
        log.Fatal(err)
    }
    defer database.Close()
    
    ctx := context.Background()
    
    // Execute synchronous query
    users, err := db.Query[User](ctx, database, "SELECT id, name, email FROM users WHERE active = ?", true)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d users\n", len(users))
    for _, user := range users {
        fmt.Printf("User: %s (%s)\n", user.Name, user.Email)
    }
}
```

### Asynchronous Queries

```go
// Execute asynchronous query
result := db.QueryAsync[User](ctx, database, "SELECT id, name, email FROM users")

// Do other work...

// Get the result when ready
users := <-result
if users.Error != nil {
    log.Fatal(err)
}

fmt.Printf("Async query returned %d users\n", len(users.Value))
```

### Struct Mapping with Tags

The library automatically maps database columns to struct fields using the `db` tag:

```go
type UserProfile struct {
    UserID      int    `db:"user_id"`
    FirstName   string `db:"first_name"`
    LastName    string `db:"last_name"`
    Email       string `db:"email"`
    CreatedAt   string `db:"created_at"`
}

// The library will automatically map columns to struct fields
profiles, err := db.Query[UserProfile](ctx, database, 
    "SELECT user_id, first_name, last_name, email, created_at FROM user_profiles")
```

### Nested Struct Support

The library supports nested structs with automatic field mapping:

```go
type Address struct {
    Street string `db:"street"`
    City   string `db:"city"`
    State  string `db:"state"`
}

type UserWithAddress struct {
    ID      int     `db:"id"`
    Name    string  `db:"name"`
    Address Address `db:"address_"` // Prefix for nested fields
}

// Maps columns like: id, name, address_street, address_city, address_state
```

## API Reference

### Query Functions

| Function | Description |
|----------|-------------|
| `Query[T any](ctx context.Context, session IDbSession, query string, args ...any) ([]T, error)` | Execute SQL query synchronously and return typed results |
| `QueryAsync[T any](ctx context.Context, session IDbSession, query string, args ...any) async.Result[[]T]` | Execute SQL query asynchronously |

### Transaction Functions

| Function | Description |
|----------|-------------|
| `ExecuteInTransaction(ctx context.Context, conn IDbConnection, opts *sql.TxOptions, fn TransactionScopeFunction) error` | Execute function within a database transaction with automatic commit/rollback |
| `ExecuteInTransactionAsync(ctx context.Context, conn IDbConnection, opts *sql.TxOptions, fn TransactionScopeFunction) async.Result[any]` | Execute transaction asynchronously |

## Error Handling

The library provides comprehensive error handling:
- Automatic transaction rollback on errors or panics
- Context cancellation support
- Proper resource cleanup

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## Acknowledgments

- Built with Go's standard `database/sql` package
- Utilizes the `github.com/uoul/go-async` library for asynchronous operations
- Designed for simplicity and performance

---
Made with ❤️ for the Go community