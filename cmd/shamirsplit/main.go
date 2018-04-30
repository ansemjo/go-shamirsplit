package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/akamensky/argparse"
	"github.com/ansemjo/go-shamirsplit/sharding"
)

func main() {

	// init parser and add flags
	parser := argparse.NewParser("shamirsplit", "Split data with Shamir secret sharing.")

	// commands
	create := parser.NewCommand("create", "split stdin into pem shards")
	combine := parser.NewCommand("combine", "reconstruct data from pem shards")

	// creation arguments
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
	nullbyte := create.Flag("0", "null", &argparse.Options{
		Required: false,
		Help:     "terminate each pem block on stdout with a null byte",
	})
	outdir := create.String("d", "directory", &argparse.Options{
		Required: false,
		Help:     "output pem blocks to files in this directory",
	})
	outname := create.String("n", "name", &argparse.Options{
		Required: false,
		Help:     "filename stub for writing to directory; will be suffixed with _000.pem, _001.pem etc.",
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

		// check if outdir is a valid directory
		if *outdir != "" {

			if *nullbyte {
				// null termination does not make sense for direct writing
				*nullbyte = false
				fmt.Fprintln(os.Stderr, "disabling nullbyte")
			}

			dir, err := os.Open(*outdir)
			if err != nil {
				log.Fatal(err)
			}
			defer dir.Close()

			stat, err := dir.Stat()
			if err != nil {
				log.Fatal(err)
			}

			if !stat.IsDir() {
				log.Fatal(fmt.Errorf("not a directory: %q", *outdir))
			}

		}

		// create shards
		shards, err := sharding.CreateShards(*threshold, *shares, stdin, *description)
		if err != nil {
			log.Fatal(err)
		}

		// serialize and output pem blocks
		for i, s := range shards {

			pem, err := s.MarshalPEM()
			if err != nil {
				log.Fatal(err)
			}

			writefile := func(f *os.File) {
				if *nullbyte {
					pem = append(pem, '\x00')
				}
				_, err = fmt.Fprintf(f, "%s", pem)
				if err != nil {
					log.Fatal(err)
				}
			}

			if *outdir == "" {

				writefile(os.Stdout)

			} else {

				var filename string
				if *outname == "" {
					filename = fmt.Sprintf("shard_%s_%03d.pem", s.UUID, i)
				} else {
					filename = fmt.Sprintf("%s_%03d.pem", *outname, i)
				}

				pathname := path.Join(*outdir, filename)
				fmt.Println(pathname)

				f, err := os.Create(pathname)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()

				writefile(f)

			}

		}

	}

	if combine.Happened() {

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
