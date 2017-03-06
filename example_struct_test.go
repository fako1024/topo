////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 by Fabian Kohn
//
// This source code is licensed under the Apache License, Version 2.0, found in
// the LICENSE file in the root directory of this source tree.
////////////////////////////////////////////////////////////////////////////////

package topo

import "testing"

// StructType is a struct type
type StructType struct {
	String string
	Int    int
	Float  float64
}

func testStructType(t *testing.T) []StructType {

	// List of all structs (to be sorted)
	var allStructs = []StructType{
		StructType{"A", 1, 1.0},
		StructType{"B", 2, 2.0},
		StructType{"C", 3, 3.0},
		StructType{"D", 4, 4.0},
		StructType{"E", 5, 5.0},
	}

	// List of all struct dependencies
	var structDependencies = []Dependency{
		Dependency{Child: StructType{"A", 1, 1.0}, Parent: StructType{"C", 3, 3.0}},
		Dependency{Child: StructType{"D", 4, 4.0}, Parent: StructType{"E", 5, 5.0}},
	}

	// Getter function to convert original elements to a generic type
	getter := func(i int) Type {
		return allStructs[i]
	}

	// Setter function to restore the original type of the data
	setter := func(i int, val Type) {
		allStructs[i] = val.(StructType)
	}

	// Perform topological sort
	if err := Sort(allStructs, structDependencies, getter, setter); err != nil {
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

		if posTo >= posFrom {
			t.Fatalf("Unexpected order, want pos(%v) < pos(%v) for %v / %v", posTo, posFrom, dependency.Child, dependency.Parent)
		}
	}

	return allStructs
}

func TestStructType(t *testing.T) {
	testStringType(t)
}

func TestStructTypeStability(t *testing.T) {
	expected := testStructType(t)

	for run := 0; run < nRunsConsistency; run++ {
		if res := testStructType(t); !testEqStruct(res, expected) {
			t.Fatalf("API stability violation, want %v, have %v", expected, res)
		}
	}
}

func testEqStruct(a, b []StructType) bool {

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
