package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	val := "hello"
	ptr := String(val)
	assert.NotNil(t, ptr)
	assert.Equal(t, val, *ptr)
}

func TestBool(t *testing.T) {
	val := true
	ptr := Bool(val)
	assert.NotNil(t, ptr)
	assert.Equal(t, val, *ptr)
}

func TestInt64(t *testing.T) {
	val := int64(42)
	ptr := Int64(val)
	assert.NotNil(t, ptr)
	assert.Equal(t, val, *ptr)
}

func TestUInt32(t *testing.T) {
	val := uint32(32)
	ptr := UInt32(val)
	assert.NotNil(t, ptr)
	assert.Equal(t, val, *ptr)
}

func TestUInt64(t *testing.T) {
	val := uint64(64)
	ptr := UInt64(val)
	assert.NotNil(t, ptr)
	assert.Equal(t, val, *ptr)
}

func TestFloat64(t *testing.T) {
	val := float64(3.14)
	ptr := Float64(val)
	assert.NotNil(t, ptr)
	assert.Equal(t, val, *ptr)
}

func TestTime(t *testing.T) {
	val := time.Now()
	ptr := Time(val)
	assert.NotNil(t, ptr)
	assert.Equal(t, val, *ptr)
}

func TestSliceWithInts(t *testing.T) {
	vals := []int{1, 2, 3}
	ptrs := Slice(vals)
	assert.NotNil(t, ptrs)
	assert.Equal(t, len(vals), len(ptrs))

	for i, ptr := range ptrs {
		assert.NotNil(t, ptr)
		assert.Equal(t, vals[i], *ptr)
	}
}

func TestSliceWithStrings(t *testing.T) {
	vals := []string{"foo", "bar", "baz"}
	ptrs := Slice(vals)
	assert.NotNil(t, ptrs)
	assert.Equal(t, len(vals), len(ptrs))

	for i, ptr := range ptrs {
		assert.NotNil(t, ptr)
		assert.Equal(t, vals[i], *ptr)
	}
}
