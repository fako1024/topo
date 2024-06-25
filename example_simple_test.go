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

func testStringType(t *testing.T) []string {

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

	// All string dependencies
	var stringDependencies = []Dependency[string]{
		{Child: "B", Parent: "A"},
		{Child: "B", Parent: "C"},
		{Child: "B", Parent: "D"},
		{Child: "A", Parent: "E"},
		{Child: "D", Parent: "C"},
	}

	// Perform topological sort
	require.Nil(t, Sort(allStrings, stringDependencies))

	// Check if all StrDependencies are fulfilled
	for _, dependency := range stringDependencies {
		posFrom, posTo := -1, -1
		for j := 0; j < len(allStrings); j++ {
			if allStrings[j] == dependency.Child {
				posFrom = j
			}
			if allStrings[j] == dependency.Parent {
				posTo = j
			}
		}

		require.Less(t, posTo, posFrom)
	}

	return allStrings
}

func TestStringType(t *testing.T) {
	testStringType(t)
}

func TestStringTypeStability(t *testing.T) {
	expected := testStringType(t)
	for run := 0; run < nRunsConsistency; run++ {
		require.EqualValues(t, expected, testStringType(t))
	}
}
