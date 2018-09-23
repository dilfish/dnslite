package main

import (
	"fmt"
	"github.com/dilfish/dnslite"
)

func main() {
	dnslite.RecordMap = make(map[string][]dnslite.TypeRecord)
	dnslite.HandleHTTP()
	err := dnslite.Handle()
	if err != nil {
		fmt.Println("error is", err)
		panic(err)
	}
}
