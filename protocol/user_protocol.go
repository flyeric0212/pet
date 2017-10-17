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

type UserProfileInfoJson struct {
    UserId          int64           `json:"user_id"`
    Name            string          `json:"name"`
    Gender          string          `json:"gender"`     // 性别，0: 无性别 1: 男 2: 女
    Phone           string          `json:"phone"`
    Email           string          `json:"email"`
    Openid          string          `json:"openid"`     // 微信用户凭证
}

// ++++++++++++++++++++ 请求参数的数据格式 ++++++++++++++++++++++
type UserRegistArgs struct {
    local.TraceParam

    Name            string          `json:"name"`
    Gender          string          `json:"gender"`     // 性别，0: 无性别 1: 男 2: 女
    Phone           string          `json:"phone"`
    Email           string          `json:"email"`
    Openid          string          `json:"openid"`     // 微信用户凭证
}

type UserRegistReply struct {
    UserId          int64           `json:"user_id"`
}