/*
 * host_test_unix.go
 *
 * Copyright 2017 Bill Zissimopoulos
 */
/*
 * This file is part of Cgofuse.
 *
 * It is licensed under the MIT license. The full license text can be found
 * in the License.txt file at the root of this project.
 */

// +build windows

package fuse

import (
	"syscall"
)

func sendInterrupt() bool {
	dll := syscall.MustLoadDLL("kernel32")
	prc := dll.MustFindProc("GenerateConsoleCtrlEvent")
	r, _, _ := p.Call(syscall.CTRL_BREAK_EVENT, ^uintptr(0))
	return 0 != r
}
