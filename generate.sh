#!/bin/sh

go run cmd/compressionid-packer/main.go -m flate samples/input.txt -o samples/encoded.flate

go run cmd/compressionid-packer/main.go -m zlib samples/input.txt -o samples/encoded.zlib

go run cmd/compressionid-packer/main.go -m lz4 samples/input.txt -o samples/encoded.lz4

go run cmd/compressionid-packer/main.go -m lz77a samples/input.txt -o samples/encoded.lz77a
go run cmd/compressionid-packer/main.go -m lz77b samples/input.txt -o samples/encoded.lz77b

go run cmd/compressionid-packer/main.go -m lzo1x samples/input.txt -o samples/encoded.lzo1x

go run cmd/compressionid-packer/main.go -m lzw-lsb8 samples/input.txt -o samples/encoded.lzw-lsb8
go run cmd/compressionid-packer/main.go -m lzw-msb8 samples/input.txt -o samples/encoded.lzw-msb8


# cargo install lzfoo
lzfoo -encode -i samples/input.txt -o samples/encoded.lzfse
