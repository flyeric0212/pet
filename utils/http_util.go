/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/9/24 21:14 
 */
package utils

import (
    "bytes"
    //"crypto/sha256"
    //"crypto/tls"
    //"encoding/hex"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    //"os"
    "reflect"
    "regexp"
    "strconv"
    "strings"
    "time"

    "third/gin"
    "third/go-local"
    "third/http_client_cluster"
    "third/httprouter"
)

const DEFAULT_API_TIMEOUT = 1 * time.Second

var g_num int

type CodoonApiResponse struct {
    Status string      `json:"status"`
    Data   interface{} `json:"data"`
    Desc   string      `json:"desc"`
}

func (this CodoonApiResponse) MarshalJSON() ([]byte, error) {
    return json.Marshal(map[string]interface{}{
        "status": this.Status,
        "data":   this.Data,
        "desc":   this.Desc,
    })
}

func NewApiResponse(status string, data interface{}, desc string) CodoonApiResponse {
    var resp CodoonApiResponse = CodoonApiResponse{"OK", nil, "success"}
    return resp
}

func CodoonGetHeader(c *gin.Context) {
    // 获取token
    r := c.Request
    Logger.Info("+++++++++++request header: %+v", r.Header)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    //io.WriteString(w, "404")
    http.Error(w, "404 page not found", http.StatusNotFound)

}

var sliceOfInts = reflect.TypeOf([]int(nil))
var sliceOfStrings = reflect.TypeOf([]string(nil))

// parse form values to struct via tag.
func ParseForm(form url.Values, obj interface{}) error {
    objT := reflect.TypeOf(obj)
    objV := reflect.ValueOf(obj)
    if !IsStructPtr(objT) {
        return fmt.Errorf("%v must be  a struct pointer", obj)
    }
    objT = objT.Elem()
    objV = objV.Elem()

    for i := 0; i < objT.NumField(); i++ {
        fieldV := objV.Field(i)
        if !fieldV.CanSet() {
            continue
        }

        fieldT := objT.Field(i)
        tags := strings.Split(fieldT.Tag.Get("form"), ",")
        var tag string
        if len(tags) == 0 || len(tags[0]) == 0 {
            tag = fieldT.Name
        } else if tags[0] == "-" {
            continue
        } else {
            tag = tags[0]
        }

        value := form.Get(tag)
        if len(value) == 0 {
            continue
        }

        switch fieldT.Type.Kind() {
        case reflect.Bool:
            if strings.ToLower(value) == "on" || strings.ToLower(value) == "1" || strings.ToLower(value) == "yes" {
                fieldV.SetBool(true)
                continue
            }
            if strings.ToLower(value) == "off" || strings.ToLower(value) == "0" || strings.ToLower(value) == "no" {
                fieldV.SetBool(false)
                continue
            }
            b, err := strconv.ParseBool(value)
            if err != nil {
                return err
            }
            fieldV.SetBool(b)
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            x, err := strconv.ParseInt(value, 10, 64)
            if err != nil {
                return err
            }
            fieldV.SetInt(x)
        case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
            x, err := strconv.ParseUint(value, 10, 64)
            if err != nil {
                return err
            }
            fieldV.SetUint(x)
        case reflect.Float32, reflect.Float64:
            x, err := strconv.ParseFloat(value, 64)
            if err != nil {
                return err
            }
            fieldV.SetFloat(x)
        case reflect.Interface:
            fieldV.Set(reflect.ValueOf(value))
        case reflect.String:
            fieldV.SetString(value)
        case reflect.Struct:
            switch fieldT.Type.String() {
            case "time.Time":
                format := time.RFC3339
                if len(tags) > 1 {
                    format = tags[1]
                }
                t, err := time.Parse(format, value)
                if err != nil {
                    return err
                }
                fieldV.Set(reflect.ValueOf(t))
            }
        case reflect.Slice:
            if fieldT.Type == sliceOfInts {
                formVals := form[tag]
                fieldV.Set(reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(int(1))), len(formVals), len(formVals)))
                for i := 0; i < len(formVals); i++ {
                    val, err := strconv.Atoi(formVals[i])
                    if err != nil {
                        return err
                    }
                    fieldV.Index(i).SetInt(int64(val))
                }
            } else if fieldT.Type == sliceOfStrings {
                formVals := form[tag]
                fieldV.Set(reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf("")), len(formVals), len(formVals)))
                for i := 0; i < len(formVals); i++ {
                    fieldV.Index(i).SetString(formVals[i])
                }
            }
        }
    }
    return nil
}

func ParseHttpBodyToArgs(r *http.Request, args interface{}) error {

    err := r.ParseForm()
    if nil != err {
        err = NewInternalError(DecodeErrCode, err)
    }
    err = ParseForm(r.Form, args)
    if nil != err {
        err = NewInternalError(DecodeErrCode, err)
    }

    var body []byte
    body, err = ioutil.ReadAll(r.Body)
    if err != nil {
        Logger.Error("UpdateUserInfo read body err : %s,%v", r.FormValue("user_id"), err)
        return err
    }
    defer r.Body.Close()
    if err := json.Unmarshal(body, args); err != nil {
        Logger.Error("Unmarshal body : %s,%s,%v", r.FormValue("user_id"), string(body), err)
        return err
    }

    return err
}

func WriteRespToBody(w http.ResponseWriter, resp interface{}) error {

    b, err := json.Marshal(resp)
    if err != nil {
        Logger.Error("Marshal json to bytes error :%v", err)
        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
        return err
    }
    w.Write(b)
    return err
}

func OptionHandler(c *gin.Context) {}
func GinCrossDomain() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Cache-Control", "no-cache")
        CheckCrossdomain(c)
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusOK)
        }
    }
}

func GinFilter() gin.HandlerFunc {
    return func(c *gin.Context) {

        if c.Request.Method == "HEAD" {
            c.AbortWithStatus(http.StatusOK)
        }
    }
}

func CheckCrossdomain(c *gin.Context) {
    c.Writer.Header().Add("Access-Control-Allow-Headers", "content-type, authorization")
    c.Writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
    c.Writer.Header().Add("Access-Control-Allow-Credentials", "true")

    origin_list := []string{"lhd.codoon.com", ".codoon.com$", ".runtopia.net$", ".blastapp.net$", "http://localhost", "192.168.\\d+"}
    for _, str := range origin_list {
        reg := regexp.MustCompile(str)
        if reg.MatchString(c.Request.Header.Get("Origin")) {
            fmt.Println("match origin: ", c.Request.Header.Get("Origin"))
            c.Writer.Header().Add("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
            break
        }
    }
}

var htmlEscape bool = true // default is true in json.Marshal

func SetHtmlEscape(b bool) {
    htmlEscape = b
}

func HTMLUnEscape(src []byte) []byte {
    src = bytes.Replace(src, []byte("\\u003c"), []byte("<"), -1)
    src = bytes.Replace(src, []byte("\\u003e"), []byte(">"), -1)
    src = bytes.Replace(src, []byte("\\u0026"), []byte("&"), -1)
    return src
}

func SendResponse(c *gin.Context, http_code int, data interface{}, err error) error {
    var resp CodoonApiResponse = CodoonApiResponse{"OK", nil, "success"}
    if err != nil {
        is_user_err, code, info := IsUserErr(err)
        if is_user_err {
            resp.Status = "Error"
            resp.Data = code
            resp.Desc = info
        } else {
            //CheckError(err)
            c.String(http_code, http.StatusText(http_code))

            // 500错误（user_error_code < 100）邮件发送
            if 500 == http_code {
                //host_name, _ := os.Hostname()
                //if IsOnline() {
                //    user_id, _ := strconv.ParseInt(c.Request.FormValue("user_id"), 10, 0)
                //    body := fmt.Sprintf("500 Code: \r\n<br> method: %s \r\n<br> uri: %s \r\n<br> hostname: %s \r\n<br> user_id: %d \r\n<br> trace_id: %s \r\n<br> err: %v", c.Request.Method, c.Request.URL.String(), host_name, user_id, local.TraceId(), err)
                //    go SendAlertMail([]string{"liangbo@codoon.com", "liucx@codoon.com"}, body)
                //}
            }

            return nil
        }
    } else {
        resp.Data = data
    }
    c.Writer.Header().Set("Content-Type", "application/json")
    c.Writer.Header().Set("ServerTime", strconv.FormatInt(time.Now().Unix(), 10))

    b, err := json.Marshal(&resp)
    if err != nil {
        Logger.Error("Marshal json to bytes error :%v", err)
    }

    c.Writer.Header().Del("Content-length")
    if 0 != len(b) {
    }

    if !htmlEscape {
        b = HTMLUnEscape(b)
    }

    // 输出结果，当大于3000个字符不在输出
    if len(b) > 3000 {
        //Logger.Info(string(b[:3000]))
    } else {
        Logger.Info("response: %s", string(b))
    }

    // 定死的变量不需要输出
    //Logger.Info("+++++++++++response header: %v ", c.Writer.Header())
    c.Writer.Write(b)

    return err
}

func SendFormRequest(http_method, urls string, header map[string]string, req_body map[string]string, resp_body interface{}) (int, error) {
    client := &http.Client{}

    var err error = nil
    form := url.Values{}
    var request *http.Request

    for key, value := range req_body {
        form.Set(key, value)
    }

    if "GET" == http_method {
        request, _ = http.NewRequest(http_method, urls+"?"+form.Encode(), nil)
    } else {
        request, _ = http.NewRequest(http_method, urls, strings.NewReader(form.Encode()))
        request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    }

    for key, value := range header {
        request.Header.Add(key, value)
    }

    response, err := client.Do(request)
    if nil != err {
        Logger.Error("send request err :%v", err)
        return 200, err
    }

    if response.StatusCode == 200 {
        defer response.Body.Close()
        decoder := json.NewDecoder(response.Body)
        decoder.UseNumber()
        decoder.Decode(&resp_body)
    }

    return response.StatusCode, err
}

func SendRequest(http_method, urls, method string, req_body interface{}) (int, error) {

    Logger.Info("http:%s,%v", urls+method, req_body)
    var err error = nil

    request, _ := http.NewRequest(http_method, urls+method, nil)

    if nil != req_body {
        b, _ := json.Marshal(req_body)
        request.Body = ioutil.NopCloser(strings.NewReader(string(b)))
    }
    /*
        if "GET" == http_method {
            request, _ = http.NewRequest(http_method, urls+method+"?"+form.Encode(), nil)
        } else {
            request, _ = http.NewRequest(http_method, urls+method, strings.NewReader(form.Encode()))
            request.Header.Set("Content-Type", "application/json")
        }
    */
    response, err := HttpClientClusterDo(request)
    if nil != err {
        err = NewInternalError(HttpErrCode, err)

        Logger.Error("send request err :%v", err)
        return 200, err
    }

    if response.StatusCode == 200 {
        defer response.Body.Close()
        body, err := ioutil.ReadAll(response.Body)
        if nil == err {
            Logger.Info("body:%v", string(body))
        }
    } else {
        err = NewInternalError(HttpErrCode, fmt.Errorf("http code :%d", response.StatusCode))
        Logger.Error("send request err :%v", err)
        return 200, err
    }

    return response.StatusCode, err
}

// modified by linagbo on 2017-08-10, 定义不被转化成int的key
var string_key map[string]int = map[string]int{
    "key":  1,
    "nick": 1,
    "club_name": 1,
    "position": 1,
}

func ParseHttpParamsToArgs(c *gin.Context, args map[string]interface{}, reply interface{}) error {
    var err error

    r := c.Request
    if nil != c.Request.Body {
        decoder := json.NewDecoder(r.Body)
        decoder.UseNumber()
        decoder.Decode(&args)

    }
    c.Request.Body.Close()
    for _, param := range c.Params {
        if 1 == string_key[param.Key] {
            args[param.Key] = param.Value
        } else {
            value_int, err := strconv.ParseInt(param.Value, 10, 0)
            if err != nil {
                args[param.Key] = param.Value
            } else {
                args[param.Key] = value_int
            }
        }
    }

    for key, value := range r.Form {
        if 1 == string_key[key] {
            args[key] = value[0]
        } else {
            value_int, err := strconv.ParseInt(value[0], 10, 0)
            if err != nil {
                args[key] = value[0]
            } else {
                args[key] = value_int
            }
        }
    }

    args["user_agent"] = r.Header.Get("User-Agent")
    Logger.Info("api request args: %+v", args)
    return err
}

//func ForwardHttpToRpc(c *gin.Context, client *RpcClient, method string, args map[string]interface{}, reply interface{}, http_code *int) error {
//
//    r := c.Request
//    if nil != c.Request.Body {
//        decoder := json.NewDecoder(r.Body)
//        decoder.UseNumber()
//        decoder.Decode(&args)
//
//    }
//    c.Request.Body.Close()
//    for _, param := range c.Params {
//
//        value_int, err := strconv.ParseInt(param.Value, 10, 0)
//        if err != nil || 1 == string_key[param.Key] {
//            args[param.Key] = param.Value
//        } else {
//            args[param.Key] = value_int
//        }
//    }
//
//    r.ParseForm()
//    for key, value := range r.Form {
//        value_int, err := strconv.ParseInt(value[0], 10, 0)
//        if err != nil || 1 == string_key[key] {
//            args[key] = value[0]
//        } else {
//            args[key] = value_int
//        }
//    }
//    args["user_agent"] = r.Header.Get("User-Agent")
//
//    // 日志输出统一都放到client.Call方法里面，重复
//    err := client.Call(method, &args, reply)
//    if nil != err {
//        is_user_err, _, _ := IsUserErr(err)
//        if !is_user_err {
//            *http_code = http.StatusInternalServerError
//        }
//    }
//    return err
//}

func GinRecovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            err := recover()
            if err != nil {
                switch err.(type) {
                case error:
                    //CheckError(err.(error))
                default:
                    err = errors.New(fmt.Sprint(err))
                    //CheckError(err)
                }

                stack := stack(3)
                Logger.Error("PANIC: %s\n%s", err, stack)

                c.Writer.WriteHeader(http.StatusInternalServerError)

                //host_name, _ := os.Hostname()
                //if IsOnline() {
                //    user_id, _ := strconv.ParseInt(c.Request.FormValue("user_id"), 10, 0)
                //    body := fmt.Sprintf("Recovery err: \r\n<br> method: %s \r\n<br> uri: %s \r\n<br> hostname: %s \r\n<br> user_id: %d \r\n<br> trace_id: %s \r\n<br> err: %v", c.Request.Method, c.Request.URL.String(), host_name, user_id, local.TraceId(), err)
                //    go SendAlertMail([]string{"liangbo@codoon.com", "liucx@codoon.com"}, body)
                //
                //    SendSmsChina("8618501790857", "500 code: "+c.Request.URL.String()+" "+host_name)
                //}
            }

        }()

        c.Next()
    }
}

func MyRecovery() {

    err := recover()
    if err != nil {
        switch err.(type) {
        case error:
            //CheckError(err.(error))
        default:
            err = errors.New(fmt.Sprint(err))
            //CheckError(err)
        }

        stack := stack(3)
        Logger.Error("PANIC: %s\n%s", err, stack)
    }

}

func GinLogger() gin.HandlerFunc {

    return func(c *gin.Context) {
        // Start timer
        start := time.Now()

        // Process request
        c.Next()

        // Stop timer
        end := time.Now()
        latency := end.Sub(start)

        clientIP := c.ClientIP()
        method := c.Request.Method
        statusCode := c.Writer.Status()
        Logger.Notice("[GIN] %v | %3d | %12v | %s | %-7s %s %s\n%s\n%s",
            end.Format("2006/01/02 - 15:04:05"),
            statusCode,
            latency,
            clientIP,
            method,
            c.Request.URL.String(),
            c.Request.URL.Opaque,
            c.Errors.String(),
            c.Keys)

        if statusCode == 500 {

        }

        if latency > DEFAULT_API_TIMEOUT {
            Logger.Error("[TIMEOUT] %v | %3d | %12v | %s | %-7s %s %s | %s",
                end.Format("2006/01/02 - 15:04:05"),
                statusCode,
                latency,
                clientIP,
                method,
                c.Request.URL.String(),
                c.Request.URL.Opaque,
                c.Errors.String())

        }
    }
}

var secretTable string = "_WY+Ytpa=A^Fm(Jl-rx@EVLl-Yx$v4+YgOhxB4s$Lqcen+BflOj_lgS3xuh5bSN-Jnhj69OSa(CmV5*91MRh8XIY423aPH_k$-u@XwaMgmPFCL1Ne-dx!kV$Q_US7f7fMV!H2CgjXmk)8aY3ftssyOrL-(c(UcW*QRd^8Fhcfs)A@qmR$8A8TFm8#)CvNE_CZ2lkvgVCC-vZaeDv^jb1QOv@W2+Ph!eQM=CtbtZPz(wX%gY)J$gdC8Rbc1L*(x6%tVO7RUutHAZF#6@sl(LzBP1DAzU7ttpHfqvKN$e5C@c!pg=@c$zL55$kg!8KJ$1SCbMbL^BYKaK9&_yxUU#XZF&GqY_tS!MN$zsWsLX*4uvCVG_EJ3-96qejb3z9m7e)BrmQMlTS9fVkA%5J5OL12BY8pzTJIeWC1z#jQaTwjnEl$cZj(sqY*LkMJG+(7l*ZNuY1rU5Tvcf6NH%5%7P8r&&yIsj=z2z4c=8VL5gelN-ZGOas$xpX8hf-qOK+MO8s"

//func CodoonCheckSign(c *gin.Context) {
//
//    // 获取token
//    var err error
//    r := c.Request
//    auth := r.Header.Get("Authorization")
//    timeStamp := r.Header.Get("TimeStamp")
//    xTable := r.Header.Get("XTable")
//    userAgent := r.Header.Get("User-Agent")
//    runSign := r.Header.Get("RunSign")
//    urlStr := r.URL.Path
//
//    if auth == "Bearer test_a670ad8689961de2c725b8b79c28956e" {
//        return
//    }
//    if "" == timeStamp || "" == runSign || "" == xTable {
//        Logger.Error("sign failed")
//        c.AbortWithStatus(http.StatusOK)
//    }
//
//    user_id, err := strconv.ParseInt(r.FormValue("user_id"), 10, 0)
//    if nil != err {
//        Logger.Error("get userid error :%v", err)
//    }
//    user_ids := []int64{user_id}
//    users, err := GetUserSummaryByIDs(user_ids, UserProfileClient)
//    if nil != err {
//        Logger.Error("get user summary by ids error :%v", err)
//    }
//    var phone string
//    if user, ok := users[user_id]; ok {
//        phone = user.Phone
//    }
//
//    xIndexStrs := strings.Split(xTable, ",")
//    var X []byte
//    for i := range xIndexStrs {
//        index, err := strconv.Atoi(xIndexStrs[i])
//        if nil != err {
//            Logger.Error("strconv.Atoi error :%s,%d", xTable, i)
//        }
//        X = append(X, secretTable[index])
//    }
//
//    sercet := userAgent + auth + phone + timeStamp + urlStr + string(X)
//
//    table := sha256.New()
//    table.Write([]byte(sercet))
//    hastSercetStr := hex.EncodeToString(table.Sum(nil))
//
//    Logger.Info("first hash :%s", hastSercetStr)
//
//    table = sha256.New()
//    table.Write([]byte(hastSercetStr + timeStamp))
//    hastSercetStr = hex.EncodeToString(table.Sum(nil))
//
//    Logger.Info("sercet:%s,runSign:%s,hash:%s", sercet, runSign, hastSercetStr)
//    if runSign == hastSercetStr {
//        Logger.Info("sign success, %d", user_id)
//        return
//    } else {
//
//        sercet := userAgent + auth + timeStamp + urlStr + string(X)
//
//        table := sha256.New()
//        table.Write([]byte(sercet))
//        hastSercetStr := hex.EncodeToString(table.Sum(nil))
//
//        Logger.Info("first hash :%s", hastSercetStr)
//
//        table = sha256.New()
//        table.Write([]byte(hastSercetStr + timeStamp))
//        hastSercetStr = hex.EncodeToString(table.Sum(nil))
//
//        if runSign == hastSercetStr {
//            Logger.Info("sign success, %d", user_id)
//            return
//        } else {
//            Logger.Info("sign fail :%s,%s, %d", runSign, hastSercetStr, user_id)
//            c.AbortWithStatus(http.StatusOK)
//            return
//        }
//    }
//    return
//}

func HttpRequest(method, url string, data []byte) (status int, body []byte, err error) {
    var data_reader io.Reader = nil
    if len(data) > 0 {
        data_reader = bytes.NewReader(data)
    }

    req, err := http.NewRequest(method, url, data_reader)
    if err != nil {
        return
    }
    local.FillTraceHttp(req)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

    resp, err := HttpClientClusterDo(req)
    if err != nil {
        return
    }
    defer resp.Body.Close()

    body, err = ioutil.ReadAll(resp.Body)
    return resp.StatusCode, body, err
}

func HttpClientClusterDo(request *http.Request) (*http.Response, error) {

    resp, err := http_client_cluster.HttpClientClusterDo(request)
    return resp, err
}

func GetTokenFromHeader(r *http.Request) string {
    var token string
    auth := r.Header.Get("Authorization")
    if "" == auth {
        cookie, err := r.Cookie("sessionid")
        if nil != err || nil == cookie {
            return token
        }
        auth = cookie.Value
        if "" == auth {
            return token
        }
    }
    auths := strings.Split(auth, " ")
    if 2 != len(auths) {
        return token
    }
    if "Bearer" == auths[0] {
        token = auths[1]
    }
    if "" == token {
        return token
    }
    return token
}

func GetGinRawPath(c *gin.Context) string {
    path := c.Request.URL.Path
    for i := range c.Params {
        path = strings.Replace(path, c.Params[i].Value, ":"+c.Params[i].Key, -1)
    }
    fmt.Println("path :%s", path)
    return path
}
