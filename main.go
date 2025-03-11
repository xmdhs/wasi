package main

import (
	"fmt"
	"unsafe"
)

var mem = make(map[int32][]byte)

//go:wasmexport malloc
func malloc(length int32) unsafe.Pointer {
	if length <= 0 {
		return nil
	}
	b := make([]byte, length)
	ptr := unsafe.Pointer(&b[0])
	mem[int32(uintptr(ptr))] = b
	return ptr
}

//go:wasmexport free
func free(ptr int32) {
	if ptr == 0 {
		return
	}
	// 从全局 map 中删除对应的内存引用
	delete(mem, int32(uintptr(ptr)))
}

//go:wasmexport print
func print(ptr int32, length int32) {
	b := mem[ptr][:length]
	fmt.Println(string(b))
}

//go:wasmexport fibonacci
func fibonacci(n int32) int32 {
	if n <= 1 {
		return n
	}
	a, b := 0, 1
	var i int32
	for i = 2; i <= n; i++ {
		a, b = b, a+b
	}
	return int32(b)
}

func main() {}
