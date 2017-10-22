/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/9/24 21:01 
 */
package models

import (
    "third/gorm"
    "fmt"
    "pet/utils"
    "pet/protocol"
)

var PET_DB *gorm.DB

func InitAllDB() error {
    var err error

    PET_DB, err = InitPetDb(utils.CommonConfig)
    if nil != err {
        utils.Logger.Warning("open pet db error")
        return err
    }

    return nil
}

func InitPetDb(config *utils.Configure) (*gorm.DB, error) {
    var err error
    PetDbSetting, ok := config.MysqlSetting["PetDbSetting"]
    if !ok {
        return nil, fmt.Errorf("can't find pet mysql settings")
    }
    PET_DB, err = utils.InitGormDbPool(&PetDbSetting)
    if err != nil {
        utils.Logger.Warning("open pet DB error")
        return nil, err
    }
    return PET_DB, nil
}

func CopyUserData(from *User, dst *protocol.UserInfoJson) error {
    //copy部分user的字段返回给api层
    dst.UserId = from.UserId
    dst.Name = from.Name
    dst.Nickname = from.Nickname
    dst.Gender = from.Gender
    dst.Phone = from.Phone
    dst.Email = from.Email
    dst.Openid = from.Openid

    return nil
}