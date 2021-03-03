package local_cache

import (
	"fmt"
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
