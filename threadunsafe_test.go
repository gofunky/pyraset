package mapset

import (
	"testing"
)

func Test_threadUnsafeSet_String(t *testing.T) {
	type fields []interface{}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name:   "empty set",
			fields: fields{},
			want:   []string{"Set{}"},
		},
		{
			name:   "unary set",
			fields: fields{"one"},
			want:   []string{"Set{one}"},
		},
		{
			name:   "binary set",
			fields: fields{"one", "two"},
			want:   []string{"Set{one, two}", "Set{two, one}"},
		},
		{
			name:   "nested set",
			fields: fields{"one", NewSet("two")},
			want:   []string{"Set{one, Set{two}}", "Set{Set{two}, one}"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSet(tt.fields...)
			for i, want := range tt.want {
				got := set.String()
				if got == want {
					break
				}
				if i == len(tt.want)-1 {
					t.Errorf("threadUnsafeSet.String() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
