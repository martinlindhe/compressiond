module github.com/martinlindhe/compressionid

go 1.19

require (
	github.com/alecthomas/kong v0.7.0
	github.com/owencmiller/LZ77 v0.0.0-20220118204303-02b62518c89a
	github.com/pierrec/lz4/v4 v4.1.17
	github.com/rasky/go-lzo v0.0.0-20200203143853-96a758eda86e
	github.com/rs/zerolog v1.28.0
)

require (
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
)

replace github.com/owencmiller/LZ77 v0.0.0-20220118204303-02b62518c89a => github.com/martinlindhe/LZ77 v0.0.0-20221025142603-82b6c92a4246
