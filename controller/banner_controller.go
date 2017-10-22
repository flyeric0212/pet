/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/10/22 15:49 
 */
package controller

import (
    "pet/protocol"
    "pet/utils"
    "pet/model"
)

// 分页获取
func GetBannerListByPage(args *protocol.BannerListArgs, reply *protocol.BannerListReply) error {
    utils.Logger.Info("[cmd:banner_list_by_page] args: %+v", args)

    if args.PageNum <= 0 {
        args.PageNum = 1
    }
    if args.PageSize <= 0 {
        args.PageSize = 10
    }

    banner_model := new(model.Banner)
    banner_list, total_num, err := banner_model.GetBannerListByPage(args.Type, args.PageNum, args.PageSize)
    if nil != err {
        return err
    }

    reply.BannerList = banner_list
    reply.TotalNum = total_num

    return nil
}
