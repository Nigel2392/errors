# github.com/Nigel2392/errors

A comprehensive utility package for error handling in Go. It provides a rich, structured error type that supports error codes, multiple related errors, and underlying causes, while seamlessly integrating with and extending the capabilities of both standard `errors` and `github.com/pkg/errors`.

## Features

- **Structured Errors:** The `Error` struct allows you to define errors with specific `GoCode` constants, clear messages, and underlying reasons.
- **Error Codes:** Easily categorize and identify errors by their code rather than relying strictly on string matching.
- **Cause Tracking:** Chain errors and preserve the original cause using `WithCause` or `Wrap`.
- **Related Errors:** Store multiple related errors alongside the primary error.
- **Drop-in Utility:** Provides top-level wrapper functions for common operations like `Is`, `As`, `Wrap`, `Unwrap`, `Join`, and `Cause` by leveraging `github.com/pkg/errors` and standard library `errors`.

## Installation

```bash
go get [github.com/Nigel2392/errors](https://github.com/Nigel2392/errors)
```

## Usage

### Creating Structured Errors

Create an error with a specific code and message:

```go
package main

import (
    "fmt"
    "os"
    myerrors "github.com/Nigel2392/errors"
)

const ErrNotFound myerrors.GoCode = "NotFound"

func main() {
    err := myerrors.New(ErrNotFound, "user could not be located in the database")
    fmt.Println(err.Error()) 
    err2 = err.Wrap("wrapped")
    fmt.Println(err.Error()) 
    fmt.Println(err2.Error()) 

    f, err := os.Open(...)
    if err != nil && errors.Is(err, os.ErrNotExist) {
        err = ErrNotFound.WithCause(err)
    }
}
```

### Adding Context and Causes

Attach underlying causes to your structured errors without losing the initial code context:

```go
package main

import (
    "database/sql"
    myerrors "github.com/Nigel2392/errors"
)

func getUser() error {
    err := sql.ErrNoRows
  
    // Create a base error, then attach the sql error as the reason
    baseErr := myerrors.New("DB_ERROR", "failed to fetch user")
    return baseErr.WithCause(err)
}
```

### Using Wrapper Utilities

You can use the package as a drop-in replacement for standard `errors` and `github.com/pkg/errors`:

```go
package main

import (
    "fmt"
    myerrors "github.com/Nigel2392/errors"
)

func main() {
    err1 := myerrors.New("ERR_1", "first error")
    err2 := myerrors.New("ERR_2", "second error")
  
    // Join multiple errors
    joined := myerrors.Join(err1, err2)
  
    // Wrap an error with a stack trace
    wrapped := myerrors.Wrap(joined, "additional context")
  
    fmt.Println(wrapped)
}
```
