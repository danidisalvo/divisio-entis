package graph

type (
	DuplicatedNodeError struct {
		err string
	}

	IllegalArgumentError struct {
		err string
	}

	NodeNotFoundError struct {
		err string
	}

	ParsingError struct {
		err string
	}
)

func (e *DuplicatedNodeError) Error() string {
	return e.err
}

func (e *IllegalArgumentError) Error() string {
	return e.err
}

func (e *NodeNotFoundError) Error() string {
	return e.err
}

func (e *ParsingError) Error() string {
	return e.err
}
