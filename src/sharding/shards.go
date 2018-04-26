package sharding

import (
	"encoding/pem"
	"errors"
	"fmt"
	"strconv"

	"github.com/ansemjo/shamir/src/cryptography"
	"github.com/ansemjo/shamir/src/util"
	proto "github.com/golang/protobuf/proto"
	"github.com/google/uuid"
)

// TODO: use random keys
const demokey = "Zxky/LE10mbSdeT4Z3cPoJVcK5Vz3A/oRIR3DcUbgM8="

// Shard is a high-level struct for construction of PEM files.
// For binary files marshalling a ProtoShard is sufficient.
type Shard struct {
	Description string
	Threshold   int
	UUID        uuid.UUID
	Proto       *ProtoShard
}

func CreateShards(threshold, shares int, message []byte, description string) (shards []*Shard, err error) {

	// message data
	uuid := util.RandomUUID()
	key, err := util.Base64decode(demokey) // TODO: use random key!
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

	ciphertext, err := cryptography.Encrypt(key, nonce, ad, message)
	//fmt.Println("Nonce:", encode(nonce))
	//fmt.Println("Ciphertext:", encode(ciphertext))

	// perform split and reed-solomon encoding
	rscoded, err := cryptography.ReedSolomonEncode(ciphertext, threshold, shares)
	if err != nil {
		return
	}

	// perform keysplit
	keyshares, err := cryptography.ShamirSplit(key, threshold, shares)
	if err != nil {
		return
	}

	// basic sanity check
	if len(keyshares) != len(rscoded) {
		err = errors.New("keyshares and reed-solomon encoded data have different lengths")
		return
	}

	// assemble and sign shards
	shards = make([]*Shard, len(keyshares))
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
		util.Fatal(err)

		// sign protobuf and amend shard
		sig, pub, err := cryptography.EdSign(key, m)
		util.Fatal(err)
		p.Pubkey = pub
		p.Signature = sig

		// finalize shard for pem construction
		shards[i] = &Shard{
			Description: description,
			Threshold:   threshold,
			UUID:        uuid,
			Proto:       p,
		}

	}
	return

}

func (s *Shard) toBlock() (block *pem.Block, err error) {

	headers := make(map[string]string)
	headers["Threshold"] = strconv.Itoa(s.Threshold)
	headers["UUID"] = fmt.Sprintf("%s", s.UUID)
	if len(s.Description) > 0 {
		headers["Description"] = s.Description
	}

	bytes, err := proto.Marshal(s.Proto)
	if err != nil {
		return
	}

	block = &pem.Block{Type: pemtype, Headers: headers, Bytes: bytes}
	return

}

// MarshalPEM marshals the contained protobuf and then marshals
// a byte representation of a PEM shard.
func (s *Shard) MarshalPEM() (p []byte, err error) {

	block, err := s.toBlock()
	if err != nil {
		return
	}

	p = pem.EncodeToMemory(block)
	return

}

// Inspect logs the internal structure to the console.
func (s *Shard) Inspect() {

	fmt.Println(util.R("Shard "+fmt.Sprint(s.UUID)+":"), "index", s.Proto.Index)
	fmt.Println(util.Y(" Threshold :"), s.Proto.Associated.Threshold)
	fmt.Println(util.Y(" Shares    :"), s.Proto.Associated.Shares)
	fmt.Println(util.G(" Keyshare  :"), util.Base64encode(s.Proto.Keyshare))
	fmt.Println(util.G(" Pubkey    :"), util.Base64encode(s.Proto.Pubkey))
	fmt.Println(util.G(" Signature :"), util.Base64encode(s.Proto.Signature))
	fmt.Println(util.B(" Data      :"), util.Base64encode(s.Proto.Data))

}

const pemtype = "SHARDED MESSAGE"
