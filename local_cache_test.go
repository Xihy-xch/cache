package local_cache

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestDefaultCache_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "测试get",
			args: args{
				key: "test_key",
			},
			want:    "test_val",
			wantErr: false,
		},
	}
	d := NewCache(WithMode(Default))
	for i := 0; i < 10; i++ {
		i := i
		go func() {
			d.Set("test_key"+strconv.Itoa(i), "test_val", WithExpiration(10*time.Second))
		}()
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := make(chan int)
			for i := 0; i < 10; i++ {
				i := i
				go func() {
					time.Sleep(1 * time.Second)
					got, err := d.Get(tt.args.key + strconv.Itoa(i))
					if (err != nil) != tt.wantErr {
						t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
						ch <- 1
						return
					}
					if !reflect.DeepEqual(got.(item).value, tt.want) {
						t.Errorf("Get() got = %v, want %v", got, tt.want)
						ch <- 1
					}
				}()
			}
			<-ch
		})
	}
}
