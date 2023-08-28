package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Printf("Welcome to 민석's Blockchain Project\n\n")
	fmt.Printf("Please use the following commands:\n\n")
	fmt.Printf("explorer:	Start the HTML Explorer\n")
	fmt.Printf("rest: 		Start the REST API (recommended)\n")
	os.Exit(0) // 프로그램 정지 및 에러 코드 // 0은 문제 없음 // 1부터 에러
}

func main() {

	// os.Args: 터미널 창에서 입력한 내용 -> os.Args[0]: 프로그램이름
	// os.Args[1부터]: 우리가 추가적으로 입력한 내용.
	// ex) go run main.go helllo -> os.Args[1] === helllo
	if len(os.Args) < 2 {
		usage()
	}

	// flag가 많을때 flagSet을 이용
	rest := flag.NewFlagSet("rest", flag.ExitOnError)                 // new flag 집합 // flag 예시 : go run main.go rest -port=9000
	portFlag := rest.Int("port", 4000, "Sets the port of the server") // rest라는 flag 집합에 넣고, 핸들링

	switch os.Args[1] {
	case "explorer":
		fmt.Println("Start Explorer")
	case "rest":
		rest.Parse(os.Args[2:])
	default: // 그 외적으로 기본 값.
		usage()
	}

	if rest.Parsed() {
		fmt.Println(portFlag)
		fmt.Println("Start server")
	}

}

// Mux : Multiplexer
// 하나의 Mux가 3000번과 4000번의 "/" 을 동시에 보고있기때문에 오류가 발생함
// port가 서로 달라서 분리된것 같지만, http.ListenAndServe 함수에서 두번째 인자에 nil을 넣으면,
// DefaultServeMux 가 default로 사용되기 때문
// 즉 우리만의 Mux를 만들어서 사용해야함
