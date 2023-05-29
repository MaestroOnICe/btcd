btcd
====
btcd is an alternative full node bitcoin implementation written in Go (golang).

This is a [SCION](https://github.com/scionproto/scion) compatible version of the [btcd](https://github.com/btcsuite/btcd) implementation.

It is *NOT compatible*, I repeat *NOT compatible*, with the official mainnet of the Bitcoin Network.

This client implementation is used to experiment within the [SEED Emulator](https://github.com/seed-labs/seed-emulator).

## What changed have been done?
- SCION compatiblity is enabled (scion.go inf scion dir)
  - dial and listen is done with [pan](https://github.com/netsec-ethz/scion-apps)
- CPU mining is enabled on our own "mainnet"
  - for this to work some of the mainnet parameters like the genisis block or checkpoint are deleted

## Requirements

[Go](http://golang.org) 1.16 or newer.

## License

btcd is licensed under the [copyfree](http://copyfree.org) ISC License.
