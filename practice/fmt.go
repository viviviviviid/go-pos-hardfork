// package main

// import "fmt"

func main() {
	name := 3124124
	fmt.Println("%b", name) // 내용 그대로 터미널에 출력
	fmt.Printf("%b", name)  // 변환된 내용을 터미널에 출력

	hi := fmt.Sprintf("%b", name) // 변환된 내용을
	fmt.Println(name, hi)
}
