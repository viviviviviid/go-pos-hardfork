// package main

// import "fmt"

// type person struct {
// 	name string
// 	age  int
// }

//  // struct의 method 만들기
// func (p person) sayHello() {  // p는 person의 앞글자로 지정 -> 관습
// 	fmt.Printf("hello! my name is %s and i'm %d", p.name, p.age)
// }

// func main() {
// 	nico := person{"nico", 12}
// 	nico.sayHello()
// }

package main

import "fmt"

type userInfo struct {
	userName string
	userAge  int
}

func (u userInfo) introducing() {
	fmt.Printf("Hi, My name is %s and my korean age is %d", u.userName, u.userAge)
}

func main() {
	minseok := userInfo{"minseok", 26}
	minseok.introducing()
}
