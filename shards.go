package main

import (
	"encoding/pem"
	"fmt"
	"strconv"

	proto "github.com/golang/protobuf/proto"
	"github.com/google/uuid"
)

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

func EncodePEM(shard *Shard) (p []byte, err error) {

	block, err := shard.toBlock()
	if err != nil {
		return
	}

	p = pem.EncodeToMemory(block)
	return

}

const pemtype = "SHAMIR SHARED MESSAGE"
