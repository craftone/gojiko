package apns

import "fmt"

type NoSuchAPNError struct{ error }

func NewNoSuchAPNError(format string, a ...interface{}) *NoSuchAPNError {
	return &NoSuchAPNError{fmt.Errorf(format, a...)}
}
