// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package header

import (
	"encoding/binary"
	"errors"
	"sync"

	"tinyrpc/compressor"
)

const (
	// MaxHeaderSize = 12 + 10 + 10 + 10 + 4 (10 refer to binary.MaxVarintLen64)
	MaxHeaderSize = 36

	Uint32Size = 4
	Uint16Size = 2
)

var ErrUnmarshal = errors.New("an error occurred in Unmarshal")

// RequestHeader request header structure looks like:
// +--------------+----------------+----------+------------+----------+
// |    Format    |      Method    |    ID    | RequestLen | Checksum |
// +--------------+----------------+----------+------------+----------+
// |    uint16    | uvarint+string |  uvarint |   uvarint  |  uint32  |
// +--------------+----------------+----------+------------+----------+
type RequestHeader struct {
	sync.RWMutex
	Format     compressor.Format
	Method     string
	ID         uint64
	RequestLen uint32
	Checksum   uint32
}

// Marshal will encode request header into a byte slice
func (r *RequestHeader) Marshal() []byte {
	r.RLock()
	defer r.RUnlock()
	idx := 0
	// MaxHeaderSize = 2 + 10 + len(string) + 10 + 10 + 4
	header := make([]byte, MaxHeaderSize+len(r.Method))

	binary.LittleEndian.PutUint16(header[idx:], uint16(r.Format))
	idx += Uint16Size

	idx += writeString(header[idx:], r.Method)
	idx += binary.PutUvarint(header[idx:], r.ID)
	idx += binary.PutUvarint(header[idx:], uint64(r.RequestLen))

	binary.LittleEndian.PutUint32(header[idx:], r.Checksum)
	idx += Uint32Size

	return header[:idx]
}

// Unmarshal will decode request header into a byte slice
func (r *RequestHeader) Unmarshal(data []byte) (err error) {
	r.Lock()
	defer r.Unlock()
	if len(data) == 0 {
		return ErrUnmarshal
	}

	defer func() {
		if r := recover(); r != nil {
			err = ErrUnmarshal
		}
	}()
	idx, size := 0, 0
	r.Format = compressor.Format(binary.LittleEndian.Uint16(data[idx:]))
	idx += Uint16Size

	r.Method, size = readString(data[idx:])
	idx += size

	r.ID, size = binary.Uvarint(data[idx:])
	idx += size

	length, size := binary.Uvarint(data[idx:])
	r.RequestLen = uint32(length)
	idx += size

	r.Checksum = binary.LittleEndian.Uint32(data[idx:])

	return
}

// get compress format
func (r *RequestHeader) GetFormat() compressor.Format {
	r.RLock()
	defer r.RUnlock()
	return r.Format
}

// reset request header
func (r *RequestHeader) ResetHeader() {
	r.Lock()
	defer r.Unlock()
	r.ID = 0
	r.Checksum = 0
	r.Method = ""
	r.Format = 0
	r.RequestLen = 0
}

// ResponseHeader request header structure looks like:
// +--------------+---------+----------------+-------------+----------+
// |    Format    |    ID   |      Error     | ResponseLen | Checksum |
// +--------------+---------+----------------+-------------+----------+
// |    uint16    | uvarint | uvarint+string |    uvarint  |  uint32  |
// +--------------+---------+----------------+-------------+----------+
type ResponseHeader struct {
	sync.RWMutex
	Format      compressor.Format
	ID          uint64
	Error       string
	ResponseLen uint32
	Checksum    uint32
}

// Marshal will encode response header into a byte slice
func (r *ResponseHeader) Marshal() []byte {
	r.RLock()
	defer r.RUnlock()
	idx := 0
	header := make([]byte, MaxHeaderSize+len(r.Error)) // prevent panic

	binary.LittleEndian.PutUint16(header[idx:], uint16(r.Format))
	idx += Uint16Size

	idx += binary.PutUvarint(header[idx:], r.ID)
	idx += writeString(header[idx:], r.Error)
	idx += binary.PutUvarint(header[idx:], uint64(r.ResponseLen))

	binary.LittleEndian.PutUint32(header[idx:], r.Checksum)
	idx += Uint32Size
	return header[:idx]
}

// Unmarshal will decode response header into a byte slice
func (r *ResponseHeader) Unmarshal(data []byte) (err error) {
	r.Lock()
	defer r.Unlock()
	if len(data) == 0 {
		return ErrUnmarshal
	}

	defer func() {
		if r := recover(); r != nil {
			err = ErrUnmarshal
		}
	}()
	idx, size := 0, 0
	r.Format = compressor.Format(binary.LittleEndian.Uint16(data[idx:]))
	idx += Uint16Size

	r.ID, size = binary.Uvarint(data[idx:])
	idx += size

	r.Error, size = readString(data[idx:])
	idx += size

	length, size := binary.Uvarint(data[idx:])
	r.ResponseLen = uint32(length)
	idx += size

	r.Checksum = binary.LittleEndian.Uint32(data[idx:])
	return
}

// get compress type
func (r *ResponseHeader) GetFormat() compressor.Format {
	r.RLock()
	defer r.RUnlock()
	return r.Format
}

// reset response header
func (r *ResponseHeader) ResetHeader() {
	r.Lock()
	defer r.Unlock()
	r.Error = ""
	r.ID = 0
	r.Format = 0
	r.Checksum = 0
	r.ResponseLen = 0
}

func readString(data []byte) (string, int) {
	idx := 0
	length, size := binary.Uvarint(data)
	idx += size
	str := string(data[idx : idx+int(length)])
	idx += len(str)

	return str, idx
}

func writeString(data []byte, str string) int {
	idx := 0
	idx += binary.PutUvarint(data, uint64(len(str)))
	copy(data[idx:], str)
	idx += len(str)

	return idx
}
