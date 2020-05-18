package bui

import (
	"github.com/lxn/win"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	mainView *WebView
	views    []*WebView
	quit     chan os.Signal
	// Life cycle
	onBeforeStart func()
	onStart       func()
	onBeforeStop  func()
	onStop        func()
}

func NewApp() *App {
	return &App{
		quit: make(chan os.Signal, 1),
	}
}

func (a *App) messageLoop() {
	signal.Notify(a.quit, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case f := <-uiCall:
			f()
		case <-a.quit:
			a.Close()
			Finalize()
			os.Exit(0)
		default:
			msg := &win.MSG{}
			if win.GetMessage(msg, 0, 0, 0) != 0 {
				win.TranslateMessage(msg)
				win.DispatchMessage(msg)
			}
		}
	}
}

func (a *App) SetMainView(view *WebView) {
	a.mainView = view
}

func (a *App) MainView() *WebView {
	return a.mainView
}

func (a *App) Start() {
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

func (a *App) Close() {
	if a.onBeforeStop != nil {
		a.onBeforeStop()
	}
	if a.mainView != nil {
		a.mainView.Close()
		a.mainView = nil
	}
	a.quit <- syscall.SIGTERM
	if a.onStop != nil {
		a.onStop()
	}
}

func (a *App) OnBeforeStart(f func()) {
	a.onBeforeStart = f
}

func (a *App) OnStart(f func()) {
	a.onStart = f
}

func (a *App) OnBeforeStop(f func()) {
	a.onBeforeStop = f
}

func (a *App) OnStop(f func()) {
	a.onStop = f
}
