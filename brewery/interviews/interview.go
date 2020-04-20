package interview

import (
	"sort"
	"fmt"
)

	// type MinHeap []IPTuple

	// func (m *MinHeap) Peek() IPTuple {
	// 	return []IPTuple(*m)[0]
	// }
	// func (m *MinHeap) Pop() IPTuple {
	// 	out :=  []IPTuple(*m)[0]
	// 	*m = []IPTuple(*m)[1:]
	// 	return out
	// }
	// func (m *MinHeap) Add(t IPTuple) {
	// 	*m = append([]IPTuple(*m),t)
	// 	sort.Slice(*m, func(i,j int) bool {
	// 		return []IPTuple(*m)[i].end.Before([]IPTuple(*m)[j].end)
	// 	})
	// }

type IPTuple struct {
	start int
	end int
}

func GetChargeTime(source []IPTuple, start int, end int) int64 {
	sort.Slice(source, func(i,j int) bool {
		return source[i].start < source[j].start
	})

	total := int64(0)
	for t:=start; t <= end; t++ {
		var active int64 = 0
		for _, ip := range source {
			if t > ip.start && t <= ip.end {
				fmt.Printf("loop\n%d\n%ds\n",t,ip.end)
					active ++
			}
		}
		fmt.Printf("active %d\n",active)
		active = active  // allow one active at a time
		if active >= 1 {
			active = active - 1
		}
		total += active
	}
	return total
}


