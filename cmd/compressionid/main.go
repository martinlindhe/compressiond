package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"

	"github.com/martinlindhe/compressionid"
)

var args struct {
	Filename string `kong:"arg" name:"filename" type:"existingfile" help:"Input file."`
	Offset   int64  `help:"Starting offset (default 0)."`
	OutFile  string `help:"Write extracted data to file." short:"o"`
}

func main() {
	compressionid.InitLogging()

	_ = kong.Parse(&args,
		kong.Name("compressionid"),
		kong.Description("A compression identifier."))

	f, err := os.Open(args.Filename)
	if err != nil {
		panic(err)
	}

	if args.Offset != 0 {
		f.Seek(args.Offset, os.SEEK_SET)
	}

	kind, v, err := compressionid.TryExtract(f)
	if err == nil {
		log.Info().Msgf("%s: %s compression detected", args.Filename, kind.String())

		if args.OutFile != "" {
			err = os.WriteFile(args.OutFile, v, 0644)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to write to %s", args.OutFile)
			}
		}
	} else {
		log.Error().Err(err).Msgf("Giving up")
	}
}
