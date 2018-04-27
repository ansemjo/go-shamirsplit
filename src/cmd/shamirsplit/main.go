package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/akamensky/argparse"
	"github.com/ansemjo/shamir/src/sharding"
)

func main() {

	// init parser and add flags
	parser := argparse.NewParser("shamirsplit", "Split data with Shamir secret sharing.")

	// commands
	create := parser.NewCommand("create", "split stdin into pem shards")
	combine := parser.NewCommand("combine", "reconstruct data from pem shards")

	threshold := create.Int("t", "threshold", &argparse.Options{
		Required: true,
		Help:     "minimum number of shares needed for reconstruction",
	})

	shares := create.Int("s", "shares", &argparse.Options{
		Required: true,
		Help:     "total number of shares to create",
	})

	description := create.String("", "description", &argparse.Options{
		Required: false,
		Help:     "add a short description to the PEM blocks",
	})

	// parse arguments and exit if necessary
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	// read stdin for input
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	// decide which command to run
	if create.Happened() {

		shards, err := sharding.CreateShards(*threshold, *shares, stdin, *description)
		if err != nil {
			log.Fatal(err)
		}

		for _, s := range shards {
			pem, err := s.MarshalPEM()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(string(pem) + "\x00")
		}

	}

	if combine.Happened() {

		// TODO: when combining files written with
		// ... | while read -d '' pem; do printf "%s" "$pem" > pem.$((i++)); done
		// later with cat pem.* | .. combine
		// there is a "index out of range" error. probably something about missing EOLs?

		shards, err := sharding.ReadAll(stdin)
		if err != nil {
			log.Fatal(err)
		}

		pshards := sharding.ExtractProtoShards(shards)

		data, err := sharding.CombineShards(pshards)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(string(data))

	}

}
