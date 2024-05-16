package graph

import (
	"backend/internal/graph/errors"
)

// FindTargetNodes returns the nodes to which the given node can be moved
func (n *Node) FindTargetNodes(node string) ([]*Node, error) {
	if node == "" {
		return nil, errors.NewIllegalArgumentError("node cannot be empty")
	}
	var candidates []*Node
	candidates = findTargetNodes(n, candidates, node)
	if len(n.Traverse()) == len(candidates) {
		// node is the tree's root
		return make([]*Node, 0), nil
	}
	return candidates, nil
}

// targetNodes recursively traverses the graph using the Depth-First Search algorithm
func findTargetNodes(node *Node, candidates []*Node, targetNode string) []*Node {
	candidates = append(candidates, &Node{Id: node.Id, Name: node.Name})
	for _, child := range node.Children {
		if child.Id != targetNode {
			candidates = findTargetNodes(child, candidates, targetNode)
		}
	}
	return candidates
}
