package errorx

import (
	"fmt"
)

// Wrap wraps an existing err into a new err
// Provides a convenient concatenation function without manual template code
// Equivalent to fmt.Errorf("%s > %w", title, err)
func Wrap(err error, desc string) error {
	return fmt.Errorf("%s > %w", desc, err)
}

// Wrapf wraps an existing err into a new err
// Provides a convenient concatenation function without manual template code
// Equivalent to fmt.Errorf("%s > %w", fmt.Sprintf(format, a...), err)
func Wrapf(err error, format string, a ...any) error {
	return fmt.Errorf("%s > %w", fmt.Sprintf(format, a...), err)
}

// Wraps returns the wrapped version of the error message as a string
func Wraps(err error, desc string) string {
	return fmt.Sprintf("%s > %v", desc, err)
}

// Wrapfs returns the Wrapf version of the error message as a string
func Wrapfs(err error, format string, a ...any) string {
	return fmt.Sprintf("%s > %v", fmt.Sprintf(format, a...), err)
}
