// Copyright (c) 554949297@qq.com . 2022-2022. All rights reserved

package game

import (
	"github.com/zhouhp1295/g3"
	"github.com/zhouhp1295/g3-game/modules/game/model"
	_ "github.com/zhouhp1295/g3-game/modules/game/routers"
	"github.com/zhouhp1295/g3/crud"
	"go.uber.org/zap"
)

func DoMigrate() {
}

func SyncTables() {
	//初始化数据结构
	tables := []interface{}{
		new(model.GameUser),
	}
	err := crud.SyncTables(crud.DbSess(), tables)
	if err != nil {
		g3.ZL().Fatal("AutoMigrate Game Database", zap.Error(err))
	}
}
