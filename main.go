package main

import (
	"github.com/viviviviviid/go-coin/cli"
	"github.com/viviviviviid/go-coin/db"
)

func main() {
	defer db.Close()
	db.InitDB()
	cli.Start()
}
