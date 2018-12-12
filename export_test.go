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
		assert.Equal(t, "1", Q(1).A())
		assert.Equal(t, "true", Q(true).A())
		assert.Equal(t, "false", Q(false).A())
	}
}

func TestB(t *testing.T) {
	{
		assert.Equal(t, true, Q(true).B())
		assert.Equal(t, false, Q(false).B())
	}
	{
		q := Q([]bool{true})
		assert.Equal(t, true, q.At(0).B())
	}
	{
		assert.Equal(t, true, Q(4).B())
		assert.Equal(t, false, Q(0).B())
	}
	{
		assert.Equal(t, true, Q("1").B())
		assert.Equal(t, true, Q("T").B())
		assert.Equal(t, true, Q("true").B())
		assert.Equal(t, true, Q("True").B())
		assert.Equal(t, true, Q("TRUE").B())
		assert.Equal(t, false, Q("0").B())
		assert.Equal(t, false, Q("F").B())
		assert.Equal(t, false, Q("false").B())
		assert.Equal(t, false, Q("False").B())
		assert.Equal(t, false, Q("FALSE").B())
	}
	{
		q, _ := FromYaml(`foo:
  bar: true`)
		assert.True(t, q.Any())

		// Exists
		assert.Equal(t, true, q.Yaml("foo.bar").O())
		assert.Equal(t, true, q.Yaml("foo.bar").B())

		// Doesn't exist
		assert.Equal(t, nil, q.Yaml("foo.foo").O())
		assert.Equal(t, false, q.Yaml("foo.foo").B())
	}
}

func TestI(t *testing.T) {
	assert.Equal(t, 2, Q(2).I())
	assert.Equal(t, 2, Q("2").I())
	assert.Equal(t, 3, Q([]int{2, 3}).At(1).O())
	assert.Equal(t, 3, Q([]int{2, 3}).At(1).I())
	assert.Equal(t, 3, Q([]string{"2", "3"}).At(1).I())
	assert.Equal(t, 1, Q(true).I())
	assert.Equal(t, 0, Q(false).I())
}

func TestM(t *testing.T) {
	{
		q, _ := FromYaml(`foo:
  bar: true`)
		assert.True(t, q.Any())

		// Key exists
		assert.Equal(t, map[string]interface{}{"bar": true}, q.Yaml("foo").M())

		// Key doesn't exist
		assert.Equal(t, map[string]interface{}{}, q.Yaml("bar").M())
	}
	{
		// Convert with simple cast
		data := map[string]interface{}{"1": "one", "2": "two", "3": "three"}
		q := Q(data)
		assert.Equal(t, 3, q.Len())
		assert.Equal(t, data, q.M())
	}
	{
		// Convert list of KeyVal into a map
		q := Q([]string{"foo=bar"})
		m := q.Map(func(x O) O {
			return A(x.(string)).Split("=").YamlKeyVal()
		}).M()
		assert.Equal(t, map[string]interface{}{"foo": "bar"}, m)
	}
}

func TestInts(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3}, Q([]int{1, 2, 3}).Ints())
}

func TestAAMap(t *testing.T) {
	{
		q, _ := FromYaml(`foo:
  bar: frodo`)
		assert.True(t, q.Any())

		// Key exists
		assert.Equal(t, map[string]string{"bar": "frodo"}, q.Yaml("foo").AAMap())

		// Key doesn't exist
		assert.Equal(t, map[string]string{}, q.Yaml("bar").AAMap())
	}
	{
		// Get map of string to string
		q, _ := FromYaml(`foobar:
  labels:
    name: one
    meta: two
    data: three`)
		assert.True(t, q.Any())
		expected := map[string]string{"name": "one", "meta": "two", "data": "three"}
		assert.Equal(t, expected, q.Yaml("foobar.labels").AAMap())
	}
	{
		// Ints as keys
		q, _ := FromYaml(`foobar:
  labels:
    1: one
    2: two
    3: 3`)
		assert.True(t, q.Any())
		expected := map[string]string{"1": "one", "2": "two", "3": "3"}
		assert.Equal(t, expected, q.Yaml("foobar.labels").AAMap())
	}
}

func TestASAMap(t *testing.T) {
	{
		q := Q(map[string][]string{"bar": []string{"frodo"}})
		assert.Equal(t, map[string][]string{"bar": []string{"frodo"}}, q.ASAMap())
	}
	{
		q, err := FromYaml("foo:\n  bar:\n  - frodo")
		assert.Nil(t, err)
		assert.True(t, q.Any())
		assert.Equal(t, map[string][]string{"bar": []string{"frodo"}}, q.Yaml("foo").ASAMap())
	}
}

func TestS(t *testing.T) {
	{
		q, _ := FromYaml(`items:
  - name: one
  - name: two
  - name: three`)
		assert.True(t, q.Any())
		expected := []interface{}{
			map[string]interface{}{"name": "one"},
			map[string]interface{}{"name": "two"},
			map[string]interface{}{"name": "three"},
		}
		assert.Equal(t, expected, q.Yaml("items").S())
	}
	{
		// Not found
		q, _ := FromYaml(`items:
  - name: one
  - name: two
  - name: three`)
		assert.True(t, q.Any())
		assert.Equal(t, N(), q.Yaml(""))
	}
}

func TestSAMap(t *testing.T) {
	{
		q, _ := FromYaml(`items:
  - name: one
  - name: two
  - name: three`)
		assert.True(t, q.Any())
		expected := []map[string]interface{}{
			{"name": "one"},
			{"name": "two"},
			{"name": "three"},
		}
		assert.Equal(t, expected, q.Yaml("items").SAMap())
	}
	{
		q, _ := FromYaml(`items:
  - name: one
  - name: two
  - name: three`)
		assert.True(t, q.Any())
		expected := []map[string]interface{}{}
		assert.Equal(t, expected, q.Yaml("foo").SAMap())
	}
}

func TestSAAMap(t *testing.T) {
	{
		// slice of string to string map
		q, _ := FromYaml(`items:
  - name: one
  - name: two
  - name: three`)
		assert.True(t, q.Any())
		expected := []map[string]string{
			{"name": "one"},
			{"name": "two"},
			{"name": "three"},
		}
		assert.Equal(t, expected, q.Yaml("items").SAAMap())
	}
	{
		// slice of string to int map
		q, _ := FromYaml(`items:
  - name: 1
  - name: 2
  - name: three`)
		assert.True(t, q.Any())
		expected := []map[string]string{
			{"name": "1"},
			{"name": "2"},
			{"name": "three"},
		}
		assert.Equal(t, expected, q.Yaml("items").SAAMap())
	}
}

func TestStrs(t *testing.T) {
	{
		// slice of string to int map
		q, _ := FromYaml(`items:
  - name: 1
  - name: 2
  - name: three`)
		assert.True(t, q.Any())
		expected := []string{}
		assert.Equal(t, expected, q.Yaml("frodo.baggins").Strs())
	}
	{
		q := Q([]string{"one"})
		assert.Equal(t, "one", q.At(0).A())
		assert.Equal(t, []string{"one"}, q.Strs())
	}
	{
		assert.Equal(t, []string{"1", "2", "3"}, Q([]interface{}{"1", "2", "3"}).Strs())
	}
}

func TestCastToTypeOf(t *testing.T) {
	//typof := []int{1, 2, 3}
	//obj := []interface{}{4}[0]
	//CastToTypeOf(typof, obj)
}
