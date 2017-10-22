/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/10/21 21:11 
 */
package main

import (
    "third/gin"
    "time"
    "net/http"
    "pet/protocol"
    "pet/utils"
    "pet/controller"
)

// 用户电话注册
func UserPhoneRegist(c *gin.Context) {
    var http_code int = http.StatusOK
    handle_start_time := time.Now()

    var args protocol.UserPhoneRegistArgs
    var reply protocol.UserPhoneRegistReply

    r := c.Request
    err := utils.ParseHttpBodyToArgs(c, &args)
    if nil != err {
        goto NOTICE
    }
    err = controller.UserPhoneRegist(&args, &reply)

NOTICE:
    g_logger.Notice("[cmd:user_phone_regist][user_id:%s][Cost:%dus][Err:%v]",
        r.FormValue("user_id"), time.Now().Sub(handle_start_time).Nanoseconds()/1000, err)

    utils.SendResponse(c, http_code, &reply.User, err)
}

// 用户会员名注册
func UserNicknameRegist(c *gin.Context) {
    var http_code int = http.StatusOK
    handle_start_time := time.Now()

    var args protocol.UserNicknameRegistArgs
    var reply protocol.UserNicknameRegistReply

    r := c.Request
    err := utils.ParseHttpBodyToArgs(c, &args)
    if nil != err {
        goto NOTICE
    }
    err = controller.UserNicknameRegist(&args, &reply)

NOTICE:
    g_logger.Notice("[cmd:user_nickname_regist][user_id:%s][Cost:%dus][Err:%v]",
        r.FormValue("user_id"), time.Now().Sub(handle_start_time).Nanoseconds()/1000, err)

    utils.SendResponse(c, http_code, &reply.User, err)
}

// 用户会员名登陆
func UserNicknameLogin(c *gin.Context) {
    var http_code int = http.StatusOK
    handle_start_time := time.Now()

    var args protocol.UserNicknameLoginArgs
    var reply protocol.UserNicknameLoginReply

    r := c.Request
    err := utils.ParseHttpBodyToArgs(c, &args)
    if nil != err {
        goto NOTICE
    }
    err = controller.UserNicknameLogin(&args, &reply)

NOTICE:
    g_logger.Notice("[cmd:user_nickname_login][user_id:%s][Cost:%dus][Err:%v]",
        r.FormValue("user_id"), time.Now().Sub(handle_start_time).Nanoseconds()/1000, err)

    utils.SendResponse(c, http_code, &reply.User, err)
}

// 分页获取banner列表
func GetBannerListByPage(c *gin.Context) {
    var http_code int = http.StatusOK
    handle_start_time := time.Now()

    var args protocol.BannerListArgs
    var reply protocol.BannerListReply

    r := c.Request
    err := utils.ParseHttpBodyToArgs(c, &args)
    if nil != err {
        goto NOTICE
    }
    err = controller.GetBannerListByPage(&args, &reply)

NOTICE:
    g_logger.Notice("[cmd:banner_list_by_page][user_id:%s][Cost:%dus][Err:%v]",
        r.FormValue("user_id"), time.Now().Sub(handle_start_time).Nanoseconds()/1000, err)

    utils.SendResponse(c, http_code, &reply, err)
}

// 分页获取banner列表
func GetArticleListByPage(c *gin.Context) {
    var http_code int = http.StatusOK
    handle_start_time := time.Now()

    var args protocol.ArticleListArgs
    var reply protocol.ArticleListReply

    r := c.Request
    err := utils.ParseHttpBodyToArgs(c, &args)
    if nil != err {
        goto NOTICE
    }
    err = controller.GetArticleListByPage(&args, &reply)

NOTICE:
    g_logger.Notice("[cmd:article_list_by_page][user_id:%s][Cost:%dus][Err:%v]",
        r.FormValue("user_id"), time.Now().Sub(handle_start_time).Nanoseconds()/1000, err)

    utils.SendResponse(c, http_code, &reply, err)
}