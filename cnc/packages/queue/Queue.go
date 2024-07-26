package queue

type Queue []interface{}

func (queue *Queue) Push(x interface{}) {
	*queue = append(*queue, x)
}

func (queue *Queue) Pop() interface{} {
	h := *queue
	var el interface{}
	l := len(h)
	el, *queue = h[0], h[1:l]
	return el
}

func (queue *Queue) Length() int {
	return len(*queue)
}

func NewQueue() *Queue {
	return &Queue{}
}
