package main

import (
	"fmt"

	"github.com/ansemjo/shamir/src/sharding"
	"github.com/ansemjo/shamir/src/util"
)

var (
	// key material
	demokey = "Zxky/LE10mbSdeT4Z3cPoJVcK5Vz3A/oRIR3DcUbgM8="
	// parameters
	threshold   = 3
	shares      = 5
	description = "Testing creation of shards."
	// just a paragraph of lorem ipsum
	lipsum = []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Ut auctor velit at urna sodales porta. Quisque tempor rutrum porttitor. Donec ac mi finibus, efficitur urna vitae, imperdiet turpis. Morbi dictum, est convallis mollis egestas, ante nunc auctor odio, in congue leo leo vitae mauris. Aliquam ornare ultricies dui vel fermentum. Sed tellus ligula, hendrerit volutpat luctus commodo, commodo a mi. Maecenas at fermentum turpis. Nullam interdum ex sed turpis venenatis, et facilisis est dignissim.")
)

func main() {

	shards, err := sharding.CreateShards(threshold, shares, lipsum, description)
	util.Fatal(err)

	for _, s := range shards {

		pem, err := s.MarshalPEM()
		util.Fatal(err)
		// 	fmt.Print(string(pem) + "\x00") // add null for easier splitting
		fmt.Print(string(pem))

		// 	//s.Inspect()

	}

	shards[2].Proto.Pubkey = []byte("\xde\xad\xbe\xef")
	shards[3].Proto.Index = 2

	pshards := sharding.ExtractProtoShards(shards)

	data, err := sharding.CombineShards(pshards)
	util.Fatal(err)

	fmt.Println(string(data))

}
