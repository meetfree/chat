package util

import (
	"bytes"
	"unsafe"
)

func BytesCombine(b ...[]byte) []byte {
	return bytes.Join(b, nil)
}

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

/**
 * 通过unsafe.Pointer伪造String的过程没有发生内存拷贝
 * 所以效率上会比发生内存拷贝的类型转换快
 * 但代价就是把底层数据暴露出来，这种做法是不安全的。
 */
func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
