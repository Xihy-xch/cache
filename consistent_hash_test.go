package local_cache

import (
	"fmt"
	"strconv"
	"testing"
)

func TestMap_Add(t *testing.T) {
	type args struct {
		keys []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "测试add",
			args: args{
				keys: []string{"test1", "test2", "test3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(3, nil)

			m.Add(tt.args.keys...)
			m.Add("test_key")
			for i := range m.keys {
				fmt.Printf("%d,", m.keys[i])
			}
		})
	}
}

func TestMap_Get(t *testing.T) {


	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "测试get",
			args: args{
				key: "7",
			},
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(3, func(data []byte) uint32 {
				i, _ := strconv.Atoi(string(data))
				return uint32(i)
			})


			m.Add("1", "2", "3")
			if got := m.Get(tt.args.key); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
