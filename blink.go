package bui

// #include "blink.h"
import "C"
import (
	"os"
	"path/filepath"
	"runtime"
	"unsafe"
)

const (
	TempDir = "bui"
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
	_ = os.MkdirAll(filepath.Join(os.TempDir(), TempDir), 0644)

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
