package main

import (
	"encoding/base64"
	"fmt"
	"log"
)

var encode = base64.StdEncoding.EncodeToString
var decode = base64.StdEncoding.DecodeString

type print func(struct{}) string

func printarray(name string, array [][]byte) {
	for i, r := range array {
		fmt.Println(fmt.Sprintf("%s[%d] =", name, i), encode(r))
	}
}

func fatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
