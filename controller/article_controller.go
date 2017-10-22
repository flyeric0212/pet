/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/10/22 19:26 
 */
package controller

import (
    "pet/protocol"
    "pet/model"
    "pet/utils"
)

// 分页获取
func GetArticleListByPage(args *protocol.ArticleListArgs, reply *protocol.ArticleListReply) error {
    utils.Logger.Info("[cmd:article_list_by_page] args: %+v", args)

    if args.PageNum <= 0 {
        args.PageNum = 1
    }
    if args.PageSize <= 0 {
        args.PageSize = 10
    }

    article_model := new(model.Article)
    article_list, total_num, err := article_model.GetArticleListByPage(args.Type, args.PageNum, args.PageSize)
    if nil != err {
        return err
    }
    if article_list == nil || len(article_list) == 0 {
        reply.ArticleList = make([]protocol.ArticleInfoJson, 0)
    } else {
        reply.ArticleList = make([]protocol.ArticleInfoJson, len(article_list))
        // format data
        for i := range article_list {
            utils.DumpStruct(&reply.ArticleList[i], &article_list[i])
        }
    }
    reply.TotalNum = total_num

    return nil
}
