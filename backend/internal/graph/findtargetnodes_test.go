package graph_test

import (
	"backend/internal/graph"
	"reflect"
	"testing"
)

func TestNode_FindTargetNodes_of_0_Success(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
	}
	expected := make([]*graph.Node, 0)

	actual, err := root.FindTargetNodes("0")
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Lists do not match")
	}
}

func TestNode_FindTargetNodes_of_D_Success(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
	}
	var expected []*graph.Node
	//expected := make([]*graph.Node, 0)
	expected = append(expected, &graph.Node{Id: "0", Name: "ens"})
	expected = append(expected, &graph.Node{Id: "id_B", Name: "B"})
	expected = append(expected, &graph.Node{Id: "id_C", Name: "C"})
	expected = append(expected, &graph.Node{Id: "id_E", Name: "E"})

	actual, err := root.FindTargetNodes("id_D")
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Lists do not match")
	}
}

func TestNode_FindTargetNodes_of_G_Success(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
	}
	var expected []*graph.Node
	expected = append(expected, &graph.Node{Id: "0", Name: "ens"})
	expected = append(expected, &graph.Node{Id: "id_B", Name: "B"})
	expected = append(expected, &graph.Node{Id: "id_C", Name: "C"})
	expected = append(expected, &graph.Node{Id: "id_D", Name: "D"})
	expected = append(expected, &graph.Node{Id: "id_F", Name: "F"})
	expected = append(expected, &graph.Node{Id: "id_E", Name: "E"})

	actual, err := root.FindTargetNodes("id_G")
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Lists do not match")
	}
}

func TestNode_FindTargetNodes_of_H_Success(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
	}
	var expected []*graph.Node
	expected = append(expected, &graph.Node{Id: "0", Name: "ens"})
	expected = append(expected, &graph.Node{Id: "id_B", Name: "B"})
	expected = append(expected, &graph.Node{Id: "id_C", Name: "C"})
	expected = append(expected, &graph.Node{Id: "id_D", Name: "D"})
	expected = append(expected, &graph.Node{Id: "id_F", Name: "F"})
	expected = append(expected, &graph.Node{Id: "id_G", Name: "G"})
	expected = append(expected, &graph.Node{Id: "id_I", Name: "I"})
	expected = append(expected, &graph.Node{Id: "id_E", Name: "E"})

	actual, err := root.FindTargetNodes("id_H")
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Lists do not match")
	}
}

func TestNode_FindTargetNodes_FailsNodeIsEmpty(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, err = root.FindTargetNodes("")
	if err == nil {
		t.Errorf("AddNode did not return an error")
		return
	}
	if err.Error() != "node cannot be empty" {
		t.Errorf("The error message does not match. Expected \"node cannot be nil\", got %s", err)
	}
}
