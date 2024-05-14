package errors

type NodeNotFoundError struct {
	err string
}

func NewNodeNotFoundError(err string) *NodeNotFoundError {
	return &NodeNotFoundError{err: err}
}

func (e *NodeNotFoundError) Error() string {
	return e.err
}
