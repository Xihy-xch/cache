package local_cache

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestNodeList_pushFront(t *testing.T) {
	type args struct {
		nodes []*Node
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "pushFront",
			args: args{
				nodes: []*Node{
					NewNode("1", item{
						value:      "val1",
						expiration: time.Now(),
					}),
					NewNode("2", item{
						value:      "val2",
						expiration: time.Now(),
					}),
				},
			},
		},
	}
	n := NewNodeList()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n.pushFront(tt.args.nodes[0])
			n.pushFront(tt.args.nodes[1])
			fmt.Println(n.front().val.value)
			//n.popBack()
			n.moveToFront(n.back())
			fmt.Println(n.front().val.value)
		})
	}

}

func TestLRUCache_Get(t *testing.T) {
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
	d := NewCache(NewLRUCache(10))
	for i := 0; i < 10; i++ {
		d.Set("test_key"+strconv.Itoa(i), "test_val", WithExpiration(10*time.Second))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := make(chan int)
			for i := 0; i < 10; i++ {
				i := i
				go func() {
					for {
						time.Sleep(1 * time.Second)
						got, err := d.Get(tt.args.key + strconv.Itoa(i))
						if (err != nil) != tt.wantErr {
							t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
							ch <- i
							return
						}
						if !reflect.DeepEqual(got, tt.want) {
							t.Errorf("Get() got = %v, want %v", got, tt.want)
							ch <- i
						}
					}

				}()
			}
			<-ch
		})
	}
}

func TestLRUCache_Get1(t *testing.T) {
	type args struct {
		key  string
		opts []OptionsFn
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "正常流程",
			args: args{
				key: "test_key",
			},
			want:    "test_val",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewCache(NewLRUCache(10))
			got, err := l.Get(tt.args.key, WithGetter(GetterFunc(func(key string) (interface{}, error) {
				return "test_val", nil
			})))
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
