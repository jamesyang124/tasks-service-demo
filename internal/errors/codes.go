package errors

// Error codes for API responses
const (
	// Task related errors (1000-1999)
	ErrCodeTaskNotFound      = 1001
	ErrCodeTaskInvalidInput  = 1002
	ErrCodeTaskNameRequired  = 1003
	ErrCodeTaskNameTooLong   = 1004
	ErrCodeTaskInvalidStatus = 1005

	// Request related errors (2000-2999)
	ErrCodeInvalidJSON   = 2001
	ErrCodeInvalidID     = 2002
	ErrCodeMissingFields = 2003

	// System related errors (5000-5999)
	ErrCodeInternalError = 5001
	ErrCodeStorageError  = 5002
)
