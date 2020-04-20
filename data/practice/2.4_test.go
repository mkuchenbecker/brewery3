package practice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test2_4(t *testing.T) {
	nodes := []*Node{
		&Node{val: 3},
		&Node{val: 5},
		&Node{val: 8},
		&Node{val: 5},
		&Node{val: 10},
		&Node{val: 2},
		&Node{val: 3},
	}
	var last *Node
	for _, n := range nodes {
		if last != nil {
			last.next = n
		}
		last = n
	}

	partition(nodes[0], 5)

	last = nodes[0]
	vals := make([]int, 0, len(nodes))
	for {
		if last == nil {
			break
		}
		vals = append(vals, last.val)
		last = last.next
	}
	assert.Equal(t, []int{3, 1, 2, 10, 5, 5, 8}, vals)
}
