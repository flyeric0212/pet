/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/10/1 11:01 
 */
package main

import (
    "github.com/chanxuehong/wechat.v2/mp/core"
    "github.com/chanxuehong/wechat.v2/mp/message/callback/request"
    "github.com/chanxuehong/wechat.v2/mp/menu"
    "github.com/chanxuehong/wechat.v2/mp/message/callback/response"
    "third/gin"
    "pet/utils"
)

const (
    //wxAppId       = "appid"
    //wxAppSecret   = "appsecret"
    //wxOriId       = "oriid"

    wxToken         = "12345678shanghaipet"
    wxEncodedAESKey = "mzTODOLGqD2aVlw53HEOMc6qMYMG4UqXst0FxHzRr2z"
)

var (
    // 下面两个变量不一定非要作为全局变量, 根据自己的场景来选择.

    msgHandler core.Handler
    msgServer  *core.Server

    accessTokenServer core.AccessTokenServer
)

func InitWeixinServer() {
    mux := core.NewServeMux()
    // 默认处理
    mux.DefaultMsgHandleFunc(defaultMsgHandler)
    mux.DefaultEventHandleFunc(defaultEventHandler)

    // 处理文本消息
    mux.MsgHandleFunc(request.MsgTypeText, textMsgHandler)
    // 处理菜单点击
    mux.EventHandleFunc(menu.EventTypeClick, menuClickEventHandler)
    msgHandler = mux

    var wxAppId         = utils.Config.External["AppId"]
    var wxAppSecret     = utils.Config.External["AppSecret"]
    var wxOriId         = utils.Config.External["OriId"]

    msgServer           = core.NewServer(wxOriId, wxAppId, wxToken, wxEncodedAESKey, msgHandler, nil)
    accessTokenServer   = core.NewDefaultAccessTokenServer(wxAppId, wxAppSecret, nil)
}

func defaultMsgHandler(ctx *core.Context) {
    g_logger.Info("收到消息:\n%s\n", ctx.MsgPlaintext)
    ctx.NoneResponse()
}

func defaultEventHandler(ctx *core.Context) {
    g_logger.Info("收到事件:\n%s\n", ctx.MsgPlaintext)
    ctx.NoneResponse()
}

// 文本消息处理
func textMsgHandler(ctx *core.Context) {
    g_logger.Info("收到文本消息:\n%s\n", ctx.MsgPlaintext)

    msg := request.GetText(ctx.MixedMsg)
    resp := response.NewText(msg.FromUserName, msg.ToUserName, msg.CreateTime, msg.Content)
    ctx.RawResponse(resp) // 明文回复
    //ctx.AESResponse(resp, 0, "", nil) // aes密文回复
}

//
func menuClickEventHandler(ctx *core.Context) {
    g_logger.Info("收到菜单 click 事件:\n%s\n", ctx.MsgPlaintext)

    event := menu.GetClickEvent(ctx.MixedMsg)
    resp := response.NewText(event.FromUserName, event.ToUserName, event.CreateTime, "收到 click 类型的事件")
    ctx.RawResponse(resp) // 明文回复
    //ctx.AESResponse(resp, 0, "", nil) // aes密文回复
}



// wxCallbackHandler 是处理回调请求的 http handler.
//  1. 不同的 web 框架有不同的实现
//  2. 一般一个 handler 处理一个公众号的回调请求(当然也可以处理多个, 这里我只处理一个)
func WxCallbackHandler(c *gin.Context) {
    query := c.Request.URL.Query()
    g_logger.Info("query: %+v \n", query)
    msgServer.ServeHTTP(c.Writer, c.Request, nil)
}
