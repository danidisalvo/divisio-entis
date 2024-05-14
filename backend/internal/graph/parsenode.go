package graph

import (
	"backend/internal/graph/errors"
	"encoding/json"
	"fmt"
)

// Parse parses a node's JSON representation
func (n *Node) Parse(bytes []byte) (*Node, error) {
	err := json.Unmarshal(bytes, n)
	if err != nil {
		return nil, errors.NewParsingError(fmt.Sprintf("failed to parse the node [%s]", err))
	}
	n.Traverse()
	return n, nil
}
