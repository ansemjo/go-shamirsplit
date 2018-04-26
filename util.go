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

func y(s string) string {
	return fmt.Sprintf("\033[33;1m%s\033[0m", s)
}

func b(s string) string {
	return fmt.Sprintf("\033[34;1m%s\033[0m", s)
}

func r(s string) string {
	return fmt.Sprintf("\033[31;1m%s\033[0m", s)
}

func g(s string) string {
	return fmt.Sprintf("\033[32;1m%s\033[0m", s)
}
