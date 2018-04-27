package util

import (
	"encoding/base64"
	"fmt"

	"github.com/ansemjo/shamir/src/sharding"
)

var enc = base64.StdEncoding.EncodeToString

// Inspect logs the internal structure to the console.
func Inspect(s *sharding.Shard) {

	fmt.Println(color(red, "Shard "+fmt.Sprint(s.UUID)+":"), "index", s.Proto.Index)
	fmt.Println(color(green, " Threshold :"), s.Proto.Associated.Threshold)
	fmt.Println(color(green, " Shares    :"), s.Proto.Associated.Shares)
	fmt.Println(color(yellow, " Keyshare  :"), enc(s.Proto.Keyshare))
	fmt.Println(color(yellow, " Pubkey    :"), enc(s.Proto.Pubkey))
	fmt.Println(color(yellow, " Signature :"), enc(s.Proto.Signature))
	fmt.Println(color(red, " Data      :"), enc(s.Proto.Data))

}

func color(c, s string) string {
	return fmt.Sprintf("\033["+c+";1m%s\033[0m", s)
}

const (
	red    = "31"
	green  = "32"
	yellow = "33"
	blue   = "34"
)
