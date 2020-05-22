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
	windowViews = make(map[C.mbWebView]*WebView)
)

type WebView struct {
	webView    C.mbWebView
	window     win.HWND
	closing    bool
	echo       *echo.Echo
	fs         http.FileSystem
	rpc        *RPC
	port       int
	tray       *Tray
	wndProcPtr uintptr
	modal      *WebView
	isModal    bool
}

type CreateViewOption struct {
	Title       string
	Width       int
	Height      int
	Transparent bool
	Fs          http.FileSystem
	Port        int
	IsModal     bool
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

	v.webView = C.createWebWindow(C.int(opt.Width), C.int(opt.Height), C.bool(opt.Transparent))
	windowViews[v.webView] = v
	v.window = win.HWND(unsafe.Pointer(C.getWindowHandle(v.webView)))

	v.wndProcPtr = win.SetWindowLongPtr(v.window, win.GWLP_WNDPROC, windows.NewCallback(v.wndProc))

	tempPath := filepath.Join(os.TempDir(), TempDir)
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
	win.SendMessage(v.window, win.WM_SETICON, 0, uintptr(icon))
	win.SendMessage(v.window, win.WM_SETICON, 1, uintptr(icon))
}

func (v *WebView) Tray() *Tray {
	if v.tray == nil {
		v.tray = NewTray(v.window)
	}
	return v.tray
}

func (v *WebView) RPC() *RPC {
	return v.rpc
}

func (v *WebView) LoadUrl(url string) {
	C.loadURL(v.webView, C.CString(url))
}

func (v *WebView) SetTitle(title string) {
	C.setWindowTitle(v.webView, C.CString(title))
}

func (v *WebView) MoveToCenter() {
	C.moveToCenter(v.webView)
}

func (v *WebView) Show() {
	C.showWindow(v.webView, C.bool(true))
}

func (v *WebView) Hide() {
	C.showWindow(v.webView, C.bool(false))
}

func (v *WebView) Minimize() {
	win.ShowWindow(v.window, win.SW_MINIMIZE)
}

func (v *WebView) Maximize() {
	win.ShowWindow(v.window, win.SW_MAXIMIZE)
}

func (v *WebView) Restore() {
	win.ShowWindow(v.window, win.SW_RESTORE)
}

func (v *WebView) ToggleMaximize() {
	if win.IsZoomed(v.window) {
		v.Restore()
	} else {
		v.Maximize()
	}
}

func (v *WebView) ShowOnTop() {
	win.SetForegroundWindow(v.window)
	v.Show()
}

func (v *WebView) Close() {
	if !v.closing {
		v.closing = true
		if v.tray != nil {
			v.tray.Dispose()
		}
		C.destroyWindow(v.webView)
		win.DestroyWindow(v.window)
	}
}

func (v *WebView) Enable() {
	win.EnableWindow(v.window, true)
}

func (v *WebView) Disable() {
	win.EnableWindow(v.window, false)
}

func (v *WebView) ShowModal(width, height int, url string) {
	modal := CreateView(CreateViewOption{
		Width:       width,
		Height:      height,
		Transparent: true,
	})
	win.EnableWindow(v.window, false)
	modal.LoadUrl(url)
	modal.Show()
	modal.OnDestroy(func() {
		win.EnableWindow(v.window, true)
	})
	v.modal = modal
}

func (v *WebView) CloseModal() {
	if v.modal != nil {
		v.modal.Close()
		v.modal = nil
	}
}

func (v *WebView) OnReady(callback func()) {
	C.onDocumentReady(v.webView, pointer.Save(&callback))
}

func (v *WebView) OnDestroy(callback func()) {
	C.onWindowDestroy(v.webView, pointer.Save(&callback))
}

func (v *WebView) onLoadUrlBegin(callback func(url string, job C.mbNetJob)) {
	C.onLoadUrlBegin(v.webView, pointer.Save(&callback))
}

func (v *WebView) onLoadUrlEnd(callback func(url string, job C.mbNetJob, buf unsafe.Pointer, length int)) {
	C.onLoadUrlEnd(v.webView, pointer.Save(&callback))
}

func (v *WebView) setLocalStorageFullPath(path string) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	C.setLocalStorageFullPath(v.webView, cPath)
}

func (v *WebView) setCookieJarFullPath(path string) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	C.setCookieJarFullPath(v.webView, cPath)
}

func (v *WebView) ShowDevTools(path string) {
	p, err := filepath.Abs(filepath.Join(path, "inspector.html"))
	if err != nil {
		log.Errorln(err)
		return
	}
	cPath := C.CString("file:///" + p)
	defer C.free(unsafe.Pointer(cPath))
	C.showDevtools(v.webView, cPath)
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
		return win.CallWindowProc(v.wndProcPtr, hWnd, msg, wParam, lParam)
	}
	return 0
}

/*
	Golang Callbacks which called by C
*/

//export goOnDocumentReady
func goOnDocumentReady(webView C.mbWebView, param, frameId unsafe.Pointer) {
	if cb := pointer.Restore(param).(*func()); cb != nil && *cb != nil {
		(*cb)()
	}
}

//export goOnWindowDestroy
func goOnWindowDestroy(webView C.mbWebView, param, _ unsafe.Pointer) int {
	if _, ok := windowViews[webView]; ok {
		delete(windowViews, webView)
	}
	if cb := pointer.Restore(param).(*func()); cb != nil && *cb != nil {
		(*cb)()
		pointer.Unref(param)
	}
	return 0
}

//export goOnLoadUrlBegin
func goOnLoadUrlBegin(webView C.mbWebView, param unsafe.Pointer, url *C.char, job C.mbNetJob) int {
	if cb := pointer.Restore(param).(*func(url string, job C.mbNetJob)); cb != nil && *cb != nil {
		(*cb)(C.GoString(url), job)
	}
	return 0
}

//export goOnLoadUrlEnd
func goOnLoadUrlEnd(webView C.mbWebView, param unsafe.Pointer, url *C.char, job C.mbNetJob, buf unsafe.Pointer, length C.int) {
	if cb := pointer.Restore(param).(*func(url string, job C.mbNetJob, buf unsafe.Pointer, length int)); cb != nil && *cb != nil {
		(*cb)(C.GoString(url), job, buf, int(length))
	}
}
