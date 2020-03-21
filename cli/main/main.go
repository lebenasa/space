package main

import (
	"log"
	"os"

	"github.com/lebenasa/space/cli"
)

func main() {
	err := cli.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
