////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 by Fabian Kohn
//
// This source code is licensed under the Apache License, Version 2.0, found in
// the LICENSE file in the root directory of this source tree.
////////////////////////////////////////////////////////////////////////////////

package topo

import (
	"fmt"
	"testing"
)

// PluginType is a dummy type for the iota below
type PluginType uint

// Simple string method to properly format the plugin name
func (p PluginType) String() string {
	return fmt.Sprintf("Plugin %d", p)
}

// List of supported plugins
const (
	A PluginType = iota
	B
	C
	D
	E
	F
	G
	H
)

func testIotaType(t *testing.T) []PluginType {
	// All plugins (to be sorted)
	var allPlugins = []PluginType{
		A,
		B,
		C,
		D,
		E,
		F,
		G,
		H,
	}

	// All plugin dependencies
	var pluginDependencies = []Dependency{
		{Child: B, Parent: A},
		{Child: B, Parent: C},
		{Child: B, Parent: D},
		{Child: A, Parent: E},
		{Child: D, Parent: C},
	}

	// Getter function to convert original elements to a generic type
	getter := func(i int) Type {
		return allPlugins[i]
	}

	// Setter function to restore the original type of the data
	setter := func(i int, val Type) {
		allPlugins[i] = val.(PluginType)
	}

	// Perform topological sort
	if err := Sort(allPlugins, pluginDependencies, getter, setter); err != nil {
		t.Fatal(err)
	}

	// Check if all dependencies are fulfilled
	for _, dependency := range pluginDependencies {
		posFrom, posTo := -1, -1
		for j := 0; j < len(allPlugins); j++ {
			if allPlugins[j] == dependency.Child {
				posFrom = j
			}
			if allPlugins[j] == dependency.Parent {
				posTo = j
			}
		}

		if posTo >= posFrom {
			t.Fatalf("Unexpected order, want pos(%v) < pos(%v) for %v / %v", posTo, posFrom, dependency.Child, dependency.Parent)
		}
	}

	return allPlugins
}

func TestIotaType(t *testing.T) {
	testIotaType(t)
}

func TestIotaTypeStability(t *testing.T) {
	expected := testIotaType(t)

	for run := 0; run < nRunsConsistency; run++ {
		if res := testIotaType(t); !testEqIota(res, expected) {
			t.Fatalf("API stability violation, want %s, have %s", expected, res)
		}
	}
}

func testEqIota(a, b []PluginType) bool {

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
