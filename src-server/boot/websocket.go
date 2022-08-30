// Copyright (c) 554949297@qq.com . 2022-2022. All rights reserved

//go:build websocket
// +build websocket

package boot

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/zhouhp1295/g3"
	"github.com/zhouhp1295/g3-game/utils"
	"github.com/zhouhp1295/g3/net"
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
	"net/http"
	"sync"
)

const UserAuthRouter = "user/auth"

var (
	worker *net.WsWorker

	wsAuthHandler WsRouterHandler

	wsRouterHandlers     map[string][]WsRouterHandler
	wsRouterHandlerMutex sync.Mutex

	wsUnAuthResp = WsResponseMsg{
		Msg: "未授权",
	}

	wsDefaultResp = WsResponseMsg{
		Msg: "Coming soon ...",
	}
)

type WsRouterHandler func(*net.WsWorker, *net.WsConn, WsRequestMsg)

func RegisterWsAuthHandler(handler WsRouterHandler) {
	wsAuthHandler = handler
}

func RegisterWsRouterHandler(router string, handler WsRouterHandler) {
	wsRouterHandlerMutex.Lock()
	if UserAuthRouter != router {
		wsRouterHandlers[router] = append(wsRouterHandlers[router], handler)
	}
	wsRouterHandlerMutex.Unlock()
}

type WsRequestMsg struct {
	Router string                 `json:"router"`
	Params map[string]interface{} `json:"params"`
}

func (wr *WsRequestMsg) Get(key string) (interface{}, error) {
	if wr.Params == nil {
		g3.ZL().Error("get value failed. params is empty",
			zap.String("key", key))
		return nil, errors.New("params is empty")
	}
	v, ok := wr.Params[key]
	if !ok {
		g3.ZL().Error("get value failed. key is not exist",
			zap.String("key", key),
			zap.Reflect("params", wr.Params))
		return nil, errors.New("key is not exist")
	}
	return v, nil
}

func (wr *WsRequestMsg) GetString(key string) (string, error) {
	v, err := wr.Get(key)
	if err != nil {
		return "", err
	}
	str, ok := v.(string)
	if !ok {
		g3.ZL().Error("get string failed. type is not incorrect",
			zap.Reflect("key", key))
		return "", errors.New("value type is incorrect")
	}
	return str, nil
}

type WsResponseMsg struct {
	Router string                 `json:"router"`
	Code   int                    `json:"code"`
	Msg    string                 `json:"msg"`
	Data   map[string]interface{} `json:"data"`
}

func init() {
	App.Name = "websocket"
	App.Identifier = "node01"
	App.Version = "v0.0.1"
	App.RunMode = "dev"
	preStart = preWebsocketStart
	start = startWebsocket
	wsRouterHandlers = make(map[string][]WsRouterHandler)
}

func loadWebsocketConfig() {
	iniPath := g3.AssetPath("conf/websocket.ini")
	if !utils.IsExist(iniPath) {
		panic("未找到配置文件: " + iniPath)
	}
	iniFile, err := ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, iniPath)

	if err != nil {
		panic(err)
	}
	// ***************************
	// ----- ServerCfg settings -----
	// ***************************
	if err = iniFile.Section("server").MapTo(&ServerCfg); err != nil {
		panic(err)
	}
}

func onWsMessage(conn *net.WsConn, msg []byte) {
	conn.Mutex.Lock()
	defer conn.Mutex.Unlock()
	reqParams := WsRequestMsg{}
	err := jsoniter.Unmarshal(msg, &reqParams)
	if err != nil {
		g3.ZL().Error("on message err .", zap.Error(err))
		worker.Close(conn)
		return
	}
	if net.WsAuthing == conn.Status {
		// 新的链接，第一个动作必须是授权验证
		if wsAuthHandler == nil || UserAuthRouter != reqParams.Router {
			g3.ZL().Error("websocket auth handler undefined.")
			conn.Failed(reqParams.Router, "auth failed")
			worker.Close(conn)
			return
		}
		wsAuthHandler(worker, conn, reqParams)
		return
	}
	if funcList, ok := wsRouterHandlers[reqParams.Router]; ok {
		for _, f := range funcList {
			f(worker, conn, reqParams)
		}
	} else {
		g3.ZL().Warn("undefined router. please check.",
			zap.String("router", reqParams.Router))
		conn.Failed(reqParams.Router, "undefined router")
	}
}

func preWebsocketStart() {
	if !IsInstalled() {
		panic("系统未安装")
	}
	// 加载配置
	loadWebsocketConfig()
	// 启动
	var err error
	worker, err = net.HandleWebsocket("/", net.WithWsMessaged(onWsMessage))
	if err != nil {
		g3.ZL().Fatal("服务启动失败", zap.Error(err))
	}
}

func startWebsocket() {
	err := http.ListenAndServe(ServerCfg.HTTPAddr+":"+ServerCfg.HTTPPort, nil)
	if err != nil {
		g3.ZL().Fatal("Application Stopped", zap.Error(err))
	}
}
