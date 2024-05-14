package graph_test

import (
	"backend/internal/graph"
	_ "embed"
	"fmt"
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
	node, err := graph.NewLexeme("new node", "new node", "")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if node.Name != "new node" || node.Color != graph.DefaultColor {
		t.Errorf("Expected {\"new node\", true, 'default color', \"some stuff\"}, got %v", node)
	}
}

func TestNewNode_Fails(t *testing.T) {
	_, err := graph.NewLexeme("id", "", "")
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

func provisionNodes() (*graph.Node, []*graph.Node, error) {
	node, err := graph.NewLexeme("0", "ens", "")
	node.SetProperty("p1", "abc")
	node.SetProperty("p2", "xyz")
	if err != nil {
		return nil, nil, err
	}
	b, err := graph.NewLexeme("id_B", "B", red)
	if err != nil {
		return nil, nil, err
	}
	node.Children = append(node.Children, b)

	c, err := graph.NewLexeme("id_C", "C", red)
	if err != nil {
		return nil, nil, err
	}
	node.Children = append(node.Children, c)

	d, err := graph.NewOpposition("id_D", "D", green)
	if err != nil {
		return nil, nil, err
	}
	node.Children = append(node.Children, d)

	e, err := graph.NewLexeme("id_E", "E", green)
	if err != nil {
		return nil, nil, err
	}
	node.Children = append(node.Children, e)

	f, err := graph.NewLexeme("id_F", "F", blu)
	if err != nil {
		return nil, nil, err
	}
	d.Children = append(d.Children, f)

	g, err := graph.NewLexeme("id_G", "G", blu)
	if err != nil {
		return nil, nil, err
	}
	d.Children = append(d.Children, g)

	h, err := graph.NewDivision("id_H", "H", yellow)
	if err != nil {
		return nil, nil, err
	}
	g.Children = append(g.Children, h)

	i, err := graph.NewLexeme("id_I", "I", yellow)
	if err != nil {
		return nil, nil, err
	}
	g.Children = append(g.Children, i)

	nodes := [9]*graph.Node{node, b, c, d, f, g, h, i, e}

	return node, nodes[:], nil
}
