package graph

import (
	"backend/internal/graph/errors"
	"fmt"
)

// UpdateNode updates a graph's node
func (n *Node) UpdateNode(parent string, targetNode *Node) (*Node, error) {
	if targetNode == nil {
		return nil, errors.NewIllegalArgumentError("targetNode cannot be nil")
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
						child.Color = DefaultColor
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
								return nil, errors.NewDuplicatedNodeError(fmt.Sprintf("duplicated ID %q", newChild.Id))
							}
						}
						child.Children = append(child.Children, newChild)
					}
				}
			}
		}
	}
	if !parentFound {
		return nil, errors.NewNodeNotFoundError(fmt.Sprintf("the parent node with ID %q was not found", parent))
	}
	if !targetFound {
		return nil, errors.NewNodeNotFoundError(fmt.Sprintf("the target node with ID %q was not found", targetNode.Id))
	}
	return n, nil
}
