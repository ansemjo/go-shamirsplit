package main

import (
	"errors"
	"fmt"
)

const (
	// key material
	demokey = "Zxky/LE10mbSdeT4Z3cPoJVcK5Vz3A/oRIR3DcUbgM8="
	// parameters
	threshold = 3
	total     = 5
	// demo data: use lipsum from lipsum.go
)

func main() {

	key, err := decode(demokey)
	fatal(err)

	shares, err := CreateKeyShares(key, threshold, total)
	fatal(err)
	//printarray("keyshares", shares)

	nonce, ciphertext, err := Encrypt(key, nil, lipsum)
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
		shards[i] = Shard{index: i, nonce: nonce, keyshare: shares[i], data: rscoded[i]}
	}

	for _, s := range shards {
		//fmt.Printf("· shard[%d] = %s\n", i, fmt.Sprintf("%v", s))
		fmt.Println(string(EncodePEM(s)))
	}

	// fmt.Println("delete shards")
	// splitdata[0] = nil
	// splitdata[1] = nil
	// printarray("splitdata", splitdata)
	// fmt.Println("recover data shards")
	// err = reed.ReconstructData(splitdata)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// printarray("splitdata", splitdata)

	// var recovered bytes.Buffer
	// reed.Join(&recovered, splitdata, K*len(splitdata[0]))
	// fmt.Println("original:", data)
	// fmt.Println("padded data:", paddeddata)
	// fmt.Println("recovered:", recovered.Bytes())

	// re := make([][]byte, 3)
	// re[0] = shares[3]
	// re[1] = shares[1]
	// re[2] = shares[4]
	// for i, r := range re {
	// 	fmt.Println(i, r[0], encode(r))
	// }

	// res, err := sss.CombineKeyshares(re)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("combined:", encode(res))

	// // hkdf testing
	// master := data
	// kdf, err := blake2b.New256(master)
	// fatal(err)

	// keys := make([][]byte, M)
	// for i := range keys {
	// 	kdf.Reset()
	// 	fmt.Println(i)
	// 	_, err = kdf.Write([]byte{byte(i)})
	// 	fatal(err)
	// 	keys[i] = kdf.Sum(nil)
	// }
	// printarray("HKDF", keys)

}
