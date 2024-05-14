package graph

import "github.com/google/uuid"

// Traverse traverses the graph starting from this node using the Depth-First Search algorithm
func (n *Node) Traverse() []*Node {
	var traversed []*Node
	traversed = traverse(n, traversed)
	return traversed
}

// traverse recursively traverses the graph using the Depth-First Search algorithm
func traverse(node *Node, traversed []*Node) []*Node {
	traversed = append(traversed, node)
	for _, child := range node.Children {
		if child.Id == "" {
			child.Id = uuid.New().String()
		}
		if child.Type == "" {
			child.Type = lexeme
		}
		if child.Properties == nil {
			child.Properties = make(map[string]string)
		}
		if child.Children == nil {
			child.Children = make([]*Node, 0)
		}
		traversed = traverse(child, traversed)
	}
	return traversed
}
