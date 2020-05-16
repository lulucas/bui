//go:generate esc -prefix web/dist -o ui.go web/dist
package main

import (
	"fmt"
	"github.com/lulucas/bui"
	"log"
	"os"
	"time"
)

func main() {
	bui.Initialize()
	app := bui.NewApp()
	app.SetMainView(bui.CreateView(bui.CreateViewOption{
		Title:       "bui",
		Width:       900,
		Height:      600,
		Transparent: false,
	}))
	if icon, err := bui.IconFromBytes(FSMustByte(false, "/favicon.ico")); err == nil {
		app.MainView().SetIcon(icon)
		app.MainView().Tray().SetIcon(icon)
		app.MainView().Tray().SetTooltip("BUI")
		app.MainView().Tray().OnLeftMouseClick(func() {
			app.MainView().Restore()
			app.MainView().ShowOnTop()
		})
	}
	app.MainView().SetFileSystem(FS(false))
	app.MainView().RPC().Register("sum", func(params []int) (int, error) {
		log.Printf("call sum: %v", params)
		return params[0] + params[1], nil
	})
	app.MainView().RPC().Register("open_url", func(params struct{ Url string }) {
		log.Printf("notify open url: %s", params.Url)
	})
	app.MainView().RPC().Register("minimize_to_tray", func() {
		bui.WkeAsyncCall(func() {
			app.MainView().Hide()
		})
	})
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			if app.MainView() == nil {
				break
			}
			app.MainView().RPC().Emit("state_changed", struct {
				State string
				Time  time.Time
			}{
				State: "start",
				Time:  time.Now(),
			})
		}
	}()
	app.MainView().ShowDevTools("devtools")
	app.OnStart(func() {
		fmt.Println("BUI started!")
	})
	app.OnStop(func() {
		// Clear
		os.Exit(0)
	})
	app.Start()
}
