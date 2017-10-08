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
)

const (
    wxAppId     = "appid"
    wxAppSecret = "appsecret"

    wxOriId         = "oriid"
    wxToken         = "12345678shanghaipet"
    wxEncodedAESKey = "mzTODOLGqD2aVlw53HEOMc6qMYMG4UqXst0FxHzRr2z"
)

var (
    // 下面两个变量不一定非要作为全局变量, 根据自己的场景来选择.
    msgHandler core.Handler
    msgServer  *core.Server
)

func init() {
    mux := core.NewServeMux()
    mux.DefaultMsgHandleFunc(defaultMsgHandler)
    mux.DefaultEventHandleFunc(defaultEventHandler)
    mux.MsgHandleFunc(request.MsgTypeText, textMsgHandler)
    mux.EventHandleFunc(menu.EventTypeClick, menuClickEventHandler)

    msgHandler = mux
    msgServer = core.NewServer(wxOriId, wxAppId, wxToken, wxEncodedAESKey, msgHandler, nil)
}

func textMsgHandler(ctx *core.Context) {
    g_logger.Info("收到文本消息:\n%s\n", ctx.MsgPlaintext)

    msg := request.GetText(ctx.MixedMsg)
    resp := response.NewText(msg.FromUserName, msg.ToUserName, msg.CreateTime, msg.Content)
    ctx.RawResponse(resp) // 明文回复
    //ctx.AESResponse(resp, 0, "", nil) // aes密文回复
}

func defaultMsgHandler(ctx *core.Context) {
    g_logger.Info("收到消息:\n%s\n", ctx.MsgPlaintext)
    ctx.NoneResponse()
}

func menuClickEventHandler(ctx *core.Context) {
    g_logger.Info("收到菜单 click 事件:\n%s\n", ctx.MsgPlaintext)

    event := menu.GetClickEvent(ctx.MixedMsg)
    resp := response.NewText(event.FromUserName, event.ToUserName, event.CreateTime, "收到 click 类型的事件")
    ctx.RawResponse(resp) // 明文回复
    //ctx.AESResponse(resp, 0, "", nil) // aes密文回复
}

func defaultEventHandler(ctx *core.Context) {
    g_logger.Info("收到事件:\n%s\n", ctx.MsgPlaintext)
    ctx.NoneResponse()
}

// wxCallbackHandler 是处理回调请求的 http handler.
//  1. 不同的 web 框架有不同的实现
//  2. 一般一个 handler 处理一个公众号的回调请求(当然也可以处理多个, 这里我只处理一个)
func WxCallbackHandler(c *gin.Context) {
    msgServer.ServeHTTP(c.Writer, c.Request, nil)
}
