package graph

import (
	_ "embed"
	"fmt"
	"reflect"
	"testing"
)

//go:embed test-graph.json
var testGraphData []byte

//go:embed test-print.txt
var testPrintData []byte

const (
	red    = "#ff0000"
	green  = "#00ff00"
	blu    = "#0000ff"
	yellow = "#00ffff"
)

func TestNewNode_Success(t *testing.T) {
	node, err := NewLexeme("new node", "new node", "")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if node.Name != "new node" || node.Color != defaultColor {
		t.Errorf("Expected {\"new node\", true, 'default color', \"some stuff\"}, got %v", node)
	}
}

func TestNewNode_Fails(t *testing.T) {
	_, err := NewLexeme("id", "", "")
	if err == nil {
		t.Errorf("NewNode did not return an error")
		return
	}
	if err.Error() != "name cannot be empty" {
		t.Errorf("Expected \"name cannot be empty\", got %s", err.Error())
	}
}

func TestNode_String(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	json, err := root.String()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	fmt.Println(json)
}

func TestNode_Parse_Success(t *testing.T) {
	root, err := NewLexeme("root", "root", "")
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

func TestNode_FindNode_Success(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	expected, err := NewLexeme("id_F", "F", blu)
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

func TestNode_AddNode_Success(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	node, err := NewLexeme("id_K", "K", "")
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
	node, err := NewLexeme("K", "K", "")
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
	node, err := NewLexeme("id_G", "G", "")
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
	actual, err := root.MoveNode("id_D", "id_G", "id_G")
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

func TestNode_UpdateNode_Success(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	targetNode, err := NewLexeme("id_F", "F", "")
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

func TestNode_UpdateNode_FailsTargetNotFound(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	targetNode, err := NewLexeme("id_Z", "Z", "")
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
	targetNode, err := NewLexeme("F", "F", "")
	_, err = root.UpdateNode("Z", targetNode)
	if err == nil {
		t.Errorf("UpdateNode did not return an error")
		return
	}
	if err.Error() != "the parent node with ID \"Z\" was not found" {
		t.Errorf("The error message does not match. Expected \"the parent node with ID \"Z\" was not found\", got %s", err)
	}
}

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

func TestNode_SimpleString_Success(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
	}
	expected := string(testPrintData)
	actual := root.SimpleString()
	if expected != actual {
		t.Errorf("strings fo not match. Expected: %s\n. Actual: %s", expected, actual)
	}
}

func provisionNodes() (*Node, []*Node, error) {
	node, err := NewLexeme("0", "ens", "")
	node.SetProperty("p1", "abc")
	node.SetProperty("p2", "xyz")
	if err != nil {
		return nil, nil, err
	}
	b, err := NewLexeme("id_B", "B", red)
	if err != nil {
	}
	node.Children = append(node.Children, b)

	c, err := NewLexeme("id_C", "C", red)
	if err != nil {
	}
	node.Children = append(node.Children, c)

	d, err := NewOpposition("id_D", "D", green)
	if err != nil {
	}
	node.Children = append(node.Children, d)

	e, err := NewLexeme("id_E", "E", green)
	if err != nil {
	}
	node.Children = append(node.Children, e)

	f, err := NewLexeme("id_F", "F", blu)
	if err != nil {
	}
	d.Children = append(d.Children, f)

	g, err := NewLexeme("id_G", "G", blu)
	if err != nil {
	}
	d.Children = append(d.Children, g)

	h, err := NewDivision("id_H", "H", yellow)
	if err != nil {
	}
	g.Children = append(g.Children, h)

	i, err := NewLexeme("id_I", "I", yellow)
	if err != nil {
	}
	g.Children = append(g.Children, i)

	nodes := [9]*Node{node, b, c, d, f, g, h, i, e}

	return node, nodes[:], nil
}
