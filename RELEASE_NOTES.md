# v0.1.0

Initial public release of `go-bitshares`.

Highlights:

- standalone BitShares Go module for client, protocol, signing, and ECC helpers
- typed wallet builders for the common user-signable operation families
- login, backup restore, memo encryption, and transaction signing support
- secret-handling hardening with byte-based secret inputs and explicit `Wipe()` helpers
- canonical binary set serialization aligned with BitShares core semantics

Notes:

- this is the first public release, so API and compatibility guarantees may still evolve
- virtual or chain-internal operations remain intentionally unwrapped as wallet builders
