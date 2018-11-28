package nub

import (
	"bytes"
	"errors"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/ghodss/yaml"
)

// Queryable provides chainable deferred execution
// and is the heart of the algorithm abstraction layer
type Queryable struct {
	O    interface{}
	Iter func() Iterator

	ref *reflect.Value
}

// Iterator provides a closure to capture the index and reset it
type Iterator func() (item interface{}, ok bool)

// KeyVal similar to C# for iterator over maps
type KeyVal struct {
	Key interface{}
	Val interface{}
}

// Load YAML/JSON from file into queryable
func Load(target string) *Queryable {
	if yamlFile, err := ioutil.ReadFile(target); err == nil {
		data := map[string]interface{}{}
		yaml.Unmarshal(yamlFile, &data)
		return Q(data)
	}
	return M()
}

// A provides a new empty Queryable string
func A() *Queryable {
	obj := string("")
	ref := reflect.ValueOf(obj)
	return &Queryable{O: obj, ref: &ref, Iter: strIter(ref, obj)}
}

func strIter(ref reflect.Value, obj interface{}) func() Iterator {
	return func() Iterator {
		i := 0
		len := ref.Len()
		return func() (item interface{}, ok bool) {
			if ok = i < len; ok {
				item = ref.Index(i).Interface()
				i++
			}
			return
		}
	}
}

// M provides a new empty Queryable map
func M() *Queryable {
	obj := map[interface{}]interface{}{}
	ref := reflect.ValueOf(obj)
	return &Queryable{O: obj, ref: &ref, Iter: mapIter(ref, obj)}
}

func mapIter(ref reflect.Value, obj interface{}) func() Iterator {
	return func() Iterator {
		i := 0
		len := ref.Len()
		keys := ref.MapKeys()
		return func() (item interface{}, ok bool) {
			if ok = i < len; ok {
				item = &KeyVal{
					Key: keys[i].Interface(),
					Val: ref.MapIndex(keys[i]).Interface(),
				}
				i++
			}
			return
		}
	}
}

// S provides a new empty Queryable slice
func S() *Queryable {
	obj := []interface{}{}
	ref := reflect.ValueOf(obj)
	return &Queryable{O: obj, ref: &ref, Iter: sliceIter(ref, obj)}
}

func sliceIter(ref reflect.Value, obj interface{}) func() Iterator {
	return func() Iterator {
		i := 0
		len := ref.Len()
		return func() (item interface{}, ok bool) {
			if ok = i < len; ok {
				item = ref.Index(i).Interface()
				i++
			}
			return
		}
	}
}

// Q provides origination for the Queryable abstraction layer
func Q(obj interface{}) *Queryable {
	if obj == nil {
		return S()
	}

	ref := reflect.ValueOf(obj)
	result := &Queryable{O: obj, ref: &ref}
	switch ref.Kind() {

	// Slice types
	case reflect.Array, reflect.Slice:
		result.Iter = sliceIter(ref, obj)

	// Handle map types
	case reflect.Map:
		result.Iter = mapIter(ref, obj)

	// Handle string types
	case reflect.String:
		result.Iter = strIter(ref, obj)

	// Chan types
	case reflect.Chan:
		panic("TODO: handle reflect.Chan")

	// Handle int types
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// No iterator should exist

	// Handle unknown types
	default:
		panic("TODO: handle custom types")
	}

	return result
}

// Append items to the end of the collection and return the queryable
// converting to a collection if necessary
func (q *Queryable) Append(obj ...interface{}) *Queryable {

	// No existing type return a new queryable
	if q.ref == nil {
		*q = *Q(obj)
		return q
	}

	// Not a collection type create a new queryable
	kind := q.ref.Kind()
	if kind != reflect.Array && kind != reflect.Slice && kind != reflect.Map {
		*q = *S().Append(q.O)
	}

	// Append to the collection type
	ref := reflect.ValueOf(obj)
	for i := 0; i < ref.Len(); i++ {
		*q.ref = reflect.Append(*q.ref, ref.Index(i))
	}
	return q
}

// At returns the item at the given index location. Allows for negative notation
func (q *Queryable) At(i int) *Queryable {
	if q.Iter != nil {
		if i < 0 {
			i = q.ref.Len() + i
		}
		if i >= 0 && i < q.ref.Len() {
			if str, ok := q.O.(string); ok {
				return Q(string(str[i]))
			}
			return Q(q.ref.Index(i).Interface())
		}
	}
	panic(errors.New("Index out of slice bounds"))
}

// Clear the queryable collection
func (q *Queryable) Clear() *Queryable {
	*q = *S()
	return q
}

// Each iterates over the queryable and executes the given action
func (q *Queryable) Each(action func(interface{})) {
	if q.Iter != nil {
		next := q.Iter()
		for x, ok := next(); ok; x, ok = next() {
			action(x)
		}
	}
}

// Join slice items as string with given delimeter
func (q *Queryable) Join(delim string) *Queryable {
	var joined bytes.Buffer
	if !q.Singular() {
		next := q.Iter()
		for x, ok := next(); ok; x, ok = next() {
			if str, ok := x.(string); ok {
				joined.WriteString(str)
				joined.WriteString(delim)
			}
		}
	} else if q.Iter != nil {
		joined.WriteString(q.O.(string))
	}
	return Q(strings.TrimSuffix(joined.String(), delim))
}

// Len of the collection type including string
func (q *Queryable) Len() int {
	if q.Iter != nil {
		return q.ref.Len()
	}
	return 1
}

// Set provides a way to set underlying object Queryable is operating on
func (q *Queryable) Set(obj interface{}) *Queryable {
	other := Q(obj)
	q.O = other.O
	q.Iter = other.Iter
	q.ref = other.ref
	return q
}

// Singular is queryable encapsulating a non-collection
func (q *Queryable) Singular() bool {
	_, strType := q.O.(string)
	if q.Iter == nil || strType {
		return true
	}
	return false
}
