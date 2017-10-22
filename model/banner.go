/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/10/22 15:28 
 */
package model

import (
    "time"
    "pet/utils"
)

// banner表
type Banner struct {
    Id              int64           `gorm:"primary_key"; sql:"AUTO_INCREMENT" json:"id"`
    Pic             string          `sql:"type:varchar(255)" json:"pic"`
    RefUrl          string          `sql:"type:varchar(255)" json:"ref_url"`
    Type            int             `sql:"type:smallint(6)" json:"type"` // 类型，1:首页 2:展商 3: 合作媒体
    CreateTime      time.Time       `sql:"type:datetime" json:"create_time"`
}

func (banner *Banner) TableName() string {
    return "pet.banner"
}

func (banner *Banner) Create() error {
    banner.CreateTime = time.Now()

    err := PET_DB.Table(banner.TableName()).Create(banner).Error
    if nil != err {
        err = utils.NewInternalError(utils.DbErrCode, err)
        utils.Logger.Error("create banner error: %v", err)
        return err
    }
    return nil
}

func (banner *Banner) Save() error {
    if banner.Id == 0 {
        banner.CreateTime = time.Now()
    }
    err := PET_DB.Table(banner.TableName()).Save(banner).Error
    if nil != err {
        err = utils.NewInternalError(utils.DbErrCode, err)
        utils.Logger.Error("save banner error: %v", err)
        return err
    }
    return nil
}

// 分页列表
func (banner *Banner) GetBannerListByPage(banner_type int, page_num, page_size int) (banner_list []Banner, total_num int, err error) {
    if page_size < 0 {
        page_size = 10
    }

    offset := (page_num - 1) * page_size
    if offset < 0 {
        offset = 0
    }

    query := PET_DB.Table(banner.TableName())
    if banner_type != 0 {
        query = query.Where("type = ?", banner_type)
    }
    if err2 := query.Count(&total_num).Error; nil != err2 {
        utils.Logger.Error("count banner list err: %v", err2)
        err = utils.NewInternalError(utils.DbErrCode, err2)
        return
    }
    query = query.Order("create_time desc").Limit(page_size).Offset(offset)

    err = query.Find(&banner_list).Error
    if nil != err {
        utils.Logger.Error("get banner list by page error :%s\n", err.Error())
        err = utils.NewInternalError(utils.DbErrCode, err)
        return
    }
    return
}