/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/10/17 22:34 
 */
package protocol

import (
    "third/go-local"
)

// +++++++++++++++++++ 通用输出VO ++++++++++++++++++++++++++++

type UserInfoJson struct {
    UserId          int64           `json:"user_id"`
    Name            string          `json:"name"`       // 姓名
    Nickname        string          `json:"nickname"`   // 昵称
    Gender          string          `json:"gender"`     // 性别，0: 无性别 1: 男 2: 女
    Phone           string          `json:"phone"`
    Email           string          `json:"email"`
    Openid          string          `json:"openid"`     // 微信用户凭证
}

// ++++++++++++++++++++ 请求参数的数据格式 ++++++++++++++++++++++
// 电话注册
type UserPhoneRegistArgs struct {
    local.TraceParam

    Name            string          `json:"name"`
    Gender          string          `json:"gender"`     // 性别，0: 无性别 1: 男 2: 女
    Phone           string          `json:"phone"`
    Email           string          `json:"email"`
    Openid          string          `json:"openid"`     // 微信用户凭证
}
type UserPhoneRegistReply struct {
    User          UserInfoJson      `json:"user_info"`
}

// 用户昵称注册
type UserNicknameRegistArgs struct {
    local.TraceParam

    Nickname        string          `json:"nickname"`
    Password        string          `json:"password"`
}
type UserNicknameRegistReply struct {
    User          UserInfoJson           `json:"user_info"`
}

// 用户昵称登陆
type UserNicknameLoginArgs struct {
    local.TraceParam

    Nickname        string          `json:"nickname"`
    Password        string          `json:"password"`
}
type UserNicknameLoginReply struct {
    User          UserInfoJson           `json:"user_info"`
}