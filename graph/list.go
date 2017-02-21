////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 by Fabian Kohn
//
// This source code is licensed under the Apache License, Version 2.0, found in
// the LICENSE file in the root directory of this source tree.
////////////////////////////////////////////////////////////////////////////////

package graph

// list is a generic structure holding a sorted array of connected vertices and
// is used to determine the correct ordering and detect cycles
type list struct {
	indices  map[Object]int
	elements ObjectList
}

// newList return a new, empty list (constructor)
func newList() *list {
	return &list{make(map[Object]int), make(ObjectList, 0)}
}

// add indicates if the element already exists and conditionally adds a new element
func (s *list) add(obj Object) bool {

	// Check if the element already exists in the list
	_, exists := s.indices[obj]

	// If not, add it to the list
	if !exists {
		s.indices[obj] = len(s.elements)
		s.elements = append(s.elements, obj)
	}

	// Indicate if the item existed
	return !exists
}

func (s *list) clone() *list {

	// Pre-allocate result map with the correct number of elements
	listCopy := &list{make(map[Object]int, len(s.indices)), make(ObjectList, len(s.elements))}

	// Populate copy indices
	for k, v := range s.indices {
		listCopy.indices[k] = v
	}

	// Populate copy elements
	for k, v := range s.elements {
		listCopy.elements[k] = v
	}

	return listCopy
}

func (s *list) findIndex(obj Object) (int, bool) {

	// Check if the element exists in the list and return its index if it does
	if i, ok := s.indices[obj]; ok {
		return i, true
	}

	// If not, indicate non-existence and return negative index
	return indexNoExist, false
}
