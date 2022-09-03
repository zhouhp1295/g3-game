// Copyright (c) 554949297@qq.com . 2022-2022. All rights reserved

//go:build websocket
// +build websocket

package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/zhouhp1295/g3"
	"github.com/zhouhp1295/g3-game/boot"
	"github.com/zhouhp1295/g3-game/modules/game/helpers"
	"github.com/zhouhp1295/g3/net"
	"go.uber.org/zap"
)

const (
	userInfoRouter = "user/info"
)

func init() {
	boot.RegisterWsAuthHandler(onUserAuth)
	boot.RegisterWsRouterHandler(userInfoRouter, onUserInfo)
}

func onUserAuth(worker *net.WsWorker, conn *net.WsConn, msg boot.WsRequestMsg) {
	token, err := msg.GetString("token")
	if err != nil {
		g3.ZL().Error("user auth failed . parse token failed. please check.",
			zap.Error(err))
		conn.Failed(msg.Router, "parse token failed")
		worker.Close(conn)
		return
	}
	claims, err := helpers.ParseWsJwtToken(token, []byte(boot.JwtCfg.WsSecret))
	if err != nil {
		g3.ZL().Error("user auth failed .token is incorrect. please check.",
			zap.Reflect("token", token),
			zap.Error(err))
		conn.Failed(msg.Router, "incorrect token")
		worker.Close(conn)
		return
	}
	conn.Data["uid"] = claims.Uid
	g3.ZL().Info("connection auth success", zap.String("uuid", conn.Uuid))
	conn.Status = net.WsConnected
	conn.Ok(msg.Router, gin.H{"uid": claims.Uid})
}

func onUserInfo(worker *net.WsWorker, conn *net.WsConn, msg boot.WsRequestMsg) {
	conn.Ok(msg.Router, msg)
}
