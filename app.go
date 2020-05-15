package bui

import (
	"github.com/lxn/win"
)

func init() {
	Initialize("")
}

type app struct {
	mainView *WebView
	views    []*WebView
	closing  chan bool

	// Life cycle
	onBeforeStart func()
	onStart       func()
	onBeforeStop  func()
	onStop        func()
}

func NewApp() *app {
	return &app{
		closing: make(chan bool, 1),
	}
}

func (a *app) messageLoop() {
	for {
		select {
		case f := <-wkeCall:
			f()
		case <-a.closing:
			return
		default:
			msg := &win.MSG{}
			if win.GetMessage(msg, 0, 0, 0) != 0 {
				win.TranslateMessage(msg)
				win.DispatchMessage(msg)
			}
		}
	}
}

func (a *app) SetMainView(view *WebView) {
	a.mainView = view
}

func (a *app) MainView() *WebView {
	return a.mainView
}

func (a *app) Start() {
	if a.onBeforeStart != nil {
		a.onBeforeStart()
	}
	a.mainView.OnDestroy(func() {
		a.Close()
	})
	a.mainView.Show()
	if a.onStart != nil {
		a.onStart()
	}
	a.messageLoop()
}

func (a *app) Close() {
	if a.onBeforeStop != nil {
		a.onBeforeStop()
	}
	if a.mainView != nil {
		a.mainView.Close()
		a.mainView = nil
	}
	a.closing <- true
	if a.onStop != nil {
		a.onStop()
	}
}

func (a *app) OnBeforeStart(f func()) {
	a.onBeforeStart = f
}

func (a *app) OnStart(f func()) {
	a.onStart = f
}

func (a *app) OnBeforeStop(f func()) {
	a.onBeforeStop = f
}

func (a *app) OnStop(f func()) {
	a.onStop = f
}
