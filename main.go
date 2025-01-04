package main

import (
	"flag"
	"fmt"
	"main/formatter"
	"os"
)

func main() {
	var file string
	flag.StringVar(&file, "file", "", "File to read from")
	flag.Parse()

	fmt.Println("formatting:" + file)

	sql, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	f := formatter.Create(string(sql))
	fmt.Println(f.Format())
}
