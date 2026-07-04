package errors_test

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/Nigel2392/errors"
)

func TestNewAndErrorf(t *testing.T) {
	err := errors.New("Code1", "test message")
	if err.Code != "Code1" {
		t.Errorf("Expected code 'Code1', got '%s'", err.Code)
	}
	if err.Message != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", err.Message)
	}

	errFmt := errors.Errorf("Code2", "test %s", "message")

	// Errorf returns an error interface, so we need to cast it
	e, ok := errFmt.(errors.Error)
	if !ok {
		t.Fatalf("Expected Errorf to return an errors.Error type")
	}

	if e.Code != "Code2" {
		t.Errorf("Expected code 'Code2', got '%s'", e.Code)
	}
	if e.Message != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", e.Message)
	}
}

func TestErrorStringFormatting(t *testing.T) {
	// Test basic format
	err1 := errors.New("AuthFailed", "invalid credentials")
	expected1 := "AuthFailed: invalid credentials"
	if err1.Error() != expected1 {
		t.Errorf("Expected %q, got %q", expected1, err1.Error())
	}

	// Test missing code fallback
	err2 := errors.New("", "something went wrong")
	expected2 := fmt.Sprintf("%s: something went wrong", errors.CodeUnknown)
	if err2.Error() != expected2 {
		t.Errorf("Expected %q, got %q", expected2, err2.Error())
	}

	// Test with reason
	baseErr := errors.New("DB_ERROR", "query failed")
	sqlErr := stderrors.New("connection timeout")
	err3 := baseErr.WithCause(sqlErr)
	expected3 := "DB_ERROR: query failed: connection timeout"
	if err3.Error() != expected3 {
		t.Errorf("Expected %q, got %q", expected3, err3.Error())
	}
}

func TestWithCause(t *testing.T) {
	baseErr := errors.New("CodeA", "base error")
	cause1 := stderrors.New("cause 1")
	cause2 := stderrors.New("cause 2")

	// Add first cause
	errWithCause := baseErr.WithCause(cause1)
	if errWithCause.Reason != cause1 {
		t.Errorf("Expected reason to be %v, got %v", cause1, errWithCause.Reason)
	}

	// Add second cause (should join them based on WithCause logic)
	errWithMultiCause := errWithCause.WithCause(cause2)

	// Check if both causes are accessible
	if !stderrors.Is(errWithMultiCause.Reason, cause1) {
		t.Errorf("Expected joined reason to contain cause1")
	}
	if !stderrors.Is(errWithMultiCause.Reason, cause2) {
		t.Errorf("Expected joined reason to contain cause2")
	}
}

func TestWrap(t *testing.T) {
	baseErr := errors.New("CodeX", "original message").WithCause(stderrors.New("root cause"))

	wrapped1 := baseErr.Wrap("new message")
	if wrapped1.Code != "CodeX" || wrapped1.Message != "new message" || wrapped1.Reason.Error() != "root cause" {
		t.Errorf("Wrap failed to map fields correctly: %+v", wrapped1)
	}

	wrapped2 := baseErr.Wrapf("new message %d", 123)
	if wrapped2.Message != "new message 123" {
		t.Errorf("Wrapf failed to format message: %s", wrapped2.Message)
	}
}

func TestIs(t *testing.T) {
	baseErr := errors.New("ErrCode", "message")
	sameCodeErr := errors.New("ErrCode", "different message")

	// 1. Matches by Code
	if !errors.Is(baseErr, sameCodeErr) {
		t.Errorf("Expected baseErr to match sameCodeErr due to identical code")
	}

	// 2. Matches by Message (when codes don't trigger match, but equals() logic says if codes match it returns true.
	// If codes differ, it compares messages).
	errNoCode1 := errors.New("", "same message")
	errNoCode2 := errors.New("", "same message")
	if !errors.Is(errNoCode1, errNoCode2) {
		t.Errorf("Expected errors to match by message when codes are empty")
	}

	// 3. Matches by Reason
	rootCause := stderrors.New("root cause")
	errWithRoot := errors.New("CodeY", "msg").WithCause(rootCause)
	if !errors.Is(errWithRoot, rootCause) {
		t.Errorf("Expected errors.Is to match the underlying Reason")
	}

	// 4. Matches by Related
	relatedErr := stderrors.New("related issue")
	errWithRelated := errors.Error{
		Code:    "CodeZ",
		Message: "main issue",
		Related: []error{relatedErr},
	}
	if !errors.Is(errWithRelated, relatedErr) {
		t.Errorf("Expected errors.Is to match the related error")
	}
}

func TestUnwrapAndCause(t *testing.T) {
	rootCause := stderrors.New("root cause")
	err := errors.New("CodeA", "msg").WithCause(rootCause)

	if err.Unwrap() != rootCause {
		t.Errorf("Expected Unwrap to return rootCause")
	}

	if err.Cause() != rootCause {
		t.Errorf("Expected Cause to return rootCause")
	}
}

func TestPkgErrorsUtilities(t *testing.T) {
	// Tests for the wrapper functions in pkg_errors.go
	stdErr := stderrors.New("standard error")

	// Wrap
	wrapped := errors.Wrap(stdErr, "wrapped context")
	if wrapped == nil {
		t.Fatal("Wrap returned nil")
	}

	// Cause
	if errors.Cause(wrapped) != stdErr {
		t.Errorf("Cause failed to extract original error")
	}

	// Join
	err1 := stderrors.New("e1")
	err2 := stderrors.New("e2")
	joined := errors.Join(err1, err2)

	// As we are using standard errors.Join under the hood, it should wrap both
	if !errors.Is(joined, err1) || !errors.Is(joined, err2) {
		t.Errorf("Join failed to include all errors")
	}
}
