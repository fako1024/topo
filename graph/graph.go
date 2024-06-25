////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 by Fabian Kohn
//
// This source code is licensed under the Apache License, Version 2.0, found in
// the LICENSE file in the root directory of this source tree.
////////////////////////////////////////////////////////////////////////////////

package graph

import "fmt"

// Indicates the non-existence in find() methods
const indexNoExist = -1

// Graph represents the relations of arbitrary objects by designating vertices and
// arcs in accordance with a discrete graph description
type Graph[T comparable] struct {
	vertices map[T]vertex[T]
	order    []T
}

// NewGraph returns a new graph representation (constructor)
func NewGraph[T comparable](objects ...T) *Graph[T] {
	gr := Graph[T]{make(map[T]vertex[T]), make([]T, 0)}

	// Optionally add all vertices already provided variadically
	for _, obj := range objects {
		gr.AddVertex(obj)
	}

	return &gr
}

// AddVertex adds a node / vertex to the graph
func (g *Graph[T]) AddVertex(obj T) {
	if _, found := g.find(obj); !found {
		g.vertices[obj] = newVertex[T]()
		g.order = append(g.order, obj)
	}
}

// AddArc adds a line / arc to the graph
func (g *Graph[T]) AddArc(arcFrom, arcTo T) error {

	// Check if the "source" vertex exists
	sourceVertex, ok := g.vertices[arcFrom]
	if !ok {
		return fmt.Errorf("source vertex %v not found in graph", arcFrom)
	}

	// Check if the "destination" vertex exists
	if _, ok := g.vertices[arcTo]; !ok {
		return fmt.Errorf("destination vertex %v not found in graph", arcTo)
	}

	// Add the arc from "source" to "destination" vertex
	sourceVertex.addArc(arcTo)

	return nil
}

// SortTopological performs a topological sort and returns the sorted list of
// arbitrary input types
func (g *Graph[T]) SortTopological() (Objects[T], error) {
	var (
		results = newList[T]()
		err     error
	)

	// Recursively check each vertex for connected vertices and construct the
	// sorted list
	for _, obj := range g.order {
		var seen = newList[T]()
		if err = g.analyze(obj, results, seen); err != nil {
			return nil, err
		}
	}

	return results.elements, nil
}

////////////////// Private methods /////////////////////////////////////////////

// Find determines if a graph contains a specific vertex
func (g *Graph[T]) find(obj T) (vertex[T], bool) {
	val, ok := g.vertices[obj]
	return val, ok
}

// analyze recursively parses all graph vertices and their connections to other
// vertices, constructing the topologically sorted list in the process
func (g *Graph[T]) analyze(obj T, results, seen *list[T]) (err error) {

	// Try to add the current vertex to the sorted list
	if isNewElement := seen.add(obj); !isNewElement {

		// Cycle detected, obtain conflicting vertex indices (we just addded the index,
		// so we can forego the chekc for its existence)
		index, _ := seen.findIndex(obj)

		// Construct cycle
		cycle := append(seen.elements[index:], obj)

		// Return descriptive error indicating the cycle
		return fmt.Errorf("cycle error: %s", cycle.String())
	}

	// Recursively analyze next layer of graph
	for _, arc := range g.vertices[obj].arcs() {
		if err = g.analyze(arc, results, seen.clone()); err != nil {
			return err
		}
	}

	// Add the current vertex to the resulting list
	results.add(obj)

	return nil
}
