/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/9/24 21:08 
 */
package main

import (
    "third/gin"
    "pet/utils"
    "os"
    "time"
    "third/go-local"
    "sync/atomic"
    "strconv"
    "math/rand"
    "fmt"
    "io/ioutil"
)

func StartHttpServer() {
    router := gin.Default()
    router.Use(GinSetTraceInfo())
    router.Use(utils.GinFilter())
    router.Use(utils.GinRecovery())
    router.Use(utils.GinLogger())
    //router.Use(util.GinCrossDomain())

    // weixin
    router.GET("/wx_callback", WxCallbackHandler)
    router.POST("/wx_callback", WxCallbackHandler)

    // user
    user_router := router.Group("/api/users")
    user_router.POST("/phone_regist", UserPhoneRegist)
    user_router.GET("/get_by_openid", GetUserByOpenid)
    user_router.POST("/send_verify_code", SendVerifyCode)

    // banner
    banner_router := router.Group("/api/banner")
    banner_router.GET("/get_banner_list", GetBannerListByPage)

    // article
    article_router := router.Group("/api/article")
    article_router.GET("get_article_list", GetArticleListByPage)


    // weixin homepage
    router.GET("/api/vistor_center_auth", VistorCenterAuth)
    router.GET("/api/vistor_center_callback", VistorCenterRedirect)

    router.GET("/MP_verify_nzlgoroX2jUMUlfT.txt", DownloadWinxinValidFile)

    router.Run(utils.Config.Listen)
}

// 访问文件
func DownloadWinxinValidFile(c *gin.Context) {
    filename := "/var/config/MP_verify_nzlgoroX2jUMUlfT.txt"

    _, err := os.Stat(filename)
    if err != nil {
        fmt.Println("ReadFile: ", err.Error())
        return
    }

    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        fmt.Println("ReadFile: ", err.Error())
        return
    }
    c.Writer.Write(bytes)

    return
}


// update the method of generating trace_id
func GinSetTraceInfo() gin.HandlerFunc {
    var pid = uint64(os.Getpid()) & 0xFFFF // max: 0xFFFF
    var incId = uint64(rand.NewSource(time.Now().Unix()).Int63()) & 0xFFFFF
    return func(c *gin.Context) {
        var traceId string
        if traceId = local.ParseTraceParamHttp(c.Request).TraceId; traceId != "" {
            g_logger.Info("Trace id already setted as %q", traceId)
            // It is the setter's responsibility to clear the trace info
        } else {
            traceIdVal := (pid << 56) |
                (uint64(time.Now().Unix()) << (7 + 16)) |
                (atomic.AddUint64(&incId, 1) & 0x7FFFFF)
            traceId = strconv.FormatUint(traceIdVal, 36)
            c.Request.Header.Set(local.TraceIdFieldName, traceId)

        }

        local.TempTraceInfo(traceId, 0, 0)
        g_logger.Info("Set trace id: %s", traceId)

        c.Next()

        local.Clear()
        g_logger.Info("Clear trace id: %s", traceId)
    }
}