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
type UserProfile struct {
    Id              int64           `gorm:"primary_key"; sql:"AUTO_INCREMENT"`
    Name            string          `sql:"type:varchar(128)"`
    Gender          int             `sql:"type:smallint(6)"`    // 性别，0: 无性别 1: 男 2: 女
    Phone           string          `sql:"type:varchar(64)"`
    Email           string          `sql:"type:varchar(128)"`
    Openid          string          `sql:"type:varchar(255)"`
    CreateTime      time.Time       `sql:"type:datetime"`
    UpdateTime      time.Time       `sql:"type:datetime"`
}

func (user_profile *UserProfile) TableName() string {
    return "pet.user_profile"
}

func (user_profile *UserProfile) Create(id *int64) error {
    user_profile.CreateTime = time.Now()

    err := PET_DB.Table(user_profile.TableName()).Create(user_profile).Error
    if nil != err {
        err = utils.NewInternalError(utils.DbErrCode, err)
        utils.Logger.Error("create user_profile error: %v", err)
        return err
    }
    return nil
}