package main

import (
	"errors"
	"fmt"

	proto "github.com/golang/protobuf/proto"
)

const (
	// key material
	demokey = "Zxky/LE10mbSdeT4Z3cPoJVcK5Vz3A/oRIR3DcUbgM8="
	// parameters
	threshold = 3
	total     = 5
	message   = "Testing"
	// demo data: use lipsum from lipsum.go
)

func main() {

	key, err := decode(demokey)
	fatal(err)

	shares, err := CreateKeyShares(key, threshold, total)
	fatal(err)
	//printarray("keyshares", shares)

	associated := &AssociatedData{Message: message, Shares: total, Threshold: threshold}
	ad, err := proto.Marshal(associated)
	fatal(err)

	nonce, ciphertext, err := Encrypt(key, ad, lipsum)
	//fmt.Println("Nonce:", encode(nonce))
	//fmt.Println("Ciphertext:", encode(ciphertext))

	ciphertext, err = Pad(ciphertext, threshold)
	fatal(err)

	rscoded, err := RSEncode(ciphertext, threshold, total)
	fatal(err)

	if len(shares) != len(rscoded) {
		fatal(errors.New("keyshares and reed-solomon encoded data have different lengths"))
	}

	//printarray("splitdata", rscoded)

	shards := make([]Shard, len(shares))
	for i := range shards {
		shards[i] = Shard{message: "Testing", index: i, nonce: nonce, keyshare: shares[i], data: rscoded[i]}
	}

	for _, s := range shards {
		//fmt.Printf("· shard[%d] = %s\n", i, fmt.Sprintf("%v", s))
		fmt.Println(string(EncodePEM(s)))
	}

}
