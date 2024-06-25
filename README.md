# Topological Sort for arbitrary slices

[![Github Release](https://img.shields.io/github/release/fako1024/topo.svg)](https://github.com/fako1024/topo/releases)
[![GoDoc](https://godoc.org/github.com/fako1024/topo?status.svg)](https://godoc.org/github.com/fako1024/topo/)
[![Go Report Card](https://goreportcard.com/badge/github.com/fako1024/topo)](https://goreportcard.com/report/github.com/fako1024/topo)
[![Build/Test Status](https://github.com/fako1024/topo/workflows/Go/badge.svg)](https://github.com/fako1024/topo/actions?query=workflow%3AGo)

Introduction
------------

The topo package implements a topological sort algorithm to facilitate dependency resolution between elements of arbitrary data types.
The topo/graph package provides a directed graph representation for arbitrary data types to perform the actual sort process by describing elements as nodes / vertices and dependencies as links / arcs between these elements.

Installation and usage
----------------------

The import path for the package is *github.com/fako1024/topo*.
To install it, run:

    go get github.com/fako1024/topo

API summary
-----------------

The API of the topo package is fairly straight-forward. The following generics-based types / methods are exposed:

```Go
// Dependency represents a dependency between one Type and another
type Dependency[T comparable] struct {
	Child  T
	Parent T
}

// Dependencies represents a list of dependencies
type Dependencies[T comparable] []Dependency[T]

// String tries to stringify a dependency. If the type of the dependency fulfills
// the Stringer interface, it will use its String() method, otherwise it will try
// to format the variable otherwise
func (d Dependency[T]) String() string

// Sort performs a topological sort on a slice and constructs a directed graph (using the
// dependency constraints) and finally converts back the resulting object list to the
// original slice (sort in place)
func Sort[T comparable](data graph.Objects[T], deps Dependencies[T]) (err error)

```
In order to perform a dependency resolution, first a slice or array containing all elements to be sorted and a list of all dependencies have to be created.
Afterwards, the actual Sort() call can be performed, causing the original slice to be sorted in-place so as to satisfy all dependencies.
Note: Sort() is a stable sort algorithm, hence the actual order of elements in the output will be deterministic. A detailed, yet simple example can be found below.

License
-------

The topo and topo/graph packages are licensed under the Apache License 2.0. Please see the LICENSE file for details.

Example
-------

```Go
package main

import (
	"fmt"
	"os"
	"github.com/fako1024/topo"
)

// List of all simple strings (to be sorted)
var stringsToSort = []string{
	"A", "B", "C", "D", "E", "F", "G", "H",
}

// List of dependencies
var stringDependencies = []topo.Dependency[string]{
	{Child: "B", Parent: "A"},
	{Child: "B", Parent: "C"},
	{Child: "B", Parent: "D"},
	{Child: "A", Parent: "E"},
	{Child: "D", Parent: "C"},
}

func main() {

	// Perform topological sort
	if err := topo.Sort(stringsToSort, stringDependencies); err != nil {
		fmt.Printf("Error performing topological sort on slice of strings: %s\n", err)
		os.Exit(1)
	}

	// Print resulting Slice in order
	fmt.Println("Sorted list of strings:", stringsToSort)
	fmt.Println("The following dependencies were taken into account:")
	for _, dep := range stringDependencies {
		fmt.Println(dep)
	}
}
```

This example will generate the following output:

```
Sorted list of strings: [E A C D B F G H]
The following dependencies were taken into account:
B depends upon A
B depends upon C
B depends upon D
A depends upon E
D depends upon C
```

Additional examples can be found in the *_example_test.go files.
