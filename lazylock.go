package herlihy

import (
	"math"
	"sync"
	"unsafe"
)


// Lazy lock differs from optimistic in no locking in contains() - just check bool marked

type nodeLazy struct {
	p   unsafe.Pointer
	next *nodeLazy
	sync.Mutex // fine grain sync
	marked bool
	val int
}

type lazyLocked struct {
	head *nodeLazy // sentinel int min
	tail *nodeLazy // sentinel int max
}

func newLazy() lazyLocked {
	head := nodeLazy{
		val: math.MinInt64,
	}
	tail := nodeLazy{
		val: math.MaxInt64,
	}
	head.next = &tail
	return lazyLocked{
		&head,
		&tail,
	}
}

func(self *lazyLocked) add(x int) bool {
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
				nodeLazy := nodeLazy{
					next: curr,
					val: x,
				}
				prev.next = &nodeLazy
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

func(self *lazyLocked) remove(x int) bool {
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
				curr.marked = true // lazy removing
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

func (self *lazyLocked)  valid(prev, curr *nodeLazy) bool {
	this := self.head
	for this != nil {
		if this == prev {
			return this.next == curr
		}
		this = this.next
	}
	return false
}

func(self *lazyLocked) contains(x int) bool {
	if x == math.MaxInt64 || x == math.MinInt64 {
		return false
	}
	curr := self.head
	for curr != nil {
		if curr.val == x {
			return !curr.marked
		}
		curr = curr.next
	}
	return false
}
