package main

import (
	"github.com/viviviviviid/go-coin/cli"
	"github.com/viviviviviid/go-coin/db"
)

func main() {
	defer db.Close()
	cli.Start()
}

// Mux : Multiplexer
// 하나의 Mux가 3000번과 4000번의 "/" 을 동시에 보고있기때문에 오류가 발생함
// port가 서로 달라서 분리된것 같지만, http.ListenAndServe 함수에서 두번째 인자에 nil을 넣으면,
// DefaultServeMux 가 default로 사용되기 때문
// 즉 우리만의 Mux를 만들어서 사용해야함

// func countToTen(c chan<- int) {
// 	for i := range [10]int{} {
// 		fmt.Printf(">> sending %d <<\n", i)
// 		c <- i
// 		fmt.Printf(">> sent %d <<\n", i)
// 	}
// 	close(c)
// }

// func receive(c <-chan int) {
// 	for {
// 		time.Sleep(10 * time.Second)
// 		a, ok := <-c
// 		if !ok {
// 			fmt.Println("we are done.")
// 			break
// 		}
// 		fmt.Printf("|| received %d ||\n", a)
// 	}
// }

// func main() {
// 	c := make(chan int, 5)
// 	// buffer channel : 누군가 꺼내기전까지 기다리지 않고 넣을 수 있는 channel
// 	// EX) make(chan int) 빈칸으로 둔다면 default size는 1
// 	// EX) make(chan int, 5) 5 개가 들어갈때까지 block 없고, 처음엔 5개까지 들어오는대로 보냄
// 	// 받을 channel에서 하나씩 받는다고 가정하면, 하나 받으면 buffer channel의 공간이 하나가 남으므로, 추가적으로 하나를 더 보냄.
// 	// 받을 채널이 받는 것에 따라서 buffer channel이 queue 형태로 작동한다고 보면됨
// 	go countToTen(c)
// 	receive(c)
// }
