package herlihy

import (
	"math"
	"sync"
)

type nodeOptimistic struct {
	next *nodeOptimistic
	sync.Mutex // fine grain sync
	val int
}


type optimisticLocked struct {
	head *nodeOptimistic // sentinel int min
	tail *nodeOptimistic // sentinel int max
}

func newOptimistic() optimisticLocked {
	head := nodeOptimistic{
		val: math.MinInt64,
	}
	tail := nodeOptimistic{
		val: math.MaxInt64,
	}
	head.next = &tail
	return optimisticLocked{
		&head,
		&tail,
	}
}

func(self *optimisticLocked) add(x int) bool {
	if x == math.MaxInt64 || x == math.MinInt64 {
		return false
	}
	startover:
	prev := self.head
	curr := prev.next
	for curr != nil {
		if curr.val == x {
			return false
		}
		if curr.val > x {
			prev.Lock()
			curr.Lock()
			if curr.val > x && self.valid(prev, curr) {
				nodeOptimistic := nodeOptimistic{
					next: curr,
					val: x,
				}
				prev.next = &nodeOptimistic
				prev.Unlock()
				curr.Unlock()
				return true
			}
			prev.Unlock()
			curr.Unlock()
			goto startover
		}
		this := curr
		curr = curr.next
		prev = this
	}
	return false
}

func(self *optimisticLocked) remove(x int) bool {
	if x == math.MaxInt64 || x == math.MinInt64 {
		return false
	}
	startover:
	prev := self.head
	curr := prev.next
	for curr != nil {
		if curr.val > x {
			return false
		}
		if curr.val == x {
			prev.Lock()
			curr.Lock()
			if curr.val == x && self.valid(prev, curr) {
				prev.next = curr.next
				prev.Unlock()
				curr.Unlock()
				return true
			}
			prev.Unlock()
			curr.Unlock()
			goto startover
		}
		this := curr
		curr = curr.next
		prev = this
	}
	return false
}

func (self *optimisticLocked)  valid(prev, curr *nodeOptimistic) bool {
	this := self.head
	for this != nil {
		if this == prev {
			return this.next == curr
		}
		this = this.next
	}
	return false
}

func(self *optimisticLocked) contains(x int) bool { // todo: why contains works without validating?
	if x == math.MaxInt64 || x == math.MinInt64 {
		return false
	}
startover:
	prev := self.head
	curr := prev.next
	for curr != nil {
		if curr.val > x {
			return false
		}
		if curr.val == x {
			prev.Lock()
			curr.Lock()
			if curr.val == x && self.valid(prev, curr) {
				prev.Unlock()
				curr.Unlock()
				return true
			}
			prev.Unlock()
			curr.Unlock()
			goto startover
		}
		this := curr
		curr = curr.next
		prev = this
	}
	return false
}
