////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2017 by Fabian Kohn
//
// This source code is licensed under the Apache License, Version 2.0, found in
// the LICENSE file in the root directory of this source tree.
////////////////////////////////////////////////////////////////////////////////

package graph

// list is a generic structure holding a sorted array of connected vertices and
// is used to determine the correct ordering and detect cycles
type list[T comparable] struct {
	indices  map[T]int
	elements Objects[T]
}

// newList return a new, empty list (constructor)
func newList[T comparable]() *list[T] {
	return &list[T]{make(map[T]int), make(Objects[T], 0)}
}

// add indicates if the element already exists and conditionally adds a new element
func (s *list[T]) add(obj T) bool {

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

func (s *list[T]) clone() *list[T] {

	// Pre-allocate result map with the correct number of elements
	listCopy := &list[T]{make(map[T]int, len(s.indices)), make(Objects[T], len(s.elements))}

	// Populate copy elements / indices
	for k, v := range s.elements {
		listCopy.elements[k] = v
		listCopy.indices[v] = s.indices[v]
	}

	return listCopy
}

func (s *list[T]) findIndex(obj T) (int, bool) {

	// Check if the element exists in the list and return its index if it does
	if i, ok := s.indices[obj]; ok {
		return i, true
	}

	// If not, indicate non-existence and return negative index
	return indexNoExist, false
}
