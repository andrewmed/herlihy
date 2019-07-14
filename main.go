package herlihy

// concurrent linked list
// chapter 9 of Prof. Herlihy's "The Art of Multiprocessor Programming"

import(
	"math"
)

var (
	startSentinel   = math.MinInt64
	endSentinel     = math.MaxInt64
	deletedSentinel = math.MinInt64 + 1
)

type list interface {
	// add returns true if x was not in the list
	add(x int) bool
	// remove returns true if x was in the list
	remove(x int) bool
	// contains return true if x was in the list
	contains(x int) bool
}
