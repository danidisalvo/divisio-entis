package graph_test

import (
	"backend/internal/graph"
	"reflect"
	"testing"
)

func TestNode_FindNode_Success(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	expected, err := graph.NewLexeme("id_F", "F", blu)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	actual, err := root.FindNode("id_F")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("The nodes do not match. Expected %v, got %v", expected, actual)
	}
}

func TestNode_FindNode_IdIsEmpty(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	_, err = root.FindNode("")
	if err == nil {
		t.Errorf("FindNode did not return an error")
		return
	}
	if err.Error() != "id cannot be empty" {
		t.Errorf("The error message does not match. Expected \"id cannot be empty\", got %s", err)
	}
}

func TestNode_FindNode_NotFound(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	_, err = root.FindNode("Z")
	if err == nil {
		t.Errorf("FindNode did not return an error")
		return
	}
	if err.Error() != "the node with ID \"Z\" was not found" {
		t.Errorf("The error message does not match. Expected \"the node with ID \"Z\" was not found\", got %s", err)
	}
}
