package main

import (
	"encoding/pem"
	"fmt"
	"strconv"

	proto "github.com/golang/protobuf/proto"
	"github.com/google/uuid"
)

// Shard is a high-level struct for construction of PEM files.
// For binary files marshalling a ProtoShard is sufficient.
type Shard struct {
	Description string
	Threshold   int
	UUID        uuid.UUID
	Proto       *ProtoShard
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

func (s *Shard) Inspect() {

	fmt.Println(r("Shard "+fmt.Sprint(s.UUID)+":"), "index", s.Proto.Index)
	fmt.Println(y(" Threshold :"), s.Proto.Associated.Threshold)
	fmt.Println(y(" Shares    :"), s.Proto.Associated.Shares)
	fmt.Println(g(" Keyshare  :"), encode(s.Proto.Keyshare))
	fmt.Println(g(" Pubkey    :"), encode(s.Proto.Pubkey))
	fmt.Println(g(" Signature :"), encode(s.Proto.Signature))
	fmt.Println(b(" Data      :"), encode(s.Proto.Data))

}

const pemtype = "SHAMIR SHARED MESSAGE"
