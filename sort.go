////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 by Fabian Kohn
//
// This source code is licensed under the Apache License, Version 2.0, found in
// the LICENSE file in the root directory of this source tree.
////////////////////////////////////////////////////////////////////////////////

package topo

import (
	"errors"
	"fmt"

	"github.com/fako1024/topo/graph"
)

var (
	// ErrUnexpectedMismatch is thrown if the sorted graph is inconsistent with the input data
	ErrUnexpectedMismatch = errors.New("unexpected mismatch between original and sorted data")
)

// Dependency represents a dependency between one Type and another
type Dependency[T comparable] struct {
	Child  T
	Parent T
}

// Dependencies represents a list of dependencies
type Dependencies[T comparable] []Dependency[T]

// String tries to stringify a dependency. If the type of the dependency fulfills
// the Stringer interface, it will use its String() method, otherwise it will try
// to format the variable otherwise
func (d Dependency[T]) String() string {
	return fmt.Sprintf("%v depends upon %v", d.Child, d.Parent)
}

// Sort performs a topological sort on a slice and constructs a directed graph (using the
// dependency constraints) and finally converts back the resulting object list to the
// original slice (sort in place)
func Sort[T comparable](data graph.Objects[T], deps Dependencies[T]) (err error) {

	// In case there are no dependencies, return immediately without action
	if len(deps) == 0 {
		return nil
	}

	// Instantiate a new (empty) graph
	gr := graph.NewGraph[T]()

	// Add all vertices (based on slice indices)
	for i := 0; i < len(data); i++ {
		gr.AddVertex(data[i])
	}

	// Add all dependencies (based on the enforced struct fields)
	for i := 0; i < len(deps); i++ {
		if err = gr.AddArc(deps[i].Child, deps[i].Parent); err != nil {
			return
		}
	}

	// Perform topological sorting, return error if e.g. a cycle is found
	var result graph.Objects[T]
	if result, err = gr.SortTopological(); err != nil {
		return
	}

	// Sanity check to make sure the resulting slice contains the same number of
	// elements as the original data
	if len(result) != len(data) {
		return ErrUnexpectedMismatch
	}

	// Copy the sorted data back to the original slice
	copy(data, result)

	return
}
