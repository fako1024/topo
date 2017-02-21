////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 by Fabian Kohn
//
// This source code is licensed under the Apache License, Version 2.0, found in
// the LICENSE file in the root directory of this source tree.
////////////////////////////////////////////////////////////////////////////////

package graph

// vertex represents a node / vertex of a graph
type vertex map[Object]struct{}

// newVertex returns a new vertex (constructor)
func newVertex() vertex {
	return make(vertex)
}

// addArc creates a new line / arc to the graph
func (v vertex) addArc(arc Object) {
	v[arc] = struct{}{}
}

// arcs returns a list of all lines / arcs a graph contains
func (v vertex) arcs() ObjectList {

	// Pre-allocate the list of arcs with the correct number of elements
	list := make([]Object, len(v))
	pos := 0

	// Populate the list of arcs
	for k := range v {
		list[pos] = k
		pos++
	}

	return list
}
