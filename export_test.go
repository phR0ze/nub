package n

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestA(t *testing.T) {
	{
		q := Q("one")
		assert.Equal(t, "o", q.At(0).A())
	}
	{
		q := Q([]string{"one"})
		assert.Equal(t, "one", q.At(0).A())
		assert.Equal(t, []string{"one"}, q.Strs())
	}
	{
		assert.Equal(t, "1", Q("1").A())
	}
}

func TestI(t *testing.T) {
	assert.Equal(t, 1, Q(1).I())
}

func TestM(t *testing.T) {
	{
		data := map[string]interface{}{"1": "one", "2": "two", "3": "three"}
		q := Q(data)
		assert.Equal(t, 3, q.Len())
		assert.Equal(t, data, q.M())
	}
}

func TestInts(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3}, Q([]int{1, 2, 3}).Ints())
}

func TestStrs(t *testing.T) {
	{
		q := Q([]string{"one"})
		assert.Equal(t, "one", q.At(0).A())
		assert.Equal(t, []string{"one"}, q.Strs())
	}
	{
		assert.Equal(t, []string{"1", "2", "3"}, Q([]interface{}{"1", "2", "3"}).Strs())
	}
}

func TestStrMap(t *testing.T) {
	{
		// Get map of string to string
		q := FromYAML(`foobar:
  labels:
    name: one
    meta: two
    data: three`)
		assert.True(t, q.Any())
		expected := map[string]string{"name": "one", "meta": "two", "data": "three"}
		assert.Equal(t, expected, q.YAML("foobar.labels").StrMap())
	}
	{
		// Ints as keys
		q := FromYAML(`foobar:
  labels:
    1: one
    2: two
    3: 3`)
		assert.True(t, q.Any())
		expected := map[string]string{"1": "one", "2": "two", "3": "3"}
		assert.Equal(t, expected, q.YAML("foobar.labels").StrMap())
	}
}

func TestCastToTypeOf(t *testing.T) {
	//typof := []int{1, 2, 3}
	//obj := []interface{}{4}[0]
	//CastToTypeOf(typof, obj)
}