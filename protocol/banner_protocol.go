/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/10/22 15:48 
 */
package protocol

import (
    "third/go-local"
    "pet/model"
)

// ++++++++++++++++++++ 请求参数的数据格式 ++++++++++++++++++++++

type BannerListArgs struct {
    local.TraceParam

    Type                int         `json:"type"`
    PageNum     		int			`json:"page_num"`
    PageSize    		int         `json:"page_size"`
}
type BannerListReply struct {
    BannerList          []model.Banner  `json:"banner_list"`
    TotalNum		    int 			`json:"total_num"`
}