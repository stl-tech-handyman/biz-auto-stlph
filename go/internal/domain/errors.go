package domain

import "fmt"

// DomainError represents a domain-level error
type DomainError struct {
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// Common error codes
const (
	ErrCodeBusinessNotFound  = "BUSINESS_NOT_FOUND"
	ErrCodePipelineNotFound  = "PIPELINE_NOT_FOUND"
	ErrCodeInvalidInput       = "INVALID_INPUT"
	ErrCodeActionFailed       = "ACTION_FAILED"
	ErrCodeCriticalFailure    = "CRITICAL_FAILURE"
)

// NewDomainError creates a new domain error
func NewDomainError(code, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

