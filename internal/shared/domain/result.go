package domain

// Result represents the outcome of an operation following the Result Pattern.
type Result[T any] struct {
	Data         *T
	Success      bool
	Message      string
	Error        error
	DetailsError string
	ErrorValids  map[string][]string
	Errors       []string
}

// Ok creates a successful Result with data and message.
func Ok[T any](data T, message string) Result[T] {
	return Result[T]{
		Data:    &data,
		Success: true,
		Message: message,
	}
}

// OkEmpty creates a successful Result without data payload.
func OkEmpty[T any](message string) Result[T] {
	return Result[T]{
		Success: true,
		Message: message,
	}
}

// Fail creates a failed Result with an error and message.
func Fail[T any](err error, message string) Result[T] {
	var details string
	if err != nil {
		details = err.Error()
	}
	return Result[T]{
		Success:      false,
		Message:      message,
		Error:        err,
		DetailsError: details,
		Errors:       []string{details},
	}
}

// FailWithValidation creates a failed Result due to input validation failures.
func FailWithValidation[T any](message string, validationErrors map[string][]string) Result[T] {
	var allErrors []string
	for _, errs := range validationErrors {
		allErrors = append(allErrors, errs...)
	}

	return Result[T]{
		Success:      false,
		Message:      message,
		Error:        ErrValidation,
		DetailsError: "validation failed",
		ErrorValids:  validationErrors,
		Errors:       allErrors,
	}
}
