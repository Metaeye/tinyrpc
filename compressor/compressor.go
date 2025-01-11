package compressor

import (
	"bytes"
	"io"
)

type Format uint16

const (
	Raw Format = iota
	Gzip
	Snappy
	Zlib
)

const NUM_FORMATS = 4

type compressFunc func([]byte) ([]byte, error)

type writer interface {
	Close() error
	Flush() error
	Write([]byte) (int, error)
}

type reader io.ReadCloser

func HasFormat(format Format) bool {
	return format < NUM_FORMATS
}

func zip(newWriter func(io.Writer) writer) compressFunc {
	return func(xs []byte) ([]byte, error) {
		buf := bytes.NewBuffer(nil)
		w := newWriter(buf)
		defer w.Close()

		if _, err := w.Write(xs); err != nil {
			return nil, err
		}

		if err := w.Flush(); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	}
}

func unzip(newReader func(io.Reader) (reader, error)) compressFunc {
	return func(xs []byte) ([]byte, error) {
		r, err := newReader(bytes.NewBuffer(xs))
		if err != nil {
			return nil, err
		}
		defer r.Close()

		ys, err := io.ReadAll(r)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return nil, err
		}

		return ys, nil
	}
}

var Zips = [NUM_FORMATS]compressFunc{
	func(xs []byte) ([]byte, error) {
		return xs, nil
	},
	zip(NewGzipWriter),
	zip(NewSnappyWriter),
	zip(NewZlibWriter),
}
var Unzips = [NUM_FORMATS]compressFunc{
	func(xs []byte) ([]byte, error) {
		return xs, nil
	},
	unzip(NewGzipReader),
	unzip(NewSnappyReader),
	unzip(NewZlibReader),
}
