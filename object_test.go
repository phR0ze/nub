package n

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestObject_ToBool(t *testing.T) {

	// w/out error
	{
		o := &Object{true}
		assert.IsType(t, true, o.ToBool())
	}

	// w/error
	{
		o := &Object{true}
		b, e := o.ToBoolE()
		assert.Nil(t, e)
		assert.IsType(t, true, b)
	}
}

func TestObject_ToTime(t *testing.T) {

	// w/out error
	{
		o := &Object{time.Time{}}
		assert.IsType(t, time.Time{}, o.ToTime())
	}

	// w/error
	{
		o := &Object{time.Time{}}
		obj, e := o.ToTimeE()
		assert.Nil(t, e)
		assert.IsType(t, time.Time{}, obj)
	}
}

func TestObject_ToDuration(t *testing.T) {

	// w/out error
	{
		o := &Object{time.Duration(0)}
		assert.IsType(t, time.Duration(0), o.ToDuration())
	}

	// w/error
	{
		o := &Object{time.Duration(0)}
		obj, e := o.ToDurationE()
		assert.Nil(t, e)
		assert.IsType(t, time.Duration(0), obj)
	}
}

func TestObject_ToFloat32(t *testing.T) {

	// w/out error
	{
		o := &Object{float32(1.0)}
		assert.IsType(t, float32(1.0), o.ToFloat32())
	}

	// w/error
	{
		o := &Object{float32(1.0)}
		obj, e := o.ToFloat32E()
		assert.Nil(t, e)
		assert.IsType(t, float32(1.0), obj)
	}
}

func TestObject_ToFloat64(t *testing.T) {

	// w/out error
	{
		o := &Object{float64(1.0)}
		assert.IsType(t, float64(1.0), o.ToFloat64())
	}

	// w/error
	{
		o := &Object{float64(1.0)}
		obj, e := o.ToFloat64E()
		assert.Nil(t, e)
		assert.IsType(t, float64(1.0), obj)
	}
}

func TestObject_ToInt(t *testing.T) {

	// w/out error
	{
		o := &Object{1}
		assert.IsType(t, 1, o.ToInt())
	}

	// w/error
	{
		o := &Object{1}
		obj, e := o.ToIntE()
		assert.Nil(t, e)
		assert.IsType(t, 1, obj)
	}
}

func TestObject_ToInt8(t *testing.T) {

	// w/out error
	{
		o := &Object{int8(1)}
		assert.IsType(t, int8(1), o.ToInt8())
	}

	// w/error
	{
		o := &Object{int8(1)}
		obj, e := o.ToInt8E()
		assert.Nil(t, e)
		assert.IsType(t, int8(1), obj)
	}
}

func TestObject_ToInt16(t *testing.T) {

	// w/out error
	{
		o := &Object{int16(1)}
		assert.IsType(t, int16(1), o.ToInt16())
	}

	// w/error
	{
		o := &Object{int16(1)}
		obj, e := o.ToInt16E()
		assert.Nil(t, e)
		assert.IsType(t, int16(1), obj)
	}
}

func TestObject_ToInt32(t *testing.T) {

	// w/out error
	{
		o := &Object{int32(1)}
		assert.IsType(t, int32(1), o.ToInt32())
	}

	// w/error
	{
		o := &Object{int32(1)}
		obj, e := o.ToInt32E()
		assert.Nil(t, e)
		assert.IsType(t, int32(1), obj)
	}
}

func TestObject_ToInt64(t *testing.T) {

	// w/out error
	{
		o := &Object{int64(1)}
		assert.IsType(t, int64(1), o.ToInt64())
	}

	// w/error
	{
		o := &Object{int64(1)}
		obj, e := o.ToInt64E()
		assert.Nil(t, e)
		assert.IsType(t, int64(1), obj)
	}
}

func TestObject_ToUInt(t *testing.T) {

	// w/out error
	{
		o := &Object{uint(1)}
		assert.IsType(t, uint(1), o.ToUint())
	}

	// w/error
	{
		o := &Object{uint(1)}
		obj, e := o.ToUintE()
		assert.Nil(t, e)
		assert.IsType(t, uint(1), obj)
	}
}

func TestObject_ToUint8(t *testing.T) {

	// w/out error
	{
		o := &Object{uint8(1)}
		assert.IsType(t, uint8(1), o.ToUint8())
	}

	// w/error
	{
		o := &Object{uint8(1)}
		obj, e := o.ToUint8E()
		assert.Nil(t, e)
		assert.IsType(t, uint8(1), obj)
	}
}

func TestObject_ToUint16(t *testing.T) {

	// w/out error
	{
		o := &Object{uint16(1)}
		assert.IsType(t, uint16(1), o.ToUint16())
	}

	// w/error
	{
		o := &Object{uint16(1)}
		obj, e := o.ToUint16E()
		assert.Nil(t, e)
		assert.IsType(t, uint16(1), obj)
	}
}

func TestObject_ToUint32(t *testing.T) {

	// w/out error
	{
		o := &Object{uint32(1)}
		assert.IsType(t, uint32(1), o.ToUint32())
	}

	// w/error
	{
		o := &Object{uint32(1)}
		obj, e := o.ToUint32E()
		assert.Nil(t, e)
		assert.IsType(t, uint32(1), obj)
	}
}

func TestObject_ToUint64(t *testing.T) {

	// w/out error
	{
		o := &Object{uint64(1)}
		assert.IsType(t, uint64(1), o.ToUint64())
	}

	// w/error
	{
		o := &Object{uint64(1)}
		obj, e := o.ToUint64E()
		assert.Nil(t, e)
		assert.IsType(t, uint64(1), obj)
	}
}

func TestObject_ToStr(t *testing.T) {

	{
		o := &Object{""}
		assert.IsType(t, (*Str)(nil), o.ToStr())
	}

	{
		o := &Object{"test"}
		obj := o.ToStr()
		assert.IsType(t, (*Str)(nil), obj)
		assert.Equal(t, A("test"), obj)
	}
}

func TestObject_ToString(t *testing.T) {

	// w/out error
	{
		o := &Object{""}
		assert.IsType(t, "", o.ToString())
	}

	// w/error
	{
		o := &Object{""}
		obj, e := o.ToStringE()
		assert.Nil(t, e)
		assert.IsType(t, "", obj)
	}
}

func TestObject_ToStringMap(t *testing.T) {

	// w/out error
	{
		o := &Object{map[string]interface{}{}}
		assert.IsType(t, (*StringMap)(nil), o.ToStringMap())
	}

	// w/error
	{
		o := &Object{map[string]interface{}{}}
		obj, e := o.ToStringMapE()
		assert.Nil(t, e)
		assert.IsType(t, (*StringMap)(nil), obj)
	}
}

func TestObject_ToStringMapG(t *testing.T) {

	// w/out error
	{
		o := &Object{map[string]interface{}{}}
		assert.IsType(t, map[string]interface{}{}, o.ToStringMapG())
	}

	// w/error
	{
		o := &Object{map[string]interface{}{}}
		obj, e := o.ToStringMapGE()
		assert.Nil(t, e)
		assert.IsType(t, map[string]interface{}{}, obj)
	}
}

func TestObject_ToStringMapString(t *testing.T) {

	// w/out error
	{
		o := &Object{map[string]string{}}
		assert.IsType(t, map[string]string{}, o.ToStringMapString())
	}

	// w/error
	{
		o := &Object{map[string]string{}}
		obj, e := o.ToStringMapStringE()
		assert.Nil(t, e)
		assert.IsType(t, map[string]string{}, obj)
	}
}

func TestObject_ToStringMapStringSlice(t *testing.T) {

	// w/out error
	{
		o := &Object{map[string][]string{}}
		assert.IsType(t, map[string][]string{}, o.ToStringMapStringSlice())
	}

	// w/error
	{
		o := &Object{map[string][]string{}}
		obj, e := o.ToStringMapStringSliceE()
		assert.Nil(t, e)
		assert.IsType(t, map[string][]string{}, obj)
	}
}

func TestObject_ToStringMapBool(t *testing.T) {

	// w/out error
	{
		o := &Object{map[string]bool{}}
		assert.IsType(t, map[string]bool{}, o.ToStringMapBool())
	}

	// w/error
	{
		o := &Object{map[string]bool{}}
		obj, e := o.ToStringMapBoolE()
		assert.Nil(t, e)
		assert.IsType(t, map[string]bool{}, obj)
	}
}

func TestObject_ToStringMapInt(t *testing.T) {

	// w/out error
	{
		o := &Object{map[string]int{}}
		assert.IsType(t, map[string]int{}, o.ToStringMapInt())
	}

	// w/error
	{
		o := &Object{map[string]int{}}
		obj, e := o.ToStringMapIntE()
		assert.Nil(t, e)
		assert.IsType(t, map[string]int{}, obj)
	}
}

func TestObject_ToStringMapInt64(t *testing.T) {

	// w/out error
	{
		o := &Object{map[string]int64{}}
		assert.IsType(t, map[string]int64{}, o.ToStringMapInt64())
	}

	// w/error
	{
		o := &Object{map[string]int64{}}
		obj, e := o.ToStringMapInt64E()
		assert.Nil(t, e)
		assert.IsType(t, map[string]int64{}, obj)
	}
}

func TestObject_ToSlice(t *testing.T) {

	// w/out error
	{
		o := &Object{[]interface{}{}}
		assert.IsType(t, []interface{}{}, o.ToSlice())
	}

	// w/error
	{
		o := &Object{[]interface{}{}}
		obj, e := o.ToSliceE()
		assert.Nil(t, e)
		assert.IsType(t, []interface{}{}, obj)
	}
}

func TestObject_ToBoolSlice(t *testing.T) {

	// w/out error
	{
		o := &Object{[]bool{}}
		assert.IsType(t, []bool{}, o.ToBoolSlice())
	}

	// w/error
	{
		o := &Object{[]bool{}}
		obj, e := o.ToBoolSliceE()
		assert.Nil(t, e)
		assert.IsType(t, []bool{}, obj)
	}
}

func TestObject_ToStringSlice(t *testing.T) {

	// w/out error
	{
		o := &Object{[]string{}}
		assert.IsType(t, &StringSlice{}, o.ToStringSlice())
	}

	// w/error
	{
		o := &Object{[]string{}}
		obj, e := o.ToStringSliceE()
		assert.Nil(t, e)
		assert.IsType(t, &StringSlice{}, obj)
	}
}

func TestObject_ToStringSliceG(t *testing.T) {

	// w/out error
	{
		o := &Object{[]string{}}
		assert.IsType(t, []string{}, o.ToStringSliceG())
	}

	// w/error
	{
		o := &Object{[]string{}}
		obj, e := o.ToStringSliceGE()
		assert.Nil(t, e)
		assert.IsType(t, []string{}, obj)
	}
}

func TestObject_ToIntSlice(t *testing.T) {

	// w/out error
	{
		o := &Object{[]int{}}
		assert.IsType(t, []int{}, o.ToIntSliceG())
	}

	// w/error
	{
		o := &Object{[]int{}}
		obj, e := o.ToIntSliceGE()
		assert.Nil(t, e)
		assert.IsType(t, []int{}, obj)
	}
}

func TestObject_ToDurationSlice(t *testing.T) {

	// w/out error
	{
		o := &Object{[]time.Duration{}}
		assert.IsType(t, []time.Duration{}, o.ToDurationSlice())
	}

	// w/error
	{
		o := &Object{[]time.Duration{}}
		obj, e := o.ToDurationSliceE()
		assert.Nil(t, e)
		assert.IsType(t, []time.Duration{}, obj)
	}
}
