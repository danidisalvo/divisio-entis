package graph

import (
	"backend/internal/graph/errors"
	"encoding/json"
	"fmt"
)

const (
	DefaultColor = "#dddddd"
	division     = NodeType("division")
	lexeme       = NodeType("lexeme")
	opposition   = NodeType("opposition")
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

// NewDivision creates a new division node
func NewDivision(id, name, color string) (*Node, error) {
	return newNode(id, name, color, division)
}

// NewLexeme creates a new LEXEME node
func NewLexeme(id, name, color string) (*Node, error) {
	return newNode(id, name, color, lexeme)
}

// NewOpposition creates a new OPPOSITION node
func NewOpposition(id, name, color string) (*Node, error) {
	return newNode(id, name, color, opposition)
}

// NewNode creates a new node
func newNode(id, name, color string, nodeType NodeType) (*Node, error) {
	if id == "" {
		return nil, errors.NewIllegalArgumentError("id cannot be empty")
	}
	if name == "" {
		return nil, errors.NewIllegalArgumentError("name cannot be empty")
	}
	if color == "" {
		color = DefaultColor
	}
	return &Node{
		Id:         id,
		Name:       name,
		Color:      color,
		Type:       nodeType,
		Properties: make(map[string]string),
		Children:   make([]*Node, 0),
	}, nil
}

// String returns the JSON representation of this node
func (n *Node) String() (string, error) {
	bytes, err := json.Marshal(n)
	if err != nil {
		return "", errors.NewParsingError(fmt.Sprintf("failed to marshal the node [%s]", err))
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
