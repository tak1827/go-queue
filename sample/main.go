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

	value := "string"
	err := q.Enqueue(value)
	if err != nil {
		panic(err.Error())
	}

	elm, _ := q.Dequeue()
	elm = elm.(string)
	if elm != value {
		panic("unexpected value")
	}

	fmt.Printf("dequeue: %v\n", elm)

	/***********
	 use `has` function
	************/
	useHasFunc = true
	size = 100
	q = queue.NewQueue(size, useHasFunc)

	entry := queue.Entry{
		Key:   "key",
		Value: "value",
	}
	err = q.Enqueue(entry)
	if err != nil {
		panic(err.Error())
	}

	has, err := q.Has(entry.Key)
	if err != nil {
		panic(err.Error())
	}
	if !has {
		panic("expected to have")
	}

	elm, _ = q.Dequeue()
	elm = elm.(queue.Entry)
	if elm != entry {
		panic("unexpected value")
	}

	fmt.Printf("dequeue: %v\n", elm)
}
