/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/9/24 20:41 
 */
package utils

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "strings"
    "time"
)

type MysqlConfig struct {
    MysqlConn            string
    MysqlConnectPoolSize int
}

var ServerName string
var Env string

func IsOnline() bool {
    if "" == Env {
        var environment = os.Getenv("GOENV")
        if environment == "" {
            environment = "online"
        } else {
            environment = strings.ToLower(environment)
        }
        fmt.Println("environment", environment)
        Env = environment
    }
    if Env == "online" {
        return true
    }
    return false
}

type RedisConfig struct {
    RedisConn      string
    RedisPasswd    string
    ReadTimeout    int
    ConnectTimeout int
    WriteTimeout   int
    IdleTimeout    int
    MaxIdle        int
    MaxActive      int
    RedisDb        string
}

type RPCSetting struct {
    Addr string
    Net  string
}

type CeleryQueue struct {
    Url   string
    Queue string
}

type OssConfig struct {
    AccessKeyId     string
    AccessKeySecret string
    Region          string
    Bucket          string
}

// statsd, circuit
type HystrixConfig struct {
    StatsdAddr                   string
    StatsdPrefix                 string
    CircuitTimeout               int
    CircuitMaxConcurrent         int
    CircuitVolumeThreshold       int
    CircuitSleepWindow           int
    CircuitErrorPercentThreshold int
}

// consul check
type AgentServiceCheckConfig struct {
    Script            string
    HTTP              string
    TCP               string
    TTL               string
    DockerContainerId string
    Shell             string

    Interval string
    Timeout  string
}

func DefaultCheckConfig() AgentServiceCheckConfig {
    var defaultCheckConfig AgentServiceCheckConfig = AgentServiceCheckConfig{
        Interval: "10s",
        Timeout:  "2s",
    }
    return defaultCheckConfig
}

// consul
type ConsulConfig struct {
    ServiceName string
    ServiceIp   string
    ServiceTags []string
    AgentAddr   string
    Check       []AgentServiceCheckConfig
}

type Configure struct {
    MysqlSetting   map[string]MysqlConfig
    RedisSetting   map[string]RedisConfig
    RpcSetting     map[string]RPCSetting
    GRpcSettings   map[string][]RPCSetting
    CelerySetting  map[string]CeleryQueue
    OssSetting     OssConfig
    HystrixSetting HystrixConfig
    ConsulSetting  ConsulConfig
    SentryUrl      string
    SearchUrl      string
    LogDir         string
    LogFile        string
    Listen         string
    RpcListen      string
    LogLevel       string
    External       map[string]string
    ExternalInt64  map[string]int64
}

var Config *Configure
var CommonConfig *Configure
var g_config_file_last_modify_time time.Time
//var g_local_conf_file string
var conf_dir = "/var/config/"

func InitConfigFileEtcd(SERVERNAME, config_file string, config *Configure) error {
    var err error

    if IsOnline() {
        ServerName = fmt.Sprintf("%s_%s", SERVERNAME, Env)
    } else {
        ServerName = fmt.Sprintf("%s_%s", SERVERNAME, Env)
    }
    if config_file != "" {
        err = InitConfigFileWithoutEnv(config_file, config)
        if err != nil {
            fmt.Println("init config file error:", err)
            return err
        }
        Config = config
        CommonConfig = config
    }
    fmt.Println("config:", *Config, time.Now())
    return nil
}

func InitConfigFileWithoutEnv(filename string, config *Configure) error {
    filename = conf_dir + filename
    fmt.Println("config file name:", filename)
    fi, err := os.Stat(filename)
    if err != nil {
        fmt.Println("ReadFile: ", err.Error())
        return err
    }

    if g_config_file_last_modify_time.Equal(fi.ModTime()) {
        return nil
    } else {
        g_config_file_last_modify_time = fi.ModTime()
    }

    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        fmt.Println("ReadFile: ", err.Error())
        return err
    }

    if err := json.Unmarshal(bytes, config); err != nil {
        err = NewInternalError(DecodeErrCode, err)
        return err
    }
    //g_local_conf_file = filename
    return nil
}

func InitConfigFile(filename string, config *Configure) error {
    var environment = os.Getenv("GOENV")

    fmt.Println("environment", environment)
    switch environment {
    case "DEV":
        filename = filename + ".dev"
    case "TEST":
        filename = filename + ".test"
    case "PRE":
        filename = filename + ".pre"
    case "ONLINE":
        filename = filename + ".online"
    }

    fmt.Println("filename", filename)
    fi, err := os.Stat(filename)
    if err != nil {
        fmt.Println("ReadFile: ", err.Error())
        return err
    }

    if g_config_file_last_modify_time.Equal(fi.ModTime()) {
        return nil
    } else {
        g_config_file_last_modify_time = fi.ModTime()
    }

    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        fmt.Println("ReadFile: ", err.Error())
        return err
    }

    if err := json.Unmarshal(bytes, config); err != nil {
        err = NewInternalError(DecodeErrCode, err)
        fmt.Println("Unmarshal: ", err.Error())
        return err
    }
    fmt.Println("conifg :", *config)
    //g_local_conf_file = filename
    Config = config
    return nil
}
