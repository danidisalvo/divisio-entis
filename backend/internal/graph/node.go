package graph

import (
	"encoding/json"
	"fmt"
)

const defaultColor = "#dddddd"

// Node represents a node of a graph which can be traversed using the Depth-First Search algorithm
type Node struct {
	Name     string  `json:"name"`
	ReadOnly bool    `json:"readOnly"`
	Color    string  `json:"color"`
	Notes    string  `json:"notes"`
	Children []*Node `json:"children"`
}

// NewNode creates a new node
func NewNode(name string, readOnly bool, color, notes string) (*Node, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if color == "" {
		color = defaultColor
	}
	return &Node{
		Name:     name,
		ReadOnly: readOnly,
		Color:    color,
		Notes:    notes,
	}, nil
}

// String returns the JSON representation of this node
func (n *Node) String() (string, error) {
	bytes, err := json.Marshal(n)
	if err != nil {
		return "", fmt.Errorf("failed to marshal the node [%s]", err)
	}
	return string(bytes), nil
}

func (n *Node) Parse(bytes []byte) (*Node, error) {
	err := json.Unmarshal(bytes, n)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the node [%s]", err)
	}
	return n, nil
}

// FindNode returns the node with the given name
func (n *Node) FindNode(name string) (*Node, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	for _, node := range n.Traverse() {
		if node.Name == name {
			return node, nil
		}
	}
	return nil, fmt.Errorf("the node %q was not found", name)
}

// AddNode adds a node to the graph
func (n *Node) AddNode(parent string, newNode *Node) (*Node, error) {
	if newNode == nil {
		return nil, fmt.Errorf("newNode cannot be nil")
	}
	// the name must be unique
	for _, node := range n.Traverse() {
		if node.Name == newNode.Name {
			return nil, fmt.Errorf("duplicated name %q", newNode.Name)
		}
	}

	nodes := n.Traverse()
	for i, node := range nodes {
		if node.Name == parent {
			node.Children = append(node.Children, newNode)
			nodes[i] = node
			return n, nil
		}
	}
	return nil, fmt.Errorf("parent %q not found", parent)
}

// RemoveNode removes a node to the graph
func (n *Node) RemoveNode(parent, target string) (*Node, error) {
	parentFound := false
	targetFound := false
	nodes := n.Traverse()
	for i, parentNode := range nodes {
		if parentNode.Name == parent {
			parentFound = true
			children := make([]*Node, 0)
			for _, child := range parentNode.Children {
				if child.Name == target {
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
		return nil, fmt.Errorf("the parent node %q was not found", parent)
	}
	if !targetFound {
		return nil, fmt.Errorf("the target node %q was not found", target)
	}
	return n, nil
}

// UpdateNode updates a graph's node
func (n *Node) UpdateNode(parent string, targetNode *Node) (*Node, error) {
	if targetNode == nil {
		return nil, fmt.Errorf("targetNode cannot be nil")
	}
	parentFound := false
	targetFound := false
	nodes := n.Traverse()
	for _, parentNode := range nodes {
		if parentNode.Name == parent {
			parentFound = true
			for _, child := range parentNode.Children {
				if child.Name == targetNode.Name {
					targetFound = true
					child.ReadOnly = targetNode.ReadOnly
					if targetNode.Color == "" {
						child.Color = defaultColor
					} else {
						child.Color = targetNode.Color
					}
					child.Notes = targetNode.Notes
					if targetNode.Children != nil && len(targetNode.Children) > 0 {
						// there should be only one child
						newChild := targetNode.Children[0]
						// the name must be unique
						for _, node := range nodes {
							if node.Name == newChild.Name {
								return nil, fmt.Errorf("duplicated name %q", newChild.Name)
							}
						}
						child.Children = append(child.Children, newChild)
					}
				}
			}
		}
	}
	if !parentFound {
		return nil, fmt.Errorf("the parent node %q was not found", parent)
	}
	if !targetFound {
		return nil, fmt.Errorf("the target node %q was not found", targetNode.Name)
	}
	return n, nil
}

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
		traversed = traverse(child, traversed)
	}
	return traversed
}
