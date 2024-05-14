package graph

import (
	"fmt"
	"strconv"
)

// Stringify returns a flat string representation of this node (see test-print.txt)
func (n *Node) Stringify() string {
	var traversed string
	var counters []int
	counters = append(counters, 1)
	return stringify(n, traversed, "", counters)
}

// stringify recursively stringifies the graph using the Depth-First Search algorithm
func stringify(node *Node, traversed string, prefix string, counters []int) string {
	traversed = fmt.Sprintf("%s%s %s\n", traversed, formatCounters(counters), node.Name)
	if len(node.Children) > 0 {
		counters = append(counters, 0)
	}
	for _, child := range node.Children {
		counters[len(counters)-1]++
		traversed = stringify(child, traversed, prefix, counters)
	}
	return traversed
}

func formatCounters(counters []int) string {
	s := ""
	for i := 0; i < len(counters); i++ {
		s += "." + strconv.Itoa(counters[i])
	}
	return s[1:]
}
