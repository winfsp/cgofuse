/*
 * trace.go
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

package shared

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	TracePattern = os.Getenv("CGOFUSE_TRACE")
)

func traceJoin(deref bool, vals []interface{}) string {
	rslt := ""
	for _, v := range vals {
		if deref {
			switch i := v.(type) {
			case *bool:
				rslt += fmt.Sprintf(", %#v", *i)
			case *int:
				rslt += fmt.Sprintf(", %#v", *i)
			case *int8:
				rslt += fmt.Sprintf(", %#v", *i)
			case *int16:
				rslt += fmt.Sprintf(", %#v", *i)
			case *int32:
				rslt += fmt.Sprintf(", %#v", *i)
			case *int64:
				rslt += fmt.Sprintf(", %#v", *i)
			case *uint:
				rslt += fmt.Sprintf(", %#v", *i)
			case *uint8:
				rslt += fmt.Sprintf(", %#v", *i)
			case *uint16:
				rslt += fmt.Sprintf(", %#v", *i)
			case *uint32:
				rslt += fmt.Sprintf(", %#v", *i)
			case *uint64:
				rslt += fmt.Sprintf(", %#v", *i)
			case *uintptr:
				rslt += fmt.Sprintf(", %#v", *i)
			case *float32:
				rslt += fmt.Sprintf(", %#v", *i)
			case *float64:
				rslt += fmt.Sprintf(", %#v", *i)
			case *complex64:
				rslt += fmt.Sprintf(", %#v", *i)
			case *complex128:
				rslt += fmt.Sprintf(", %#v", *i)
			case *string:
				rslt += fmt.Sprintf(", %#v", *i)
			default:
				rslt += fmt.Sprintf(", %#v", v)
			}
		} else {
			rslt += fmt.Sprintf(", %#v", v)
		}
	}
	if len(rslt) > 0 {
		rslt = rslt[2:]
	}
	return rslt
}

func Trace(skip int, prfx string, vals ...interface{}) func(vals ...interface{}) {
	if "" == TracePattern {
		return func(vals ...interface{}) {
		}
	}
	pc, _, _, ok := runtime.Caller(skip + 1)
	name := "<UNKNOWN>"
	if ok {
		fn := runtime.FuncForPC(pc)
		name = fn.Name()
		if m, _ := filepath.Match(TracePattern, name); !m {
			return func(vals ...interface{}) {
			}
		}
	}
	if "" != prfx {
		prfx = prfx + ": "
	}
	args := traceJoin(false, vals)
	return func(vals ...interface{}) {
		f := "%v%v(%v) = (%v)"
		if len(vals) == 1 {
			f = "%v%v(%v) = %v"
		}
		log.Printf(f, prfx, name, args, traceJoin(true, vals))
	}
}
