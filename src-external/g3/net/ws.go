// Copyright (c) 554949297@qq.com . 2022-2022 . All rights reserved

package net

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/zhouhp1295/g3"
	"github.com/zhouhp1295/g3/helpers"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

var workers map[string]bool

type WsConnStatus int

const (
	WsErrorFailed = -1
	WsErrorOk     = 0

	WsErrorOKMsg = "OK"
)

const (
	WsConnecting WsConnStatus = iota
	WsAuthing
	WsConnected
	WsClosed
)

type WsConn struct {
	Uuid     string
	Conn     *websocket.Conn
	Query    map[string]string
	Data     map[string]interface{}
	CreateAt time.Time
	Status   WsConnStatus
	Mutex    sync.Mutex
}

func (wc *WsConn) Failed(router string, msg string) {
	wc.WriteJSON(router, WsErrorFailed, msg, gin.H{})
}

func (wc *WsConn) Ok(router string, data interface{}) {
	wc.WriteJSON(router, WsErrorOk, WsErrorOKMsg, data)
}

func (wc *WsConn) WriteJSON(router string, code int, msg string, data interface{}) {
	err := wc.Conn.WriteJSON(gin.H{
		"router": router,
		"code":   code,
		"msg":    msg,
		"data":   data,
	})
	if err != nil {
		g3.ZL().Error("failed to write to websocket connection",
			zap.String("router", router),
			zap.Error(err))
	}
}

type WsWorker struct {
	Router string

	upgrader websocket.Upgrader

	onConnected func(conn *WsConn)
	onClosed    func(conn *WsConn)
	onError     func(conn *WsConn, err error)
	onMessage   func(conn *WsConn, message []byte)

	connHandler WsConnHandler

	connections map[string]*WsConn
	rwMutex     sync.RWMutex
}

func (w *WsWorker) listen(conn *WsConn) {
	g3.ZL().Info("start listen", zap.String("uuid", conn.Uuid))
	for {
		_, message, err := conn.Conn.ReadMessage()
		if err != nil {
			g3.ZL().Debug("on message err", zap.Error(err))
			if w.onError != nil {
				w.onError(conn, err)
			}
			break
		}
		g3.ZL().Debug("on message", zap.String("uuid", conn.Uuid))
		if w.onMessage != nil {
			w.onMessage(conn, message)
		} else {
			if WsAuthing == conn.Status {
				// 没注册自定义消息处理, 不进行auth验证
				conn.Status = WsConnected
				_ = conn.Conn.WriteJSON(map[string]interface{}{
					"msg": "Login Success!",
				})
				continue
			}
			_ = conn.Conn.WriteJSON(map[string]interface{}{
				"msg": "Message From Server : " + string(message),
			})
		}
	}
}

// handleConnect
func (w *WsWorker) handleConnect(conn *WsConn) bool {
	g3.ZL().Info("connected",
		zap.String("uuid", conn.Uuid),
		zap.Reflect("query", conn.Query),
	)
	w.rwMutex.Lock()
	w.connections[conn.Uuid] = conn
	conn.Status = WsAuthing
	if w.onConnected != nil {
		w.onConnected(conn)
	}
	w.rwMutex.Unlock()
	return true
}

func (w *WsWorker) closeAndDelete(conn *WsConn) {
	w.rwMutex.Lock()
	if WsClosed != conn.Status {
		if w.onClosed != nil {
			w.onClosed(conn)
		}
		if conn != nil {
			_ = conn.Conn.Close()
		}
		delete(w.connections, conn.Uuid)
	}
	g3.ZL().Info("connection closed",
		zap.String("uuid", conn.Uuid),
	)
	w.rwMutex.Unlock()
}

func (w *WsWorker) closeConn(conn *WsConn) {
	g3.ZL().Info("auto closing connection",
		zap.String("uuid", conn.Uuid),
	)
	w.closeAndDelete(conn)
}

func (w *WsWorker) Close(conn *WsConn) {
	g3.ZL().Info("manual closing",
		zap.String("uuid", conn.Uuid),
	)
	w.closeAndDelete(conn)
}

// WsConnHandler WsConnHandler
type WsConnHandler func(conn *WsConn)

// WsWorkerOption WsWorkerOption
type WsWorkerOption func(*WsWorker)

// WithWsUpgrader 自定义upgrader
func WithWsUpgrader(up websocket.Upgrader) WsWorkerOption {
	return func(worker *WsWorker) {
		worker.upgrader = up
	}
}

// WithWsConnHandler 自定义conn处理
func WithWsConnHandler(handler WsConnHandler) WsWorkerOption {
	return func(worker *WsWorker) {
		worker.connHandler = handler
	}
}

// WithWsConnected 自定义建立链接处理
func WithWsConnected(handler func(conn *WsConn)) WsWorkerOption {
	return func(worker *WsWorker) {
		worker.onConnected = handler
	}
}

// WithWsMessaged 自定义消息处理
func WithWsMessaged(handler func(conn *WsConn, message []byte)) WsWorkerOption {
	return func(worker *WsWorker) {
		worker.onMessage = handler
	}
}

// WithWsError 自定义错误处理
func WithWsError(handler func(conn *WsConn, err error)) WsWorkerOption {
	return func(worker *WsWorker) {
		worker.onError = handler
	}
}

// defaultUpgrader defaultUpgrader
func defaultUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

func init() {
	workers = make(map[string]bool)
}

func onError(w *WsWorker, conn *WsConn, err error) {
	g3.ZL().Error("connected",
		zap.String("uuid", conn.Uuid),
		zap.Reflect("query", conn.Query),
		zap.Error(err),
	)
	if conn.Conn != nil {
		_ = conn.Conn.Close()
	}
	if w.onError != nil {
		w.onError(conn, err)
	}
}

// HandleWebsocket HandleWebsocket
func HandleWebsocket(router string, opts ...WsWorkerOption) (*WsWorker, error) {
	if _, exist := workers[router]; exist {
		return nil, errors.New("Create Websocket Router Failed :" + router)
	}
	workers[router] = true
	worker := new(WsWorker)
	worker.Router = router
	worker.upgrader = defaultUpgrader()
	worker.connections = make(map[string]*WsConn)
	for _, opt := range opts {
		opt(worker)
	}
	http.HandleFunc(router, func(writer http.ResponseWriter, request *http.Request) {
		var conn *websocket.Conn
		g3Conn := new(WsConn)
		g3Conn.Status = WsConnecting
		g3Conn.Uuid = fmt.Sprintf("%v", uuid.New())
		g3Conn.Query = helpers.ParseQueryString(request.RequestURI)
		g3Conn.Data = make(map[string]interface{})
		g3Conn.CreateAt = time.Now()
		conn, err := worker.upgrader.Upgrade(writer, request, nil)
		if err != nil {
			onError(worker, g3Conn, err)
			return
		}
		if worker.connHandler != nil {
			worker.connHandler(g3Conn)
			return
		}
		defer worker.closeConn(g3Conn)
		g3Conn.Conn = conn
		if !worker.handleConnect(g3Conn) {
			return
		}
		worker.listen(g3Conn)
	})
	return worker, nil
}
