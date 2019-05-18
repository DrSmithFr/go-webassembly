# Golang - Web Assembly

Example of web assembly stateless API (without DOM modification).

All APIs get wrapper within a javascript object to allow embedded wasm on modern Javascript frameworks.

## rules followed by API:

- MUST NOT modify DOM
- MUST be stateless
- MUST return pure Javascript object
- MUST use Go Object within logic
- SHOULD use gpu based capabilities when available
