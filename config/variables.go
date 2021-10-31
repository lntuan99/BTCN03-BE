package config

import (
    "fmt"
    "github.com/spf13/viper"
    "net/url"
    "os"
)

const (
    Local      = "local"
    Production = "production"
)

type SystemConfig struct {
    // Special configs
    CreateNecessaryConfigsForTables string

    // Postgres config
    PostgresConnectionString string

    // Redis config
    RedisConnectionAddress  string
    RedisConnectionPassword string
    RedisURL                string
    RedisPage               int

    // SMTP config
    SMTPHost      string
    SMTPPort      int
    SMTPUsername  string
    SMTPPassword  string
    SMTPFromEmail string
    SMTPFromName  string

    // Worker quantity
    WorkerQuantity int

    // Log/debug/swagger config
    PostgresLogMode bool
    DebugApi        bool
    QORSecretKey    string

    // App & Company config
    AppName          string
    OrganizationName string

    // Portal & api port & domain config
    CdnUrl                     string
    Domain                     string
    MediaDir                   string
    JobsPort                   string
    ApiPort                    string

    // Firebase config
    FirebaseFileName string
    FirebaseName     string
}

var Config *SystemConfig

func FetchEnvironmentVariables() {
    EnvType := os.Getenv("ENV")
    Config = NewSystemConfig(EnvType)
}

func NewSystemConfig(env string) *SystemConfig {
    cf := SystemConfig{}
    fileConfig := cf.GetConfigFile(env)
    fmt.Printf("Load Config File: %s \n", fileConfig)
    viper.SetConfigFile(fileConfig)
    err := viper.ReadInConfig()
    if err != nil {
        panic(err)
    }

    // Special configs
    cf.CreateNecessaryConfigsForTables = os.Getenv("CREATE_NECESSARY_CONFIGS_FOR_TABLES")

    // Postgres config
    cf.PostgresConnectionString = viper.GetString("postgres_connection_string")

    // Redis config
    cf.RedisConnectionAddress = fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.GetString("redis.port"))
    cf.RedisConnectionPassword = viper.GetString("redis.password")
    cf.RedisURL = fmt.Sprintf("redis://:%s@%s:%s", url.QueryEscape(viper.GetString("redis.password")), viper.GetString("redis.host"), viper.GetString("redis.port"))
    cf.RedisPage = viper.GetInt("redis.page")

    // SMTP config
    cf.SMTPHost = viper.GetString("smtp.host")
    cf.SMTPPort = viper.GetInt("smtp.port")
    cf.SMTPUsername = viper.GetString("smtp.username")
    cf.SMTPPassword = viper.GetString("smtp.password")
    cf.SMTPFromEmail = viper.GetString("smtp.from_email")
    cf.SMTPFromName = viper.GetString("smtp.from_name")

    // Worker quantity
    cf.WorkerQuantity = viper.GetInt("worker_quantity")

    // Log/debug/swagger config
    cf.PostgresLogMode = viper.GetBool("postgres_log_mode")
    cf.DebugApi = viper.GetBool("debug_api")
    cf.QORSecretKey = viper.GetString("qor_secret_key")

    // App & Company config
    cf.AppName = viper.GetString("app_name")
    cf.OrganizationName = viper.GetString("organization_name")
    if cf.AppName == "" || cf.OrganizationName == "" {
        panic("Don't let app-name and organization-name be empty!")
    }

    // Portal & api port & domain config
    cf.CdnUrl = viper.GetString("cdn_url")

    cf.Domain = viper.GetString("domain")
    if cf.Domain == "" {
        cf.Domain = "http://127.0.0.1:9002"
    }

    cf.MediaDir = viper.GetString("media_dir")
    cf.JobsPort = viper.GetString("jobs_port")
    cf.ApiPort = viper.GetString("api_port")

    // Firebase config
    cf.FirebaseFileName = viper.GetString("firebase_file_name")
    cf.FirebaseName = viper.GetString("firebase_name")

    return &cf
}

func (config *SystemConfig) GetConfigFile(env string) string {
    fileF := "zzz/config/%s_config.json"
    switch env {
    case Local, Production:
        return fmt.Sprintf(fileF, env)
    default:
        return fmt.Sprintf(fileF, Local)
    }
}
