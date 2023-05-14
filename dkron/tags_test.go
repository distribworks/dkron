package dkron

import (
	"reflect"
	"testing"
)

func Test_cleanTags(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	type args struct {
		tags map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		want1   int
		wantErr bool
	}{
		{
			name:  "Clean Tags",
			args:  args{map[string]string{"key1": "value1", "key2": "value2"}},
			want:  map[string]string{"key1": "value1", "key2": "value2"},
			want1: maxInt,
		},
		{
			name:  "With Cardinality",
			args:  args{map[string]string{"key1": "value1", "key2": "value2:5"}},
			want:  map[string]string{"key1": "value1", "key2": "value2"},
			want1: 5,
		},
		{
			name:  "With Multiple Cardinalities",
			args:  args{map[string]string{"key1": "value1:2", "key2": "value2:5"}},
			want:  map[string]string{"key1": "value1", "key2": "value2"},
			want1: 2,
		},
		{
			name:  "With String Cardinality",
			args:  args{map[string]string{"key1": "value1", "key2": "value2:cardinality"}},
			want:  map[string]string{"key1": "value1", "key2": "value2"},
			want1: 0,
		},
	}
	logger := getTestLogger()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := cleanTags(tt.args.tags, logger)
			if !reflect.DeepEqual(got, tt.want) {
				t.Logf("got map: %#v", got)
				t.Logf("want map: %#v", tt.want)
				t.Errorf("cleanTags() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("cleanTags() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
