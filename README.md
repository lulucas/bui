# BUI

Golang UI base on https://github.com/weolar/miniblink49, only for windows.

## Quick Start

Download miniblink release, copy `ndoe.dll` and `mb.dll` to work directory.

If your application is 64bit, copy `miniblink_x64.dll` and `mb_x64.dll` to work directory.

```
https://github.com/weolar/miniblink49/releases
```

## Life Cycle

### App life cycle

1. OnBeforeStart()
1. OnStart()
1. OnBeforeStop()
1. OnStop()

## Examples

### Simple app
```go
package main

import "github.com/lulucas/bui"

func main() {
    bui.Initialize()
    app := bui.NewApp()
    app.SetMainView(bui.CreateView(bui.CreateViewOption{
        Title:       "bui",
        Width:       900,
        Height:      600,
        Transparent: false,
    }))
    app.Start()
}
```

### Tray

```
if icon, err := bui.IconFromBytes(FSMustByte(false, "/favicon.ico")); err == nil {
    app.MainView().SetIcon(icon)
    app.MainView().Tray().SetIcon(icon)
    app.MainView().Tray().SetTooltip("BUI")
    app.MainView().Tray().OnLeftMouseClick(func() {
        app.MainView().Restore()
        app.MainView().ShowOnTop()
    })
}
```

### RPC

BUI uses websocket as rpc channel.

In web client side, you can use [rpc-websockets](https://www.npmjs.com/package/rpc-websockets).

BUI serves build-in rpc which is compatible with rpc-websockets. 

#### Golang side example

```go
package main

import (
	"github.com/lulucas/bui"
	"fmt"
	"log"
	"time"
)

func main() {
    bui.Initialize()
    app := bui.NewApp()
    app.SetMainView(bui.CreateView(bui.CreateViewOption{}))
    app.MainView().RPC().Register("sum", func(params []int) (int, error) {
        log.Printf("call sum: %v", params)
        return params[0] + params[1], nil
    })
    app.MainView().RPC().Register("open_url", func(params struct{Url string}) {
        log.Printf("notify open url: %s", params.Url)
    })
    time.AfterFunc(5*time.Second, func() {
        // emit a event to who subscribes this event name
        app.MainView().RPC().Emit("state_changed", struct {
            State string
            Time  time.Time
        }{
            State: "start",
            Time:  time.Now(),
        })
    })
    app.Start()
}
```

#### Web side example

```javascript
const WebSocket = require('rpc-websockets').Client
const ws = new WebSocket(`ws://127.0.0.1:${window.location.port}/rpc`)
ws.on('open', () => {
  ws.call('sum', [5, 3]).then(result => {
    console.log(`sum result: ${result}`)
  })

  ws.notify('open_url', {url: "https://www.google.com"})

  ws.subscribe('state_changed')

  ws.on('state_changed', state => {
    console.log(`stateChanged ${state}`)
  })
})
```

#### Main Thread Call Problem For Wke

Miniblink function must be called at the main thread,
you can use `bui.SyncCall` or `bui.AsyncCall` to make calling in the main thread.

```
bui.AsyncCall(func() {
    app.MainView().Hide()
})
```

### Devtools

Extract miniblink `front_end` folder to disk, 
then `ShowDevTools(path string)` of WebView to show dev tools.

### Application Update

#### TBD

Client side

```go
```

Server side metadata

```json
{
  "Version": "1.1.0",
  "Force": false
}
```


## Build

First, install pack tools.
 
* https://github.com/akavel/rsrc
* https://github.com/mjibson/esc

### Icon and UAC

Currently, rsrc has a manifest bug with CGO build.

Please make a launcher.exe to start bui app.

### Web UI assets

You can custom your path or package name. 

```
esc -pkg ui -prefix ui_folder -o ui/ui.go ui_folder
```

Set file system to app webview

```go
package main

import (
    "github.com/lulucas/bui"
    "ui"
)

func main() {
    bui.Initialize()
    app := bui.NewApp()
    app.SetMainView(bui.CreateView(bui.CreateViewOption{}))
    app.MainView().SetFileSystem(ui.FS(false))
    app.Start()
}
```

### Hide console window

```
go build -ldflags "-w -s -H=windowsgui"
```
