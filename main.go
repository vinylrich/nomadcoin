package main

import (
	"fmt"
	"sync"
	"time"
)

var wait sync.WaitGroup

func countToTen(c chan<- int) {
	for i := range [10]int{} {
		fmt.Println(i, c)
		c <- i
		time.Sleep(1 * time.Second)
	}
	close(c)
}

func receive(c <-chan int) { //받기 전용 채널
	for {
		a, ok := <-c

		if !ok {
			fmt.Println("NO VALUE")
			break
		}
		fmt.Println(a)

	}

}

func main() {
	c := make(chan int)
	go countToTen(c)
	receive(c)
}
