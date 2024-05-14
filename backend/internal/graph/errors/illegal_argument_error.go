package errors

type IllegalArgumentError struct {
	err string
}

func NewIllegalArgumentError(err string) *IllegalArgumentError {
	return &IllegalArgumentError{err: err}
}

func (e *IllegalArgumentError) Error() string {
	return e.err
}
