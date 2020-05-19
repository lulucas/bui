package bui

// #include "blink.h"
import "C"
import (
	"runtime"
	"unsafe"
)

var (
	uiCall = make(chan func(), 1)
)

func AsyncCall(f func()) {
	uiCall <- f
}

func SyncCall(f func()) {
	resolve := make(chan interface{}, 1)
	uiCall <- func() {
		f()
		resolve <- true
	}
	<-resolve
}

func Initialize() {
	InitializeByDllPath("", "")
}

func InitializeByDllPath(nodeDll string, mbDll string) {
	runtime.LockOSThread()

	if nodeDll == "" {
		nodeDll = FindNodeDLL()
	}

	cNodePath := C.CString(nodeDll)
	C.mbSetMbMainDllPath(cNodePath)
	C.free(unsafe.Pointer(cNodePath))

	if mbDll == "" {
		mbDll = FindMbDLL()
	}
	cMbPath := C.CString(mbDll)
	C.mbSetMbDllPath(cMbPath)
	C.free(unsafe.Pointer(cMbPath))

	C.mbInitialize()
}

func Finalize() {
	C.mbFinalize()
}
