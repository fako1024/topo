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

	"github.com/stretchr/testify/require"
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
	var pluginDependencies = []Dependency[PluginType]{
		{Child: B, Parent: A},
		{Child: B, Parent: C},
		{Child: B, Parent: D},
		{Child: A, Parent: E},
		{Child: D, Parent: C},
	}

	// Perform topological sort
	require.Nil(t, Sort(allPlugins, pluginDependencies))

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

		require.Less(t, posTo, posFrom)
	}

	return allPlugins
}

func TestIotaType(t *testing.T) {
	testIotaType(t)
}

func TestIotaTypeStability(t *testing.T) {
	expected := testIotaType(t)

	for run := 0; run < nRunsConsistency; run++ {
		require.EqualValues(t, expected, testIotaType(t))
	}
}
