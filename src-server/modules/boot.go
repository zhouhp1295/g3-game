// Copyright (c) 554949297@qq.com . 2022-2022 . All rights reserved

package modules

import (
	"github.com/zhouhp1295/g3-game/boot"
	_ "github.com/zhouhp1295/g3-game/modules/common/routers"
	"github.com/zhouhp1295/g3-game/modules/content"
	"github.com/zhouhp1295/g3-game/modules/game"
	_ "github.com/zhouhp1295/g3-game/modules/install/routers"
	_ "github.com/zhouhp1295/g3-game/modules/render"
	"github.com/zhouhp1295/g3-game/modules/system"
	"github.com/zhouhp1295/g3-game/modules/system/dao"
)

func init() {
	boot.RegisterAfterInstallFunction(func() {
		// 系统模块
		system.SyncTables()
		system.DoMigrate()
		// 内容管理模块
		content.SyncTables()
		content.DoMigrate()
		// 游戏模块
		game.SyncTables()
		game.DoMigrate()
		// 创建超级管理员
		dao.SysUserDao.CreateSuperUser("admin", "123456")
		// 初始化权限
		dao.SysRoleDao.RefreshRolePerms()
	})
}
