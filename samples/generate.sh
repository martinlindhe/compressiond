# cargo install lzfoo
lzfoo -encode -i input.txt -o encoded.lzfse

go run cmd/compressionid-packer/main.go -m zlib samples/input.txt -o samples/encoded.zlib

go run cmd/compressionid-packer/main.go -m flate samples/input.txt -o samples/encoded.flate


# XXX generate LZO1X

# XXX generate LZ4

# XXX generate LZW (LSB, 8-bit)

