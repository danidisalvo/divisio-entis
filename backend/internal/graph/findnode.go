package graph

import (
	"backend/internal/graph/errors"
	"fmt"
)

// FindNode returns the node with the given id
func (n *Node) FindNode(id string) (*Node, error) {
	if id == "" {
		return nil, errors.NewIllegalArgumentError("id cannot be empty")
	}
	for _, node := range n.Traverse() {
		if node.Id == id {
			return node, nil
		}
	}
	return nil, errors.NewNodeNotFoundError(fmt.Sprintf("the node with ID %q was not found", id))
}
