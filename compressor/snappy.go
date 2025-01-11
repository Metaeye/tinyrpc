package compressor

import (
	"io"

	"github.com/golang/snappy"
)

func NewSnappyWriter(w io.Writer) writer {
	return snappy.NewBufferedWriter(w)
}

type snappyReader snappy.Reader

func (r *snappyReader) Close() error {
	return nil
}

func (r *snappyReader) Read(p []byte) (int, error) {
	return (*snappy.Reader)(r).Read(p)
}

func NewSnappyReader(r io.Reader) (reader, error) {
	return (*snappyReader)(snappy.NewReader(r)), nil
}
