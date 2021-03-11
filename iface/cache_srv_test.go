package iface

import (
	"context"
	"fmt"
	"local-cache/local_cache"
	cache "local-cache/proto"
	"reflect"
	"testing"
)

func TestCacheSrv_Get(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *cache.CacheGetRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *cache.CacheGetResponse
		wantErr bool
	}{
		{
			name: "",
			args: args{
				ctx:     context.Background(),
				request: &cache.CacheGetRequest{Key: "test_key"},
			},
			want:    nil,
			wantErr: false,
		},
	}

	person := Person{Name: "xiaoMing"}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CacheSrv{
				cache: local_cache.NewLRUCache(1024),
			}
			val, _ := marshal(person)
			c.Set(tt.args.ctx, &cache.CacheSetRequest{
				Key:   "test_key",
				Value: val,
			})

			got, err := c.Get(tt.args.ctx, tt.args.request)
			person1 := Person{}
			unmarshal(got.GetValue(), &person1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(person1, person) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}

			fmt.Println(person, person1)
		})
	}
}

type Person struct {
	Name string
}

func Test_marshal(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "解码测试",
			args: args{
				v: Person{Name: "xiaoMing"},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := marshal(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var person interface{}
			err = unmarshal(got, &person)
			if (err != nil) != tt.wantErr {
				t.Errorf("marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(person, tt.args.v) {
				t.Errorf("marshal() got = %v, want %v", person, tt.args.v)
			}
			fmt.Printf("marshal() got = %v, want %v", person, tt.args.v)
		})
	}
}
