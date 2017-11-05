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
    "strconv"
    "math/rand"
    "third/go-local"
    "fmt"
    "time"
)

// 用户电话注册
func UserPhoneRegist(args *protocol.UserPhoneRegistArgs, reply *protocol.UserPhoneRegistReply) error {
    local.TempTraceInfoArgs(args)
    defer local.Clear()

    utils.Logger.Info("[cmd:user_phone_regist] args: %+v", args)

    var err error

    // 参数校验
    if args.Name == "" || args.Phone == "" || args.RegistType == 0 {
        err = utils.NewInternalErrorByStr(utils.ParameterErrCode, "参数不全")
        utils.Logger.Error("UserPhoneRegist failed, param err: %s \n", err.Error())
        return err
    }

    // 电话号码校验
    if !utils.PhoneValid(args.Phone) {
        err = utils.NewInternalErrorByStr(utils.ParameterErrCode, "电话号码错误")
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

// 获取用户信息 by openid
func GetUserByOpenid(args *protocol.GetUserByOpenidArgs, reply *protocol.GetUserByOpenidReply) error {
    local.TempTraceInfoArgs(args)
    defer local.Clear()

    utils.Logger.Info("[cmd:get_user_by_openid] args: %+v", args)

    if args.Openid == "" {
        return nil
    }

    user_model := new(model.User)
    err := user_model.GetUserByOpenid(args.Openid)
    if nil != err {
        return err
    }
    // 复制数据，输出到api
    model.CopyUserData(user_model, &reply.User)

    return nil
}

func GenerateVerifyCode() string {
    var code = ""
    for i := 0; i < 4; i++ {
        code = code + strconv.Itoa(rand.Intn(9))
    }
    return code
}

// 发送验证码
func SendVerifyCode(args *protocol.SendVerifyCodeArgs, reply *protocol.SendVerifyCodeReply) error {
    local.TempTraceInfoArgs(args)
    defer local.Clear()

    utils.Logger.Info("[cmd:send_verify_code] args: %+v", args)

    var err error

    if args.Phone == "" || !utils.PhoneValid(args.Phone) {
        err = utils.NewInternalErrorByStr(utils.ParameterErrCode, "电话号码错误")
        utils.Logger.Error("SendVerifyCode failed, param err: %s \n", err.Error())
        return err
    }

    // redis 判断请求频率（5秒）
    res, err := g_cache.Get(args.Phone)
    if res != nil {
        err = utils.NewInternalErrorByStr(utils.HighFrequencyErrCode, "验证码请求频率太快")
        utils.Logger.Error("SendVerifyCode HighFrequencyErrCode: %v", err)
        return err
    }

    // redis 判断每天请求次数
    var times int
    res, err = g_cache.Get(fmt.Sprintf("%s:%s", args.Phone, time.Now().Format("2006-01-02")))
    if nil != err {
        times = 0
    } else {
        times, _ = strconv.Atoi(string(res))
    }
    if times > 10 {
        err = utils.NewInternalErrorByStr(utils.DayMaxTimeErrCode, "验证码超过每天次数")
        utils.Logger.Error("SendVerifyCode DayMaxTimeErrCode: %v", err)
        return err
    }

    // redis 设置验证码，过期时间15分钟
    code := GenerateVerifyCode()
    err = g_cache.Set(code + ":" + args.Phone, 0, 15*60)    // 10分钟
    if nil != err {
        utils.Logger.Error("set code cache err")
        err = utils.NewInternalError(utils.CacheErrCode, err)
        return err
    }

    // redis 设置
    err = g_cache.Set(args.Phone, code, 5)
    if nil != err {
        utils.Logger.Error("set phone expire key err")
        err = utils.NewInternalError(utils.CacheErrCode, err)
        return err
    }

    // TODO: 发送验证码
    err = utils.YpSendSms(args.Phone, code)
    if nil != err {
        utils.Logger.Error("SendVerifyCode call yunpian client failed, err: %v", err)
        return err
    }

    _, err = g_cache.Incr(fmt.Sprintf("%s:%s", args.Phone, time.Now().Format("2006-01-02")))
    if nil != err {
        utils.Logger.Error("set code max time err")
        err = utils.NewInternalError(utils.CacheErrCode, err)
        return err
    }
    g_cache.Expire(fmt.Sprintf("%s:%s", args.Phone, time.Now().Format("2006-01-02")), 24*3600)

    return nil
}

/**
 * 验证验证码
 *
 * state, 0: 可注册 ·1：微信已注册  2: 官网已注册
 */
func CheckVerifyCode(args *protocol.CheckVerifyCodeArgs, reply *protocol.CheckVerifyCodeReply) error {
    local.TempTraceInfoArgs(args)
    defer local.Clear()

    utils.Logger.Info("[cmd:check_verify_code] args: %+v", args)

    var err error
    if args.Phone == "" || !utils.PhoneValid(args.Phone) {
        err = utils.NewInternalErrorByStr(utils.ParameterErrCode, "电话号码错误")
        utils.Logger.Error("CheckVerifyCode failed, param err: %s \n", err.Error())
        return err
    }

    if args.VerifyCode == "" {
        err = utils.NewInternalErrorByStr(utils.ParameterErrCode, "验证码不能为空")
        utils.Logger.Error("CheckVerifyCode failed, param err: %s \n", err.Error())
        return err
    }

    res, err := g_cache.Get(args.VerifyCode + ":" + args.Phone)
    if nil != err {
        utils.Logger.Error("get verify code cache err: %v", err)
        err = utils.NewInternalError(utils.CacheErrCode, err)
        return err
    }

    if res == nil {
        err = utils.NewInternalErrorByStr(utils.VerifyCodeWrong, "验证码错误")
        return err
    }

    user_model := new(model.User)
    err = user_model.GetUserByPhone(args.Phone)
    if nil != err {
        return err
    }
    if user_model.UserId == 0 {
        reply.State = 0
    } else {
        if user_model.RegistType == 1 {
            reply.State = 1
        } else if user_model.RegistType == 2 {
            reply.State = 2
        }
        model.CopyUserData(user_model, &reply.User)
    }
    return nil
}