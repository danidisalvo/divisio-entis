package graph

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

const (
	defaultColor = "#dddddd"
	DIVISION     = NodeType("division")
	LEXEME       = NodeType("lexeme")
	OPPOSITION   = NodeType("opposition")
)

type NodeType string

// Node represents a node of a graph which can be traversed using the Depth-First Search algorithm
type Node struct {
	Id         string            `json:"id"`
	Name       string            `json:"name"`
	Type       NodeType          `json:"type"`
	Color      string            `json:"color"`
	Properties map[string]string `json:"properties"`
	Children   []*Node           `json:"children"`
}

// NewDivision creates a new DIVISION node
func NewDivision(id, name, color string) (*Node, error) {
	return newNode(id, name, color, DIVISION)
}

// NewLexeme creates a new LEXEME node
func NewLexeme(id, name, color string) (*Node, error) {
	return newNode(id, name, color, LEXEME)
}

// NewOpposition creates a new OPPOSITION node
func NewOpposition(id, name, color string) (*Node, error) {
	return newNode(id, name, color, OPPOSITION)
}

// NewNode creates a new node
func newNode(id, name, color string, nodeType NodeType) (*Node, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if color == "" {
		color = defaultColor
	}
	return &Node{
		Id:         id,
		Name:       name,
		Color:      color,
		Type:       nodeType,
		Properties: make(map[string]string),
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

// GetProperty returns a property
func (n *Node) GetProperty(name string) string {
	return n.Properties[name]
}

// SetProperty sets a new property
func (n *Node) SetProperty(name, value string) {
	n.Properties[name] = value
}

// Parse parses a node's JSON representation
func (n *Node) Parse(bytes []byte) (*Node, error) {
	err := json.Unmarshal(bytes, n)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the node [%s]", err)
	}
	n.Traverse() // TODO
	return n, nil
}

// FindNode returns the node with the given id
func (n *Node) FindNode(id string) (*Node, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}
	for _, node := range n.Traverse() {
		if node.Id == id {
			return node, nil
		}
	}
	return nil, fmt.Errorf("the node with ID %q was not found", id)
}

// AddNode adds a node to the graph
func (n *Node) AddNode(parent string, newNode *Node) (*Node, error) {
	if newNode == nil {
		return nil, fmt.Errorf("newNode cannot be nil")
	}
	// the id must be unique
	for _, node := range n.Traverse() {
		if node.Id == newNode.Id {
			return nil, fmt.Errorf("duplicated ID %q", newNode.Id)
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
	return nil, fmt.Errorf("parent %q not found", parent)
}

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
		return nil, fmt.Errorf("the parent node with ID %q was not found", parent)
	}
	if !targetFound {
		return nil, fmt.Errorf("the target node with ID %q was not found", target)
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
		if parentNode.Id == parent {
			parentFound = true
			for _, child := range parentNode.Children {
				if child.Id == targetNode.Id {
					targetFound = true
					if targetNode.Color == "" {
						child.Color = defaultColor
					} else {
						child.Color = targetNode.Color
					}
					child.Properties = targetNode.Properties
					if targetNode.Children != nil && len(targetNode.Children) > 0 {
						// there should be only one child
						newChild := targetNode.Children[0]
						// the id must be unique
						for _, node := range nodes {
							if node.Id == newChild.Id {
								return nil, fmt.Errorf("duplicated ID %q", newChild.Id)
							}
						}
						child.Children = append(child.Children, newChild)
					}
				}
			}
		}
	}
	if !parentFound {
		return nil, fmt.Errorf("the parent node with ID %q was not found", parent)
	}
	if !targetFound {
		return nil, fmt.Errorf("the target node with ID %q was not found", targetNode.Id)
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
		if child.Id == "" {
			child.Id = uuid.New().String()
		}
		if child.Type == "" {
			child.Type = LEXEME
		}
		traversed = traverse(child, traversed)
	}
	return traversed
}
