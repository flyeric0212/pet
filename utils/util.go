/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/9/24 21:28 
 */
package utils

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/md5"
    "crypto/rand"
    "errors"
    "io"
    "reflect"
    "encoding/hex"
    "regexp"
    ypclnt "github.com/yunpian/yunpian-go-sdk/sdk"
)

var (
    dunno     = []byte("???")
    centerDot = []byte("·")
    dot       = []byte(".")
    slash     = []byte("/")
)

var (
    AESBlock       cipher.Block
    ErrAESTextSize = errors.New("ciphertext is not a multiple of the block size")
    ErrAESPadding  = errors.New("cipher padding size err")

    // 云片sms
    ypApiKey            string
    ypClient            ypclnt.YunpianClient
)

var ses_res_chan chan struct{} = make(chan struct{}, 10) // aws ses send limit is 14 per second

const (
    AESTable = "shanghaipet1464345----1577808000"
)

func init() {
    var err error
    AESBlock, err = aes.NewCipher([]byte(AESTable))
    if err != nil {
        panic(err)
    }
    loggerInit()

    // init ses_res
    for i := 0; i < 10; i++ {
        ses_res_chan <- struct{}{}
    }
}

func InitYpClient() {
    // 云片client
    ypApiKey            = Config.External["ypApiKey"]
    ypClient            = ypclnt.New(ypApiKey)
}

// aes 加密
func AesEncrypt(srcStr string) string {
    src := []byte(srcStr)
    paddingLen := aes.BlockSize - (len(src) % aes.BlockSize)
    for i := 0; i < paddingLen; i++ {
        src = append(src, byte(paddingLen))
    }
    srcLen := len(src)

    encryptText := make([]byte, srcLen+aes.BlockSize)
    iv := encryptText[srcLen:]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        Logger.Error("aesencrypt err: %v", err)
        return ""
    }

    mode := cipher.NewCBCEncrypter(AESBlock, iv)
    mode.CryptBlocks(encryptText[:srcLen], src)
    return string(encryptText)
}

// ase 解密
func AesDecrypt(srcStr string) string {
    src := []byte(srcStr)
    // 长度不小于aes.Blocksize * 2
    if len(src) < aes.BlockSize*2 || len(src)%aes.BlockSize != 0 {
        Logger.Error("aesdecrypt err: %v", ErrAESTextSize)
        return ""
    }

    srcLen := len(src) - aes.BlockSize
    decyptText := make([]byte, srcLen)
    iv := src[srcLen:]
    mode := cipher.NewCBCDecrypter(AESBlock, iv)
    mode.CryptBlocks(decyptText, src[:srcLen])
    paddingLen := int(decyptText[srcLen-1])
    if paddingLen > 16 || paddingLen <= 0 {
        Logger.Error("aesdecrypt err: %v", ErrAESPadding)
        return ""
    }

    return string(decyptText[:(srcLen - paddingLen)])
}

func MD5(p []byte) string {
    sum := md5.Sum(p)
    return hex.EncodeToString(sum[:])
}


func IsStructPtr(t reflect.Type) bool {
    return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

//Be careful to use, from,to must be pointer
func DumpStruct(to interface{}, from interface{}) {
    fromv := reflect.ValueOf(from)
    tov := reflect.ValueOf(to)
    if fromv.Kind() != reflect.Ptr || tov.Kind() != reflect.Ptr {
        return
    }

    from_val := reflect.Indirect(fromv)
    to_val := reflect.Indirect(tov)

    for i := 0; i < from_val.Type().NumField(); i++ {
        fdi_from_val := from_val.Field(i)
        fd_name := from_val.Type().Field(i).Name
        fdi_to_val := to_val.FieldByName(fd_name)

        if fdi_to_val.IsValid() && fdi_from_val.Type() == fdi_to_val.Type() {
            fdi_to_val.Set(fdi_from_val)
        }
    }
}

// 校验电话号码
func PhoneValid(phone string) bool {
    var ret bool

    reg := `^1([38][0-9]|14[57]|5[^4])\d{8}$`
    rgx := regexp.MustCompile(reg)
    ret = rgx.MatchString(phone)

    return ret
}

// 云片发送验证码
func YpSendSms(phone, code string) error {
    param := ypclnt.NewParam(2)
    param[ypclnt.MOBILE] = phone
    param[ypclnt.TEXT] = "您的验证码是 " + code + "，15分钟内有效"
    result := ypClient.Sms().SingleSend(param)

    if result != nil && result.Code != 0 {
        Logger.Error("Yunpian send sms failed, result: %+v", result)
    }

    return nil
}