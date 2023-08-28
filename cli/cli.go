package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/viviviviviid/go-coin/explorer"
	"github.com/viviviviviid/go-coin/rest"
)

func usage() {
	fmt.Printf("Welcome to 민석's Blockchain Project\n\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-port:	Set the PORT of the server\n")
	fmt.Printf("-mode:	Choose between 'html' and 'rest'\n")
	os.Exit(0) // 프로그램 정지 및 에러 코드 // 0은 문제 없음 // 1부터 에러
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}

	port := flag.Int("port", 4000, "Set port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")
	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	default:
		usage()
	}

	fmt.Println(*port, *mode)
}
