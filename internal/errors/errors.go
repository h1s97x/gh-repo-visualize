package errors

import "fmt"

// ErrorCode represents an error code
type ErrorCode string

const (
	ErrCodeNotGitRepo    ErrorCode = "NOT_GIT_REPO"
	ErrCodeNoCommits     ErrorCode = "NO_COMMITS"
	ErrCodeInvalidBranch ErrorCode = "INVALID_BRANCH"
	ErrCodeGitFailed     ErrorCode = "GIT_FAILED"
	ErrCodeInvalidInput  ErrorCode = "INVALID_INPUT"
	ErrCodeNoAuthor      ErrorCode = "NO_AUTHOR"
	ErrCodeRenderFailed  ErrorCode = "RENDER_FAILED"
)

// Sentinel errors for common cases
var (
	ErrNotGitRepo    = NewError(ErrCodeNotGitRepo, "not a git repository (or any of the parent directories)", nil)
	ErrNoCommits     = NewError(ErrCodeNoCommits, "no commits found in repository", nil)
	ErrInvalidBranch = NewError(ErrCodeInvalidBranch, "branch not found", nil)
	ErrGitFailed     = NewError(ErrCodeGitFailed, "git command failed", nil)
	ErrInvalidInput  = NewError(ErrCodeInvalidInput, "invalid input parameter", nil)
	ErrNoAuthor      = NewError(ErrCodeNoAuthor, "no commits found for author", nil)
	ErrRenderFailed  = NewError(ErrCodeRenderFailed, "failed to render graph", nil)
)

// VisualizeError represents a structured error
type VisualizeError struct {
	Code    ErrorCode
	Message string
	Cause   error
	Detail  string
}

// Error implements the error interface
func (e *VisualizeError) Error() string {
	if e.Cause != nil {
		if e.Detail != "" {
			return fmt.Sprintf("%s: %s (%s, caused by: %v)", e.Code, e.Message, e.Detail, e.Cause)
		}
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause
func (e *VisualizeError) Unwrap() error {
	return e.Cause
}

// NewError creates a new VisualizeError
func NewError(code ErrorCode, message string, cause error) *VisualizeError {
	return &VisualizeError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// NewErrorWithDetail creates a new VisualizeError with detail
func NewErrorWithDetail(code ErrorCode, message string, detail string, cause error) *VisualizeError {
	return &VisualizeError{
		Code:    code,
		Message: message,
		Detail:  detail,
		Cause:   cause,
	}
}

// GitError creates a git command error
func GitError(command string, output string, cause error) *VisualizeError {
	return NewErrorWithDetail(ErrCodeGitFailed, "git command failed", command+": "+output, cause)
}

// NotGitRepo creates a not git repo error
func NotGitRepo(path string) *VisualizeError {
	return NewErrorWithDetail(ErrCodeNotGitRepo, "not a git repository", path, nil)
}

// NoCommits creates a no commits error
func NoCommits(branch string) *VisualizeError {
	if branch != "" {
		return NewErrorWithDetail(ErrCodeNoCommits, "no commits found", "branch: "+branch, nil)
	}
	return ErrNoCommits
}

// InvalidBranch creates an invalid branch error
func InvalidBranch(branch string) *VisualizeError {
	return NewErrorWithDetail(ErrCodeInvalidBranch, "branch not found", branch, nil)
}

// InvalidInput creates an invalid input error
func InvalidInput(param string, value string) *VisualizeError {
	return NewErrorWithDetail(ErrCodeInvalidInput, "invalid input parameter", param+"="+value, nil)
}

// IsGitError checks if error is a git error
func IsGitError(err error) bool {
	if e, ok := err.(*VisualizeError); ok {
		return e.Code == ErrCodeGitFailed
	}
	return false
}

// IsNotGitRepo checks if error is a not git repo error
func IsNotGitRepo(err error) bool {
	if e, ok := err.(*VisualizeError); ok {
		return e.Code == ErrCodeNotGitRepo
	}
	return false
}
