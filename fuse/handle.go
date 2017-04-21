/*
 * handle.go
 *
 * Copyright 2017 Bill Zissimopoulos
 */
/*
 * This file is part of Cgofuse.
 *
 * You can redistribute it and/or modify it under the terms of the GNU
 * General Public License version 3 as published by the Free Software
 * Foundation.
 */

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

func newHandleForInterface(i interface{}) unsafe.Pointer {
	if nil == i {
		return nil
	}
	p := C.malloc(1)
	guard.Lock()
	defer guard.Unlock()
	table[p] = i
	return p
}

func delHandleForInterface(p unsafe.Pointer) interface{} {
	guard.Lock()
	defer guard.Unlock()
	if i, ok := table[p]; ok {
		delete(table, p)
		C.free(p)
		return i
	}
	return nil
}

func getInterfaceForHandle(p unsafe.Pointer) interface{} {
	guard.Lock()
	defer guard.Unlock()
	if i, ok := table[p]; ok {
		return i
	}
	return nil
}
