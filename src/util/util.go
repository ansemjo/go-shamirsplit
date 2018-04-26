package util

import (
	"encoding/base64"
	"fmt"
	"log"
)

var Base64encode = base64.StdEncoding.EncodeToString
var Base64decode = base64.StdEncoding.DecodeString

type print func(struct{}) string

func PrintArray(name string, array [][]byte) {
	for i, r := range array {
		fmt.Println(fmt.Sprintf("%s[%d] =", name, i), Base64encode(r))
	}
}

func Fatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func Y(s string) string {
	return fmt.Sprintf("\033[33;1m%s\033[0m", s)
}

func B(s string) string {
	return fmt.Sprintf("\033[34;1m%s\033[0m", s)
}

func R(s string) string {
	return fmt.Sprintf("\033[31;1m%s\033[0m", s)
}

func G(s string) string {
	return fmt.Sprintf("\033[32;1m%s\033[0m", s)
}
