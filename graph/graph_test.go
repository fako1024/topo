////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 by Fabian Kohn
//
// This source code is licensed under the Apache License, Version 2.0, found in
// the LICENSE file in the root directory of this source tree.
////////////////////////////////////////////////////////////////////////////////

package graph

import (
	"fmt"
	"strings"
	"testing"
)

const nRunTestTable = 1000

type testArc struct {
	from interface{}
	to   interface{}
}

type testCase struct {
	graph  *Graph
	arcs   []testArc
	result ObjectList
}

func (c *testCase) init() {
	for _, arc := range c.arcs {
		c.graph.AddArc(arc.from, arc.to)
	}
}

func (c *testCase) assessOrder() error {
	for _, arc := range c.arcs {
		posSource, _ := c.result.find(arc.from)
		posDestination, _ := c.result.find(arc.to)

		if posDestination >= posSource {
			return fmt.Errorf("Unexpected order, want pos(%s) < pos(%s)", arc.to, arc.from)
		}
	}

	return nil
}

var testTable = []testCase{
	testCase{
		graph: NewGraph("a"),
	},
	testCase{
		graph: NewGraph("a", "b", "c", "d", "e"),
	},
	testCase{
		graph: NewGraph("a", "b", "c", "d", "e"),
		arcs:  []testArc{testArc{"a", "b"}, testArc{"b", "c"}, testArc{"c", "d"}},
	},
	testCase{
		graph: NewGraph("a", "b", "c", "d", "e"),
		arcs:  []testArc{testArc{"a", "b"}, testArc{"c", "b"}, testArc{"e", "d"}},
	},
	testCase{
		graph: NewGraph(1, 2, 3, 4, 5, 6, 7, 8, 9, 10),
		arcs:  []testArc{testArc{7, 8}, testArc{7, 1}, testArc{7, 3}, testArc{2, 7}},
	},
	testCase{
		graph: NewGraph(1, 2.63535, 3, 4, 5, 6, "A", 8, 9, 10),
		arcs:  []testArc{testArc{"A", 8}, testArc{"A", 1}, testArc{"A", 3}, testArc{5, "A"}, testArc{2.63535, 1}},
	},
}

func TestGraphInteraction(t *testing.T) {

	var dummyVtxList = make(ObjectList, 0)
	if idx, ok := dummyVtxList.find("doesnotexist"); idx != -1 || ok {
		t.Fatal("Expected negative index, but got:", idx, ok)
	}

	var dummyList = newList()
	if idx, ok := dummyList.findIndex("doesnotexist"); idx != -1 || ok {
		t.Fatal("Expected negative index, but got:", idx, ok)
	}

	graph := NewGraph()
	result, err := graph.SortTopological()
	if len(result) != 0 {
		t.Fatal("Expected empty result, got", result.String())
	}
	_ = result.String()

	graph = NewGraph("a", "b", "c", "d", "e")

	// Try successful addition of vertex
	graph.AddVertex("f")

	// Try successful addition of arc
	if err = graph.AddArc("d", "a"); err != nil {
		t.Fatal("Graph interaction error:", err)
	}

	// Try failed addition of arc
	if err = graph.AddArc("d", "doesnotexist"); err == nil {
		t.Fatal("Expected graph interaction error, but got none")
	}

	// Try failed addition of arc
	if err = graph.AddArc("dontexist", "a"); err == nil {
		t.Fatal("Expected graph interaction error, but got none")
	}
}

func TestGraphTable(t *testing.T) {
	var err error
	for nRun := 0; nRun < nRunTestTable; nRun++ {
		for _, test := range testTable {
			test.init()

			test.result, err = test.graph.SortTopological()

			if err != nil {
				t.Fatal(err)
			}

			if len(test.result) != len(test.graph.vertices) {
				t.Fatal("Number of elements does not match expectation:", test.result, "vs.", test.graph.vertices)
			}

			if err = test.assessOrder(); err != nil {
				t.Fatal("Sort order mismatch detected:", err, test.result)
			}
		}
	}
}

func TestGraphCyclic(t *testing.T) {
	cyclicGraph := NewGraph("a", "b", "c", "d")

	cyclicGraph.AddArc("a", "b")
	cyclicGraph.AddArc("b", "c")
	cyclicGraph.AddArc("c", "a")

	if _, err := cyclicGraph.SortTopological(); err == nil {
		t.Fatal("Expected cyclic error not seen")
	} else if !strings.Contains(err.Error(), "Cycle error:") {
		t.Errorf("Unexpected error message: %s", err)
	}
}
