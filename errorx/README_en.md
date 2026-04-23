# errorx package

The `errorx` package provides convenient utility functions for error wrapping and formatting.

## Design Goals

When returning errors across multiple layers, if you directly return the underlying function's error each time, it often lacks context information from the upper layer's business logic.

For example, a database operation might return `sql.ErrNoRows`, but for the caller, it would be more helpful to see a complete information chain like "failed to get user info > record not found".

This package adds context to errors in a simple way while preserving the original error's type information, so upper-layer code can use `errors.Is` or `errors.As` for checking.

## Core Functions

### `Wrap(err error, desc string) error`

The `Wrap` function prepends a descriptive string `desc` to an existing error `err`, generating a new error.
It is equivalent to `fmt.Errorf("%s > %w", desc, err)`, "wrapping" the original `err` while preserving its type information.

#### Usage

```go
package main

import (
    "database/sql"
    "errors"
    "log"
    "fmt"
    "github.com/l0neman/go-kit/errorx"
)

// mock a data access function
func fetchUserFromDB(id int) error {
    // Simulate database record not found
    return sql.ErrNoRows
}

// service layer function
func getUser(id int) error {
    err := fetchUserFromDB(id)
    if err != nil {
        // Use Wrap to add context information
        return errorx.Wrap(err, fmt.Sprintf("failed to get user with id %d", id))
    }
    
    return nil
}

func main() {
    err := getUser(123)

    if err != nil {
        // Print the complete error chain
        // Output: failed to get user with id 123 > sql: no rows in result set
        log.Println(err)

        // Since Wrap preserves the original error, we can use errors.Is to check
        if errors.Is(err, sql.ErrNoRows) {
            log.Println("Confirmed: The user was not found in the database.")
        }
    }
}
```

### `Wrapf(err error, format string, a ...any) error`

The `Wrapf` function is similar to `Wrap`, but allows using a format string to generate the description. It is equivalent to `fmt.Errorf("%s > %w", fmt.Sprintf(format, a...), err)`.

#### Usage

```go
package main

import (
    "database/sql"
    "errors"
    "fmt"
    "log"
    "github.com/l0neman/go-kit/errorx"
)

func getUserWithAge(id int, age int) error {
    // Simulate database error
    err := sql.ErrNoRows
    if err != nil {
        // Use Wrapf to add context information, supports formatting
        return errorx.Wrapf(err, "failed to get user with id %d and age %d", id, age)
    }
    
    return nil
}

func main() {
    err := getUserWithAge(123, 30)

    if err != nil {
        // Print the complete error chain
        // Output: failed to get user with id 123 and age 30 > sql: no rows in result set
        log.Println(err)

        // Still can use errors.Is for type checking
        if errors.Is(err, sql.ErrNoRows) {
            log.Println("Confirmed: The user was not found in the database.")
        }
    }
}
```

### `Wraps(err error, desc string) string`

The `Wraps` function generates a formatted **string**, without error wrapping. It is equivalent to `fmt.Sprintf("%s > %v", desc, err)`.

When you only need an error description with context (for example, for logging), but don't need to preserve the original error type, you can use this function.

#### Usage

```go
package main

import (
    "fmt"
    "log"
    "github.com/l0neman/go-kit/errorx"
)

func main() {
    originalErr := fmt.Errorf("timeout")
    
    // Generate a string for logging
    logMessage := errorx.Wraps(originalErr, "network operation failed")
    
    // Output: network operation failed > timeout
    log.Println(logMessage)
}
```

### `Wrapfs(err error, format string, a ...any) string`

The `Wrapfs` function is similar to `Wraps`, but allows using a format string to generate the description. It is equivalent to `fmt.Sprintf("%s > %v", fmt.Sprintf(format, a...), err)`.

#### Usage

```go
package main

import (
    "fmt"
    "log"
    "github.com/l0neman/go-kit/errorx"
)

func main() {
    originalErr := fmt.Errorf("timeout")
    
    // Generate a formatted string for logging
    logMessage := errorx.Wrapfs(originalErr, "network operation failed for host %s", "example.com")
    
    // Output: network operation failed for host example.com > timeout
    log.Println(logMessage)
}
```