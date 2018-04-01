package domain

import "fmt"

type InvalidGtpSessionStateError struct {
	expected GtpSessionStatus
	current  GtpSessionStatus
}

func NewInvalidGtpSessionStateError(expected, current GtpSessionStatus) *InvalidGtpSessionStateError {
	return &InvalidGtpSessionStateError{expected, current}
}

func (e *InvalidGtpSessionStateError) Error() string {
	return fmt.Sprintf("The session's status is not %s, it is %s", e.expected.String(), e.current.String())
}
