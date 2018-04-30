package sharding

import (
	"fmt"
	"os"

	"github.com/jinzhu/copier"

	"github.com/ansemjo/shamir/cryptography"

	proto "github.com/golang/protobuf/proto"
)

// CombineShards attempts to reconstruct the encoded and encrypted data
// from a slice of ProtoShards. The first shard in the slice is assumed
// to be 'ours' and some of its values are used to verify the rest of
// the shards.
func CombineShards(shards []*ProtoShard) (data []byte, err error) {

	// shards must not be empty
	if len(shards) == 0 {
		err = fmt.Errorf("shards slice cannot be empty")
		return
	}

	// assume the first shard to be our and use for later reference
	my := &ProtoShard{}
	err = copier.Copy(my, shards[0])
	if err != nil {
		return
	}
	threshold := int(my.Associated.Threshold)
	shares := int(my.Associated.Shares)
	uuid := my.Associated.Uuid

	// verify and collect valid shards
	verified := shards[:0]
	for i, s := range shards {
		if ok, e := verifyShardSignature(s, my.Pubkey); ok {
			verified = append(verified, s)
		} else if e != nil {
			err = e
			return
		} else {
			// TODO: only on verbose flag
			fmt.Fprintf(os.Stderr, "invalid signature on shard %d/%d\n", i+1, len(shards))
		}
	}

	// check number of shards
	if len(verified) < threshold {
		err = fmt.Errorf("not enough valid shards given (%d/%d)", len(verified), threshold)
		return
	}

	// collect reed-solomon blocks and keyshares
	rsdata := make([][]byte, shares)
	keyshares := make([][]byte, 0, shares)
	for _, v := range verified {
		// reed-solomon requires correct length and ordering
		rsdata[v.Index] = v.Data
		// while shamir shares may not be nil
		keyshares = append(keyshares, v.Keyshare)
	}

	// reconstruct ciphertext
	ciphertext, err := cryptography.ReedSolomonReconstruct(rsdata, threshold, shares)
	if err != nil {
		return
	}

	// recombine keyshares
	key, err := cryptography.ShamirCombine(keyshares)
	if err != nil {
		return
	}

	// serialize associated data
	ad, err := proto.Marshal(my.Associated)
	if err != nil {
		return
	}

	// authenticate and decrypt to cleartext
	data, err = cryptography.Decrypt(key, uuid[:12], ad, ciphertext)
	return

}

// ExtractProtoShards simply extracts ProtoShard structs to a new slice.
func ExtractProtoShards(shards []*Shard) (ps []*ProtoShard) {

	ps = make([]*ProtoShard, len(shards))
	for i, s := range shards {
		ps[i] = s.Proto
	}
	return

}

// verifySignature uses a given Ed25519 public key to verify the signature
// on a ProtoShard is correct. It does not modify the given shard.
func verifyShardSignature(s *ProtoShard, pubkey []byte) (ok bool, err error) {

	// copy to a temporary object
	newshard := &ProtoShard{}
	err = copier.Copy(newshard, s)
	if err != nil {
		return
	}

	// marshal newshard w/o signature and pubkey
	newshard.Pubkey, newshard.Signature = []byte{}, []byte{}
	serialized, err := proto.Marshal(newshard)
	if err != nil {
		return
	}

	// verify signature
	ok = cryptography.EdVerify(pubkey, serialized, s.Signature)
	return

}
