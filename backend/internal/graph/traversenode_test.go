package graph_test

import (
	"reflect"
	"testing"
)

func TestNode_Traverse_Success(t *testing.T) {
	root, expected, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
	}
	actual := root.Traverse()
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Lists do not match")
	}
}
