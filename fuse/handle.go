package fuse

/*
#include <stdlib.h>
*/
import "C"

import (
	"sync"
	"unsafe"
)

var (
	guard = sync.Mutex{}
	table = map[unsafe.Pointer]interface{}{}
)

func getPointerForInterface(i interface{}) unsafe.Pointer {
	if i == nil {
		return nil
	}
	p := C.malloc(1)
	guard.Lock()
	defer guard.Unlock()
	table[p] = i
	return p
}

func getInterfaceForPointer(p unsafe.Pointer) interface{} {
	guard.Lock()
	defer guard.Unlock()
	if i, ok := table[p]; ok {
		return i
	}
	return nil
}

func delInterfaceFromPointer(p unsafe.Pointer) interface{} {
	guard.Lock()
	defer guard.Unlock()
	if i, ok := table[p]; ok {
		delete(table, p)
		C.free(p)
		return i
	}
	return nil
}
