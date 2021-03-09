package local_cache

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_getDefaultOptions(t *testing.T) {
	tests := []struct {
		name string
		want *Options
	}{
		// TODO: Add test cases.
		{
			name: "正常流程",
			want: getDefaultOptions(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getDefaultOptions()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDefaultOptions() = %v, want %v", got, tt.want)
			}

			fmt.Println(got.getter, "aa")
		})
	}
}
