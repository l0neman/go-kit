# stringx package

The `stringx` package provides convenient helper functions for checking empty strings, as a supplement to Go's standard `strings` package.

## `IsEmpty(str string) bool`

The `IsEmpty` function checks whether the length of a string is 0. It is a simple alias for `len(str) == 0`.

**Note:** This function does not consider strings containing only whitespace characters (such as spaces, tabs) as empty.

## `HasEmpty(strings ...string) bool`

The `HasEmpty` function takes one or more strings as parameters, and returns `true` if **any** of the strings is empty (according to the definition of `IsEmpty`).

It is very useful for validating multiple input fields at once.

## Usage Example

```go
package main

import (
    "log"
    "github.com/l0neman/go-kit/stringx"
)

func main() {
    // --- IsEmpty example ---
    log.Printf("'' is empty: %v\n", stringx.IsEmpty("")) // true
    log.Printf("' ' is empty: %v\n", stringx.IsEmpty(" "))  // false
    log.Printf("'hello' is empty: %v\n", stringx.IsEmpty("hello")) // false

    log.Println("---")

    // --- HasEmpty example ---
    // Check a set of valid inputs
    validInputs := []string{"user", "pass123", "token"}
    if stringx.HasEmpty(validInputs...) {
        log.Println("Error: inputs contain empty fields")
    } else {
        log.Println("All inputs are valid") // Will print this
    }

    // Check a set of inputs containing empty strings
    invalidInputs := []string{"user", "", "token"} // The second one is an empty string
    if stringx.HasEmpty(invalidInputs...) {
        log.Println("Error: inputs contain empty fields") // Will print this
    } else {
        log.Println("All inputs are valid")
    }
}
```