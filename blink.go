package bui

// #include "blink.h"
import "C"
import (
	"runtime"
	"unsafe"
)

var (
	wkeCall = make(chan func(), 1)
)

func WkeAsyncCall(f func()) {
	wkeCall <- f
}

func WkeSyncCall(f func()) {
	resolve := make(chan interface{}, 1)
	wkeCall <- func() {
		f()
		resolve <- true
	}
	<-resolve
}

func Initialize(path string) {
	runtime.LockOSThread()
	if path == "" {
		path = FindDLL()
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	C.wkeSetWkeDllPath(cPath)
	C.wkeInitialize()
	C.bindPort()
}
