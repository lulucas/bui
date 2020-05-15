package bui

import (
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
	"math/rand"
	"syscall"
	"time"
	"unsafe"
)

type Tray struct {
	hwnd              win.HWND
	guid              syscall.GUID
	onLeftMouseClick  func()
	onRightMouseClick func()
}

const TrayMsg = win.WM_APP + 1

func NewTray(hwnd win.HWND) *Tray {
	tray := &Tray{hwnd: hwnd, guid: newGUID()}
	nid := tray.createNid()
	nid.UFlags |= win.NIF_MESSAGE
	nid.UCallbackMessage = TrayMsg
	if !win.Shell_NotifyIcon(win.NIM_ADD, nid) {
		return nil
	}
	return tray
}

func (t *Tray) createNid() *win.NOTIFYICONDATA {
	var nid win.NOTIFYICONDATA
	nid.CbSize = uint32(unsafe.Sizeof(nid))
	nid.UFlags = win.NIF_GUID
	nid.HWnd = t.hwnd
	nid.GuidItem = t.guid
	return &nid
}

func (t *Tray) SetIcon(icon win.HICON) {
	nid := t.createNid()
	nid.UFlags |= win.NIF_ICON
	nid.HIcon = icon
	win.Shell_NotifyIcon(win.NIM_MODIFY, nid)
}

func (t *Tray) SetTooltip(tooltip string) {
	nid := t.createNid()
	nid.UFlags |= win.NIF_TIP
	copy(nid.SzTip[:], windows.StringToUTF16(tooltip))
	win.Shell_NotifyIcon(win.NIM_MODIFY, nid)
}

func (t *Tray) Dispose() {
	win.Shell_NotifyIcon(win.NIM_DELETE, t.createNid())
}

func (t *Tray) OnLeftMouseClick(f func()) {
	t.onLeftMouseClick = f
}

func (t *Tray) OnRightMouseClick(f func()) {
	t.onRightMouseClick = f
}

func newGUID() syscall.GUID {
	var buf [16]byte
	rand.Read(buf[:])
	return *(*syscall.GUID)(unsafe.Pointer(&buf[0]))
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
