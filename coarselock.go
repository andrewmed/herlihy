package herlihy

import (
	"math"
	"sync"
)


type nodeCoarse struct {
	next *nodeCoarse
	val int
}

type coarseLocked struct {
	sync.Mutex // coarse grain sync
	head *nodeCoarse // sentinel int min
	tail *nodeCoarse // sentinel int max
}

func newCoarse() coarseLocked {
	head := nodeCoarse{
		val: math.MinInt64,
	}
	tail := nodeCoarse{
		val: math.MaxInt64,
	}
	head.next = &tail
	return coarseLocked{
		head: &head,
		tail: &tail,
	}
}

func(self *coarseLocked) add(x int) bool {
	if x == math.MaxInt64 || x == math.MinInt64 {
		return false
	}
	self.Lock()
	defer self.Unlock()
	prev := self.head
	curr := prev.next
	for curr != nil {
		if curr.val == x {
			return false
		}
		if curr.val > x {
			nodeCoarse := nodeCoarse{
				next: curr,
				val: x,
			}
			prev.next = &nodeCoarse
			return true
		}
		this := curr
		curr = curr.next
		prev = this
	}
	return false
}

func(self *coarseLocked) remove(x int) bool {
	if x == math.MaxInt64 || x == math.MinInt64 {
		return false
	}
	self.Lock()
	defer self.Unlock()
	prev := self.head
	curr := prev.next
	for curr != nil {
		if curr.val > x {
			return false
		}
		if curr.val == x {
			prev.next = curr.next
			return true
		}
		this := curr
		curr = curr.next
		prev = this
	}
	return false
}

func(self *coarseLocked) contains(x int) bool {
	if x == math.MaxInt64 || x == math.MinInt64 {
		return false
	}
	self.Lock()
	defer self.Unlock()
	curr := self.head
	for curr != nil {
		if curr.val == x {
			return true
		}
		curr = curr.next
	}
	return false
}
