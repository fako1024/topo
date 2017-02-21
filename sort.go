////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 by Fabian Kohn
//
// This source code is licensed under the Apache License, Version 2.0, found in
// the LICENSE file in the root directory of this source tree.
////////////////////////////////////////////////////////////////////////////////

package topo

import (
	"fmt"
	"reflect"

	"github.com/fako1024/topo/graph"
)

// Type defines a generic data type
type Type interface{}

// Dependency represents a dependency between one Type and another
type Dependency struct {
	Child  Type
	Parent Type
}

// String tries to stringify a dependency. If the type of the dependency fulfills
// the Stringer interface, it will use its String() method, otherwise it will try
// to format the variable otherwise
func (d Dependency) String() string {
	return fmt.Sprint(d.Child) + " depends upon " + fmt.Sprint(d.Parent)
}

// Sort performs a topological sort on a slice using a functional approach to generalize
// the input data, construct a directed graph (using the dependency constraints) and
// finally converting back the resulting object list to the original slice (sort in place)
func Sort(data interface{}, deps []Dependency, getter func(i int) Type, setter func(i int, val Type)) (err error) {

	// Obtain the number of elements in the original data slice using reflection
	nObj := reflect.ValueOf(data).Len()

	// Instantiate a new (empty) graph
	gr := graph.NewGraph()

	// Add all vertices (based on slice indices) using the 1st class getter function
	for i := 0; i < nObj; i++ {
		gr.AddVertex(getter(i))
	}

	// Add all dependencies (based on the enforced struct fields)
	for i := 0; i < len(deps); i++ {
		if err = gr.AddArc(deps[i].Child, deps[i].Parent); err != nil {
			return
		}
	}

	// Perform topological sorting, return error if e.g. a cycle is found
	var result graph.ObjectList
	if result, err = gr.SortTopological(); err != nil {
		return
	}

	// Sanity check to make sure the resulting slice contains the same number of
	// elements as the original data
	if len(result) != nObj {
		panic("Unexpected mismatch between original and sorted data")
	}

	// Convert the generic data back to the original slice type using the 1st class
	// setter function
	for i, val := range result {
		setter(i, val)
	}

	return
}
