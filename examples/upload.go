package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gifs/gifs-go"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		return
	}

	res, err := gifs.Upload(f.Name(), f)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%+v", res.Success.Page)
}
