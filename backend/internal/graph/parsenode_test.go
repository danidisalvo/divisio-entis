package graph_test

import (
	"backend/internal/graph"
	"reflect"
	"testing"
)

func TestNode_Parse_Success(t *testing.T) {
	root, err := graph.NewLexeme("root", "root", "")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, expected, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	root, err = root.Parse(testGraphData)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	traversed := root.Traverse()
	if !reflect.DeepEqual(expected, traversed) {
		t.Error("Lists do not match")
	}
}
