# btcd

This is the mvp testing branch. SCION on this branch. This branch is used to test the btcd client in the SEED emulator with a SCION Connection.

## Changes

- connmanager debug messages disabled for visibility
  - lines 210-212, 296 and 350-351
- checkpoints for "mainnet" deletes because this an enclosed "new" mainnet
- CPU mining support enabled on mainnet
- uses the simnet genesis block
- no DNS seeding
- allowing RFC 1918 ipv4 address (e.g. 10.0.0.0/8) becaus this is the default addressing scheme in SEED
- AddressCache returns addresses without fisher-yates shuffle
  - this caused problems with a few participants in the network, returned an empty slice
- AddressCache is dumped to disk every minute

## SCION changes

-
