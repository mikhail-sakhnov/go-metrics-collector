package agent

import "fmt"

// ErrTooManyFailures error for the too many failures case
type ErrTooManyFailures struct {
	Agent string
}

// Error interface implementation
func (e ErrTooManyFailures) Error() string {
	return fmt.Sprintf("Agent `%s` has failed too many times", e.Agent)
}
