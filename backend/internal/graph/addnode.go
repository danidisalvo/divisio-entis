package graph

import (
	"backend/internal/graph/errors"
	"fmt"
)

// AddNode adds a node to the graph
func (n *Node) AddNode(parent string, newNode *Node) (*Node, error) {
	if newNode == nil {
		return nil, errors.NewIllegalArgumentError("newNode cannot be nil")
	}
	// the id must be unique
	for _, node := range n.Traverse() {
		if node.Id == newNode.Id {
			return nil, errors.NewDuplicatedNodeError(fmt.Sprintf("duplicated ID %q", newNode.Id))
		}
	}

	nodes := n.Traverse()
	for i, node := range nodes {
		if node.Id == parent {
			node.Children = append(node.Children, newNode)
			nodes[i] = node
			return n, nil
		}
	}
	return nil, errors.NewNodeNotFoundError(fmt.Sprintf("parent %q not found", parent))
}
