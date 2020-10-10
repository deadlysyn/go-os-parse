package main

import (
	"github.com/deadlysyn/go-os-parse/detector"

	"fmt"
	"log"
)

func main() {
	p, err := detector.PackageManagerCmd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Detected package manager: %q\n", p)
}
