package graph_test

import "testing"

func TestNode_SimpleString_Success(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
	}
	expected := string(testPrintData)
	actual := root.Stringify()
	if expected != actual {
		t.Errorf("strings fo not match. Expected: %s\n. Actual: %s", expected, actual)
	}
}
