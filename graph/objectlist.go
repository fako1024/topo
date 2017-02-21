////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 by Fabian Kohn
//
// This source code is licensed under the Apache License, Version 2.0, found in
// the LICENSE file in the root directory of this source tree.
////////////////////////////////////////////////////////////////////////////////

package graph

import "fmt"

// ObjectList represents a simple list of arbitrary objects
type ObjectList []Object

// string returns a generic string denoting the connection between contained vertices
func (l ObjectList) String() string {

	// Return empty string if the vertex list is empty
	if len(l) == 0 {
		return ""
	}

	// Join all elements into chain
	objString := fmt.Sprint(l[0])
	for i := 1; i < len(l); i++ {
		objString += " -> " + fmt.Sprint(l[i])
	}

	return objString
}

// find determines if a VertexList contains a specific element
func (l ObjectList) find(obj Object) (int, bool) {

	// Check if the vertex exists in the list and return its index if it does
	for i := 0; i < len(l); i++ {
		if l[i] == obj {
			return i, true
		}
	}

	// If not, indicate non-existence and return negative index
	return indexNoExist, false
}
