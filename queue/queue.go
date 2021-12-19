package queue

import (
	"errors"
	"sync"
)

const (
	DefaultQueusSize = 1 << 24
)

var (
	ErrOverflow            = errors.New("queue overflow")
	ErrInvalidtype         = errors.New("invalid type")
	ErrHasFuncDisactivated = errors.New("has function is disactivated")
)

type Entry struct {
	Key   string
	Value interface{}
}

type Queue struct {
	sync.Mutex

	elements   []interface{}
	head       int
	tail       int
	headIsLeft bool

	useHasFunc bool            // whether using `has` function
	elementMap map[string]bool // for checking having
}

func NewQueue(size int, useHasFunc bool) (q Queue) {
	if size == 0 {
		size = DefaultQueusSize
	}
	q.elements = make([]interface{}, size)
	q.head = 0
	q.tail = 0
	q.headIsLeft = true

	q.useHasFunc = useHasFunc
	if useHasFunc {
		q.elementMap = make(map[string]bool, size)
	}

	return
}

func (q *Queue) Enqueue(element interface{}) error {
	q.Lock()
	defer q.Unlock()

	if !q.headIsLeft && q.tail == q.head {
		return ErrOverflow
	}

	q.elements[q.tail] = element
	q.tail += 1
	if q.tail == cap(q.elements) {
		q.tail = 0
		q.headIsLeft = !q.headIsLeft
	}

	if !q.useHasFunc {
		return nil
	}

	entry, ok := element.(Entry)
	if !ok {
		return ErrInvalidtype
	}
	q.elementMap[entry.Key] = true

	return nil
}

func (q *Queue) Dequeue() (element interface{}, isEmpty bool) {
	q.Lock()
	defer q.Unlock()

	if q.headIsLeft && q.tail-q.head == 0 {
		isEmpty = true
		return
	}

	element = q.elements[q.head]
	q.head += 1
	if q.head == cap(q.elements) {
		q.head = 0
		q.headIsLeft = !q.headIsLeft
	}

	if q.useHasFunc {
		entry := element.(Entry)
		delete(q.elementMap, entry.Key)
	}

	return
}

func (q *Queue) Has(key string) (bool, error) {
	q.Lock()
	defer q.Unlock()

	if !q.useHasFunc {
		return false, ErrHasFuncDisactivated
	}

	return q.elementMap[key], nil
}

func (q *Queue) IsEmpty() bool {
	q.Lock()
	defer q.Unlock()

	return q.headIsLeft && q.tail-q.head == 0
}

func (q *Queue) Len() int {
	q.Lock()
	defer q.Unlock()

	if q.headIsLeft {
		return q.tail - q.head
	}

	return q.tail + cap(q.elements) - q.head
}
