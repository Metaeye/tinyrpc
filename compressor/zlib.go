package compressor

import (
	"compress/zlib"
	"io"
)

func NewZlibWriter(w io.Writer) writer {
	return zlib.NewWriter(w)
}

func NewZlibReader(r io.Reader) (reader, error) {
	return zlib.NewReader(r)
}
