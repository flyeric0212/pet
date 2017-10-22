/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/10/22 19:18 
 */
package protocol

import "third/go-local"

type ArticleInfoJson struct {
    Id              int64           `json:"id"`
    Title           string          `json:"title"`
    Content         string          `json:"content"`
    Type            int             `json:"type"` // 类型，1:会展动态
}

// ++++++++++++++++++++ 请求参数的数据格式 ++++++++++++++++++++++

type ArticleListArgs struct {
    local.TraceParam

    Type                int         `json:"type"`
    PageNum     		int			`json:"page_num" mapstructure:"page_num"`
    PageSize    		int         `json:"page_size" mapstructure:"page_size"`
}
type ArticleListReply struct {
    ArticleList         []ArticleInfoJson       `json:"article_list"`
    TotalNum		    int 			        `json:"total_num"`
}
