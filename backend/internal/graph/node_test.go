package graph

import (
	_ "embed"
	"fmt"
	"reflect"
	"testing"
)

//go:embed test-graph.json
var testGraphData []byte

const (
	red    = "#ff0000"
	green  = "#00ff00"
	blu    = "#0000ff"
	yellow = "#00ffff"
)

func TestNewNode_Success(t *testing.T) {
	node, err := NewNode("new node", true, "", "first division", "some stuff")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if node.Name != "new node" || !node.ReadOnly ||
		node.Color != defaultColor || node.Division != "first division" || node.Notes != "some stuff" {
		t.Errorf("Expected {\"new node\", true, 'default color', \"first division\", \"some stuff\"}, got %v", node)
	}
}

func TestNewNode_Fails(t *testing.T) {
	_, err := NewNode("", true, "", "", "")
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
	root, err := NewNode("root", true, "", "", "")
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

	expected, err := NewNode("F", false, blu, "division 3", "notes F")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	actual, err := root.FindNode("F")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("The nodes do not match. Expected %v, got %v", expected, actual)
	}
}

func TestNode_FindNode_NameIsEmpty(t *testing.T) {
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
	if err.Error() != "name cannot be empty" {
		t.Errorf("The error message does not match. Expected \"name cannot be empty\", got %s", err)
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
	if err.Error() != "the node \"Z\" was not found" {
		t.Errorf("The error message does not match. Expected \"the node \"Z\" was not found\", got %s", err)
	}
}

func TestNode_AddNode_Success(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	node, err := NewNode("K", false, "", "", "")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	root, err = root.AddNode("F", node)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	found, err := root.FindNode("F")
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
	node, err := NewNode("K", false, "", "", "")
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
	node, err := NewNode("G", false, "", "", "")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, err = root.AddNode("B", node)
	if err == nil {
		t.Errorf("AddNode did not return an error")
		return
	}
	if err.Error() != "duplicated name \"G\"" {
		t.Errorf("The error message does not match. Expected \"duplicated name \"G\"\", got %s", err)
	}
}

func TestNode_RemoveNode_Success(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	root, err = root.RemoveNode("D", "F")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, err = root.FindNode("F")
	if err == nil {
		t.Errorf("FindNode did not return an error")
		return
	}
	if err.Error() != "the node \"F\" was not found" {
		t.Errorf("The error message does not match. Expected \"the node \"F\" was not found\", got %s", err)
	}
}

func TestNode_RemoveNode_FailsTargetNotFound(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, err = root.RemoveNode("D", "Z")
	if err == nil {
		t.Errorf("RemoveNode did not return an error")
		return
	}
	if err.Error() != "the target node \"Z\" was not found" {
		t.Errorf("The error message does not match. Expected \"the target node \"Z\" was not found\", got %s", err)
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
	if err.Error() != "the parent node \"Z\" was not found" {
		t.Errorf("The error message does not match. Expected \"the parent node \"Z\" was not found\", got %s", err)
	}
}

func TestNode_UpdateNode_Success(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	targetNode, err := NewNode("F", false, "", "New division F", "bla bla bla")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	root, err = root.UpdateNode("D", targetNode)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	found, err := root.FindNode("F")
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
	targetNode, err := NewNode("Z", false, "", "", "")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, err = root.UpdateNode("D", targetNode)
	if err == nil {
		t.Errorf("UpdateNode did not return an error")
		return
	}
	if err.Error() != "the target node \"Z\" was not found" {
		t.Errorf("The error message does not match. Expected \"the target node \"Z\" was not found\", got %s", err)
	}
}

func TestNode_UpdateNode_FailsParentNotFound(t *testing.T) {
	root, _, err := provisionNodes()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	targetNode, err := NewNode("F", false, "", "", "")
	_, err = root.UpdateNode("Z", targetNode)
	if err == nil {
		t.Errorf("UpdateNode did not return an error")
		return
	}
	if err.Error() != "the parent node \"Z\" was not found" {
		t.Errorf("The error message does not match. Expected \"the parent node \"Z\" was not found\", got %s", err)
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

func provisionNodes() (*Node, []*Node, error) {
	node, err := NewNode("ens", false, "", "", "")
	if err != nil {
		return nil, nil, err
	}
	b, err := NewNode("B", false, red, "division 1", "notes B")
	if err != nil {
	}
	node.Children = append(node.Children, b)

	c, err := NewNode("C", false, red, "division 1", "notes C")
	if err != nil {
	}
	node.Children = append(node.Children, c)

	d, err := NewNode("D", false, green, "division 2", "notes D")
	if err != nil {
	}
	node.Children = append(node.Children, d)

	e, err := NewNode("E", false, green, "division 2", "notes E")
	if err != nil {
	}
	node.Children = append(node.Children, e)

	f, err := NewNode("F", false, blu, "division 3", "notes F")
	if err != nil {
	}
	d.Children = append(d.Children, f)

	g, err := NewNode("G", false, blu, "division 3", "notes G")
	if err != nil {
	}
	d.Children = append(d.Children, g)

	h, err := NewNode("H", false, yellow, "division 4", "notes H")
	if err != nil {
	}
	g.Children = append(g.Children, h)

	i, err := NewNode("I", false, yellow, "division 4", "notes I")
	if err != nil {
	}
	g.Children = append(g.Children, i)

	nodes := [9]*Node{node, b, c, d, f, g, h, i, e}

	return node, nodes[:], nil
}
