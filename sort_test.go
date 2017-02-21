////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 by Fabian Kohn
//
// This source code is licensed under the Apache License, Version 2.0, found in
// the LICENSE file in the root directory of this source tree.
////////////////////////////////////////////////////////////////////////////////

package topo

import (
	"strings"
	"testing"
)

func TestDepString(t *testing.T) {
	if pluginDependencies[0].String() != "Plugin 1 depends upon Plugin 0" {
		t.Fatalf("Unexpected dependency string, want \"Plugin 1 depends upon Plugin 0\", got \"%s\"", pluginDependencies[0].String())
	}
	if stringDependencies[0].String() != "B depends upon A" {
		t.Fatalf("Unexpected dependency string, want \"B depends upon A\", got \"%s\"", stringDependencies[0].String())
	}
	if structDependencies[0].String() != "{A 1 1} depends upon {C 3 3}" {
		t.Fatalf("Unexpected dependency string, want \"{A 1 1} depends upon {C 3 3}\", got \"%s\"", structDependencies[0].String())
	}
}

func TestSortInline(t *testing.T) {

	// Based on example_struct_test.go

	// Perform topological sort (inline)
	err := Sort(allStructs, structDependencies, func(i int) Type { return allStructs[i] }, func(i int, val Type) { allStructs[i] = val.(StructType) })

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

		if posTo >= posFrom {
			t.Fatalf("Unexpected order, want pos(%v) < pos(%v) for %v / %v", posTo, posFrom, dependency.Child, dependency.Parent)
		}
	}
}

func TestSortCyclic(t *testing.T) {

	// Based on example_simple_test.go
	var stringCyclicDependencies = []Dependency{
		Dependency{"B", "A"},
		Dependency{"B", "C"},
		Dependency{"B", "D"},
		Dependency{"A", "E"},
		Dependency{"D", "C"},
		Dependency{"C", "B"},
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
	if err := Sort(allStrings, stringCyclicDependencies, getter, setter); err == nil {
		t.Fatal("Expected cyclic error not seen")
	} else if !strings.Contains(err.Error(), "Cycle error:") {
		t.Errorf("Unexpected error message: %s", err)
	}
}

func TestSortNonExistVertex(t *testing.T) {

	// Based on example_simple_test.go
	var stringNonExistVertexDependencies = []Dependency{
		Dependency{"B", "A"},
		Dependency{"B", "C"},
		Dependency{"B", "D"},
		Dependency{"A", "E"},
		Dependency{"D", "C"},
		Dependency{"Z", "B"},
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
	if err := Sort(allStrings, stringNonExistVertexDependencies, getter, setter); err == nil {
		t.Fatal("Expected error (non-existing vertex) not seen")
	} else if !strings.Contains(err.Error(), "Source vertex Z not found in graph") {
		t.Errorf("Unexpected error message: %s", err)
	}
}
