module github.com/martinlindhe/compressionid

go 1.19

require (
	github.com/alecthomas/kong v0.7.0
	github.com/fbonhomm/LZSS v0.0.0-20200907090355-ba1a01a92989
	github.com/owencmiller/LZ77 v0.0.0-20220118204303-02b62518c89a
	github.com/pierrec/lz4/v4 v4.1.17
	github.com/rasky/go-lzo v0.0.0-20200203143853-96a758eda86e
	github.com/rs/zerolog v1.28.0
	github.com/ulikunitz/xz v0.5.10
	github.com/writingtoole/pdb v0.0.0-20190310153406-4473c8eabb5e
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	golang.org/x/sys v0.1.0 // indirect
)

replace (
	github.com/owencmiller/LZ77 v0.0.0-20220118204303-02b62518c89a => ../LZ77 // github.com/martinlindhe/LZ77 v0.0.0-20221025182440-d2eec1a23269

	github.com/writingtoole/pdb v0.0.0-20190310153406-4473c8eabb5e => ../pdb // github.com/martinlindhe/pdb

	github.com/fbonhomm/LZSS v0.0.0-20200907090355-ba1a01a92989 => ../LZSS // github.com/martinlindhe/LZSS
)
