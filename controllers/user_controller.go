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

// 用户注册
func UserRegist(args *protocol.UserRegistArgs, reply *protocol.UserRegistReply) error {
    var err error

    user_model := new(models.UserProfile)
    utils.DumpStruct(user_model, &args)

    err = user_model.Create(&reply.UserId)
    if nil != err {
        return err
    }
    return nil
}