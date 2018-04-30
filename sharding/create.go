package sharding

import (
	"fmt"

	"github.com/ansemjo/go-shamirsplit/cryptography"
	proto "github.com/golang/protobuf/proto"
	"github.com/google/uuid"
)

// TODO: use random keys
const demokey = "Zxky/LE10mbSdeT4Z3cPoJVcK5Vz3A/oRIR3DcUbgM8="

// CreateShards creates key/data shards
// TODO: write specification / description
func CreateShards(threshold, shares int, message []byte, description string) (shards []*Shard, err error) {

	// message data
	uuid, key := uuid.New(), cryptography.RandomKey()
	if err != nil {
		return
	}

	// assemble and serialize associated data
	associated := &AssociatedData{
		Uuid:      uuid[:],
		Shares:    int32(shares),
		Threshold: int32(threshold),
	}
	ad, err := proto.Marshal(associated)
	if err != nil {
		return
	}

	// use the first 12 bytes of uuid as nonce
	nonce := uuid[:12]

	// encrypt the message
	ciphertext, err := cryptography.Encrypt(key, nonce, ad, message)
	if err != nil {
		return
	}

	// split ciphertext and perform erasure coding
	rscoded, err := cryptography.ReedSolomonEncode(ciphertext, threshold, shares)
	if err != nil {
		return
	}

	// split key material with shamir secret sharing
	keyshares, err := cryptography.ShamirSplit(key, threshold, shares)
	if err != nil {
		return
	}

	// basic sanity check
	if len(keyshares) != len(rscoded) {
		err = fmt.Errorf("keyshare array (%d) and data block array (%d) have different lengths",
			len(keyshares), len(rscoded))
		return
	}

	// assemble and sign shards
	shards = make([]*Shard, len(keyshares))
	for i := range shards {

		// assemble protobuf shard
		p := &ProtoShard{
			Associated: associated,
			Index:      int32(i),
			Keyshare:   keyshares[i],
			Data:       rscoded[i],
		}

		// marshal shard w/o signature and pubkey
		m, e := proto.Marshal(p)
		if e != nil {
			err = e
			return
		}

		// sign resulting protobuf
		sig, pub, e := cryptography.EdSign(key, m)
		if e != nil {
			err = e
			return
		}

		// amend shard with generated signature
		p.Pubkey, p.Signature = pub, sig

		// assemble high-level shard for pem construction
		shards[i] = &Shard{
			Description: description,
			Threshold:   threshold,
			UUID:        uuid,
			Proto:       p,
		}

	}
	return

}
