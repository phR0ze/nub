package n

import (
	"fmt"
	"reflect"
)

// Collecting functions that return external Go types here

// A exports queryable into a string
func (q *Queryable) A() string {
	return q.v.Interface().(string)
}

// I exports queryable into an int
func (q *Queryable) I() int {
	return q.v.Interface().(int)
}

// Ints exports queryable into an int slice
func (q *Queryable) Ints() []int {
	result := []int{}
	next := q.Iter()
	for x, ok := next(); ok; x, ok = next() {
		result = append(result, x.(int))
	}
	return result
}

// M exports queryable into a map
func (q *Queryable) M() (result map[string]interface{}) {
	if v, ok := q.O().(map[string]interface{}); ok {
		result = v
	} else {
		result = map[string]interface{}{}
		if q.TypeMap() {
			next := q.Iter()
			for x, ok := next(); ok; x, ok = next() {
				pair := x.(KeyVal)
				result[fmt.Sprint(pair.Key)] = pair.Val
			}
		}
	}
	return result
}

// O exports queryable into a interface{}
func (q *Queryable) O() interface{} {
	return q.v.Interface()
}

// Strs exports queryable into an string slice
func (q *Queryable) Strs() []string {
	result := []string{}
	if q.TypeSlice() {
		next := q.Iter()
		for x, ok := next(); ok; x, ok = next() {
			result = append(result, fmt.Sprint(x))
		}
	}
	return result
}

// StrMap exports queryable into an string to string map
func (q *Queryable) StrMap() (result map[string]string) {
	if v, ok := q.O().(map[string]string); ok {
		result = v
	} else {
		result = map[string]string{}
		if q.TypeMap() {
			next := q.Iter()
			for x, ok := next(); ok; x, ok = next() {
				pair := x.(KeyVal)
				result[fmt.Sprint(pair.Key)] = fmt.Sprint(pair.Val)
			}
		}
	}
	return result
}

// CastToTypeOf casts the obj to the type of the typof
func CastToTypeOf(typof interface{}, obj interface{}) *reflect.Value {
	panic("TODO: experimenting with reflection")
	typ := reflect.TypeOf(typof)
	switch typ.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		targetType := typ.Elem()
		originType := reflect.TypeOf(obj)
		fmt.Println(targetType)
		fmt.Println(originType)
	default:
	}

	return nil
}