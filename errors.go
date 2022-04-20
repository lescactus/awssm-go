package main

import (
	"fmt"
)

// Key related errors
var (
	ErrKeyAlreadyExists  = &keyError{"Key already exists"}
	ErrKeyDoesNotExists  = &keyError{"Key does not exists"}
	ErrFlagSecretMissing = &flagError{"Missing flag '-secret'"}
	ErrFlagKeyMissing    = &flagError{"Missing flag '-key'"}
	ErrFlagValueMissing  = &flagError{"Missing flag '-value'"}
)

// Custom Error type for key related errors
type keyError struct {
	message string
}

// Custom Error type for flag related errors
type flagError struct {
	message string
}

// Implement the Error() method to satisfy the Error interface
func (err keyError) Error() string {
	return fmt.Sprintf("Error: %s in secret \"%s\": \"%s\"", err.message, *secretName, *secretKey)
}

// Implement the Error() method to satisfy the Error interface
func (err flagError) Error() string {
	return fmt.Sprintf("Error: %s", err.message)
}
