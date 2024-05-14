package graph

import (
	"backend/internal/graph/errors"
	"fmt"
)

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
							return nil, errors.NewIllegalArgumentError(msg)
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
		return nil, errors.NewNodeNotFoundError(fmt.Sprintf("the parent node with ID %q was not found", parentId))
	}
	if target == nil {
		return nil, errors.NewNodeNotFoundError(fmt.Sprintf("the target node with ID %q was not found", targetId))
	}
	if newParent == nil {
		return nil, errors.NewNodeNotFoundError(fmt.Sprintf("the new parent node with ID %q was not found", newParentId))
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
