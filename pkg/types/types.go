package types

import "time"

func String(v string) *string {
	return &v
}

func Bool(v bool) *bool {
	return &v
}

func Int64(v int64) *int64 {
	return &v
}

func UInt32(v uint32) *uint32 {
	return &v
}

func UInt64(v uint64) *uint64 {
	return &v
}

func Float64(v float64) *float64 {
	return &v
}

func Time(v time.Time) *time.Time {
	return &v
}

func Slice[T any](s []T) []*T {
	res := make([]*T, len(s))
	for i := range s {
		res[i] = &s[i]
	}
	return res
}
