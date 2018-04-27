package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"
	"github.com/ansemjo/shamir/src/sharding"
)

var (
	// just a paragraph of lorem ipsum
	lipsum = []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Ut auctor velit at urna sodales porta. Quisque tempor rutrum porttitor. Donec ac mi finibus, efficitur urna vitae, imperdiet turpis. Morbi dictum, est convallis mollis egestas, ante nunc auctor odio, in congue leo leo vitae mauris. Aliquam ornare ultricies dui vel fermentum. Sed tellus ligula, hendrerit volutpat luctus commodo, commodo a mi. Maecenas at fermentum turpis. Nullam interdum ex sed turpis venenatis, et facilisis est dignissim.")
)

func main() {

	// init parser and add flags
	parser := argparse.NewParser("shamirsplit", "Split data with Shamir secret sharing.")

	// TODO: mode, create / combine

	threshold := parser.Int("t", "threshold", &argparse.Options{
		Required: true,
		Help:     "minimum number of shares needed for reconstruction",
	})

	shares := parser.Int("s", "shares", &argparse.Options{
		Required: true,
		Help:     "total number of shares to create",
	})

	description := parser.String("", "description", &argparse.Options{
		Required: false,
		Help:     "add a short description to the PEM blocks",
	})

	// parse arguments and exit if necessary
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	///////////////////////////////

	shards, err := sharding.CreateShards(*threshold, *shares, lipsum, *description)
	if err != nil {
		log.Fatal(err)
	}

	pemcollect := make([]byte, 0)

	for _, s := range shards {
		pem, err := s.MarshalPEM()
		if err != nil {
			log.Fatal(err)
		}
		pemcollect = append(pemcollect, pem...)
	}

	fmt.Print(string(pemcollect))

	// shards[2].Proto.Pubkey = []byte("\xde\xad\xbe\xef")
	// shards[3].Proto.Index = 2

	readshards, err := sharding.ReadAll(pemcollect)
	if err != nil {
		log.Fatal(err)
	}

	pshards := sharding.ExtractProtoShards(readshards)

	data, err := sharding.CombineShards(pshards)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))

}
