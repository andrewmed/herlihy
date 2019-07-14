package herlihy

import (
	"math"
	"sync"
)

type nodeFine struct {
	next *nodeFine
	sync.Mutex // fine grain sync
	val int
}

type fineLocked struct {
	head *nodeFine // sentinel int min
	tail *nodeFine // sentinel int max
}

func newFineLocked() fineLocked {
	head := nodeFine{
		val: math.MinInt64,
	}
	tail := nodeFine{
		val: math.MaxInt64,
	}
	head.next = &tail
	return fineLocked{
		&head,
		&tail,
	}
}

func(self *fineLocked) add(x int) bool {
	if x == math.MaxInt64 || x == math.MinInt64 {
		return false
	}
	self.head.Lock()
	prev := self.head
	curr := prev.next
	for curr != nil {
		curr.Lock()
		if curr.val == x {
			prev.Unlock()
			curr.Unlock()
			return false
		}
		if curr.val > x {
			nodeFine := nodeFine{
				next: curr,
				val: x,
			}
			prev.next = &nodeFine
			prev.Unlock()
			curr.Unlock()
			return true
		}
		this := curr
		prev.Unlock()
		curr = curr.next
		prev = this
	}
	prev.Unlock()
	return false
}

func(self *fineLocked) remove(x int) bool {
	if x == math.MaxInt64 || x == math.MinInt64 {
		return false
	}
	self.head.Lock()
	prev := self.head
	curr := prev.next
	for curr != nil {
		if curr.val > x {
			prev.Unlock()
			return false
		}
		curr.Lock()
		if curr.val == x {
			prev.next = curr.next
			prev.Unlock()
			curr.Unlock()
			return true
		}
		this := curr
		prev.Unlock()
		curr = curr.next
		prev = this
	}
	prev.Unlock()
	return false
}

func(self *fineLocked) contains(x int) bool { // fixme: needs locking?
	if x == math.MaxInt64 || x == math.MinInt64 {
		return false
	}
	curr := self.head
	for curr != nil {
		curr.Lock()
		if curr.val == x {
			curr.Unlock()
			return true
		}
		curr.Unlock()
		curr = curr.next
	}
	return false
}
