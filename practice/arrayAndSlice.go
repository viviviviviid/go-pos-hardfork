// package main

// import "fmt"

// func main() {

// 	// array는 크기를 지정해줘야함
// 	foods := [3]string{"potato", "pizza", "pasta"}

// 	// range 방식
// 	for _, food := range foods {
// 		fmt.Println(food)
// 	}

// 	// for in 방식
// 	for i := 0; i < len(foods); i++ {
// 		fmt.Println(foods[i])
// 	}

// 	// slice는 javascript의 array 처럼 내용을 계속 더할 수 있음
// 	foods := []string{"potato", "pizza", "pasta"}
// 	fmt.Printf("%v\n", foods)
// 	foods = append(foods, "tomato") // foods로 받아야 비로소 업데이트
// 	fmt.Printf("%v\n", foods)

// }
