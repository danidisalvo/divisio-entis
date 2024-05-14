package errors

type ParsingError struct {
	err string
}

func NewParsingError(err string) *ParsingError {
	return &ParsingError{err: err}
}

func (e *ParsingError) Error() string {
	return e.err
}
