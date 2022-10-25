package compressionid

import (
	"bytes"
	"compress/flate"
	"compress/lzw"
	"compress/zlib"
	"fmt"
	"io"
	"os"

	"github.com/pierrec/lz4/v4"
	"github.com/rasky/go-lzo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogging() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
}

type CompressionKind int

const (
	ZLib CompressionKind = iota
	Flate
	LZO1X
	LZ4
	LZW_LSB8 // LSB, 8-bit
	LZW_MSB8 // MSB, 8-bit
)

func (k CompressionKind) String() string {
	switch k {
	case ZLib:
		return "ZLib"
	case Flate:
		return "Flate"
	case LZO1X:
		return "LZO1x"
	case LZ4:
		return "LZ4"
	case LZW_LSB8:
		return "LZW, LSB, 8 bit"
	case LZW_MSB8:
		return "LZW, MSB, 8 bit"
	default:
		panic(k)
	}
}

func TryExtract(r io.Reader) (CompressionKind, []byte, error) {

	data, err := io.ReadAll(r)
	if err != nil {
		return 0, nil, err
	}

	var b bytes.Buffer
	// ZLIB
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err == nil {
		defer reader.Close()
		_, err = io.Copy(&b, reader)
		if err == nil {
			return ZLib, b.Bytes(), nil
		}
		log.Error().Err(err).Msgf("ZLIB extraction failed")
	} else {
		log.Error().Err(err).Msgf("ZLIB reading failed")
	}

	// FLATE
	flateDec := flate.NewReader(bytes.NewReader(data))
	defer flateDec.Close()
	_, err = io.Copy(&b, flateDec)
	if err == nil {
		return Flate, b.Bytes(), nil
	}
	log.Error().Err(err).Msgf("DEFLATE extraction failed")

	// LZO1X
	expanded, err := lzo.Decompress1X(bytes.NewReader(data), 0, 0)
	if err == nil {
		return LZO1X, expanded, nil
	}
	log.Error().Err(err).Msgf("LZO extraction failed")

	// LZ4
	lz4Dec := lz4.NewReader(bytes.NewReader(data))
	_, err = io.Copy(&b, lz4Dec)
	if err == nil {
		return LZ4, b.Bytes(), nil
	}
	log.Error().Err(err).Msgf("LZ4 extraction failed")

	// LZW-LSB-8
	lzwDec := lzw.NewReader(bytes.NewReader(data), lzw.LSB, 8)
	output := make([]byte, 1024*1024)
	count, err := lzwDec.Read(output)
	if err == nil {
		pct := (float64(count) / float64(len(data))) * 100
		if pct < 50 {
			// we maybe had some error
			log.Warn().Msgf("LZW-LSB-8 extracted %d of %d bytes (%.0f%%)", count, len(data), pct)

			fmt.Printf("output: %#v\n", string(output[:count]))
		} else {
			return LZW_LSB8, output[:count], nil
		}
	}
	log.Error().Err(err).Msgf("LZW-LSB-8 extraction failed")

	// LZW-MSB-8
	lzwDec = lzw.NewReader(bytes.NewReader(data), lzw.MSB, 8)
	output = make([]byte, 1024*1024)
	count, err = lzwDec.Read(output)
	if err == nil {
		pct := (float64(count) / float64(len(data))) * 100
		if pct < 50 {
			// we maybe had some error
			log.Warn().Msgf("LZW-MSB-8 extracted %d of %d bytes (%.0f%%)", count, len(data), pct)

			fmt.Printf("output: %#v\n", string(output[:count]))
		} else {
			return LZW_MSB8, output[:count], nil
		}
	}
	log.Error().Err(err).Msgf("LZW-MSB-8 extraction failed")

	// lzfse  - used by apple in xcode?

	return 0, nil, fmt.Errorf("no compression recognized")
}
