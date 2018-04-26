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
	pemtype = "SHARDED MESSAGE"
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
	headers["Threshold"] = strconv.Itoa(s.Threshold)
	headers["UUID"] = fmt.Sprintf("%s", s.UUID)
	if len(s.Description) > 0 {
		headers["Description"] = s.Description
	}

	// marshal underlying protobuf
	bytes, err := proto.Marshal(s.Proto)
	if err != nil {
		return
	}

	// marshal pem to byteslice
	armor = pem.EncodeToMemory(&pem.Block{
		Type:    pemtype,
		Headers: headers,
		Bytes:   bytes,
	})
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
