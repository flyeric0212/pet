/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/10/17 22:30 
 */
package controller

import (
    "pet/protocol"
    "pet/model"
    "pet/utils"
)

// 用户电话注册
func UserPhoneRegist(args *protocol.UserPhoneRegistArgs, reply *protocol.UserPhoneRegistReply) error {
    utils.Logger.Info("[cmd:user_phone_regist] args: %+v", args)

    var err error

    // 参数校验
    if args.Name == "" || args.Phone == "" || args.RegistType == 0 {
        err = utils.NewInternalErrorByStr(utils.ParameterErrCode, "参数不全")
        utils.Logger.Error("UserPhoneRegist failed, param err: %s \n", err.Error())
        return err
    }

    if args.RegistType == 1 {
        if args.Openid == "" {
            err = utils.NewInternalErrorByStr(utils.ParameterErrCode, "openid不能为空")
            utils.Logger.Error("UserPhoneRegist failed, param err: %s \n", err.Error())
            return err
        }
    }

    // 判断电话号码是否存在
    err, ok, _ := model.CheckPhoneExist(args.Phone)
    if nil != err {
        return err
    }
    if ok {
        err = utils.NewInternalErrorByStr(utils.PhoneRepeatErrCode, "电话号码重复")
        utils.Logger.Error("UserPhoneRegist failed, err: %s \n", err.Error())
        return err
    }

    // 参数拷贝，插入数据库
    user_model := new(model.User)
    utils.DumpStruct(user_model, args)

    var user_id int64
    err = user_model.Create(&user_id)
    if nil != err {
        return err
    }

    // 复制数据，输出到api
    model.CopyUserData(user_model, &reply.User)

    return nil
}

// 发送验证码
func SendVerifyCode() {

}

// 验证验证码
func CheckVerifyCode() {

}

// 用户用户名注册
//func UserNicknameRegist(args *protocol.UserNicknameRegistArgs, reply *protocol.UserNicknameRegistReply) error {
//    utils.Logger.Info("[cmd:user_nickname_regist] args: %+v", args)
//
//    var err error
//
//    // 判断参数是否异常
//    if args.Nickname == "" || args.Password == "" {
//        err = utils.NewInternalErrorByStr(utils.ParameterErrCode, "用户名或密码为空")
//        utils.Logger.Error("UserNicknameRegist failed, param err: %s \n", err.Error())
//        return err
//    }
//
//    // 判断昵称是否存在
//    err, ok, _ := model.CheckNicknameExist(args.Nickname)
//    if nil != err {
//        return err
//    }
//    if ok {
//        err = utils.NewInternalErrorByStr(utils.NicknameRepeatErrCode, "用户名重复")
//        utils.Logger.Error("UserNicknameRegist failed, err: %s \n", err.Error())
//        return err
//    }
//
//    // 复制用户数据
//    user_model := new(model.User)
//    user_model.Nickname = args.Nickname
//    user_model.Password = utils.AesEncrypt(args.Password)
//    var user_id int64
//    err = user_model.Create(&user_id)
//    if nil != err {
//        return err
//    }
//    // 复制数据
//    model.CopyUserData(user_model, &reply.User)
//
//    return nil
//}
//
//func UserNicknameLogin(args *protocol.UserNicknameLoginArgs, reply *protocol.UserNicknameLoginReply) error {
//    utils.Logger.Info("[cmd:user_nickname_login] args: %+v", args)
//
//    var err error
//
//    // 判断参数是否异常
//    if args.Nickname == "" || args.Password == "" {
//        err = utils.NewInternalErrorByStr(utils.ParameterErrCode, "用户名或密码为空")
//        utils.Logger.Error("UserNicknameLogin failed, param err: %s \n", err.Error())
//        return err
//    }
//
//    // 判断昵称是否存在
//    err, ok, user_info := model.CheckNicknameExist(args.Nickname)
//    if nil != err {
//        return err
//    }
//    if !ok {
//        err = utils.NewInternalErrorByStr(utils.NicknameNotFoundErrCode, "用户名不存在")
//        utils.Logger.Error("UserNicknameLogin failed, err: %s \n", err.Error())
//        return err
//    }
//
//    // 判断密码
//    if args.Password != utils.AesDecrypt(user_info.Password) {
//        err = utils.NewInternalErrorByStr(utils.PasswordWrongErrCode, "密码错误")
//        utils.Logger.Error("UserNicknameLogin failed, err: %s \n", err.Error())
//        return err
//    }
//
//    user_info.LastLogin = time.Now()
//    err = model.PET_DB.Save(&user_info).Error
//    if nil != err {
//        err = utils.NewInternalError(utils.DbErrCode, err)
//        utils.Logger.Error("update user login time error: %v", err)
//        return err
//    }
//    // 复制数据
//    model.CopyUserData(user_info, &reply.User)
//
//    return nil
//}