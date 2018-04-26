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
	shares    = 5
	message   = "Testing"
	// demo data: use lipsum from lipsum.go
)

func main() {

	// message data
	uuid := RandomUUID()
	key, err := decode(demokey)
	fatal(err)

	// assemble and serialize associated data
	associated := &AssociatedData{
		Uuid:      uuid[:],
		Shares:    shares,
		Threshold: threshold,
	}
	ad, err := proto.Marshal(associated)
	fatal(err)

	// use the first 12 bytes of uuid as nonce
	nonce := uuid[:12]

	ciphertext, err := Encrypt(key, nonce, ad, lipsum)
	//fmt.Println("Nonce:", encode(nonce))
	//fmt.Println("Ciphertext:", encode(ciphertext))

	ciphertext, err = Pad(ciphertext, threshold)
	fatal(err)

	// perform split and reed-solomon encoding
	rscoded, err := RSEncode(ciphertext, threshold, shares)
	fatal(err)

	// perform keysplit
	keyshares, err := CreateKeyShares(key, threshold, shares)
	fatal(err)

	// basic sanity check
	if len(keyshares) != len(rscoded) {
		fatal(errors.New("keyshares and reed-solomon encoded data have different lengths"))
	}

	// assemble and sign shards
	shards := make([]*Shard, len(keyshares))
	for i := range shards {

		// assemble shard
		p := &ProtoShard{
			Associated: associated,
			Index:      int32(i),
			Keyshare:   keyshares[i],
			Data:       rscoded[i],
		}

		// marshal shard w/o signature and pubkey
		m, err := proto.Marshal(p)
		fatal(err)

		// sign protobuf and amend shard
		sig, pub, err := Sign(key, m)
		fatal(err)
		p.Pubkey = pub
		p.Signature = sig

		// finalize shard for pem construction
		shards[i] = &Shard{
			Threshold: threshold,
			UUID:      uuid,
			Proto:     p,
		}

	}

	for _, s := range shards {

		pem, err := EncodePEM(s)
		fatal(err)
		fmt.Println(string(pem))

	}

}
