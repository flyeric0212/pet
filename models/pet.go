/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/10/11 23:15 
 */
package models

// 宠物表
type Pet struct {
    Id              int64           `gorm:"primary_key"; sql:"AUTO_INCREMENT"`
    Name            string          `sql:"type:varchar(128)"`   // 宠物名称
}

func (pet *Pet) TableName() string {
    return "pet.pet"
}