package workstealing

import (
	"sync/atomic"
)

const bufferSize = 1024

// Implemented the lock-free bounded dequeue using atomic operations
type Deque struct {
	buffer [bufferSize]string
	head   uint32
	tail   uint32
}

// PushBottom inserts an item at the bottom of the dequeue.
func (d *Deque) PushBottom(item string) bool {
	tail := atomic.LoadUint32(&d.tail)
	nextTail := (tail + 1) % bufferSize
	if nextTail == atomic.LoadUint32(&d.head) {
		return false // Queue is full
	}

	d.buffer[tail] = item
	atomic.StoreUint32(&d.tail, nextTail)
	return true
}

// PopBottom removes and returns an item from the bottom of the dequeue.
func (d *Deque) PopBottom() (string, bool) {
	tail := atomic.LoadUint32(&d.tail)
	if tail == atomic.LoadUint32(&d.head) {
		return "", false // Queue is empty
	}
	prevTail := (tail - 1 + bufferSize) % bufferSize
	item := d.buffer[prevTail]
	atomic.StoreUint32(&d.tail, prevTail)
	return item, true
}

// PopTop removes and returns an item from the top of the queue.
func (d *Deque) PopTop() (string, bool) {
	head := atomic.LoadUint32(&d.head)
	if head == atomic.LoadUint32(&d.tail) {
		return "", false // Queue is empty
	}
	item := d.buffer[head]
	nextHead := (head + 1) % bufferSize
	atomic.StoreUint32(&d.head, nextHead)
	return item, true
}

// IsEmpty checks if the dequeue is empty.
func (d *Deque) IsEmpty() bool {
	return atomic.LoadUint32(&d.head) == atomic.LoadUint32(&d.tail)
}
