package main

import (
	"fmt"
	"time"
)

// func main() {
// 	defer db.Close()
// 	cli.Start()
// }

// Mux : Multiplexer
// 하나의 Mux가 3000번과 4000번의 "/" 을 동시에 보고있기때문에 오류가 발생함
// port가 서로 달라서 분리된것 같지만, http.ListenAndServe 함수에서 두번째 인자에 nil을 넣으면,
// DefaultServeMux 가 default로 사용되기 때문
// 즉 우리만의 Mux를 만들어서 사용해야함

func countToTen(c chan int) {
	for i := range [10]int{} {
		c <- i
		time.Sleep(1 * time.Second)
	}
	close(c)
}

func receive(c <-chan int) {
	for {
		a, ok := <-c
		if !ok {
			fmt.Println("we are done.")
			break
		}
		fmt.Println(a, ok)
	}
}

func main() {
	c := make(chan int)
	go countToTen(c)
	receive(c)
}
