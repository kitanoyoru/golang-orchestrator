package queue

type (
	Queue struct {
		start, end *node
		length     int
	}
	node struct {
		value interface{}
		next  *node
	}
)

func New() *Queue {
	return &Queue{nil, nil, 0}
}

func (this *Queue) Dequeue() interface{} {
	if this.length == 0 {
		return nil
	}
	n := this.start
	if this.length == 1 {
		this.start = nil
		this.end = nil
	} else {
		this.start = this.start.next
	}
	this.length--
	return n.value
}

func (this *Queue) Enqueue(value interface{}) {
	n := &node{value, nil}
	if this.length == 0 {
		this.start = n
		this.end = n
	} else {
		this.end.next = n
		this.end = n
	}
	this.length++
}

func (this *Queue) Len() int {
	return this.length
}

func (this *Queue) Peek() interface{} {
	if this.length == 0 {
		return nil
	}
	return this.start.value
}
