/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/9/24 21:01 
 */
package models

import (
    "third/gorm"

    "blast/common/util"
    "fmt"
)

var PET_DB *gorm.DB

func InitAllDB() error {
    var err error

    PET_DB, err = InitPetDb(util.CommonConfig)
    if nil != err {
        util.Logger.Warning("open pet db error")
        return err
    }

    return nil
}

func InitPetDb(config *util.Configure) (*gorm.DB, error) {
    var err error
    PetDbSetting, ok := config.MysqlSetting["PetDbSetting"]
    if !ok {
        return nil, fmt.Errorf("can't find pet mysql settings")
    }
    PET_DB, err = util.InitGormDbPool(&PetDbSetting)
    if err != nil {
        util.Logger.Warning("open pet DB error")
        return nil, err
    }
    return PET_DB, nil
}