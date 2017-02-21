# Topological Sort for arbitrary slices

[ ![Build Status](https://app.codeship.com/projects/6b3f9840-da5f-0134-d6f9-3e892a3f83ae/status?branch=master)](https://app.codeship.com/projects/203641)

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

The API of the topo package is fairly straight-forward. The following types / methods are exposed:

```Go
// Type defines a generic data type
type Type interface{}

// Dependency represents a dependency between one Type and another
type Dependency struct {
	Child  Type
	Parent Type
}

// String tries to stringify a dependency. If the type of the dependency fulfills
// the Stringer interface, it will use its String() method, otherwise it will try
// to format the variable otherwise
func (d Dependency) String() string

// Sort performs a topological sort on a slice using a functional approach to generalize
// the input data, construct a directed graph (using the dependency constraints) and
// finally converting back the resulting object list to the original slice (sort in place)
func Sort(data interface{}, deps []Dependency, getter func(i int) Type, setter func(i int, val Type)) (err error)

```
In order to perform a dependency resolution, first a slice or array containing all elements to be sorted and a list of all dependencies have to be created.
Afterwards, a "Getter" and a "Setter" function have to be defined in order to perform the actual type conversion for the type in question.
Finally, the actual Sort() call can be performed, causing the original slice to be sorted in-place so as to satisfy all dependencies.
Note that Sort() is not a stable sort algorithm, hence the actual order of elements in the output may vary. A detailed, yet simple example can be found below.

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
var stringDependencies = []topo.Dependency{
	topo.Dependency{Child: "B", Parent: "A"},
	topo.Dependency{Child: "B", Parent: "C"},
	topo.Dependency{Child: "B", Parent: "D"},
	topo.Dependency{Child: "A", Parent: "E"},
	topo.Dependency{Child: "D", Parent: "C"},
}

func main() {
	// Getter function to convert original elements to a generic type
	getter := func(i int) topo.Type {
		return stringsToSort[i]
	}

	// Setter function to restore the original type of the data
	setter := func(i int, val topo.Type) {
		stringsToSort[i] = val.(string)
	}

	// Perform topological sort
	if err := topo.Sort(stringsToSort, stringDependencies, getter, setter); err != nil {
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

}
```

This example will generate the following output (note that Sort() is not a stable
sort algorithm, hence the actual order of elements in the output may vary while
still satisfying all dependencies):

```
Sorted list of strings: [G H E A C D B F]
The following dependencies were taken into account:
B depends upon A
B depends upon C
B depends upon D
A depends upon E
D depends upon C
```

Additional examples can be found in the *_example_test.go files.
