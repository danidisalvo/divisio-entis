package graph

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strconv"
)

const (
	defaultColor = "#dddddd"
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
		return nil, &IllegalArgumentError{"id cannot be empty"}
	}
	if name == "" {
		return nil, &IllegalArgumentError{"name cannot be empty"}
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
		Children:   make([]*Node, 0),
	}, nil
}

// String returns the JSON representation of this node
func (n *Node) String() (string, error) {
	bytes, err := json.Marshal(n)
	if err != nil {
		return "", &ParsingError{fmt.Sprintf("failed to marshal the node [%s]", err)}
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
		return nil, &ParsingError{fmt.Sprintf("failed to parse the node [%s]", err)}
	}
	n.Traverse()
	return n, nil
}

// FindNode returns the node with the given id
func (n *Node) FindNode(id string) (*Node, error) {
	if id == "" {
		return nil, &IllegalArgumentError{"id cannot be empty"}
	}
	for _, node := range n.Traverse() {
		if node.Id == id {
			return node, nil
		}
	}
	return nil, &NodeNotFoundError{fmt.Sprintf("the node with ID %q was not found", id)}
}

// AddNode adds a node to the graph
func (n *Node) AddNode(parent string, newNode *Node) (*Node, error) {
	if newNode == nil {
		return nil, &IllegalArgumentError{"newNode cannot be nil"}
	}
	// the id must be unique
	for _, node := range n.Traverse() {
		if node.Id == newNode.Id {
			return nil, &DuplicatedNodeError{fmt.Sprintf("duplicated ID %q", newNode.Id)}
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
	return nil, &NodeNotFoundError{fmt.Sprintf("parent %q not found", parent)}
}

// MoveNode moves a node from its parent to a new parent
func (n *Node) MoveNode(parentId, targetId, newParentId string) (*Node, error) {
	// trivial case, nothing to be done
	if parentId == newParentId {
		return n, nil
	}

	var parent *Node
	var target *Node
	var newParent *Node

	for _, node := range n.Traverse() {
		if node.Id == parentId {
			parent = node
			for _, child := range node.Children {
				if child.Id == targetId {
					// newParent cannot be a target's child
					for _, c := range child.Children {
						if c.Id == newParentId {
							msg := fmt.Sprintf("the node %q cannot be moved to its child %q", targetId, newParentId)
							return nil, &IllegalArgumentError{msg}
						}
					}
					target = child
				}
			}
		} else if node.Id == newParentId {
			newParent = node
		}
		if parent != nil && target != nil && newParent != nil {
			break
		}
	}

	if parent == nil {
		return nil, &NodeNotFoundError{fmt.Sprintf("the parent node with ID %q was not found", parentId)}
	}
	if target == nil {
		return nil, &NodeNotFoundError{fmt.Sprintf("the target node with ID %q was not found", targetId)}
	}
	if newParent == nil {
		return nil, &NodeNotFoundError{fmt.Sprintf("the new parent node with ID %q was not found", newParentId)}
	}

	// add the target to the new parent's children
	newParent.Children = append(newParent.Children, target)

	// remove the target from the parent's children
	children := make([]*Node, 0)
	for _, child := range parent.Children {
		if child.Id != targetId {
			children = append(children, child)
		}
	}
	parent.Children = children

	return n, nil
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
		return nil, &NodeNotFoundError{fmt.Sprintf("the parent node with ID %q was not found", parent)}
	}
	if !targetFound {
		return nil, &NodeNotFoundError{fmt.Sprintf("the target node with ID %q was not found", target)}
	}
	return n, nil
}

// UpdateNode updates a graph's node
func (n *Node) UpdateNode(parent string, targetNode *Node) (*Node, error) {
	if targetNode == nil {
		return nil, &IllegalArgumentError{"targetNode cannot be nil"}
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
					child.Name = targetNode.Name
					child.Type = targetNode.Type
					if targetNode.Color == "" {
						child.Color = defaultColor
					} else {
						child.Color = targetNode.Color
					}
					child.Properties = targetNode.Properties
					if targetNode.Children != nil && len(targetNode.Children) == 1 {
						// there should be only one child
						newChild := targetNode.Children[0]
						// the id must be unique
						for _, node := range nodes {
							if node.Id == newChild.Id {
								return nil, &DuplicatedNodeError{fmt.Sprintf("duplicated ID %q", newChild.Id)}
							}
						}
						child.Children = append(child.Children, newChild)
					}
				}
			}
		}
	}
	if !parentFound {
		return nil, &NodeNotFoundError{fmt.Sprintf("the parent node with ID %q was not found", parent)}
	}
	if !targetFound {
		return nil, &NodeNotFoundError{fmt.Sprintf("the target node with ID %q was not found", targetNode.Id)}
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

// SimpleString returns a simplified string representation of this node (see test-print.txt)
func (n *Node) SimpleString() string {
	var traversed string
	var counters []int
	counters = append(counters, 1)
	return simpleString(n, traversed, "", counters)
}

// traverse recursively traverses the graph using the Depth-First Search algorithm
func simpleString(node *Node, traversed string, prefix string, counters []int) string {
	traversed = fmt.Sprintf("%s%s %s\n", traversed, formatCounters(counters), node.Name)
	if len(node.Children) > 0 {
		counters = append(counters, 0)
	}
	for _, child := range node.Children {
		counters[len(counters)-1]++
		traversed = simpleString(child, traversed, prefix, counters)
	}
	return traversed
}

func formatCounters(counters []int) string {
	s := ""
	for i := 0; i < len(counters); i++ {
		s += "." + strconv.Itoa(counters[i])
	}
	return s[1:]
}
