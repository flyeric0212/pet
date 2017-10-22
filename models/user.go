/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/10/11 22:58 
 */
package models

import (
    "time"
    "pet/utils"
    "third/gorm"
)

// 用户信息表
type User struct {
    UserId          int64           `gorm:"primary_key"; sql:"AUTO_INCREMENT"`
    Name            string          `sql:"type:varchar(128)"`   // 姓名
    Nickname        string          `sql:"type:varchar(128)"`  // 昵称
    Password        string          `sql:"type:varchar(128)"`
    Gender          string          `sql:"type:smallint(6)"`    // 性别，0: 无性别 1: 男 2: 女
    Phone           string          `sql:"type:varchar(64)"`
    Email           string          `sql:"type:varchar(128)"`
    Openid          string          `sql:"type:varchar(255)"`   // 微信公共号的用户标志
    CreateTime      time.Time       `sql:"type:datetime"`
    UpdateTime      time.Time       `sql:"type:datetime"`
    LastLogin       time.Time       `sql:"type:datetime"`
}

func (user_info *User) TableName() string {
    return "pet.user"
}

// 新建用户
func (user_info *User) Create(id *int64) error {
    user_info.CreateTime = time.Now()
    user_info.UpdateTime = time.Now()

    err := PET_DB.Table(user_info.TableName()).Create(user_info).Error
    if nil != err {
        err = utils.NewInternalError(utils.DbErrCode, err)
        utils.Logger.Error("create user error: %v", err)
        return err
    }
    id = &user_info.UserId
    return nil
}

// 判断昵称是否存在
func CheckNicknameExist(nickname string) (err error, flag bool, user_info *User) {
    user_info = new(User)
    err = PET_DB.Table("pet.user").Where("nickname = ?", nickname).Limit(1).Find(user_info).Error
    if err == gorm.RecordNotFound {
        flag = false
        err = nil
    } else if err == nil {
        flag = true
    } else {
        err = utils.NewInternalError(utils.DbErrCode, err)
        utils.Logger.Error("CheckNicknameExist failed, error: %v", err)
        return
    }
    return
}