package graph_test

import (
	"backend/internal/graph"
	"testing"
)

func TestNode_AddNode_Success(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	node, err := graph.NewLexeme("id_K", "K", "")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	root, err = root.AddNode("id_F", node)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	found, err := root.FindNode("id_F")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if found.Children[0].Name != node.Name {
		t.Errorf("Nodes do not match. Expected %v, got %v", node, found.Children[0])
	}
}

func TestNode_AddNode_FailsNodeIsNil(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, err = root.AddNode("F", nil)
	if err == nil {
		t.Errorf("AddNode did not return an error")
		return
	}
	if err.Error() != "newNode cannot be nil" {
		t.Errorf("The error message does not match. Expected \"node cannot be nil\", got %s", err)
	}
}

func TestNode_AddNode_FailsNameIsNotUnique(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	node, err := graph.NewLexeme("K", "K", "")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, err = root.AddNode("Z", node)
	if err == nil {
		t.Errorf("validate did not return an error")
		return
	}
	if err.Error() != "parent \"Z\" not found" {
		t.Errorf("The error message does not match. Expected \"parent \"Z\" not found\", got %s", err)
	}
}

func TestNode_AddNode_FailsParentNotFound(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	node, err := graph.NewLexeme("id_G", "G", "")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, err = root.AddNode("0", node)
	if err == nil {
		t.Errorf("AddNode did not return an error")
		return
	}
	if err.Error() != "duplicated ID \"id_G\"" {
		t.Errorf("The error message does not match. Expected \"duplicated ID \"id_G\"\", got %s", err)
	}
}
