package herlihy

import (
	"sync/atomic"
	"unsafe"
)

// we support invariant that deleted node (ref) cant be updated, to do so add must use traverse() which physically removes marked nodes

type nodeAtomic struct {
	p   unsafe.Pointer
	val int
}

type atomicList struct {
	head  *nodeAtomic // sentinel int min
}

func newAtomic() atomicList {
	head := nodeAtomic{
		val: startSentinel,
	}
	tail := nodeAtomic{
		val: endSentinel,
	}
	head.p = unsafe.Pointer(&tail)
	return atomicList{
		&head,
	}
}

func (self *nodeAtomic) isDeleted() bool {
	next := (*nodeAtomic)(self.p)
	if next == nil {
		return false
	}
	return next.val == deletedSentinel
}

func (self *atomicList) debugPrint() {
	this := self.head
	for this != nil {
		print(this.val, " ")
		this = (*nodeAtomic)(this.p)
	}
	println()
}

func (self *atomicList) traverse(x int) (prev, curr *nodeAtomic) {
startover:
	prev = self.head
	curr = (*nodeAtomic)(prev.p)
	for curr != nil {
		// invariant: we CAS pointer only if next nodeAtomic is not a delete sentinel
		p := prev.p
		if curr != (*nodeAtomic)(p) {
			goto startover
		}
		next := (*nodeAtomic)(curr.p)
		if next != nil && next.val == deletedSentinel {
			if !atomic.CompareAndSwapPointer(&prev.p, p, next.p) {
				goto startover
			}
			curr = (*nodeAtomic)(next.p)
			continue
		}
		if curr.val >= x {
			return prev, curr
		}
		prev = curr
		curr = (*nodeAtomic)(curr.p)
	}
	panic(prev.val)
}

func (self *atomicList) add(x int) bool {
	for {
		prev, curr := self.traverse(x)
		p := prev.p
		if curr == (*nodeAtomic)(p) {
			if curr.val == x {
				return false
			}
			if self.addDo(prev, p, x) {
				return true
			}
		}
	}
}

func (self *atomicList) addDo(ref *nodeAtomic, p unsafe.Pointer, x int) bool {
	this := nodeAtomic{
		p:   ref.p,
		val: x,
	}
	return atomic.CompareAndSwapPointer(&ref.p, p, unsafe.Pointer(&this))
}

func (self *atomicList) remove(x int) bool {
	for {
		_, curr := self.traverse(x)
		if curr.val != x {
			return false
		}
		if self.addDo(curr, curr.p, deletedSentinel) {
			return true
		}
	}
}

func (self *atomicList) contains(x int) bool {
	_, curr := self.traverse(x)
	return curr.val == x
}
