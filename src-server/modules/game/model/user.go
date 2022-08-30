// Copyright (c) 554949297@qq.com . 2022-2022. All rights reserved

package model

import "github.com/zhouhp1295/g3/crud"

type GameUser struct {
	crud.BaseModel
	Username string `gorm:"INDEX;TYPE:VARCHAR(36);COMMENT:登录账户" json:"username" form:"username" query:"like"`
	Nickname string `gorm:"TYPE:VARCHAR(20);COMMENT:昵称" json:"nickname" form:"nickname" query:"like"`
	Avatar   string `gorm:"TYPE:VARCHAR(255);COMMENT:头像" json:"avatar" form:"avatar"`
	Sex      string `gorm:"TYPE:CHAR(1);NOT NULL;DEFAULT:0;COMMENT:性别" json:"sex" form:"sex" query:"like"`
	crud.TailColumns
}

// Table 返回表名
func (*GameUser) Table() string {
	return "game_user"
}

// NewModel 返回实例
func (*GameUser) NewModel() crud.ModelInterface {
	return new(GameUser)
}

// NewModels 返回实例数组
func (*GameUser) NewModels() interface{} {
	return make([]GameUser, 0)
}
