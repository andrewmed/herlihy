package herlihy

import (
	"math"
)

type nodeRaces struct {
	next *nodeRaces
	val int
}

type Races struct {
	head *nodeRaces // sentinel int min
	tail *nodeRaces // sentinel int max
}

func newRaces() Races {
	head := nodeRaces{
		val: math.MinInt64,
	}
	tail := nodeRaces{
		val: math.MaxInt64,
	}
	head.next = &tail
	return Races{
		&head,
		&tail,
	}
}

func(self *Races) add(x int) bool {
	if x == math.MaxInt64 || x == math.MinInt64 {
		return false
	}
	prev := self.head
	curr := prev.next
	for curr != nil {
		if curr.val == x {
			return false
		}
		if curr.val > x {
			nodeRaces := nodeRaces{
				next: curr,
				val: x,
			}
			prev.next = &nodeRaces
			return true
		}
		this := curr
		curr = curr.next
		prev = this
	}
	return false
}

func(self *Races) remove(x int) bool {
	if x == math.MaxInt64 || x == math.MinInt64 {
		return false
	}
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

func(self *Races) contains(x int) bool {
	if x == math.MaxInt64 || x == math.MinInt64 {
		return false
	}
	curr := self.head
	for curr != nil {
		if curr.val == x {
			return true
		}
		curr = curr.next
	}
	return false
}
