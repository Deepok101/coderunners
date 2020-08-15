package utils

import (
	"errors"
)

// Queue represents the queue data structure
type Queue interface {
	Enqueue(interface{}) error
	Dequeue() error
}

// Code represents a data structure holding a programming language and its code text
type Code struct {
	language string
	content  string
}

type CodeQueue struct {
	list   []Code
	length int
}

// NewCodeQueue instantiates a codeQueue reference
func NewCodeQueue() *CodeQueue {
	return &CodeQueue{
		list:   []Code{},
		length: 0,
	}
}

func (c *CodeQueue) Enqueue(el interface{}) error {
	codeEl, ok := el.(Code)
	if ok == false {
		return errors.New("Element to be enqueued is not of type code")
	}

	newList := append(c.list, codeEl)
	c.list = newList
	c.length = len(newList)
	return nil
}

func (c *CodeQueue) Dequeue() error {
	if c.length == 0 {
		return errors.New("Nothing to deqeue")
	}
	newList := c.list[1:]
	c.list = newList
	c.length = len(newList)

	return nil
}
