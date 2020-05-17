package bui

//#include "webview.h"
import "C"
import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lxn/win"
	"github.com/mattn/go-pointer"
	"golang.org/x/sys/windows"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"unsafe"
)

var (
	windowViews = make(map[C.wkeWebView]*WebView)
)

type WebView struct {
	window     C.wkeWebView
	handle     win.HWND
	closing    bool
	echo       *echo.Echo
	fs         http.FileSystem
	rpc        *RPC
	port       int
	tray       *Tray
	wkeWndProc uintptr
}

type CreateViewOption struct {
	Title       string
	Width       int
	Height      int
	Transparent bool
	Fs          http.FileSystem
	Port        int
}

func CreateView(opt CreateViewOption) *WebView {
	v := &WebView{
		echo: echo.New(),
		rpc:  NewRPC(),
		port: opt.Port,
	}

	v.echo.HidePort = true
	v.echo.HideBanner = true

	v.echo.Use(middleware.CORS())

	v.echo.GET("/rpc", v.rpc.websocket)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", v.port))
	if err != nil {
		log.Fatalf("http file system listen error: %s", err.Error())
		return nil
	}
	v.port = l.Addr().(*net.TCPAddr).Port
	v.echo.Listener = l
	go func() {
		if err := v.echo.Start(""); err != nil {
			log.Fatalf("http file system listen error: %s", err.Error())
		}
	}()

	v.window = C.createWebWindow(C.int(opt.Width), C.int(opt.Height), C.bool(opt.Transparent))
	windowViews[v.window] = v
	v.handle = win.HWND(unsafe.Pointer(C.getWindowHandle(v.window)))

	v.wkeWndProc = win.SetWindowLongPtr(v.handle, win.GWLP_WNDPROC, windows.NewCallback(v.wndProc))
	v.tray = NewTray(v.handle)

	tempPath := filepath.Join(os.TempDir(), "bui")
	v.setLocalStorageFullPath(tempPath)
	v.setCookieJarFullPath(filepath.Join(tempPath, "data.dat"))

	v.SetFileSystem(opt.Fs)
	v.SetTitle(opt.Title)
	v.MoveToCenter()
	return v
}

func (v *WebView) SetFileSystem(fs http.FileSystem) {
	v.fs = fs
	if v.fs != nil {
		v.echo.GET("/*", echo.WrapHandler(http.FileServer(v.fs)))
		url := fmt.Sprintf("http://127.0.0.1:%d", v.port)
		log.Debugf("file system serves at %s", url)
		v.LoadUrl(url)
	}
}

func (v *WebView) SetIcon(icon win.HICON) {
	win.SendMessage(v.handle, win.WM_SETICON, 0, uintptr(icon))
	win.SendMessage(v.handle, win.WM_SETICON, 1, uintptr(icon))
}

func (v *WebView) Tray() *Tray {
	return v.tray
}

func (v *WebView) RPC() *RPC {
	return v.rpc
}

func (v *WebView) LoadUrl(url string) {
	C.loadURL(v.window, C.CString(url))
}

func (v *WebView) SetTitle(title string) {
	C.setWindowTitle(v.window, C.CString(title))
}

func (v *WebView) MoveToCenter() {
	C.moveToCenter(v.window)
}

func (v *WebView) Show() {
	C.showWindow(v.window, C.bool(true))
}

func (v *WebView) Hide() {
	C.showWindow(v.window, C.bool(false))
}

func (v *WebView) Minimize() {
	win.ShowWindow(v.handle, win.SW_MINIMIZE)
}

func (v *WebView) Maximize() {
	win.ShowWindow(v.handle, win.SW_MAXIMIZE)
}

func (v *WebView) Restore() {
	win.ShowWindow(v.handle, win.SW_RESTORE)
}

func (v *WebView) ShowOnTop() {
	win.SetForegroundWindow(v.handle)
	v.Show()
}

func (v *WebView) Close() {
	if !v.closing {
		v.closing = true
		v.tray.Dispose()
		C.destroyWindow(v.window)
	}
}

func (v *WebView) OnReady(callback func()) {
	C.onDocumentReady(v.window, pointer.Save(&callback))
}

func (v *WebView) OnDestroy(callback func()) {
	C.onWindowDestroy(v.window, pointer.Save(&callback))
}

func (v *WebView) onLoadUrlBegin(callback func(url string, job C.wkeNetJob)) {
	C.onLoadUrlBegin(v.window, pointer.Save(&callback))
}

func (v *WebView) onLoadUrlEnd(callback func(url string, job C.wkeNetJob, buf unsafe.Pointer, length int)) {
	C.onLoadUrlEnd(v.window, pointer.Save(&callback))
}

func (v *WebView) setLocalStorageFullPath(path string) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	C.setLocalStorageFullPath(v.window, cPath)
}

func (v *WebView) setCookieJarFullPath(path string) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	C.setCookieJarFullPath(v.window, cPath)
}

func (v *WebView) ShowDevTools(path string) {
	p, err := filepath.Abs(filepath.Join(path, "inspector.html"))
	if err != nil {
		log.Errorln(err)
		return
	}
	cPath := C.CString("file:///" + p)
	defer C.free(unsafe.Pointer(cPath))
	C.showDevtools(v.window, cPath)
}

func (v *WebView) wndProc(hWnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case TrayMsg:
		switch nMsg := win.LOWORD(uint32(lParam)); nMsg {
		case win.WM_LBUTTONDOWN:
			if v.tray != nil && v.tray.onLeftMouseClick != nil {
				v.tray.onLeftMouseClick()
			}
		case win.WM_RBUTTONDOWN:
			if v.tray != nil && v.tray.onRightMouseClick != nil {
				v.tray.onRightMouseClick()
			}
		}
	default:
		return win.CallWindowProc(v.wkeWndProc, hWnd, msg, wParam, lParam)
	}
	return 0
}

/*
	Golang Callbacks which called by C
*/

//export goGetBuiPort
func goGetBuiPort(window C.wkeWebView) C.int {
	if v, ok := windowViews[window]; ok {
		return C.int(v.port)
	}
	return 0
}

//export goOnDocumentReady
func goOnDocumentReady(window C.wkeWebView, param unsafe.Pointer) {
	if cb := pointer.Restore(param).(*func()); cb != nil {
		(*cb)()
	}
}

//export goOnWindowDestroy
func goOnWindowDestroy(window C.wkeWebView, param unsafe.Pointer) {
	if _, ok := windowViews[window]; ok {
		delete(windowViews, window)
	}
	if cb := pointer.Restore(param).(*func()); cb != nil {
		(*cb)()
		pointer.Unref(param)
	}
}

//export goOnLoadUrlBegin
func goOnLoadUrlBegin(window C.wkeWebView, param unsafe.Pointer, url *C.char, job C.wkeNetJob) {
	if cb := pointer.Restore(param).(*func(url string, job C.wkeNetJob)); cb != nil {
		(*cb)(C.GoString(url), job)
	}
}

//export goOnLoadUrlEnd
func goOnLoadUrlEnd(window C.wkeWebView, param unsafe.Pointer, url *C.char, job C.wkeNetJob, buf unsafe.Pointer, length C.int) {
	if cb := pointer.Restore(param).(*func(url string, job C.wkeNetJob, buf unsafe.Pointer, length int)); cb != nil {
		(*cb)(C.GoString(url), job, buf, int(length))
	}
}
