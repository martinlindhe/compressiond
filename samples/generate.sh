# cargo install lzfoo
lzfoo -encode -i input.txt -o encoded.lzfse

# XXX generate zlib file

# XXX generate DEFLATE. need custom writer it seems (or gzip file minus headers)
go run cmd/compressionid-packer/main.go -m flate samples/input.txt -o samples/encoded.flate


# XXX generate LZO1X

# XXX generate LZ4

# XXX generate LZW (LSB, 8-bit)

