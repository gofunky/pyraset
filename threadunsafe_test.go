package mapset

import (
	"testing"
)

func Test_threadUnsafeSet_String(t *testing.T) {
	type fields []interface{}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "empty set",
			fields: fields{},
			want:   "Set{}",
		},
		{
			name:   "unary set",
			fields: fields{"one"},
			want:   "Set{one}",
		},
		{
			name:   "binary set",
			fields: fields{"one", "two"},
			want:   "Set{one, two}",
		},
		{
			name:   "nested set",
			fields: fields{"one", NewSet("two")},
			want:   "Set{one, Set{two}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.fields...)
			if got := set.String(); got != tt.want {
				t.Errorf("threadUnsafeSet.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
