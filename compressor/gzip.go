package compressor

import (
	"compress/gzip"
	"io"
)

func NewGzipWriter(w io.Writer) writer {
	return gzip.NewWriter(w)
}

func NewGzipReader(r io.Reader) (reader, error) {
	return gzip.NewReader(r)
}
