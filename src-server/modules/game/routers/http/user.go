// Copyright (c) 554949297@qq.com . 2022-2022. All rights reserved

//go:build http
// +build http

package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zhouhp1295/g3"
	"github.com/zhouhp1295/g3-game/boot"
	"github.com/zhouhp1295/g3-game/modules/game/dao"
	"github.com/zhouhp1295/g3-game/modules/game/helpers"
	"github.com/zhouhp1295/g3-game/modules/game/model"
	"github.com/zhouhp1295/g3/crud"
	"github.com/zhouhp1295/g3/net"
	"go.uber.org/zap"
	"net/http"
)

type _gameUserApi struct {
	net.BaseApi
}

var GameUserApi = &_gameUserApi{
	net.BaseApi{Dao: dao.GameUserDao},
}

const (
	PermGameUserList   = "game:user:list"
	PermGameUserQuery  = "game:user:query"
	PermGameUserAdd    = "game:user:add"
	PermGameUserEdit   = "game:user:edit"
	PermGameUserRemove = "game:user:remove"
)

func init() {
	boot.RegisterAfterInstallFunction(func() {
		g3.GetGin().Group("/api").
			Bind(http.MethodGet, "/admin/game/user/page", GameUserApi.HandlePage, PermGameUserQuery)
		g3.GetGin().Group("/api").
			Bind(http.MethodGet, "/admin/game/user/get", GameUserApi.HandleGet, PermGameUserQuery)
		g3.GetGin().Group("/api").
			Bind(http.MethodPost, "/admin/game/user/insert", GameUserApi.HandleInsert, PermGameUserAdd)
		g3.GetGin().Group("/api").
			Bind(http.MethodPut, "/admin/game/user/update", GameUserApi.HandleUpdate, PermGameUserEdit)
		g3.GetGin().Group("/api").
			Bind(http.MethodPut, "/admin/game/user/status", GameUserApi.HandleUpdateStatus, PermGameUserEdit)
		g3.GetGin().Group("/api").
			Bind(http.MethodDelete, "/admin/game/user/delete", GameUserApi.HandleDelete, PermGameUserRemove)

		// 游戏接口
		g3.GetGin().Group("/api/game").MakeOpen("/user/fast")
		g3.GetGin().Group("/api/game").
			Bind(http.MethodGet, "/user/fast", onGameUserFastLogin)
		g3.GetGin().Group("/api/game").
			Bind(http.MethodPost, "/user/fast", onGameUserFastLogin)
	})
}

func onGameUserFastLogin(ctx *gin.Context) {
	user := new(model.GameUser)
	user.Username = fmt.Sprintf("%v", uuid.New())
	user.Nickname = "游客"
	crud.DbSess().Create(user)
	httpToken, err := g3.GetGin().Group("/api/game").NewJwtToken(user.Id, "")
	if err != nil {
		g3.ZL().Error("create game token failed. please check")
		net.FailedServerError(ctx, err.Error(), "")
		return
	}
	wsToken, err := helpers.NewWsJwtToken([]byte(boot.JwtCfg.WsSecret), user.Id)
	if err != nil {
		g3.ZL().Error("create websocket token failed", zap.Error(err))
		net.FailedMessage(ctx, "create jwt wsToken failed")
		return
	}
	net.SuccessData(ctx, gin.H{
		"httpToken":      httpToken,
		"websocketToken": wsToken,
	})
}
