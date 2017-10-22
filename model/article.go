/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/10/22 18:21 
 */
package model

import (
    "time"
    "pet/utils"
)

// 文章
type Article struct {
    Id              int64           `gorm:"primary_key"; sql:"AUTO_INCREMENT"`
    Title           string          `sql:"type:varchar(128)"`
    Content         string          `sql:"type:text"`
    Type            int             `sql:"type:smallint(6)"` // 类型，1:展会动态
    CreateTime      time.Time       `sql:"type:datetime"`
}

func (article *Article) TableName() string {
    return "pet.article"
}

func (article *Article) Create() error {
    article.CreateTime = time.Now()

    err := PET_DB.Table(article.TableName()).Create(article).Error
    if nil != err {
        err = utils.NewInternalError(utils.DbErrCode, err)
        utils.Logger.Error("create article error: %v", err)
        return err
    }
    return nil
}

func (article *Article) Save() error {
    if article.Id == 0 {
        article.CreateTime = time.Now()
    }

    err := PET_DB.Table(article.TableName()).Save(article).Error
    if nil != err {
        err = utils.NewInternalError(utils.DbErrCode, err)
        utils.Logger.Error("save article error: %v", err)
        return err
    }
    return nil
}

func (article *Article) GetArticleListByPage(article_type int, page_num, page_size int) (article_list []Article, total_num int, err error) {
    if page_size < 0 {
        page_size = 10
    }

    offset := (page_num - 1) * page_size
    if offset < 0 {
        offset = 0
    }

    query := PET_DB.Table(article.TableName())
    if article_type != 0 {
        query = query.Where("type = ?", article_type)
    }
    if err2 := query.Count(&total_num).Error; nil != err2 {
        utils.Logger.Error("count article list err: %v", err2)
        err = utils.NewInternalError(utils.DbErrCode, err2)
        return
    }

    query = query.Order("create_time desc").Limit(page_size).Offset(offset)

    err = query.Find(&article_list).Error
    if nil != err {
        utils.Logger.Error("get article list by page error :%s\n", err.Error())
        err = utils.NewInternalError(utils.DbErrCode, err)
        return
    }
    return
}
