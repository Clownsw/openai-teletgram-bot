package util

import "unsafe"

func StringToByteSlice(str string) []byte {
	return *(*[]byte)(unsafe.Pointer(&str))
}

func ByteSliceToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
