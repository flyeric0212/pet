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
    Nickname        string          `json:"nickname"`   // 微信昵称
    Avatar          string          `json:"avatar"`     // 微信头像
    Gender          string          `json:"gender"`     // 性别，0: 无性别 1: 男 2: 女
    Phone           string          `json:"phone"`
    Email           string          `json:"email"`
    Openid          string          `json:"openid"`     // 微信用户凭证
}

// ++++++++++++++++++++ 请求参数的数据格式 ++++++++++++++++++++++
// 电话注册参数
type UserPhoneRegistArgs struct {
    local.TraceParam

    Name            string              `json:"name"`
    Gender          string              `json:"gender"`         // 性别，0: 无性别 1: 男 2: 女
    Phone           string              `json:"phone"`
    Email           string              `json:"email"`

    RegistType      int                 `json:"regist_type"`    // 1: 微信   2：官网
    Nickname        string              `json:"nickname"`       // 微信昵称
    Avatar          string              `json:"avatar"`         // 微信头像
    Openid          string              `json:"openid"`         // 微信用户凭证
}
type UserPhoneRegistReply struct {
    User            UserInfoJson        `json:"user_info"`
}

// openid获取用户信息
type GetUserByOpenidArgs struct {
    local.TraceParam

    Openid          string              `json:"openid"`
}
type GetUserByOpenidReply struct {
    User            UserInfoJson        `json:"user_info"`
}

// 发送验证码
type SendVerifyCodeArgs struct {
    local.TraceParam

    Phone           string              `json:""`
}
type SendVerifyCodeReply struct {

}