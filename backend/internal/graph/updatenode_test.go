package graph_test

import (
	"backend/internal/graph"
	"reflect"
	"testing"
)

func TestNode_UpdateNode_Success(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	targetNode, err := graph.NewLexeme("id_F", "F", "")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	root, err = root.UpdateNode("id_D", targetNode)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	found, err := root.FindNode("id_F")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if !reflect.DeepEqual(targetNode, found) {
		t.Errorf("The nodes do not match. Expected %v, got %v", targetNode, found)
	}
}

func TestNode_UpdateNode_FailsTargetIsNil(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, err = root.UpdateNode("id_D", nil)
	if err == nil {
		t.Errorf("UpdateNode did not return an error")
		return
	}
	if err.Error() != "targetNode cannot be nil" {
		t.Errorf("The error message does not match. Expected \"targetNode cannot be nil\", got %s", err)
	}
}

func TestNode_UpdateNode_FailsTargetNotFound(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	targetNode, err := graph.NewLexeme("id_Z", "Z", "")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, err = root.UpdateNode("id_D", targetNode)
	if err == nil {
		t.Errorf("UpdateNode did not return an error")
		return
	}
	if err.Error() != "the target node with ID \"id_Z\" was not found" {
		t.Errorf("The error message does not match. Expected \"the target node with ID \"id_Z\" was not found\", got %s", err)
	}
}

func TestNode_UpdateNode_FailsParentNotFound(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	targetNode, err := graph.NewLexeme("F", "F", "")
	_, err = root.UpdateNode("Z", targetNode)
	if err == nil {
		t.Errorf("UpdateNode did not return an error")
		return
	}
	if err.Error() != "the parent node with ID \"Z\" was not found" {
		t.Errorf("The error message does not match. Expected \"the parent node with ID \"Z\" was not found\", got %s", err)
	}
}
