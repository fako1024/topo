////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 by Fabian Kohn
//
// This source code is licensed under the Apache License, Version 2.0, found in
// the LICENSE file in the root directory of this source tree.
////////////////////////////////////////////////////////////////////////////////

package topo

import "testing"

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
	var stringDependencies = []Dependency{
		Dependency{Child: "B", Parent: "A"},
		Dependency{Child: "B", Parent: "C"},
		Dependency{Child: "B", Parent: "D"},
		Dependency{Child: "A", Parent: "E"},
		Dependency{Child: "D", Parent: "C"},
	}

	// Getter function to convert original elements to a generic type
	getter := func(i int) Type {
		return allStrings[i]
	}

	// Setter function to restore the original type of the data
	setter := func(i int, val Type) {
		allStrings[i] = val.(string)
	}

	// Perform topological sort
	if err := Sort(allStrings, stringDependencies, getter, setter); err != nil {
		t.Fatal(err)
	}

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

		if posTo >= posFrom {
			t.Fatalf("Unexpected order, want pos(%v) < pos(%v) for %v / %v", posTo, posFrom, dependency.Child, dependency.Parent)
		}
	}

	return allStrings
}

func TestStringType(t *testing.T) {
	testStringType(t)
}

func TestStringTypeStability(t *testing.T) {
	expected := testStringType(t)

	for run := 0; run < nRunsConsistency; run++ {
		if res := testStringType(t); !testEqString(res, expected) {
			t.Fatalf("API stability violation, want %s, have %s", expected, res)
		}
	}
}

func testEqString(a, b []string) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
