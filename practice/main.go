package main

import (
	"fmt"

	"github.com/viviviviviid/go-coin/practice/person"
)

func main() {
	nico := person.Person{}
	nico.SetDetails("nico", 12) // 이걸 했다하더라도, 복사본이 변경된거지 여기 nico는 변하지 않았음
	fmt.Println("Main 'nico", nico)
}
