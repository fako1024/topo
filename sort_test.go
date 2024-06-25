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

const nRunsConsistency = 1000

func TestDepString(t *testing.T) {

	// All plugin dependencies
	var pluginDependencies = []Dependency[PluginType]{
		{Child: B, Parent: A},
		{Child: B, Parent: C},
		{Child: B, Parent: D},
		{Child: A, Parent: E},
		{Child: D, Parent: C},
	}

	// All string dependencies
	var stringDependencies = []Dependency[string]{
		{Child: "B", Parent: "A"},
		{Child: "B", Parent: "C"},
		{Child: "B", Parent: "D"},
		{Child: "A", Parent: "E"},
		{Child: "D", Parent: "C"},
	}

	// List of all struct dependencies
	var structDependencies = []Dependency[StructType]{
		{Child: StructType{"A", 1, 1.0}, Parent: StructType{"C", 3, 3.0}},
		{Child: StructType{"D", 4, 4.0}, Parent: StructType{"E", 5, 5.0}},
	}

	require.Equal(t, pluginDependencies[0].String(), "Plugin 1 depends upon Plugin 0")
	require.Equal(t, stringDependencies[0].String(), "B depends upon A")
	require.Equal(t, structDependencies[0].String(), "{A 1 1} depends upon {C 3 3}")
}

func TestSortInline(t *testing.T) {

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

	// Perform topological sort (inline)
	err := Sort(allStructs, structDependencies)
	if err != nil {
		t.Fatal(err)
	}

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
}

func TestSortNoDeps(t *testing.T) {

	// List of all simple strings (to be sorted)
	var allStrings = []string{
		"A",
		"B",
		"C",
		"D",
		"E",
		"F",
		"G",
		"H",
	}

	// No dependencies
	var stringEmptyDependencies = []Dependency[string]{}

	// Save original data
	allStringsOld := make([]string, len(allStrings))
	copy(allStringsOld, allStrings)

	// Perform topological sort
	require.Nil(t, Sort(allStrings, stringEmptyDependencies))
	require.EqualValues(t, allStrings, allStringsOld)
}

func TestSortCyclic(t *testing.T) {

	// List of all simple strings (to be sorted)
	var allStrings = []string{
		"A",
		"B",
		"C",
		"D",
		"E",
		"F",
		"G",
		"H",
	}

	// Based on example_simple_test.go
	var stringCyclicDependencies = []Dependency[string]{
		{"B", "A"},
		{"B", "C"},
		{"B", "D"},
		{"A", "E"},
		{"D", "C"},
		{"C", "B"},
	}

	// Perform topological sort
	require.ErrorContains(t, Sort(allStrings, stringCyclicDependencies), "cycle error")
}

func TestSortNonExistVertex(t *testing.T) {

	// List of all simple strings (to be sorted)
	var allStrings = []string{
		"A",
		"B",
		"C",
		"D",
		"E",
		"F",
		"G",
		"H",
	}

	// Based on example_simple_test.go
	var stringNonExistVertexDependencies = []Dependency[string]{
		{"B", "A"},
		{"B", "C"},
		{"B", "D"},
		{"A", "E"},
		{"D", "C"},
		{"Z", "B"},
	}

	// Perform topological sort
	require.ErrorContains(t, Sort(allStrings, stringNonExistVertexDependencies), "source vertex Z not found in graph")
}
