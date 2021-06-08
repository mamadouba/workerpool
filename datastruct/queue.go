package datastruct

type Queue []interface{}

func (q *Queue) IsEmpty() bool {
	return len(*q) == 0
}

func (q *Queue) Enque(item interface{}) {
	*q = append(*q, item)
}

func (q *Queue) Deque() (interface{}, bool) {
	if q.IsEmpty() {
		return nil, false
	}
	item := (*q)[0]
	*q = (*q)[1:]
	return item, true
}

func (q *Queue) Front() (interface{}, bool) {
	if q.IsEmpty() {
		return nil, false
	}
	return (*q)[0], true
}
