/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/10/22 15:48 
 */
package protocol

import (
    "third/go-local"
)

type BannerInfoJson struct {
    Id              int64           `json:"id"`
    Pic             string          `json:"pic"`
    RefUrl          string          `json:"ref_url"`
    Type            int             `json:"type"` // 类型，1:首页 2:展商 3: 合作媒体
}

// ++++++++++++++++++++ 请求参数的数据格式 ++++++++++++++++++++++

type BannerListArgs struct {
    local.TraceParam

    Type                int         `json:"type"`
    PageNum     		int			`json:"page_num"`
    PageSize    		int         `json:"page_size"`
}
type BannerListReply struct {
    BannerList          []BannerInfoJson    `json:"banner_list"`
    TotalNum		    int 			    `json:"total_num"`
}