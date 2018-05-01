# shamirsplit

This is a commandline program to split secrets into a number of shares where a
certain threshold amount is necessary to reconstruct the original secret.

## Idea

After watching [Daan Sprenkels](https://dsprenkels.com/)' 34C3 talk
["We should share our secrets"](https://dsprenkels.com/sss-34c3.html) I was
fascinated how simple yet powerful Shamir secret sharing is. (If you need
technical details on Shamir secret sharing, go watch that talk and read the
documentation to the accompanying projects on [GitHub](https://github.com/dsprenkels/sss)).

However, I soon realized that using hexadeciaml encoding and appending the entire
ciphertext to every keyshare was wasteful and
[wanted to use](https://github.com/dsprenkels/sss-cli/issues/8) erasure coding
with the same threshold values to minimize the space used by shares.

Sometime in the future this tool may also support binary protobuf-serialized
files to further minimize memory usage.

## Description

Informally speaking, the input data is encrypted with an authenticated
encryption construction using a random key. The key is then split using Shamir
secret sharing and the ciphertext is split and encoded using Reed-Solomon
erasure coding. The resulting shares are then signed using a signature
algorithm on an elliptic curve with a key derived from the same initial random
key.

The combination of keyshare, data share and signature is then encoded in a PEM
formatted file, which also includes a unique UUID for easier identification of
shards that go together. (I call the combination of key- and data shares a
'shard'.)

When creating shards there are two main arguments, threshold and shares. The
latter is the total number of shards created and threshold is the minimum number
of shards needed to reconstruct the original data.

For an attempted visualization of the process and data flows see
[`specification.txt`](/specification.txt).

In the spirit of "don't roll your own crypto", I didn't implement any of the
primitives used within. I merely plugged them together and hope that I didn't
make any grave mistakes in doing so. Some of the pieces used within:

- [dspenkels/sss](https://github.com/dsprenkels/sss) - Shamir secret sharing
- [golang/chacha20poly1305](https://godoc.org/golang.org/x/crypto/chacha20poly1305) - Authenticated Encryption with ChaCha20Poly1305
- [golang/ed25519](https://godoc.org/golang.org/x/crypto/ed25519) - Elliptic curve signatures on Ed25519
- [klauspost/reedsolomon](https://github.com/klauspost/reedsolomon) - Reed-Solomon erasure coding
- [alexbakker/go-pkcs7](https://github.com/alexbakker/go-pkcs7) - PKCS7 padding
- [golang/protobuf](github.com/golang/protobuf) - Google protobuf for efficient serialization
- [golang/pem](https://godoc.org/encoding/pem) - PEM encoding

## Installation

* use `go get github.com/ansemjo/go-shamirsplit/cmd/shamirsplit`
* or clone the repository and run `make && sudo make install`

## Usage

Most up-to-date usage information should always be available with
`shamirsplit --help`. Generally there are two subcommands, `create` and
`combine` and usage data for these subcommands is available with
`shamirsplit <subcommand> --help`.

### Creation

In its simplest form, data is piped into the command and the shards are output
on the console:

```bash
$ echo 'Hello, World!' | shamirsplit create -t 2 -s 3
-----BEGIN SHARDED MESSAGE-----
Threshold: 2
UUID: d6204844-90ef-4f5c-8e52-f3d9f071c89f

ChYKENYgSESQ709cjlLz2fBxyJ8QAhgDGiEBBCYpTker/wuy0QKYvSB9KZRtwWxu
A9W5bEIohzFVoXoiINmBXEWdOJdQ3OjZGuujkWS8mRAtxXbIL+s1ZJ0YO3/xKkBj
Gr9hIztF665uKdnlGiFUGJm1ithjp9agwPhALcEy4MXhSxaXwxZwBtK319j6XGOm
A3EQNjF/yNLbN8fyItsOMhBNeB6VxuaX5AhX2V4hiXSZ
-----END SHARDED MESSAGE-----
-----BEGIN SHARDED MESSAGE-----
Threshold: 2
UUID: d6204844-90ef-4f5c-8e52-f3d9f071c89f

ChYKENYgSESQ709cjlLz2fBxyJ8QAhgDEAEaIQJ2hiK0vaAHvZqv9Oe31BfWh/3V
ptSYd/sEUS7reHXLiiIg2YFcRZ04l1Dc6Nka66ORZLyZEC3Fdsgv6zVknRg7f/Eq
QLLaIpMWENe6U4OyamI+3ywJ4pAZ5fIVxVF2Dn1KI4Hd+dj45Rvg0JkuQyn5w+19
1XhlX7aVY/xP82/AGZMxeAIyEMl49vZR53jMP8Vu4z+RAgI=
-----END SHARDED MESSAGE-----
-----BEGIN SHARDED MESSAGE-----
Threshold: 2
UUID: d6204844-90ef-4f5c-8e52-f3d9f071c89f

ChYKENYgSESQ709cjlLz2fBxyJ8QAhgDEAIaIQNY5tLi61CmJoKFpjuxcTGDf43Z
4EsY4MXVqSzPtpzt2iIg2YFcRZ04l1Dc6Nka66ORZLyZEC3Fdsgv6zVknRg7f/Eq
QOG1AV6ie5rfTFEXtIgIAtizHslpdVJSNyBE7QqQXXnXHOyhaRfrYJQXk0bF5Q8Y
I8FXtOtno8Jf7XBAqxjhgw4yEFh401P15FS0Zm6qOR25mLI=
-----END SHARDED MESSAGE-----
```

To make splitting easier, a null-byte can be appended after each PEM block with
`-0/--null`. While-read loops or GNU split can then delimit on that null-byte:

```bash
$ echo 'Hello, World!' | shamirsplit create -t 2 -s 3 -0 |\
 while read -d $'\0' shard; do echo "$shard" > shard.$((i++)); done
```

```bash
$ echo 'Hello, World!' | shamirsplit create -t 2 -s 3 -0 | split -l1 -t'\0' -d - "shard."
```

_Note however_ that you need to remove the null-bytes from the output files
manually in case you use `split`. Otherwise `shamirsplit combine` will only read
the first shard and then complain about missing shards.

Alternatively use the `-d/--directory` and `-n/--name` arguments to write the
shards to files in a directory directly:

```bash
$ echo 'Hello, World!' | shamirsplit create -t 2 -s 3 -d /tmp -n testingshard
/tmp/testingshard_000.pem
/tmp/testingshard_001.pem
/tmp/testingshard_002.pem
```

### Combination

To combine shards at a later point in time, concatenate and pipe all the shards
that you have in any order. If you have more valid shards than the threshold,
the original data will be reconstructed.

___Note:___ upon verification the first shard given is assumed to be "your" shard
and its embedded public key will be used to verify the signatures on all shards.
Since this provides integrity you should take care to always pipe your own shard
first.

```bash
$ rm /tmp/testingshard_001.pem 
rm: remove regular file '/tmp/testingshard_001.pem'? y
$ cat /tmp/testingshard_00{2,0}.pem | shamirsplit combine
Hello, World!
```