/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/10/11 22:58 
 */
package models

import (
    "time"
    "pet/utils"
)

// 用户信息表
type User struct {
    UserId          int64           `gorm:"primary_key"; sql:"AUTO_INCREMENT"`
    Name            string          `sql:"type:varchar(128)"`
    Gender          int             `sql:"type:smallint(6)"`    // 性别，0: 无性别 1: 男 2: 女
    Phone           string          `sql:"type:varchar(64)"`
    Email           string          `sql:"type:varchar(128)"`
    Openid          string          `sql:"type:varchar(255)"`   //
    CreateTime      time.Time       `sql:"type:datetime"`
    UpdateTime      time.Time       `sql:"type:datetime"`
}

func (user *User) TableName() string {
    return "pet.user"
}

func (user *User) Create(id *int64) error {
    user.CreateTime = time.Now()

    err := PET_DB.Table(user.TableName()).Create(user).Error
    if nil != err {
        err = utils.NewInternalError(utils.DbErrCode, err)
        utils.Logger.Error("create user error: %v", err)
        return err
    }
    id = &user.UserId
    return nil
}