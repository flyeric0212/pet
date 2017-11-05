/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/10/11 22:58 
 */
package model

import (
    "time"
    "pet/utils"
    "third/gorm"
)

// 用户信息表
type User struct {
    UserId          int64           `gorm:"primary_key"; sql:"AUTO_INCREMENT"`
    Phone           string          `sql:"type:varchar(64)"`    // 电话号码
    Name            string          `sql:"type:varchar(128)"`   // 姓名
    Gender          string          `sql:"type:smallint(6)"`    // 性别，0: 无性别 1: 男 2: 女
    Email           string          `sql:"type:varchar(128)"`

    RegistType      int             `sql:"type:"smallint(6)"`   // 1：微信  2：官网
    Nickname        string          `sql:"type:varchar(128)"`   // 微信昵称
    Avatar          string          `sql:"type:varchar(128)"`   // 微信头像
    Openid          string          `sql:"type:varchar(255)"`   // 微信公共号的用户标志

    CreateTime      time.Time       `sql:"type:datetime"`
    UpdateTime      time.Time       `sql:"type:datetime"`

    LastLogin       time.Time       `sql:"type:datetime"`
    Password        string          `sql:"type:varbinary(128)"`
}

func (user_info *User) TableName() string {
    return "pet.user"
}

// 新建用户
func (user_info *User) Create(id *int64) error {
    now := time.Now()
    user_info.CreateTime = now
    user_info.UpdateTime = now
    user_info.LastLogin = now

    err := PET_DB.Table(user_info.TableName()).Create(user_info).Error
    if nil != err {
        err = utils.NewInternalError(utils.DbErrCode, err)
        utils.Logger.Error("create user error: %v", err)
        return err
    }
    id = &user_info.UserId
    return nil
}

// 通过微信用户标志拉取用户标志
func (user_info *User) GetUserByOpenid(openid string) error {
    err := PET_DB.Table(user_info.TableName()).Where("openid = ?", openid).Limit(1).Find(user_info).Error
    if gorm.RecordNotFound == err {
        utils.Logger.Warning("user not found by openid, openid: %s", openid)
        err = nil
    } else if nil != err {
        err = utils.NewInternalError(utils.DbErrCode, err)
        utils.Logger.Error("get user by openid failed, openid: %s, error: %v", openid, err)
        return err
    }
    return nil
}

// 判断电话号码是否存在
func CheckPhoneExist(phone string) (err error, flag bool, user_info *User) {
    user_info = new(User)
    err = PET_DB.Table("pet.user").Where("phone = ?", phone).Limit(1).Find(user_info).Error
    if err == gorm.RecordNotFound {
        flag = false
        err = nil
    } else if err == nil {
        flag = true
    } else {
        err = utils.NewInternalError(utils.DbErrCode, err)
        utils.Logger.Error("CheckPhoneExist failed, error: %v", err)
        return
    }
    return
}