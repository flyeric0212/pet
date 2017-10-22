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
    "time"
)

// 用户电话注册
func UserPhoneRegist(args *protocol.UserPhoneRegistArgs, reply *protocol.UserPhoneRegistReply) error {
    utils.Logger.Info("[cmd:user_phone_regist] args: %+v", args)

    var err error

    user_model := new(models.User)
    utils.DumpStruct(user_model, &args)

    var user_id int64
    err = user_model.Create(&user_id)
    if nil != err {
        return err
    }
    // 复制数据
    models.CopyUserData(user_model, &reply.User)

    return nil
}

// 用户用户名注册
func UserNicknameRegist(args *protocol.UserNicknameRegistArgs, reply *protocol.UserNicknameRegistReply) error {
    utils.Logger.Info("[cmd:user_nickname_regist] args: %+v", args)

    var err error

    // 判断参数是否异常
    if args.Nickname == "" || args.Password == "" {
        err = utils.NewInternalErrorByStr(utils.ParameterErrCode, "用户名或密码为空")
        utils.Logger.Error("UserNicknameRegist failed, param err: %s \n", err.Error())
        return err
    }

    // 判断昵称是否存在
    err, ok, _ := models.CheckNicknameExist(args.Nickname)
    if nil != err {
        return err
    }
    if ok {
        err = utils.NewInternalErrorByStr(utils.NicknameRepeatErrCode, "用户名重复")
        utils.Logger.Error("UserNicknameRegist failed, err: %s \n", err.Error())
        return err
    }

    // 复制用户数据
    user_model := new(models.User)
    user_model.Nickname = args.Nickname
    user_model.Password = utils.AesEncrypt(args.Password)
    var user_id int64
    err = user_model.Create(&user_id)
    if nil != err {
        return err
    }
    // 复制数据
    models.CopyUserData(user_model, &reply.User)

    return nil
}

func UserNicknameLogin(args *protocol.UserNicknameLoginArgs, reply *protocol.UserNicknameLoginReply) error {
    utils.Logger.Info("[cmd:user_nickname_login] args: %+v", args)

    var err error

    // 判断参数是否异常
    if args.Nickname == "" || args.Password == "" {
        err = utils.NewInternalErrorByStr(utils.ParameterErrCode, "用户名或密码为空")
        utils.Logger.Error("UserNicknameLogin failed, param err: %s \n", err.Error())
        return err
    }

    // 判断昵称是否存在
    err, ok, user_info := models.CheckNicknameExist(args.Nickname)
    if nil != err {
        return err
    }
    if !ok {
        err = utils.NewInternalErrorByStr(utils.NicknameNotFoundErrCode, "用户名不存在")
        utils.Logger.Error("UserNicknameLogin failed, err: %s \n", err.Error())
        return err
    }

    // 判断密码
    if args.Password != utils.AesDecrypt(user_info.Password) {
        err = utils.NewInternalErrorByStr(utils.PasswordWrongErrCode, "密码错误")
        utils.Logger.Error("UserNicknameLogin failed, err: %s \n", err.Error())
        return err
    }

    user_info.LastLogin = time.Now()
    err = models.PET_DB.Save(&user_info).Error
    if nil != err {
        err = utils.NewInternalError(utils.DbErrCode, err)
        utils.Logger.Error("update user login time error: %v", err)
        return err
    }
    // 复制数据
    models.CopyUserData(user_info, &reply.User)

    return nil
}