/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/9/24 20:34 
 */
package utils

import (
    "fmt"
    "log"
    "os"
    "third/go-logging"
)

var Logger *logging.Logger
var MysqlLogger *log.Logger
var backend_info_leveld logging.LeveledBackend

func PathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

func loggerInit() {

    Logger = logging.MustGetLogger("log_not_configed")

    info_log_fp := os.Stdout
    err_log_fp := os.Stderr
    format_str := "%{color}%{level:.4s}:%{time:2006-01-02 15:04:05.000}[%{id:03x}][%{goroutineid}/%{goroutinecount}][trace:%{traceid}][span:%{spanid}][parent:%{parentid}] %{shortfile}%{color:reset} %{message}"

    backend_info := logging.NewLogBackend(info_log_fp, "", 0)
    backend_err := logging.NewLogBackend(err_log_fp, "", 0)
    format := logging.MustStringFormatter(format_str)
    backend_info_formatter := logging.NewBackendFormatter(backend_info, format)
    backend_err_formatter := logging.NewBackendFormatter(backend_err, format)

    backend_info_leveld = logging.AddModuleLevel(backend_info_formatter)
    backend_info_leveld.SetLevel(logging.NOTICE, "")

    backend_err_leveld := logging.AddModuleLevel(backend_err_formatter)
    backend_err_leveld.SetLevel(logging.WARNING, "")

    logging.SetBackend(backend_info_leveld, backend_err_leveld)
}

func InitLogger(process_name string) (*logging.Logger, error) {

    Logger = logging.MustGetLogger(process_name)

    ok, _ := PathExists(Config.LogDir)
    if !ok {
        err := os.MkdirAll(Config.LogDir, 0777)
        if nil != err {
            fmt.Println("can't make dir : %s, %v", Config.LogDir, err)
            return nil, err
        }
    }

    sql_log_fp, err := os.OpenFile(Config.LogDir+"/"+process_name+".log.mysql", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        fmt.Println("open file[%s.mysql] failed[%s]", Config.LogFile, err)
        return nil, err
    }

    MysqlLogger = log.New(sql_log_fp, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

    info_log_fp, err := os.OpenFile(Config.LogDir+"/"+process_name+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        fmt.Println("open file[%s] failed[%s]", Config.LogFile, err)
        return nil, err
    }

    err_log_fp, err := os.OpenFile(Config.LogDir+"/"+process_name+".log.wf", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        fmt.Println("open file[%s.wf] failed[%s]", Config.LogFile, err)
        return nil, err
    }

    backend_info := logging.NewLogBackend(info_log_fp, "", 0)
    backend_err := logging.NewLogBackend(err_log_fp, "", 0)
    //format := logging.MustStringFormatter("%{level}: [%{time:2006-01-02 15:04:05.000}][%{pid}][goroutine:%{goroutinecount}][trace:%{traceid}][span:%{spanid}][parent:%{parentid}][%{module}][%{shortfile}][%{shortfunc}][%{message}]")
    format := logging.MustStringFormatter("%{level}: [%{time:2006-01-02 15:04:05.000}][trace:%{traceid}][%{module}][%{shortfile}][%{shortfunc}][%{message}]")
    backend_info_formatter := logging.NewBackendFormatter(backend_info, format)
    backend_err_formatter := logging.NewBackendFormatter(backend_err, format)

    backend_info_leveld = logging.AddModuleLevel(backend_info_formatter)
    switch Config.LogLevel {
    case "ERROR":
        backend_info_leveld.SetLevel(logging.ERROR, "")
    case "WARNING":
        backend_info_leveld.SetLevel(logging.WARNING, "")
    case "NOTICE":
        backend_info_leveld.SetLevel(logging.NOTICE, "")
    case "INFO":
        backend_info_leveld.SetLevel(logging.INFO, "")
    case "DEBUG":
        backend_info_leveld.SetLevel(logging.DEBUG, "")
    default:
        backend_info_leveld.SetLevel(logging.ERROR, "")
    }

    backend_err_leveld := logging.AddModuleLevel(backend_err_formatter)
    backend_err_leveld.SetLevel(logging.ERROR, "")

    logging.SetBackend(backend_info_leveld, backend_err_leveld)

    return Logger, err
}

func ChangeLogLevel(LogLevel string) {
    switch LogLevel {
    case "ERROR":
        backend_info_leveld.SetLevel(logging.ERROR, "")
    case "WARNING":
        backend_info_leveld.SetLevel(logging.WARNING, "")
    case "NOTICE":
        backend_info_leveld.SetLevel(logging.NOTICE, "")
    case "INFO":
        backend_info_leveld.SetLevel(logging.INFO, "")
    case "DEBUG":
        backend_info_leveld.SetLevel(logging.DEBUG, "")
    default:
        backend_info_leveld.SetLevel(logging.ERROR, "")
    }
}
