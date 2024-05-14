package graph

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

// Graph contains the graph's root node
type Graph struct {
	Root     *Node
	Filename string
}

// NewGraph create a new graph
func NewGraph(name, filename string) (*Graph, error) {
	root, err := NewLexeme("0", name, DefaultColor)
	if err != nil {
		return nil, err
	}
	return &Graph{Root: root, Filename: filename}, nil
}

// Clear reset this graph
func (g *Graph) Clear() *Graph {
	g.Root.Children = make([]*Node, 0)
	return g
}

// Load loads the graph as a JSON file from disk
func (g *Graph) Load() {
	bytes, err := os.ReadFile(g.Filename)
	if err != nil {
		msg := fmt.Sprintf("Failed to read the file [%s]", err)
		log.Error(msg)
		return
	}

	log.Debugf("Read %d bytes", len(bytes))
	g.Root, err = g.Root.Parse(bytes)
	if err != nil {
		msg := fmt.Sprintf("Failed to read the file [%s]", err)
		log.Error(msg)
	}
}

// Save saves the graph as a JSON file to disk
func (g *Graph) Save() {
	json, err := g.Root.String()
	if err != nil {
		msg := fmt.Sprintf("Failed to generate the JSON string [%s]", err)
		log.Error(msg)
		return
	}
	bytes := []byte(json)
	err = os.WriteFile(g.Filename, bytes, 0600)
	if err != nil {
		msg := fmt.Sprintf("Failed to save the file [%s]", err)
		log.Error(msg)
		return
	}
	log.Debugf("Written %d bytes", len(bytes))
}
