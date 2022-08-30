// Copyright (c) 554949297@qq.com . 2022-2022 . All rights reserved

package dao

import (
	"github.com/zhouhp1295/g3-game/modules/game/model"
	"github.com/zhouhp1295/g3/crud"
)

type gameUserDAO struct {
	crud.BaseDao
}

var GameUserDao = &gameUserDAO{
	crud.BaseDao{Model: new(model.GameUser)},
}

func (dao *gameUserDAO) BeforeInsert(m crud.ModelInterface) (ok bool, msg string) {
	if _m, _ok := m.(*model.GameUser); _ok {
		if dao.CountByColumn("username", _m.Username) > 0 {
			msg = "用户名已存在"
			return
		}
		ok = true
	}
	return
}

func (dao *gameUserDAO) BeforeUpdate(m crud.ModelInterface) (ok bool, msg string) {
	if _m, _ok := m.(*model.GameUser); _ok {
		if dao.CountByColumn("username", _m.Username) > 1 {
			msg = "用户名已存在"
			return
		}
		ok = true
	}
	return
}
