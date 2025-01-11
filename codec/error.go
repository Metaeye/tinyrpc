// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codec

import "errors"

var (
	ErrInvalidSequence        = errors.New("invalid sequence number in response")
	ErrUnexpectedChecksum     = errors.New("unexpected checksum")
	ErrNotFoundCompressFormat = errors.New("not found compressor")
	ErrCompressorTypeMismatch = errors.New("request and response Compressor type mismatch")
)
