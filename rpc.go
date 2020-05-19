package bui

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/sourcegraph/jsonrpc2"
	"net/http"
	"reflect"
	"sync"
)

var upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type RPC struct {
	methodMtx sync.RWMutex
	methods   map[string]Handler
	eventMtx  sync.RWMutex
	event     map[string]map[*websocket.Conn]bool
}

func NewRPC() *RPC {
	rpc := &RPC{
		methods: map[string]Handler{},
		event:   map[string]map[*websocket.Conn]bool{},
	}
	rpc.registerEventMethod()
	return rpc
}

func (rpc *RPC) registerEventMethod() {
	// The 'rpc.on' and 'rpc.off' methods are expanded by rpc-websockets
	// which is used to subscribe and unsubscribe event.
	rpc.Register("rpc.on", func(conn *websocket.Conn, events []string) {
		rpc.eventMtx.Lock()
		defer rpc.eventMtx.Unlock()
		for _, evt := range events {
			if rpc.event[evt] == nil {
				rpc.event[evt] = map[*websocket.Conn]bool{}
			}
			rpc.event[evt][conn] = true
		}
	})
	rpc.Register("rpc.off", func(conn *websocket.Conn, events []string) {
		rpc.eventMtx.Lock()
		defer rpc.eventMtx.Unlock()
		for _, evt := range events {
			if rpc.event[evt] != nil {
				delete(rpc.event[evt], conn)
			}
		}
	})
}

func (rpc *RPC) Register(method string, handle interface{}) {
	rpc.methodMtx.Lock()
	defer rpc.methodMtx.Unlock()
	rpc.methods[method] = NewHandler(handle)
}

func (rpc *RPC) Emit(event string, params interface{}) {
	rpc.eventMtx.RLock()
	defer rpc.eventMtx.RUnlock()

	connMap, exists := rpc.event[event]
	if !exists {
		return
	}

	notification := struct {
		Notification string          `json:"notification"`
		Params       json.RawMessage `json:"params"`
	}{
		Notification: event,
	}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		log.Errorf("emit %s event, error: %s", event, err.Error())
		return
	}
	notification.Params = paramsBytes

	for conn, _ := range connMap {
		err = conn.WriteJSON(notification)
		if err != nil {
			log.Errorf("emit %s event to %s error: %s", event, conn.RemoteAddr().String(), err.Error())
		}
	}
}

func (rpc *RPC) response(conn *websocket.Conn, id jsonrpc2.ID, result *json.RawMessage) {
	res := jsonrpc2.Response{
		ID:     id,
		Result: result,
	}
	if err := conn.WriteJSON(res); err != nil {
		log.Errorf("websocket write json error: %s", err.Error())
		return
	}
}

func (rpc *RPC) error(conn *websocket.Conn, id jsonrpc2.ID, code int64, err error) {
	res := jsonrpc2.Response{
		ID: id,
		Error: &jsonrpc2.Error{
			Code:    code,
			Message: err.Error(),
		},
	}
	if err := conn.WriteJSON(res); err != nil {
		log.Errorf("websocket write json error: %s", err.Error())
		return
	}
}

func (rpc *RPC) dispatch(conn *websocket.Conn, req *jsonrpc2.Request) {
	if req == nil {
		return
	}
	rpc.methodMtx.RLock()
	defer rpc.methodMtx.RUnlock()

	if m, ok := rpc.methods[req.Method]; ok {
		resBytes, err := m.Invoke(conn, *req.Params)
		if err != nil {
			rpc.error(conn, req.ID, jsonrpc2.CodeInvalidRequest, err)
			return
		}
		if !req.Notif {
			rpc.response(conn, req.ID, (*json.RawMessage)(&resBytes))
		}
	} else {
		if !req.Notif {
			rpc.error(conn, req.ID, jsonrpc2.CodeMethodNotFound, fmt.Errorf("method %s not found", req.Method))
		}
	}
}

func (rpc *RPC) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("websocket conn error: %s", err.Error())
		return
	}
	for {
		var req jsonrpc2.Request
		if err := conn.ReadJSON(&req); err != nil {
			rpc.clearConn(conn)
			break
		}
		rpc.dispatch(conn, &req)
	}
}

func (rpc *RPC) websocket(c echo.Context) error {
	rpc.wsHandler(c.Response(), c.Request())
	return nil
}

func (rpc *RPC) clearConn(conn *websocket.Conn) {
	rpc.eventMtx.Lock()
	for _, connections := range rpc.event {
		for c, _ := range connections {
			if c == conn {
				delete(connections, c)
			}
		}
	}
	rpc.eventMtx.Unlock()
}

type Handler interface {
	Invoke(conn *websocket.Conn, payload []byte) ([]byte, error)
}

type genericHandler func(*websocket.Conn, []byte) (interface{}, error)

func (handler genericHandler) Invoke(conn *websocket.Conn, payload []byte) ([]byte, error) {
	response, err := handler(conn, payload)
	if err != nil {
		return nil, err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return responseBytes, nil
}

func errorHandler(e error) genericHandler {
	return func(conn *websocket.Conn, event []byte) (interface{}, error) {
		return nil, e
	}
}

func validateArguments(handler reflect.Type) (bool, error) {
	handlerTakesConn := false
	if handler.NumIn() > 2 {
		return false, fmt.Errorf("handlers may not take more than two arguments, but handler takes %d", handler.NumIn())
	} else if handler.NumIn() > 0 {
		connType := reflect.TypeOf((*websocket.Conn)(nil))
		argumentType := handler.In(0)
		handlerTakesConn = argumentType == connType
		if handler.NumIn() > 1 && !handlerTakesConn {
			return false, fmt.Errorf("handler takes two arguments, but the first is not conn. got %s", argumentType.Kind())
		}
	}

	return handlerTakesConn, nil
}

func validateReturns(handler reflect.Type) error {
	errorType := reflect.TypeOf((*error)(nil)).Elem()

	switch n := handler.NumOut(); {
	case n > 2:
		return fmt.Errorf("handler may not return more than two values")
	case n > 1:
		if !handler.Out(1).Implements(errorType) {
			return fmt.Errorf("handler returns two values, but the second does not implement error")
		}
	case n == 1:
		if !handler.Out(0).Implements(errorType) {
			return fmt.Errorf("handler returns a single value, but it does not implement error")
		}
	}

	return nil
}

func NewHandler(handlerFunc interface{}) Handler {
	if handlerFunc == nil {
		return errorHandler(fmt.Errorf("handler is nil"))
	}
	handler := reflect.ValueOf(handlerFunc)
	handlerType := reflect.TypeOf(handlerFunc)
	if handlerType.Kind() != reflect.Func {
		return errorHandler(fmt.Errorf("handler kind %s is not %s", handlerType.Kind(), reflect.Func))
	}

	takesConn, err := validateArguments(handlerType)
	if err != nil {
		return errorHandler(err)
	}

	if err := validateReturns(handlerType); err != nil {
		return errorHandler(err)
	}

	return genericHandler(func(conn *websocket.Conn, payload []byte) (interface{}, error) {
		// construct arguments
		var args []reflect.Value
		if takesConn {
			args = append(args, reflect.ValueOf(conn))
		}
		if (handlerType.NumIn() == 1 && !takesConn) || handlerType.NumIn() == 2 {
			eventType := handlerType.In(handlerType.NumIn() - 1)
			event := reflect.New(eventType)

			if err := json.Unmarshal(payload, event.Interface()); err != nil {
				return nil, err
			}
			args = append(args, event.Elem())
		}

		response := handler.Call(args)

		// convert return values into (interface{}, error)
		var err error
		if len(response) > 0 {
			if errVal, ok := response[len(response)-1].Interface().(error); ok {
				err = errVal
			}
		}
		var val interface{}
		if len(response) > 1 {
			val = response[0].Interface()
		}

		return val, err
	})
}
