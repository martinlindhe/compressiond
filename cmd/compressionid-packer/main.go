package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"compress/flate"
	"compress/lzw"
	"compress/zlib"

	"github.com/rasky/go-lzo"

	"github.com/alecthomas/kong"
	"github.com/pierrec/lz4/v4"
	"github.com/rs/zerolog/log"

	"github.com/martinlindhe/compressionid"
)

var args struct {
	Filename string `kong:"arg" name:"filename" type:"existingfile" help:"Input file."`
	Method   string `help:"Compression method ('deflate')." short:"m" required:""`
	OutFile  string `help:"Write compressed data to file." short:"o" required:""`
}

func main() {

	compressionid.InitLogging()

	_ = kong.Parse(&args,
		kong.Name("compressionid"),
		kong.Description("A compression identifier."))

	r, err := os.Open(args.Filename)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	var b bytes.Buffer

	switch args.Method {
	case "flate":
		w, err := flate.NewWriter(&b, flate.DefaultCompression)
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(w, r)
		if err != nil {
			panic(err)
		}
		w.Close()

	case "zlib":
		w := zlib.NewWriter(&b)
		_, err = io.Copy(w, r)
		if err != nil {
			panic(err)
		}
		w.Close()

	case "lz4":
		w := lz4.NewWriter(&b)
		_, err = io.Copy(w, r)
		if err != nil {
			panic(err)
		}
		w.Close()

	case "lzo1x":
		data, err := ioutil.ReadAll(r)
		if err != nil {
			panic(err)
		}
		b.Write(lzo.Compress1X(data))

	case "lzw-lsb8":
		w := lzw.NewWriter(&b, lzw.LSB, 8)
		_, err = io.Copy(w, r)
		if err != nil {
			panic(err)
		}
		w.Close()

	case "lzw-msb8":
		w := lzw.NewWriter(&b, lzw.MSB, 8)
		_, err = io.Copy(w, r)
		if err != nil {
			panic(err)
		}
		w.Close()

	default:
		log.Error().Msgf("Unrecognized compression method '%s'", args.Method)
	}

	err = os.WriteFile(args.OutFile, b.Bytes(), 0644)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to write to %s", args.OutFile)
	}
}
