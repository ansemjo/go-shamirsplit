# shamirsplit

This is a program to split secrets into a number of shares using Shamir secret
sharing.

## Description

Informally speaking, the input data is encrypted with an authenticated
encryption construction using a random key. The key is then split using Shamir
secret sharing and the ciphertext is split and encoded using Reed-Solomon
erasure coding. The resulting shares are then signed using a signature
algorithm on an elliptic curve with a key derived from the same initial random
key.

The combination of keyshare, data share and signature is then encoded in a PEM
formatted file, which also includes a unique UUID for easier identification of
shards that go together.

When creating shards there are two arguments, threshold and shares. The latter
is the total number of shards created and threshold is the minimum number of
shards needed to reconstruct the original data.

## Installation

TODO: simple go get command?

`make && sudo make install`

## Usage

See `shamirsplit -h`.

Generally there are two subcommands, `create` and `combine`. Data is piped into
the command via stdin and output is given on stdout.