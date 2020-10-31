package queue

import "errors"

var ErrEmptyQueue = errors.New("cannot pop from empty queue")

type Queue struct {
	data []interface{}
}

func NewQueue() (q *Queue) {
	return &Queue{
		data: []interface{}{},
	}
}

func (q *Queue) Push(i interface{}) {
	q.data = append(q.data, i)
}

func (q *Queue) Pop() (i interface{}, err error) {
	if len(q.data) == 0 {
		return i, ErrEmptyQueue
	}
	item, items := q.data[0], q.data[1:]
	q.data = items
	items = nil
	return item, nil
}
