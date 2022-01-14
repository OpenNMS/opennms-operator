/**
 Copied from https://github.com/lightbend/akka-cluster-operator/blob/a2e764272030151940171031a9a7a3168e47b15b/pkg/controller/akkacluster/subset.go#L1

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/
package subsets

import (
	"reflect"
	"unsafe"
)

// Problem statement: we have a baseline Kubernetes resource, and a live resource, and need
// to decide if they are out of sync, meaning the live resource doesn't match our baseline.
// The live resource has timestamps, uids, and other downstream ephemera that we want to ignore
// for purposes of deciding "is the live resource the same as the baseline?" We want the subset.
//
// Alternative approach to this problem is to convert objects to JSON--getting their exported
// and normalized form--and compare those. A challenge with that approach is dealing with
// default values which may be hard to detect once converted, like time.Time.

// SubsetEqual (A,B) returns true if A is a subset of B.
// This allows us to focus on a smaller set of required fields and ignore other fields
// that have downstream mutations to objects like creationTimestamp, uid, resourceVersion.
// The algorithm here is similar to reflect.DeepCopy, in that we use reflection to walk
// a potentially recursive tree. For comparison, we ignore empty or zero value fields in A.
func SubsetEqual(subset, superset interface{}) bool {
	t := newTreeWalk()
	return t.subsetValueEqual(reflect.ValueOf(subset), reflect.ValueOf(superset))
}

type treeWalk struct {
	matches int
	visited map[node]bool
}

func newTreeWalk() *treeWalk {
	t := treeWalk{}
	t.visited = make(map[node]bool)
	return &t
}

// track visited nodes to short circuit recursion
type node struct {
	a1  unsafe.Pointer
	a2  unsafe.Pointer
	typ reflect.Type
}

// Tests for subset equality using reflected types.
func (t *treeWalk) subsetValueEqual(subset, superset reflect.Value) bool {
	// if subset side is undefined, then nothing to compare
	if !subset.IsValid() {
		return true
	}

	// sanity check, rest of code assume same type on both sides
	if !superset.IsValid() {
		return false
	}

	if subset.Type() != superset.Type() {
		return false
	}

	// short circuit references already seen
	isRef := func(k reflect.Kind) bool {
		switch k {
		case reflect.Map, reflect.Slice, reflect.Ptr, reflect.Interface:
			return true
		}
		return false
	}
	if subset.CanAddr() && superset.CanAddr() && isRef(subset.Kind()) {
		n := node{
			unsafe.Pointer(subset.UnsafeAddr()),
			unsafe.Pointer(superset.UnsafeAddr()),
			subset.Type(),
		}
		if t.visited[n] {
			return true
		}
		t.visited[n] = true
	}

	// walk tree
	switch subset.Kind() {
	case reflect.Array, reflect.Slice:
		// recursive subset: superset may have extra elements at the end, only the subset members must match
		// Some array types (e.g. uint8) need to be fully compared
		if superset.Len() > 0 && superset.Index(0).Type().String() == "uint8" {
			// Full comparison - not subset
			if superset.Len() != subset.Len() {
				return false
			}
		} else if superset.Len() < subset.Len() {
			// Case for allowing subset comparison
			return false
		}
		for i := 0; i < subset.Len(); i++ {
			if !t.subsetValueEqual(subset.Index(i), superset.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Interface, reflect.Ptr:
		return t.subsetValueEqual(subset.Elem(), superset.Elem())
	case reflect.Struct:
		for i, n := 0, subset.NumField(); i < n; i++ {
			if !t.subsetValueEqual(subset.Field(i), superset.Field(i)) {
				return false
			}
		}
		return true
	case reflect.Map:
		for _, k := range subset.MapKeys() {
			if !t.subsetValueEqual(subset.MapIndex(k), superset.MapIndex(k)) {
				return false
			}
		}
		return true
	default:
		// Leaf node: if exported, non-default value, compare subset with superset.
		if subset.CanInterface() {
			// Ignore default values, like empty string, bool false, zero int... :-|
			// If needed, could be more selective in the kinds of zeros ignored.
			if subset.Interface() == reflect.Zero(subset.Type()).Interface() {
				return true
			}
			if subset.Interface() == superset.Interface() {
				t.matches++
				return true
			}
			return false
		}
		return true // Ignore non-exported opaque internal fields, like time.Time{}.
	}
}
