package main

import (
	"fmt"
	"sync"

	"github.com/sloan-dog/linkchecker/structures"
)

func main() {
	grph := structures.ConfigGraphDefault()
	counter := structures.NewCounter()
	queue := structures.NewQueue()
	tracker := structures.NewTracker()
	fmt.Printf("entries: %v\nnodes: %v\n", grph.EntryPoints, grph.Nodes)
	wg := sync.WaitGroup{}
	fmt.Println(len(grph.EntryPoints))
	l := len(grph.EntryPoints)
	for i := 0; i < l; i++ {
		wg.Add(1)
		go func(index int) {
			structures.CountLinks(&grph, counter, queue, tracker, grph.EntryPoints[index], index)
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Printf("counter: %v", counter.GetCounts())
}
