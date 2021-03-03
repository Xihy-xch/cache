package local_cache

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Test_cache_Get(t *testing.T) {
	type fields struct {
		valueMap map[string]item
	}
	type args struct {
		key string
		val interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    interface{}
	}{
		// TODO: Add test cases.
		{
			name: "get",
			args: args{
				key: "key",
				val: struct {
					Name string
				}{Name: "xiaoWang"},
			},
			wantErr: false,
			want: struct {
				Name string
			}{Name: "xiaoWang"},
		},
	}
	c := NewCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.Set(tt.args.key, tt.args.val)
			for i := 1; i < 2000; i++ {
				time.Sleep(1 * time.Second)
				i := i
				go func() {
					got, err := c.Get(tt.args.key)
					fmt.Println(got)
					if (err != nil) != tt.wantErr {
						t.Errorf("Get() error = %v, wantErr %v %d 次", err, tt.wantErr, i)
						return
					}
					if !reflect.DeepEqual(got, tt.want) {
						t.Errorf("Get() got = %v, want %v", got, tt.want)
					}
				}()
			}

		})
	}
}