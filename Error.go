package db

import "fmt"

// ----------------------------------------------------------------------
// ErrInvalidDataType
// ----------------------------------------------------------------------
type ErrInvalidDataType struct {
	Message string
}

// Error implements error.
func (e ErrInvalidDataType) Error() string {
	return fmt.Sprintf("ErrInvalidDataType: %s", e.Message)
}

func NewErrInvalidDataType(format string, args ...any) error {
	return &ErrInvalidDataType{
		Message: fmt.Sprintf(format, args...),
	}
}
