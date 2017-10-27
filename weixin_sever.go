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

    mpoauth2 "github.com/chanxuehong/wechat.v2/mp/oauth2"
    "github.com/chanxuehong/wechat.v2/oauth2"
    "github.com/chanxuehong/session"
    "net/http"
    "github.com/chanxuehong/sid"
    "io"
    "github.com/chanxuehong/rand"
    "fmt"
    "net/url"
)

const (
    wxToken             = "12345678shanghaipet"
    wxEncodedAESKey     = "mzTODOLGqD2aVlw53HEOMc6qMYMG4UqXst0FxHzRr2z"

    oauth2RedirectURI   = "http://mp.petfair.cc/api/vistor_center_callback"
    oauth2Scope         = "snsapi_userinfo"

    vistorCenterHomeURI = "http://wx.petfair.cc/"
)

var (
    wxAppId string
    wxAppSecret string
    wxOriId string

    msgHandler core.Handler
    msgServer  *core.Server

    accessTokenServer core.AccessTokenServer
    wechatClient      *core.Client

    oauth2Endpoint oauth2.Endpoint
    sessionStorage      = session.New(20*60, 60*60)
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

    wxAppId         = utils.Config.External["AppId"]
    wxAppSecret     = utils.Config.External["AppSecret"]
    wxOriId         = utils.Config.External["OriId"]

    msgServer           = core.NewServer(wxOriId, wxAppId, wxToken, wxEncodedAESKey, msgHandler, nil)

    accessTokenServer   = core.NewDefaultAccessTokenServer(wxAppId, wxAppSecret, nil)
    wechatClient        = core.NewClient(accessTokenServer, nil)

    oauth2Endpoint      = mpoauth2.NewEndpoint(wxAppId, wxAppSecret)
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

func AuthPage(w http.ResponseWriter, r *http.Request) {
    sid := sid.New()
    state := string(rand.NewHex())

    if err := sessionStorage.Add(sid, state); err != nil {
        io.WriteString(w, err.Error())
        utils.Logger.Error("add session err: %v", err)
        return
    }

    cookie := http.Cookie{
        Name:       "sid",
        Value:      sid,
        HttpOnly:   true,
    }
    http.SetCookie(w, &cookie)

    AuthCodeURL := mpoauth2.AuthCodeURL(wxAppId, oauth2RedirectURI, oauth2Scope, state)
    utils.Logger.Info("AuthCodeURL: %s", AuthCodeURL)

    http.Redirect(w, r, AuthCodeURL, http.StatusFound)
}

func AuthCallback(w http.ResponseWriter, r *http.Request) {
    utils.Logger.Info(r.RequestURI)

    cookie, err := r.Cookie("sid")
    if err != nil {
        io.WriteString(w, err.Error())
        utils.Logger.Error("get cookie err: %v", err)
        return
    }

    session, err := sessionStorage.Get(cookie.Value)
    if err != nil {
        io.WriteString(w, err.Error())
        utils.Logger.Error("get session err: %v", err)
        return
    }

    savedState := session.(string) // 一般是要序列化的, 这里保存在内存所以可以这么做

    queryValues, err := url.ParseQuery(r.URL.RawQuery)
    if err != nil {
        io.WriteString(w, err.Error())
        utils.Logger.Error("parse url err: %v", err)
        return
    }

    code := queryValues.Get("code")
    if code == "" {
        utils.Logger.Error("用户禁止授权")
        return
    }

    queryState := queryValues.Get("state")
    if queryState == "" {
        utils.Logger.Error("state 参数为空")
        return
    }
    if savedState != queryState {
        str := fmt.Sprintf("state 不匹配, session 中的为 %q, url 传递过来的是 %q", savedState, queryState)
        io.WriteString(w, str)
        utils.Logger.Error(str)
        return
    }

    oauth2Client := oauth2.Client{
        Endpoint: oauth2Endpoint,
    }
    token, err := oauth2Client.ExchangeToken(code)
    if err != nil {
        io.WriteString(w, err.Error())
        utils.Logger.Error("exchange token err: %v", err)
        return
    }
    utils.Logger.Info("token: %+v\r\n", token)

    userinfo, err := mpoauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
    if err != nil {
        io.WriteString(w, err.Error())
        utils.Logger.Error("get weixin user info err: %v", err)
        return
    }
    utils.Logger.Info("userinfo: %+v\r\n", userinfo)

    vistorCenterURL := vistorCenterHomeURI + "?openid=" + url.QueryEscape(userinfo.OpenId) +
        "&nickname=" + url.QueryEscape(userinfo.Nickname) + "&headurl=" + url.QueryEscape(userinfo.HeadImageURL)

    http.Redirect(w, r, vistorCenterURL, http.StatusFound)

    return
}
