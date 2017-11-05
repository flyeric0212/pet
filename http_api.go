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

    err := utils.ParseHttpBodyToArgs(c, &args)
    if nil != err {
        goto NOTICE
    }
    err = controller.UserPhoneRegist(&args, &reply)

NOTICE:
    g_logger.Notice("[cmd:user_phone_regist][Cost:%dus][Err:%v]",
        time.Now().Sub(handle_start_time).Nanoseconds()/1000, err)

    utils.SendResponse(c, http_code, &reply.User, err)
}

// 用户会员名注册
func GetUserByOpenid(c *gin.Context) {
    var http_code int = http.StatusOK
    handle_start_time := time.Now()

    var args protocol.GetUserByOpenidArgs
    var reply protocol.GetUserByOpenidReply

    err := utils.ParseHttpBodyToArgs(c, &args)
    if nil != err {
        goto NOTICE
    }
    err = controller.GetUserByOpenid(&args, &reply)

NOTICE:
    g_logger.Notice("[cmd:get_user_by_openid][Cost:%dus][Err:%v]",
        time.Now().Sub(handle_start_time).Nanoseconds()/1000, err)

    utils.SendResponse(c, http_code, &reply.User, err)
}

// 发送验证码
func SendVerifyCode(c *gin.Context) {
    var http_code int = http.StatusOK
    handle_start_time := time.Now()

    var args protocol.SendVerifyCodeArgs
    var reply protocol.SendVerifyCodeReply

    err := utils.ParseHttpBodyToArgs(c, &args)
    if nil != err {
        goto NOTICE
    }
    err = controller.SendVerifyCode(&args, &reply)

NOTICE:
    g_logger.Notice("[cmd:send_verify_code][Cost:%dus][Err:%v]",
        time.Now().Sub(handle_start_time).Nanoseconds()/1000, err)

    utils.SendResponse(c, http_code, &reply, err)
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

    err := utils.ParseHttpBodyToArgs(c, &args)
    if nil != err {
        goto NOTICE
    }
    err = controller.GetArticleListByPage(&args, &reply)

NOTICE:
    g_logger.Notice("[cmd:article_list_by_page][Cost:%dus][Err:%v]",
        time.Now().Sub(handle_start_time).Nanoseconds()/1000, err)

    utils.SendResponse(c, http_code, &reply, err)
}

// 观众中心授权
func VistorCenterAuth(c *gin.Context) {
    AuthPage(c.Writer, c.Request)
}

// 观众中心授权回调页面
func VistorCenterRedirect(c *gin.Context) {
    AuthCallback(c.Writer, c.Request)
}