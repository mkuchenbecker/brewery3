package practice

type Node struct {
	next *Node
	val  int
}

func partition(root *Node, partition int) *Node {
	var lower, higher Node
	partitionInternal(root, &lower, &higher, partition)
	curr := &lower
	for {
		if curr.next == nil {
			break
		}
		curr = curr.next
	}
	curr.next = &higher
	return &lower
}

func partitionInternal(root *Node, lower *Node, higher *Node, partition int) {
	if root == nil {
		return
	}
	if root.val >= partition {
		if higher == nil {
			*higher = *root
			higher.next = nil
			partitionInternal(root.next, lower, higher, partition)
			return
		}
		higher.next = root
		higher.next.next = nil
		partitionInternal(root.next, lower, higher.next, partition)
		return
	}
	if lower == nil {
		*lower = *root
		lower.next = nil
		partitionInternal(root.next, lower, higher, partition)
		return
	}
	lower.next = root
	lower.next.next = nil
	partitionInternal(root.next, lower.next, higher, partition)
	return
}
