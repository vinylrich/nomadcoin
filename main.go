package main

import (
	"fmt"
	"sync"
	"time"
)

var wait sync.WaitGroup

func countToTen(c chan<- int) {
	for i := range [10]int{} {
		fmt.Println(">>sending<<", i)
		c <- i
	}
	close(c)
}

func receive(c <-chan int) { //받기 전용 채널
	for {
		time.Sleep(10 * time.Second)
		a, ok := <-c

		if !ok {
			fmt.Println("NO VALUE")
			break
		}
		fmt.Println("receive", a)

	}

}

//메세지를 엄청 빠르게 보내고 천천히 처리할 수 있음
func main() {
	c := make(chan int, 1) //blocking 없이 최대 n개의 값을 보낼 수 있다.
	go countToTen(c)
	receive(c)
}
