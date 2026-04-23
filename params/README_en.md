# params package

The `params` package provides a struct tag-based declarative parameter validator.
It allows defining validation rules directly on struct fields, greatly simplifying data validation logic.

## Core Function `CheckFields`

The core function is `CheckFields(v interface{}) error`.
It takes a struct (or a pointer to a struct) as a parameter and recursively validates all fields based on the `check` tags on the fields.

- If all fields pass validation, the function returns `nil`.
- If any field fails validation, the function returns an `error` describing the first failure reason.

## Usage

Here is an example of validating an API request body:

```go
package main

import (
    "fmt"
    "github.com/l0neman/go-kit/params"
)

// Profile is a nested struct
type Profile struct {
    Score float64 `check:">0,<100"`
}

// CreateUserRequest defines the user's request body
// Use the `check` tag to define validation rules for each field
type CreateUserRequest struct {
    Username string  `check:"not_empty"`
    Age      int     `check:">=18,<120"`
    Email    string  `check:"not_empty"`
    Profile  *Profile // Supports nested struct pointers
}

func main() {
    // --- Case 1: A valid request ---
    validRequest := &CreateUserRequest{
        Username: "johndoe",
        Age:      30,
        Email:    "johndoe@example.com",
        Profile:  &Profile{Score: 99.5},
    }

    if err := params.CheckFields(validRequest); err != nil {
        fmt.Printf("Valid request failed validation: %v\n", err)
    } else {
        fmt.Println("Valid request passed validation!")
    }

    // --- Case 2: An invalid request (age doesn't match) ---
    invalidRequest := &CreateUserRequest{
        Username: "jane",
        Age:      17, // Age is less than 18
        Email:    "jane@example.com",
        Profile:  &Profile{Score: 80},
    }

    if err := params.CheckFields(invalidRequest); err != nil {
        // Error output will be very clear, e.g.: "field 'Age' value '17' does not satisfy condition '>=18'"
        fmt.Printf("Invalid request failed validation: %v\n", err)
    }
}
```

## Supported Rules

Multiple rules can be separated by commas `,`.

| Rule | Description | Type |
| :--- | :--- | :--- |
| `not_empty` | Field cannot be an empty string. | `string` |
| `is_empty` | Field must be an empty string. | `string` |
| `>value` | Must be greater than `value`. | Numeric types (`int`, `float64`, etc.) |
| `<value` | Must be less than `value`. | Numeric types |
| `>=value` | Must be greater than or equal to `value`. | Numeric types |
| `<=value` | Must be less than or equal to `value`. | Numeric types |
| `==value` | Must equal `value`. | Numeric types |
| `!=value` | Must not equal `value`. | Numeric types |