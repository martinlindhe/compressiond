package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"

	"github.com/martinlindhe/compressionid"
)

var args struct {
	Filename string `kong:"arg" name:"filename" type:"existingfile" help:"Input file."`
	Method   string `help:"Compression method." enum:"flate,zlib,lz4,lz77a,lz77b,lzo1x,lzw-lsb8,lzw-msb8" short:"m" required:""`
	OutFile  string `help:"Write compressed data to file." short:"o" required:""`
}

func main() {
	compressionid.InitLogging()

	_ = kong.Parse(&args,
		kong.Name("compressionid"),
		kong.Description("A compression identifier."))

	r, err := os.Open(args.Filename)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to open file")
		return
	}
	defer r.Close()

	b, err := compressionid.CompressFromReader(args.Method, r)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to compress data")
		return
	}

	err = os.WriteFile(args.OutFile, b, 0644)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to write to %s", args.OutFile)
	}
}
