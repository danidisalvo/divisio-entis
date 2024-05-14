package graph_test

import "testing"

func TestNode_RemoveNode_Success(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	root, err = root.RemoveNode("id_D", "id_F")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, err = root.FindNode("id_F")
	if err == nil {
		t.Errorf("FindNode did not return an error")
		return
	}
	if err.Error() != "the node with ID \"id_F\" was not found" {
		t.Errorf("The error message does not match. Expected \"the node with ID \"id_F\" was not found\", got %s", err)
	}
}

func TestNode_RemoveNode_FailsTargetNotFound(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, err = root.RemoveNode("id_D", "Z")
	if err == nil {
		t.Errorf("RemoveNode did not return an error")
		return
	}
	if err.Error() != "the target node with ID \"Z\" was not found" {
		t.Errorf("The error message does not match. Expected \"the target node with ID \"Z\" was not found\", got %s", err)
	}
}

func TestNode_RemoveNode_FailsParentNotFound(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, err = root.RemoveNode("Z", "F")
	if err == nil {
		t.Errorf("RemoveNode did not return an error")
		return
	}
	if err.Error() != "the parent node with ID \"Z\" was not found" {
		t.Errorf("The error message does not match. Expected \"the parent node with ID \"Z\" was not found\", got %s", err)
	}
}
