////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 by Fabian Kohn
//
// This source code is licensed under the Apache License, Version 2.0, found in
// the LICENSE file in the root directory of this source tree.
////////////////////////////////////////////////////////////////////////////////

package topo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// StructType is a struct type
type StructType struct {
	String string
	Int    int
	Float  float64
}

func testStructType(t *testing.T) []StructType {

	// List of all structs (to be sorted)
	var allStructs = []StructType{
		{"A", 1, 1.0},
		{"B", 2, 2.0},
		{"C", 3, 3.0},
		{"D", 4, 4.0},
		{"E", 5, 5.0},
	}

	// List of all struct dependencies
	var structDependencies = []Dependency[StructType]{
		{Child: StructType{"A", 1, 1.0}, Parent: StructType{"C", 3, 3.0}},
		{Child: StructType{"D", 4, 4.0}, Parent: StructType{"E", 5, 5.0}},
	}

	// Perform topological sort
	require.Nil(t, Sort(allStructs, structDependencies))

	// Check if all StrDependencies are fulfilled
	for _, dependency := range structDependencies {
		posFrom, posTo := -1, -1
		for j := 0; j < len(allStructs); j++ {
			if allStructs[j] == dependency.Child {
				posFrom = j
			}
			if allStructs[j] == dependency.Parent {
				posTo = j
			}
		}

		require.Less(t, posTo, posFrom)
	}

	return allStructs
}

func TestStructType(t *testing.T) {
	testStringType(t)
}

func TestStructTypeStability(t *testing.T) {
	expected := testStructType(t)
	for run := 0; run < nRunsConsistency; run++ {
		require.EqualValues(t, expected, testStructType(t))
	}
}
