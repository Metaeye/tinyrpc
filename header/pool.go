package header

import "sync"

var RequestPool = sync.Pool{New: func() any { return &RequestHeader{} }}
var ResponsePool = sync.Pool{New: func() any { return &ResponseHeader{} }}
