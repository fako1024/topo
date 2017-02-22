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
	_NUMPlugins // KEEP AT BOTTOM (AND DO NOT ADD TO PluginStrings)
)

func TestIotaType(t *testing.T) {

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
		Dependency{Child: B, Parent: A},
		Dependency{Child: B, Parent: C},
		Dependency{Child: B, Parent: D},
		Dependency{Child: A, Parent: E},
		Dependency{Child: D, Parent: C},
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
}
