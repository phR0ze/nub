package n

import (
	"sort"
	"strings"

	"github.com/pkg/errors"
)

// MapSlice implements the Slice interface providing a generic way to work with slice types
// including convenience methods on par with rapid development languages.
type MapSlice []map[string]interface{}

// NewMapSlice creates a new *MapSlice
func NewMapSlice(slice interface{}) *MapSlice {
	return ToMapSlice(slice)
}

// NewMapSliceV creates a new *MapSlice from the given variadic elements. Always returns
// at least a reference to an empty MapSlice.
func NewMapSliceV(elems ...interface{}) *MapSlice {
	return ToMapSlice(elems)
}

// A is an alias to String for brevity
func (p *MapSlice) A() string {
	return p.String()
}

// All tests if this Slice is not empty or optionally if it contains
// all of the given variadic elements. Incompatible types will return false.
func (p *MapSlice) All(elems ...interface{}) bool {
	if p == nil || len(*p) == 0 {
		return false
	}

	// Not looking for anything
	if len(elems) == 0 {
		return true
	}

	// Looking for something specific returns false if incompatible type
	return p.AllS(elems)
}

// AllS tests if this Slice contains all of the given Slice's elements;
// Incompatible types will return false;
// Supports MapSlice, *MapSlice, []string or *[]string
func (p *MapSlice) AllS(slice interface{}) bool {
	if p == nil || len(*p) == 0 {
		return false
	}

	panic("NOT IMPLEMENTED")
	// if elems, err := ToMapSliceE(slice); err == nil {
	// 	for i := range *elems {
	// 		for j := range *p {
	// 			if (*p)[j] == (*elems)[i] {
	// 				return true
	// 			}
	// 		}
	// 	}
	// }
	// return false
}

// Any tests if this Slice is not empty or optionally if it contains
// any of the given variadic elements. Incompatible types will return false.
func (p *MapSlice) Any(elems ...interface{}) bool {
	if p == nil || len(*p) == 0 {
		return false
	}

	// Not looking for anything
	if len(elems) == 0 {
		return true
	}

	// Looking for something specific returns false if incompatible type
	return p.AnyS(elems)
}

// AnyS tests if this Slice contains any of the given Slice's elements.
// Incompatible types will return false.
// Supports MapSlice, *MapSlice, []string or *[]string
func (p *MapSlice) AnyS(slice interface{}) bool {
	if p == nil || len(*p) == 0 {
		return false
	}

	panic("NOT IMPLEMENTED")
	// if elems, err := ToMapSliceE(slice); err == nil {
	// 	for i := range *elems {
	// 		for j := range *p {
	// 			if (*p)[j] == (*elems)[i] {
	// 				return true
	// 			}
	// 		}
	// 	}
	// }
	// return false
}

// AnyW tests if this Slice contains any that match the lambda selector.
func (p *MapSlice) AnyW(sel func(O) bool) bool {
	return p.CountW(sel) != 0
}

// Append an element to the end of this Slice and returns a reference to this Slice.
func (p *MapSlice) Append(elem interface{}) ISlice {
	if p == nil {
		p = NewMapSliceV()
	}
	*p = append(*p, *ToStringMap(elem))
	return p
}

// AppendV appends the variadic elements to the end of this Slice and returns a reference to this Slice.
func (p *MapSlice) AppendV(elems ...interface{}) ISlice {
	if p == nil {
		p = NewMapSliceV()
	}
	for _, elem := range elems {
		*p = append(*p, *ToStringMap(elem))
	}
	return p
}

// At returns the element at the given index location. Allows for negative notation.
func (p *MapSlice) At(i int) (elem *Object) {
	elem = &Object{}
	if p == nil {
		return
	}
	if i = absIndex(len(*p), i); i == -1 {
		return
	}
	elem.o = (*p)[i]
	return
}

// Clear modifies this Slice to clear out all elements and returns a reference to this Slice.
func (p *MapSlice) Clear() ISlice {
	if p == nil {
		p = NewMapSliceV()
	} else {
		p.Drop()
	}
	return p
}

// Concat returns a new Slice by appending the given Slice to this Slice using variadic expansion.
// Supports MapSlice, *MapSlice, []string or *[]string
func (p *MapSlice) Concat(slice interface{}) (new ISlice) {
	return p.Copy().ConcatM(slice)
}

// ConcatM modifies this Slice by appending the given Slice using variadic expansion and returns a reference to this Slice.
// Supports MapSlice, *MapSlice, []string or *[]string
func (p *MapSlice) ConcatM(slice interface{}) ISlice {
	if p == nil {
		p = NewMapSliceV()
	}
	if elems, err := ToMapSliceE(slice); err == nil {
		*p = append(*p, *elems...)
	}
	return p
}

// Copy returns a new Slice with the indicated range of elements copied from this Slice.
// Expects nothing, in which case everything is copied, or two indices i and j, in which
// case positive and negative notation is supported and uses an inclusive behavior such
// that Slice(0, -1) includes index -1 as opposed to Go's exclusive behavior. Out of
// bounds indices will be moved within bounds.
//
// An empty Slice is returned if indicies are mutually exclusive or nothing can be returned.
func (p *MapSlice) Copy(indices ...int) (new ISlice) {
	if p == nil || len(*p) == 0 {
		return NewMapSliceV()
	}

	// Handle index manipulation
	i, j, err := absIndices(len(*p), indices...)
	if err != nil {
		return NewMapSliceV()
	}

	// Copy elements over to new Slice
	x := make([]map[string]interface{}, j-i, j-i)
	copy(x, (*p)[i:j])
	return NewMapSlice(x)
}

// Count the number of elements in this Slice equal to the given element.
func (p *MapSlice) Count(elem interface{}) (cnt int) {
	if y, ok := elem.(string); ok {
		cnt = p.CountW(func(x O) bool { return ExB(x.(string) == y) })
	}
	return
}

// CountW counts the number of elements in this Slice that match the lambda selector.
func (p *MapSlice) CountW(sel func(O) bool) (cnt int) {
	if p == nil || len(*p) == 0 {
		return
	}
	for i := range *p {
		if sel((*p)[i]) {
			cnt++
		}
	}
	return
}

// Drop modifies this Slice to delete the indicated range of elements and returns a referece to this Slice.
// Expects nothing, in which case everything is dropped, or two indices i and j, in which case positive and
// negative notation is supported and uses an inclusive behavior such that DropAt(0, -1) includes index -1
// as opposed to Go's exclusive behavior. Out of bounds indices will be moved within bounds.
func (p *MapSlice) Drop(indices ...int) ISlice {
	if p == nil || len(*p) == 0 {
		return p
	}

	// Handle index manipulation
	i, j, err := absIndices(len(*p), indices...)
	if err != nil {
		return p
	}

	// Execute
	n := j - i
	if i+n < len(*p) {
		*p = append((*p)[:i], (*p)[i+n:]...)
	} else {
		*p = (*p)[:i]
	}
	return p
}

// DropAt modifies this Slice to delete the element at the given index location. Allows for negative notation.
// Returns a reference to this Slice.
func (p *MapSlice) DropAt(i int) ISlice {
	return p.Drop(i, i)
}

// DropFirst modifies this Slice to delete the first element and returns a reference to this Slice.
func (p *MapSlice) DropFirst() ISlice {
	return p.Drop(0, 0)
}

// DropFirstN modifies this Slice to delete the first n elements and returns a reference to this Slice.
func (p *MapSlice) DropFirstN(n int) ISlice {
	if n == 0 {
		return p
	}
	return p.Drop(0, abs(n)-1)
}

// DropLast modifies this Slice to delete the last element and returns a reference to this Slice.
func (p *MapSlice) DropLast() ISlice {
	return p.Drop(-1, -1)
}

// DropLastN modifies thi Slice to delete the last n elements and returns a reference to this Slice.
func (p *MapSlice) DropLastN(n int) ISlice {
	if n == 0 {
		return p
	}
	return p.Drop(absNeg(n), -1)
}

// DropW modifies this Slice to delete the elements that match the lambda selector and returns a reference to this Slice.
// The slice is updated instantly when lambda expression is evaluated not after DropW completes.
func (p *MapSlice) DropW(sel func(O) bool) ISlice {
	if p == nil || len(*p) == 0 {
		return p
	}
	l := len(*p)
	for i := 0; i < l; i++ {
		if sel((*p)[i]) {
			p.DropAt(i)
			l--
			i--
		}
	}
	return p
}

// Each calls the given lambda once for each element in this Slice, passing in that element
// as a parameter. Returns a reference to this Slice
func (p *MapSlice) Each(action func(O)) ISlice {
	if p == nil {
		return p
	}
	for i := range *p {
		action((*p)[i])
	}
	return p
}

// EachE calls the given lambda once for each element in this Slice, passing in that element
// as a parameter. Returns a reference to this Slice and any error from the lambda.
func (p *MapSlice) EachE(action func(O) error) (ISlice, error) {
	var err error
	if p == nil {
		return p, err
	}
	for i := range *p {
		if err = action((*p)[i]); err != nil {
			return p, err
		}
	}
	return p, err
}

// EachI calls the given lambda once for each element in this Slice, passing in the index and element
// as a parameter. Returns a reference to this Slice
func (p *MapSlice) EachI(action func(int, O)) ISlice {
	if p == nil {
		return p
	}
	for i := range *p {
		action(i, (*p)[i])
	}
	return p
}

// EachIE calls the given lambda once for each element in this Slice, passing in the index and element
// as a parameter. Returns a reference to this Slice and any error from the lambda.
func (p *MapSlice) EachIE(action func(int, O) error) (ISlice, error) {
	var err error
	if p == nil {
		return p, err
	}
	for i := range *p {
		if err = action(i, (*p)[i]); err != nil {
			return p, err
		}
	}
	return p, err
}

// EachR calls the given lambda once for each element in this Slice in reverse, passing in that element
// as a parameter. Returns a reference to this Slice
func (p *MapSlice) EachR(action func(O)) ISlice {
	if p == nil {
		return p
	}
	for i := len(*p) - 1; i >= 0; i-- {
		action((*p)[i])
	}
	return p
}

// EachRE calls the given lambda once for each element in this Slice in reverse, passing in that element
// as a parameter. Returns a reference to this Slice and any error from the lambda.
func (p *MapSlice) EachRE(action func(O) error) (ISlice, error) {
	var err error
	if p == nil {
		return p, err
	}
	for i := len(*p) - 1; i >= 0; i-- {
		if err = action((*p)[i]); err != nil {
			return p, err
		}
	}
	return p, err
}

// EachRI calls the given lambda once for each element in this Slice in reverse, passing in that element
// as a parameter. Returns a reference to this Slice
func (p *MapSlice) EachRI(action func(int, O)) ISlice {
	if p == nil {
		return p
	}
	for i := len(*p) - 1; i >= 0; i-- {
		action(i, (*p)[i])
	}
	return p
}

// EachRIE calls the given lambda once for each element in this Slice in reverse, passing in that element
// as a parameter. Returns a reference to this Slice and any error from the lambda.
func (p *MapSlice) EachRIE(action func(int, O) error) (ISlice, error) {
	var err error
	if p == nil {
		return p, err
	}
	for i := len(*p) - 1; i >= 0; i-- {
		if err = action(i, (*p)[i]); err != nil {
			return p, err
		}
	}
	return p, err
}

// Empty tests if this Slice is empty.
func (p *MapSlice) Empty() bool {
	if p == nil || len(*p) == 0 {
		return true
	}
	return false
}

// First returns the first element in this Slice as Object.
// Object.Nil() == true will be returned when there are no elements in the slice.
func (p *MapSlice) First() (elem *Object) {
	return p.At(0)
}

// FirstN returns the first n elements in this slice as a Slice reference to the original.
// Best effort is used such that as many as can be will be returned up until the request is satisfied.
func (p *MapSlice) FirstN(n int) ISlice {
	if n == 0 {
		return NewMapSliceV()
	}
	return p.Slice(0, abs(n)-1)
}

// G returns the underlying data structure as a builtin Go type
func (p *MapSlice) G() []map[string]interface{} {
	return p.O().([]map[string]interface{})
}

// Index returns the index of the first element in this Slice where element == elem
// Returns a -1 if the element was not not found.
func (p *MapSlice) Index(elem interface{}) (loc int) {
	loc = -1
	if p == nil || len(*p) == 0 {
		return
	}

	panic("NOT IMPLEMENTED")
	// x := ToStringMap(elem)
	// for i := range *p {
	// 	if (*p)[i] == x {
	// 		return i
	// 	}
	// }
	// return
}

// Insert modifies this Slice to insert the given elements before the element(s) with the given index.
// Negative indices count backwards from the end of the slice, where -1 is the last element. If a
// negative index is used, the given element will be inserted after that element, so using an index
// of -1 will insert the element at the end of the slice. If a Slice is given all elements will be
// inserted starting from the beging until the end. Slice is returned for chaining. Invalid
// index locations will not change the slice.
func (p *MapSlice) Insert(i int, obj interface{}) ISlice {
	if p == nil || len(*p) == 0 {
		return p.ConcatM(obj)
	}

	// Insert the item before j if pos and after j if neg
	j := i
	if j = absIndex(len(*p), j); j == -1 {
		return p
	}
	if i < 0 {
		j++
	}
	if elems, err := ToMapSliceE(obj); err == nil {
		if j == 0 {
			*p = append(*elems, *p...)
		} else if j < len(*p) {
			*p = append(*p, *elems...)           // ensures enough space exists
			copy((*p)[j+len(*elems):], (*p)[j:]) // shifts right elements drop added
			copy((*p)[j:], *elems)               // set new in locations vacated
		} else {
			*p = append(*p, *elems...)
		}
	}
	return p
}

// InterSlice returns true if the underlying implementation is a RefSlice
func (p *MapSlice) InterSlice() bool {
	return false
}

// Join converts each element into a string then joins them together using the given separator or comma by default.
func (p *MapSlice) Join(separator ...string) (str *Object) {
	if p == nil || len(*p) == 0 {
		str = &Object{""}
		return
	}
	sep := ","
	if len(separator) > 0 {
		sep = separator[0]
	}

	var builder strings.Builder
	for i := range *p {
		builder.WriteString(ToString((*p)[i]))
		if i+1 < len(*p) {
			builder.WriteString(sep)
		}
	}
	str = &Object{builder.String()}
	return
}

// Last returns the last element in this Slice as an Object.
// Object.Nil() == true will be returned if there are no elements in the slice.
func (p *MapSlice) Last() (elem *Object) {
	return p.At(-1)
}

// LastN returns the last n elements in this Slice as a Slice reference to the original.
// Best effort is used such that as many as can be will be returned up until the request is satisfied.
func (p *MapSlice) LastN(n int) ISlice {
	if n == 0 {
		return NewMapSliceV()
	}
	return p.Slice(absNeg(n), -1)
}

// Len returns the number of elements in this Slice
func (p *MapSlice) Len() int {
	if p == nil {
		return 0
	}
	return len(*p)
}

// Less returns true if the element indexed by i is less than the element indexed by j.
func (p *MapSlice) Less(i, j int) bool {
	panic("NOT IMPLEMENTED")
}

// Map creates a new slice with the modified elements from the lambda.
func (p *MapSlice) Map(mod func(O) O) ISlice {
	var slice ISlice
	if p == nil || len(*p) == 0 {
		return NewMapSliceV()
	}
	for i := range *p {
		v := mod((*p)[i])
		if slice == nil {
			slice = Slice(v)
		} else {
			slice.Append(v)
		}
	}
	return slice
}

// Nil tests if this Slice is nil
func (p *MapSlice) Nil() bool {
	if p == nil {
		return true
	}
	return false
}

// O returns the underlying data structure as is
func (p *MapSlice) O() interface{} {
	if p == nil {
		return []map[string]interface{}{}
	}
	return []map[string]interface{}(*p)
}

// Pair simply returns the first and second Slice elements as Objects
func (p *MapSlice) Pair() (first, second *Object) {
	first, second = &Object{}, &Object{}
	if p == nil {
		return
	}
	if len(*p) > 0 {
		first = p.At(0)
	}
	if len(*p) > 1 {
		second = p.At(1)
	}
	return
}

// Pop modifies this Slice to remove the last element and returns the removed element as an Object.
func (p *MapSlice) Pop() (elem *Object) {
	elem = p.Last()
	p.DropLast()
	return
}

// PopN modifies this Slice to remove the last n elements and returns the removed elements as a new Slice.
func (p *MapSlice) PopN(n int) (new ISlice) {
	if n == 0 {
		return NewMapSliceV()
	}
	new = p.Copy(absNeg(n), -1)
	p.DropLastN(n)
	return
}

// Prepend modifies this Slice to add the given element at the begining and returns a reference to this Slice.
func (p *MapSlice) Prepend(elem interface{}) ISlice {
	return p.Insert(0, elem)
}

// RefSlice returns true if the underlying implementation is a RefSlice
func (p *MapSlice) RefSlice() bool {
	return false
}

// Reverse returns a new Slice with the order of the elements reversed.
func (p *MapSlice) Reverse() (new ISlice) {
	if p == nil || len(*p) < 2 {
		return p.Copy()
	}
	return p.Copy().ReverseM()
}

// ReverseM modifies this Slice reversing the order of the elements and returns a reference to this Slice.
func (p *MapSlice) ReverseM() ISlice {
	if p == nil || len(*p) == 0 {
		return p
	}
	for i, j := 0, len(*p)-1; i < j; i, j = i+1, j-1 {
		p.Swap(i, j)
	}
	return p
}

// S is an alias to ToStringSlice
func (p *MapSlice) S() (slice *StringSlice) {
	return ToStringSlice(p.O())
}

// Select creates a new slice with the elements that match the lambda selector.
func (p *MapSlice) Select(sel func(O) bool) (new ISlice) {
	slice := NewMapSliceV()
	if p == nil || len(*p) == 0 {
		return slice
	}
	for i := range *p {
		if sel((*p)[i]) {
			*slice = append(*slice, (*p)[i])
		}
	}
	return slice
}

// Set the element(s) at the given index location to the given element(s). Allows for negative notation.
// Returns a reference to this Slice and swallows any errors.
func (p *MapSlice) Set(i int, elem interface{}) ISlice {
	slice, _ := p.SetE(i, elem)
	return slice
}

// SetE the element(s) at the given index location to the given element(s). Allows for negative notation.
// Returns a referenc to this Slice and an error if out of bounds or elem is the wrong type.
func (p *MapSlice) SetE(i int, elems interface{}) (ISlice, error) {
	var err error
	if p == nil {
		return p, err
	}
	if i = absIndex(len(*p), i); i == -1 {
		err = errors.Errorf("slice assignment is out of bounds")
		return p, err
	}

	// Account for length of elems
	if x, err := ToMapSliceE(elems); err == nil {
		if len(*x) > 0 {
			copy((*p)[i:], *x)
		}
	} else {
		err = errors.Wrapf(err, "can't set type '%T' in '%T'", elems, p)
	}
	return p, err
}

// Shift modifies this Slice to remove the first element and returns the removed element as an Object.
func (p *MapSlice) Shift() (elem *Object) {
	elem = p.First()
	p.DropFirst()
	return
}

// ShiftN modifies this Slice to remove the first n elements and returns the removed elements as a new Slice.
func (p *MapSlice) ShiftN(n int) (new ISlice) {
	if n == 0 {
		return NewMapSliceV()
	}
	new = p.Copy(0, abs(n)-1)
	p.DropFirstN(n)
	return
}

// Single reports true if there is only one element in this Slice.
func (p *MapSlice) Single() bool {
	return p.Len() == 1
}

// Slice returns a range of elements from this Slice as a Slice reference to the original. Allows for negative notation.
// Expects nothing, in which case everything is included, or two indices i and j, in which case an inclusive behavior
// is used such that Slice(0, -1) includes index -1 as opposed to Go's exclusive behavior. Out of bounds indices will
// be moved within bounds.
//
// An empty Slice is returned if indicies are mutually exclusive or nothing can be returned.
//
// e.g. NewMapSliceV(1,2,3).Slice(0, -1) == [1,2,3] && NewMapSliceV(1,2,3).Slice(1,2) == [2,3]
func (p *MapSlice) Slice(indices ...int) ISlice {
	if p == nil || len(*p) == 0 {
		return NewMapSliceV()
	}

	// Handle index manipulation
	i, j, err := absIndices(len(*p), indices...)
	if err != nil {
		return NewMapSliceV()
	}

	slice := MapSlice((*p)[i:j])
	return &slice
}

// Sort returns a new Slice with sorted elements.
func (p *MapSlice) Sort() (new ISlice) {
	if p == nil || len(*p) < 2 {
		return p.Copy()
	}
	return p.Copy().SortM()
}

// SortM modifies this Slice sorting the elements and returns a reference to this Slice.
func (p *MapSlice) SortM() ISlice {
	if p == nil || len(*p) < 2 {
		return p
	}
	sort.Sort(p)
	return p
}

// SortReverse returns a new Slice sorting the elements in reverse.
func (p *MapSlice) SortReverse() (new ISlice) {
	if p == nil || len(*p) < 2 {
		return p.Copy()
	}
	return p.Copy().SortReverseM()
}

// SortReverseM modifies this Slice sorting the elements in reverse and returns a reference to this Slice.
func (p *MapSlice) SortReverseM() ISlice {
	if p == nil || len(*p) < 2 {
		return p
	}
	sort.Sort(sort.Reverse(p))
	return p
}

// String returns a string representation of this Slice, implements the Stringer interface
func (p *MapSlice) String() string {
	var builder strings.Builder
	builder.WriteString("[")
	if p != nil {
		for i := range *p {
			builder.WriteString(ToString((*p)[i]))
			if i+1 < len(*p) {
				builder.WriteString(" ")
			}
		}
	}
	builder.WriteString("]")
	return builder.String()
}

// Swap modifies this Slice swapping the indicated elements.
func (p *MapSlice) Swap(i, j int) {
	if p == nil || len(*p) < 2 || i < 0 || j < 0 || i >= len(*p) || j >= len(*p) {
		return
	}
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

// Take modifies this Slice removing the indicated range of elements from this Slice and returning them as a new Slice.
// Expects nothing, in which case everything is taken, or two indices i and j, in which case positive and negative
// notation is supported and uses an inclusive behavior such that Take(0, -1) includes index -1 as opposed to Go's
// exclusive behavior. Out of bounds indices will be moved within bounds.
func (p *MapSlice) Take(indices ...int) (new ISlice) {
	new = p.Copy(indices...)
	p.Drop(indices...)
	return
}

// TakeAt modifies this Slice removing the elemement at the given index location and returns the removed element as an Object.
// Allows for negative notation.
func (p *MapSlice) TakeAt(i int) (elem *Object) {
	elem = p.At(i)
	p.DropAt(i)
	return
}

// TakeW modifies this Slice removing the elements that match the lambda selector and returns them as a new Slice.
func (p *MapSlice) TakeW(sel func(O) bool) (new ISlice) {
	slice := NewMapSliceV()
	if p == nil || len(*p) == 0 {
		return slice
	}
	l := len(*p)
	for i := 0; i < l; i++ {
		if sel((*p)[i]) {
			*slice = append(*slice, (*p)[i])
			p.DropAt(i)
			l--
			i--
		}
	}
	return slice
}

// ToInts converts the underlying slice into a []int
func (p *MapSlice) ToInts() (slice []int) {
	return ToIntSlice(p.O()).G()
}

// ToIntSlice converts the underlying slice into a *IntSlice
func (p *MapSlice) ToIntSlice() (slice *IntSlice) {
	return ToIntSlice(p.O())
}

// ToInterSlice converts the given slice to a generic []interface{} slice
func (p *MapSlice) ToInterSlice() (slice []interface{}) {
	return ToInterSlice(p.O()).G()
}

// ToStringSlice converts the underlying slice into a *StringSlice
func (p *MapSlice) ToStringSlice() (slice *StringSlice) {
	return ToStringSlice(p.O())
}

// ToStrs converts the underlying slice into a []string slice
func (p *MapSlice) ToStrs() (slice []string) {
	return ToStrs(p.O())
}

// Union returns a new Slice by joining uniq elements from this Slice with uniq elements from the given Slice while preserving order.
// Supports MapSlice, *MapSlice, []string or *[]string
func (p *MapSlice) Union(slice interface{}) (new ISlice) {
	return p.Copy().UnionM(slice)
}

// UnionM modifies this Slice by joining uniq elements from this Slice with uniq elements from the given Slice while preserving order.
// Supports MapSlice, *MapSlice, []string or *[]string
func (p *MapSlice) UnionM(slice interface{}) ISlice {
	return p.ConcatM(slice).UniqM()
}

// Uniq returns a new Slice with all non uniq elements removed while preserving element order.
// Cost for this call vs the UniqM is roughly the same, this one is appending that one dropping.
func (p *MapSlice) Uniq() (new ISlice) {
	if p == nil || len(*p) < 2 {
		return p.Copy()
	}
	m := NewStringMapBool()
	slice := NewMapSliceV()
	for i := range *p {
		if ok := m.Set((*p)[i], true); ok {
			slice.Append((*p)[i])
		}
	}
	return slice
}

// UniqM modifies this Slice to remove all non uniq elements while preserving element order.
// Cost for this call vs the Uniq is roughly the same, this one is dropping that one appending.
func (p *MapSlice) UniqM() ISlice {
	if p == nil || len(*p) < 2 {
		return p
	}
	m := NewStringMapBool()
	l := len(*p)
	for i := 0; i < l; i++ {
		if ok := m.Set((*p)[i], true); !ok {
			p.DropAt(i)
			l--
			i--
		}
	}
	return p
}
