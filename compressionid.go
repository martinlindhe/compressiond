package compressionid

import (
	"bytes"
	"compress/flate"
	"compress/lzw"
	"compress/zlib"
	"fmt"
	"io"
	"os"

	"github.com/fbonhomm/LZSS/source"
	lzss "github.com/fbonhomm/LZSS/source"
	"github.com/pierrec/lz4/v4"
	"github.com/rasky/go-lzo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ulikunitz/xz/lzma"

	lz77b "github.com/owencmiller/LZ77"
	lz77a "github.com/writingtoole/pdb/lz77"
)

func InitLogging() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
}

type CompressionKind int

const (
	ZLib  CompressionKind = iota
	Flate                 // (de)FLATE compression
	Lzma
	Lzma2
	Lzss
	Lzo1X
	Lz4
	Lz77a    // Lz77, variant #a
	Lz77b    // Lz77, variant #b
	Lzw_LSB8 // LSB, 8-bit
	Lzw_MSB8 // MSB, 8-bit
)

func (k CompressionKind) String() string {
	switch k {
	case ZLib:
		return "ZLib"
	case Flate:
		return "Flate"
	case Lzma:
		return "Lzma"
	case Lzma2:
		return "Lzma2"
	case Lzss:
		return "Lzss"
	case Lzo1X:
		return "Lzo1x"
	case Lz4:
		return "Lz4"
	case Lz77a:
		return "Lz77 (a)"
	case Lz77b:
		return "Lz77 (b)"
	case Lzw_LSB8:
		return "Lzw, LSB, 8 bit"
	case Lzw_MSB8:
		return "Lzw, MSB, 8 bit"
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
	log.Error().Err(err).Msgf("FLATE extraction failed")

	// LZO1X
	expanded, err := lzo.Decompress1X(bytes.NewReader(data), 0, 0)
	if err == nil {
		return Lzo1X, expanded, nil
	}
	log.Error().Err(err).Msgf("Lzo extraction failed")

	// LZ4
	lz4Dec := lz4.NewReader(bytes.NewReader(data))
	_, err = io.Copy(&b, lz4Dec)
	if err == nil {
		return Lz4, b.Bytes(), nil
	}
	log.Error().Err(err).Msgf("Lz4 extraction failed")

	// LZMA
	lzmaDec, err := lzma.NewReader(bytes.NewReader(data))
	if err == nil {
		_, err = io.Copy(&b, lzmaDec)
		if err == nil {
			return Lzma, b.Bytes(), nil
		}
	}
	log.Error().Err(err).Msgf("Lzma extraction failed")

	// LZMA2
	lzma2Dec, err := lzma.NewReader2(bytes.NewReader(data))
	if err == nil {
		_, err = io.Copy(&b, lzma2Dec)
		if err == nil {
			return Lzma2, b.Bytes(), nil
		}
	}
	log.Error().Err(err).Msgf("Lzma2 extraction failed")

	// LZSS
	lzssMode0 := source.LZSS{}
	expanded, err = lzssMode0.Decompress(data)
	if err == nil {
		return Lzss, expanded, nil
	}
	log.Error().Err(err).Msgf("Lzss extraction failed")

	// Lz77 (a)
	expanded, err = lz77a.Decompress(data)
	if err == nil {
		return Lz77a, expanded, nil
	}
	log.Error().Err(err).Msgf("Lz77 (a) extraction failed")

	// Lz77 (b)
	expanded, err = lz77b.Decompress(data)
	if err == nil {
		return Lz77b, expanded, nil
	}
	log.Error().Err(err).Msgf("Lz77 (b) extraction failed")

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
			return Lzw_LSB8, output[:count], nil
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
			return Lzw_MSB8, output[:count], nil
		}
	}
	log.Error().Err(err).Msgf("LZW-MSB-8 extraction failed")

	return 0, nil, fmt.Errorf("no compression recognized")
}

func CompressFromReader(method string, r io.Reader) ([]byte, error) {
	var b bytes.Buffer

	switch method {
	case "flate":
		w, err := flate.NewWriter(&b, flate.DefaultCompression)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(w, r)
		if err != nil {
			return nil, err
		}
		w.Close()

	case "zlib":
		w := zlib.NewWriter(&b)
		_, err := io.Copy(w, r)
		if err != nil {
			return nil, err
		}
		w.Close()

	case "lz4":
		w := lz4.NewWriter(&b)
		_, err := io.Copy(w, r)
		if err != nil {
			return nil, err
		}
		w.Close()

	case "lz77a":
		data, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		compressed, err := lz77a.Compress(data)
		if err != nil {
			return nil, err
		}
		b.Write(compressed)

	case "lz77b":
		data, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		compressed := lz77b.Compress(data)
		b.Write(compressed)

	case "lzma":
		w, err := lzma.NewWriter(&b)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(w, r)
		if err != nil {
			return nil, err
		}
		w.Close()

	case "lzma2":
		w, err := lzma.NewWriter2(&b)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(w, r)
		if err != nil {
			return nil, err
		}
		w.Close()

	case "lzss":
		data, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		LzssMode0 := lzss.LZSS{}
		b.Write(LzssMode0.Compress(data))

	case "lzo1x":
		data, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		b.Write(lzo.Compress1X(data))

	case "lzw-lsb8":
		w := lzw.NewWriter(&b, lzw.LSB, 8)
		_, err := io.Copy(w, r)
		if err != nil {
			return nil, err
		}
		w.Close()

	case "lzw-msb8":
		w := lzw.NewWriter(&b, lzw.MSB, 8)
		_, err := io.Copy(w, r)
		if err != nil {
			return nil, err
		}
		w.Close()
	}

	return b.Bytes(), nil
}
