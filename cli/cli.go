package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/viviviviviid/go-coin/pos"
	"github.com/viviviviviid/go-coin/rest"
)

// cli 명령어 기본 가이드 (Ex. go run main.go -mode=rest -port=4000)
func usage() {
	fmt.Printf("Welcome to 민석's Blockchain Project\n\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-port:	Set the PORT of the server\n")
	fmt.Printf("-mode:	Choose between 'auto' and 'rest'\n")
	os.Exit(0)
}

// cli 명령어를 감지하여 auto 또는 rest 모드로 실행
func Start() {
	if len(os.Args) == 1 {
		usage()
	}
	port := flag.Int("port", 4000, "Set port of the server")
	mode := flag.String("mode", "rest", "Choose between 'auto' and 'rest'")
	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "auto":
		pos.PoS(*port)
	default:
		usage()
	}

}
