package sharding

import (
	"encoding/pem"
	"fmt"
	"strconv"

	"github.com/ansemjo/shamir/src/util"
	proto "github.com/golang/protobuf/proto"
	"github.com/google/uuid"
)

const (
	// expected begin and end string
	pemType    = "SHARDED MESSAGE"
	headerDesc = "Description"
	headerUUID = "UUID"
	headerThre = "Threshold"
)

// Shard is a high-level struct for construction of PEM files.
// For binary files marshalling a ProtoShard is sufficient.
type Shard struct {
	Description string
	Threshold   int
	UUID        uuid.UUID
	Proto       *ProtoShard
}

// MarshalPEM marshals the contained protobuf and then marshals
// a byte representation of a PEM armored shard.
func (s *Shard) MarshalPEM() (armor []byte, err error) {

	// prepare headers, omit description if none given
	// TODO: throw if threshold or uuid is missing
	headers := make(map[string]string)
	headers[headerThre] = strconv.Itoa(s.Threshold)
	headers[headerUUID] = fmt.Sprintf("%s", s.UUID)
	if len(s.Description) > 0 {
		headers[headerDesc] = s.Description
	}

	// marshal underlying protobuf
	bytes, err := proto.Marshal(s.Proto)
	if err != nil {
		return
	}

	// marshal pem to byteslice
	armor = pem.EncodeToMemory(&pem.Block{
		Type:    pemType,
		Headers: headers,
		Bytes:   bytes,
	})
	return

}

// UnmarshalPEM attempts to parse a Shard from the given byteslice. It
// returns the Shard and any remaining input in rest. If no PEM block
// was found, shard will be nil.
func UnmarshalPEM(armor []byte) (shard *Shard, rest []byte, err error) {

	// decode pem
	block, rest := pem.Decode(armor)
	if block == nil {
		//err = fmt.Errorf("no pem blocks found")
		return
	}
	if block.Type != pemType {
		err = fmt.Errorf("unexpected pem type: %s", block.Type)
		return
	}

	// unmarshal data content
	ps := &ProtoShard{}
	err = proto.Unmarshal(block.Bytes, ps)
	if err != nil {
		return
	}

	// type conversions
	t, err := strconv.Atoi(block.Headers[headerThre])
	if err != nil {
		return
	}
	u, err := uuid.Parse(block.Headers[headerUUID])
	if err != nil {
		return
	}

	// assemble struct
	shard = &Shard{Description: block.Headers[headerDesc], Threshold: t, UUID: u, Proto: ps}
	return

}

// ReadAll reads multiple PEM-armored shards from a given byteslice.
func ReadAll(input []byte) (shards []*Shard, err error) {

	shards = []*Shard{}
	for len(input) > 0 {
		if shard, rest, e := UnmarshalPEM(input); shard != nil {
			shards, input = append(shards, shard), rest
		} else if e != nil {
			err = e
			return
		} else {
			// both shard and err are nil => no more pem blocks
			return
		}
	}
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
