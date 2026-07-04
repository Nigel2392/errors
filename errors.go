package errors

import (
	"errors"
	"fmt"
)

type GoCode string

const (
	CodeUnknown GoCode = "Unknown"
)

type Error struct {
	Code    GoCode
	Reason  error
	Message string
	Related []error
}

func New(code GoCode, message string, related ...error) Error {
	return Error{
		Code:    code,
		Message: message,
		Related: related,
	}
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
func Errorf(code GoCode, format string, args ...interface{}) error {
	return Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

func (e Error) Error() string {
	var reasonStr string
	if e.Reason != nil {
		var reason = e.Reason.Error()
		var rBytes = make([]byte, len(reason)+2)
		rBytes[0] = ':'
		rBytes[1] = ' '
		copy(rBytes[2:], reason)
		reasonStr = string(rBytes)
	}

	var code = e.Code
	if code == "" {
		code = CodeUnknown
	}

	return fmt.Sprintf("%s: %s%s", code, e.Message, reasonStr)
}

func (e Error) WithCause(reason error) Error {
	if e.Reason != nil {
		return Error{
			Code:    e.Code,
			Message: e.Message,
			Reason:  errors.Join(e.Reason, reason),
			Related: e.Related,
		}
	}

	//	if diff, ok := reason.(Error); ok && diff.Code == e.Code &&
	//		// If the reason is already an Error with the same code, we can just return it.
	//		if e.Message != diff.Message && e.Message != "" {
	//			if diff.Reason == nil {
	//				diff.Reason = errors.New(e.Message)
	//			} else {
	//				diff.Reason = Wrap(diff.Reason, e.Message)
	//			}
	//		}
	//		return diff
	//	}

	return Error{
		Code:    e.Code,
		Message: e.Message,
		Reason:  reason,
		Related: e.Related,
	}
}

func (e Error) Wrap(message string) Error {
	return Error{
		Code:    e.Code,
		Message: message,
		Reason:  e.Reason,
		Related: e.Related,
	}
}

func (e Error) Wrapf(format string, args ...any) Error {
	return Error{
		Code:    e.Code,
		Message: fmt.Sprintf(format, args...),
		Reason:  e.Reason,
		Related: e.Related,
	}
}

func (e Error) equals(other Error) bool {
	// If the codes are the same, we consider them equal.
	if e.Code == other.Code && e.Code != "" {
		return true
	}
	return e.Message == other.Message
}

func (e Error) Is(chk error) bool {
	if e2, ok := chk.(Error); ok && e.equals(e2) {
		return true
	}
	if chk == nil {
		return false
	}

	if errors.Is(chk, e.Reason) {
		return true
	}

	for _, rel := range e.Related {
		if errors.Is(chk, rel) {
			return true
		}
	}

	return false
}

func (e Error) Unwrap() error {
	return e.Reason
}

func (e Error) Cause() error {
	return e.Reason
}
