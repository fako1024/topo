////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 by Fabian Kohn
//
// This source code is licensed under the Apache License, Version 2.0, found in
// the LICENSE file in the root directory of this source tree.
////////////////////////////////////////////////////////////////////////////////

package graph

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const nRunTestTable = 1000

type testArc struct {
	from interface{}
	to   interface{}
}

type testCase struct {
	graph  *Graph[any]
	arcs   []testArc
	result Objects[any]
}

func (c *testCase) init() error {
	for _, arc := range c.arcs {
		if err := c.graph.AddArc(arc.from, arc.to); err != nil {
			return err
		}
	}

	return nil
}

func (c *testCase) assessOrder(t *testing.T) error {
	for _, arc := range c.arcs {
		posSource, _ := c.result.find(arc.from)
		posDestination, _ := c.result.find(arc.to)

		require.Less(t, posDestination, posSource)
	}

	return nil
}

// find determines if a VertexList contains a specific element
func (l Objects[T]) find(obj T) (int, bool) {

	// Check if the vertex exists in the list and return its index if it does
	for i := 0; i < len(l); i++ {
		if l[i] == obj {
			return i, true
		}
	}

	// If not, indicate non-existence and return negative index
	return indexNoExist, false
}

var testTable = []testCase{
	{
		graph: NewGraph[any]("a"),
	},
	{
		graph: NewGraph[any]("a", "b", "c", "d", "e"),
	},
	{
		graph: NewGraph[any]("a", "b", "c", "d", "e"),
		arcs:  []testArc{{"a", "b"}, {"b", "c"}, {"c", "d"}},
	},
	{
		graph: NewGraph[any]("a", "b", "c", "d", "e"),
		arcs:  []testArc{{"a", "b"}, {"c", "b"}, {"e", "d"}},
	},
	{
		graph: NewGraph[any](1, 2, 3, 4, 5, 6, 7, 8, 9, 10),
		arcs:  []testArc{{7, 8}, {7, 1}, {7, 3}, {2, 7}},
	},
	{
		graph: NewGraph[any](1, 2.63535, 3, 4, 5, 6, "A", 8, 9, 10),
		arcs:  []testArc{{"A", 8}, {"A", 1}, {"A", 3}, {5, "A"}, {2.63535, 1}},
	},
}

func TestGraphInteraction(t *testing.T) {

	dummyVtxList := make(Objects[any], 0)
	idx, ok := dummyVtxList.find("doesnotexist")
	require.False(t, ok)
	require.Equal(t, idx, -1)

	dummyList := newList[any]()
	idx, ok = dummyList.findIndex("doesnotexist")
	require.False(t, ok)
	require.Equal(t, idx, -1)

	graph := NewGraph[string]()
	result, err := graph.SortTopological()
	require.Nil(t, err)
	require.Zero(t, len(result))

	graph = NewGraph("a", "b", "c", "d", "e")

	// Try successful addition of vertex
	graph.AddVertex("f")

	// Try successful addition of arc
	require.Nil(t, graph.AddArc("d", "a"))

	// Try failed addition of arcs
	require.Error(t, graph.AddArc("d", "doesnotexist"))
	require.Error(t, graph.AddArc("dontexist", "a"))
}

func TestGraphTable(t *testing.T) {
	var err error
	for nRun := 0; nRun < nRunTestTable; nRun++ {
		for _, test := range testTable {
			require.Nil(t, test.init())
			test.result, err = test.graph.SortTopological()
			require.Nil(t, err)
			require.Equal(t, len(test.result), len(test.graph.vertices))
			require.Nil(t, test.assessOrder(t))
		}
	}
}

func TestGraphCyclic(t *testing.T) {
	cyclicGraph := NewGraph("a", "b", "c", "d")

	require.Nil(t, cyclicGraph.AddArc("a", "b"))
	require.Nil(t, cyclicGraph.AddArc("b", "c"))
	require.Nil(t, cyclicGraph.AddArc("c", "a"))

	_, err := cyclicGraph.SortTopological()
	require.ErrorContains(t, err, "cycle error")
}
