package main

import (
	"fmt"

	"github.com/tak1827/go-queue/queue"
)

func main() {
	/***********
	 no `has` function
	************/
	useHasFunc := false
	size := 100
	q := queue.NewQueue(size, useHasFunc)

	e := &queue.Entry{
		Key:   "key",
		Value: []byte("value"),
	}
	err := q.Enqueue(e)
	if err != nil {
		panic(err.Error())
	}

	elm, _ := q.Dequeue()
	if elm.Key != e.Key {
		panic("unexpected key")
	}

	fmt.Printf("dequeue: %v\n", elm)

	/********************
	 use `has` function
	*********************/
	useHasFunc = true
	size = 100
	q = queue.NewQueue(size, useHasFunc)

	e = &queue.Entry{
		Key:   "key2",
		Value: []byte("value2"),
	}
	err = q.Enqueue(e)
	if err != nil {
		panic(err.Error())
	}

	has, err := q.Has(e.Key)
	if err != nil {
		panic(err.Error())
	}
	if !has {
		panic("expected to have")
	}

	elm, _ = q.Dequeue()
	if elm.Key != e.Key {
		panic("unexpected value")
	}

	fmt.Printf("dequeue: %v\n", elm)
}
