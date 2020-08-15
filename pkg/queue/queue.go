package queue

type Queue interface {
	enqeue(interface{}) error
	dequeue() error
}

type code struct {
	language string
	content  string
}

type codeQueue struct {
	list   []code
	length int
}

func NewCodeQueue() *codeQueue {
	return codeQueue{
		list:   []code{},
		length: 0,
	}
}

func (c *codeQueue) enqueue(code code) {
	newList := append(c.list, code)
	c.list = newList
	c.length = len(newList)
}

func (c *codeQueue) dequeue() {
	newList := c.list[1:]
	c.list = newList
	c.length = len(newList)
}
