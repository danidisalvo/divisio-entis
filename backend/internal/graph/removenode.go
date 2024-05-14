package graph

import (
	"backend/internal/graph/errors"
	"fmt"
)

// RemoveNode removes a node from the graph
func (n *Node) RemoveNode(parent, target string) (*Node, error) {
	parentFound := false
	targetFound := false
	nodes := n.Traverse()
	for i, parentNode := range nodes {
		if parentNode.Id == parent {
			parentFound = true
			children := make([]*Node, 0)
			for _, child := range parentNode.Children {
				if child.Id == target {
					targetFound = true
				} else {
					children = append(children, child)
				}
			}
			nodes[i].Children = children
			break
		}
	}
	if !parentFound {
		return nil, errors.NewNodeNotFoundError(fmt.Sprintf("the parent node with ID %q was not found", parent))
	}
	if !targetFound {
		return nil, errors.NewNodeNotFoundError(fmt.Sprintf("the target node with ID %q was not found", target))
	}
	return n, nil
}
