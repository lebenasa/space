package main

import (
	"log"
	"os"

	"github.com/lebenasa/cli"
)

func main() {
	err := cli.Run(func() []string { return os.Args })
	if err != nil {
		log.Fatalln(err)
	}
}
