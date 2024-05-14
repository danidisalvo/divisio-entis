package graph_test

import (
	"reflect"
	"testing"
)

func TestNode_MoveNode(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	root, err = root.MoveNode("id_D", "id_G", "id_B")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	found, err := root.FindNode("id_B")
	if len(found.Children) == 0 {
		t.Errorf("The node B has no children")
		return
	}
	found, err = root.FindNode("id_D")
	for _, n := range found.Children {
		if n.Id == "id_G" {
			t.Errorf("The node G is a child of node D")
			return
		}
	}
}

func TestNode_MoveNode_NoAction(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	actual, err := root.MoveNode("id_D", "id_G", "id_D")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if !reflect.DeepEqual(root, actual) {
		t.Errorf("The tree has changed")
		return
	}
}

func TestNode_MoveNode_FailsMoveToChild(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, err = root.MoveNode("id_D", "id_G", "id_H")
	if err == nil {
		t.Errorf("MoveNode did not return an error")
		return
	}
	if err.Error() != "the node \"id_G\" cannot be moved to its child \"id_H\"" {
		t.Errorf("The error message does not match. Expected \"the node \"id_G\" cannot be moved to its child \"id_H\"\", got %s", err)
	}
}

func TestNode_MoveNode_FailsParentNotFound(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, err = root.MoveNode("id_Z", "id_G", "id_B")
	if err == nil {
		t.Errorf("MoveNode did not return an error")
		return
	}
	if err.Error() != "the parent node with ID \"id_Z\" was not found" {
		t.Errorf("The error message does not match. Expected \"the parent node with ID \"id_Z\" was not found\", got %s", err)
	}
}

func TestNode_MoveNode_FailsTargetNotFound(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, err = root.MoveNode("id_D", "id_Z", "id_B")
	if err == nil {
		t.Errorf("MoveNode did not return an error")
		return
	}
	if err.Error() != "the target node with ID \"id_Z\" was not found" {
		t.Errorf("The error message does not match. Expected \"the target node with ID \"id_Z\" was not found\", got %s", err)
	}
}

func TestNode_MoveNode_FailsNewParentNotFound(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, err = root.MoveNode("id_D", "id_G", "id_Z")
	if err == nil {
		t.Errorf("MoveNode did not return an error")
		return
	}
	if err.Error() != "the new parent node with ID \"id_Z\" was not found" {
		t.Errorf("The error message does not match. Expected \"the new parent node with ID \"id_Z\" was not found\", got %s", err)
	}
}
