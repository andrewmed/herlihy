package herlihy

import (
	"math"
)

type chanList struct {
	cAdd, cContains, cRemove chan int
	response                 chan bool
	head *nodeChan // sentinel int min
	tail *nodeChan // sentinel int max
}

type nodeChan struct {
	next *nodeChan
	val int
}

func newChan() chanList {
	head := nodeChan{
		val: math.MinInt64,
	}
	tail := nodeChan{
		val: math.MaxInt64,
	}
	head.next = &tail
	list :=  chanList{
		make(chan int),
		make(chan int),
		make(chan int),
		make(chan bool),
		&head,
		&tail,
	}
	go list.serve()
	return list
}

func (self *chanList) serve() {
	var response bool
	var prev, curr *nodeChan
	for {
		select {
		case val := <-self.cAdd:
			for prev, curr = self.head, self.head.next ; curr != nil; prev, curr = curr, curr.next {
				if curr.val == val {
					response = false
					break
				}
				if curr.val > val {
					this := nodeChan{
						next: curr,
						val:  val,
					}
					prev.next = &this
					response = true
					break
				}
			}
		case val := <-self.cContains:
			for curr := self.head; curr != nil; curr = curr.next {
				if curr.val == val {
					response = true
					break
				}
				if curr.val > val {
					response = false
					break
				}
			}
		case val := <- self.cRemove:
			for prev, curr = self.head, self.head.next ; curr != nil; prev, curr = curr, curr.next {
				if curr.val == val {
					prev.next = curr.next
					response = true
					break
				}
				if curr.val > val {
					response = false
					break
				}
			}
		}
		self.response <- response
	}
}

func (self *chanList) add(x int) bool {
	self.cAdd <- x
	return <-self.response
}

func (self *chanList) remove(x int) bool {
	self.cRemove <- x
	return <-self.response
}

func (self *chanList) contains(x int) bool {
	self.cContains <- x
	return <-self.response
}
