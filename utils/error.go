/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/9/24 20:41
 */
package utils

import (
    "fmt"
    "strings"
    "strconv"
)

type ErrCode int

const (
    //需要返回http500的错误码
    OKCode            		ErrCode = 0
    InternalErrorCode 		ErrCode = 1
    DbErrCode         		ErrCode = 2
    CacheErrCode      		ErrCode = 3
    DecodeErrCode     		ErrCode = 4

    //400以后禁用
    UserNotFoundCode    	ErrCode = 100


    //直接抛给用户的错误码
    ParameterErrCode      	ErrCode = 501 	// 参数异常
    PhoneRepeatErrCode      ErrCode = 502   // 电话号码重复

    //NicknameRepeatErrCode   ErrCode = 502   // 昵称重复
    //NicknameNotFoundErrCode ErrCode = 503   // 昵称不存在
    //PasswordWrongErrCode    ErrCode = 504   // 密码错误

    MaxUserError 			ErrCode = 9999
)

type InternalError struct {
    Code       ErrCode
    Error_info error
}

func NewInternalError(code ErrCode, err error) error {
    //if code < UserNotFoundCode && code > ErrCode(0) {
    //    stack := stack(3)
    //    Logger.Error("PANIC: %v\n%v", err, string(stack))
    //    CheckError(&InternalError{Code: ErrCode(code), Error_info: err})
    //}
    //panic(err)
    return &InternalError{Code: ErrCode(code), Error_info: err}
}

func NewInternalErrorByStr(code ErrCode, err string) error {
    //if code < UserNotFoundCode && code > ErrCode(0) {
    //    CheckError(&InternalError{Code: ErrCode(code), Error_info: fmt.Errorf(err)})
    //}
    //panic(err)
    return &InternalError{Code: ErrCode(code), Error_info: fmt.Errorf(err)}
}

func NewInternalErrByStrDefault(code ErrCode) error {
    var err_str string
    //var ok bool
    //if err_str, ok = ErrStr[code]; !ok {
    //    err_str = fmt.Sprintf("error code is %d", code)
    //}
    //if code < UserNotFoundCode && code > ErrCode(0) {
    //    CheckError(&InternalError{Code: ErrCode(code), Error_info: fmt.Errorf(err_str)})
    //}
    //panic(err)
    return &InternalError{Code: ErrCode(code), Error_info: fmt.Errorf(err_str)}
}

func (err *InternalError) Error() string {
    return fmt.Sprintf("%d:%s", err.Code, err.Error_info)
}

func IsUserErr(original_err error) (bool, int, string) {
    infos := strings.Split(original_err.Error(), ":")
    if len(infos) < 2 {
        return false, 0, infos[0]
    }
    code, err := strconv.Atoi(infos[0])
    if nil != err {
        return false, 0, infos[1]
    }
    if code >= int(UserNotFoundCode) && code <= int(MaxUserError) {
        return true, code, infos[1]
    }
    return false, 0, infos[1]
}

func GetErrInfo(original_err error) string {
    if original_err == nil {
        return ""
    }
    var code int
    var info string
    fmt.Sscanf(original_err.Error(), "%d:%s", &code, &info)
    return info
}

func GetErrCode(original_err error) ErrCode {
    if original_err == nil {
        return 0
    }
    var code int
    var info string
    fmt.Sscanf(original_err.Error(), "%d:%s", &code, &info)
    return ErrCode(code)
}