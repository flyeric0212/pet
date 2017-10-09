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
)

func StartHttpServer() {
    router := gin.Default()
    router.Use(GinSetTraceInfo())
    //router.Use(util.GinCrossDomain())

    router.GET("/wx_callback", WxCallbackHandler)
    router.POST("/wx_callback", WxCallbackHandler)

    router.Run(utils.Config.Listen)
    //router.Run(":80")
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