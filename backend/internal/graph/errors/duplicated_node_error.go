package errors

type DuplicatedNodeError struct {
	err string
}

func NewDuplicatedNodeError(err string) *DuplicatedNodeError {
	return &DuplicatedNodeError{err: err}
}

func (e *DuplicatedNodeError) Error() string {
	return e.err
}
