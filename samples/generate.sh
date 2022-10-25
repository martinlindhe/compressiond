# cargo install lzfoo
lzfoo -encode -i input.txt -o encoded.lzfse

go run cmd/compressionid-packer/main.go -m zlib samples/input.txt -o samples/encoded.zlib

go run cmd/compressionid-packer/main.go -m flate samples/input.txt -o samples/encoded.flate

# lz4 samples/input.txt samples/encoded.lz4
go run cmd/compressionid-packer/main.go -m lz4 samples/input.txt -o samples/encoded.lz4

go run cmd/compressionid-packer/main.go -m lzo1x samples/input.txt -o samples/encoded.lzo1x

# XXX generate LZW (LSB, 8-bit)

