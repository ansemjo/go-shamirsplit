package sharding

import (
	"fmt"

	"github.com/ansemjo/shamir/src/cryptography"
	"github.com/ansemjo/shamir/src/util"

	proto "github.com/golang/protobuf/proto"
)

func verifyShard(s *ProtoShard, pubkey []byte) (ok bool, err error) {

	// save signature
	sig := s.GetSignature()

	// marshal shard w/o signature and pubkey
	s.Pubkey, s.Signature = []byte{}, []byte{}
	m, err := proto.Marshal(s)
	if err != nil {
		return
	}

	// verify signature
	ok = cryptography.EdVerify(pubkey, m, sig)
	return

}

func CombineShards(shards []*ProtoShard) (data []byte, err error) {

	// first verify uuids
	// assume the first shard to be ours and use that public key
	mypub := shards[0].GetPubkey()

	// verify and collect valid shards
	verified := make([]*ProtoShard, 0, len(shards))
	for _, s := range shards {
		if ok, e := verifyShard(s, mypub); ok {
			fmt.Println("VERIFIED:", s.Index)
			verified = append(verified, s)
		} else if e != nil {
			err = e
			return
		}
	}

	// TODO: save our / one verified shard for easier access to e.g. associated data

	// collect reed-solomon encoded blocks in correct order
	rsdata := make([][]byte, verified[0].Associated.Shares)
	for _, v := range verified {
		rsdata[v.Index] = v.Data
	}
	ciphertext, err := cryptography.ReedSolomonReconstruct(rsdata, int(verified[0].Associated.Threshold), int(verified[0].Associated.Shares))
	if err != nil {
		return
	}

	fmt.Println(util.Base64encode(ciphertext))

	// combine shamir secret
	keyshares := make([][]byte, 0, len(verified))
	// TODO: combine collection into a single loop
	for _, k := range verified {
		keyshares = append(keyshares, k.Keyshare)
	}
	key, err := cryptography.ShamirCombine(keyshares)
	if err != nil {
		return
	}

	// decrypt
	ad, err := proto.Marshal(verified[0].Associated)
	if err != nil {
		return
	}
	message, err := cryptography.Decrypt(key, verified[0].Associated.Uuid[:12], ad, ciphertext)
	if err != nil {
		return
	}

	data = message
	fmt.Println(string(data))

	return

}
