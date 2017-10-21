/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/10/17 22:30 
 */
package controllers

import (
    "pet/protocol"
    "pet/models"
    "pet/utils"
)

// 用户电话注册
func UserPhoneRegist(args *protocol.UserPhoneRegistArgs, reply *protocol.UserPhoneRegistReply) error {
    var err error

    user_model := new(models.User)
    utils.DumpStruct(user_model, &args)

    err = user_model.Create(&reply.UserId)
    if nil != err {
        return err
    }
    return nil
}