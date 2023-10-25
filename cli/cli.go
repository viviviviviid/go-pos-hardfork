package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/viviviviviid/go-coin/rest"
)

func usage() {
	fmt.Printf("Welcome to 민석's Blockchain Project\n\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-port:	Set the PORT of the server\n")
	fmt.Printf("-mode:	Choose between 'html' and 'rest'\n")
	os.Exit(0) // 프로그램 정지 및 에러 코드 // 0은 문제 없음 // 1부터 에러
	// 설명보면 defer 이후에 사용해야한다고 나와있음 // defer도 죽이기때문아닐까 // 그래서 runtime 패키지 이용하는거고
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}
	port := flag.Int("port", 3000, "Set port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html', 'rest' and 'both' (both mean html and rest)")
	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	// case "html":
	// 	explorer.Start(*port)
	// case "both":
	// 	go explorer.Start(*port)
	// 	rest.Start(*port + 1)
	// 	// fmt.Scanln() // go가 받을때 main 함수가 먼저 종료되지 않게 대기하기 위해서는 fmt.Scanln()를 입력해줘야 한다.
	default:
		usage()
	}

}
