package config

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Port :nodoc:
func Port() string {
	return viper.GetString("port")
}

// Env :nodoc:
func Env() string {
	return viper.GetString("env")
}

// PapertrailHost :nodoc:
func PapertrailHost() string {
	return viper.GetString("papertrail.host")
}

// PapertrailPort :nodoc:
func PapertrailPort() int {
	return viper.GetInt("papertrail.port")
}

// PapertrailAppName :nodoc:
func PapertrailAppName() string {
	return viper.GetString("papertrail.app_name")
}

// PapertrailLogLevel :nodoc:
func PapertrailLogLevel() string {
	return viper.GetString("papertrail.log_level")
}

// CockroachHost :nodoc:
func CockroachHost() string {
	return viper.GetString("cockroach.host")
}

// CockroachDatabase :nodoc:
func CockroachDatabase() string {
	return viper.GetString("cockroach.database")
}

// CockroachUsername :nodoc:
func CockroachUsername() string {
	return viper.GetString("cockroach.username")
}

// CockroachPassword :nodoc:
func CockroachPassword() string {
	return viper.GetString("cockroach.password")
}

// DatabaseDSN :nodoc:
func DatabaseDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		CockroachUsername(),
		CockroachPassword(),
		CockroachHost(),
		CockroachDatabase())
}

// RedisCacheHost :nodoc:
func RedisCacheHost() string {
	return viper.GetString("redis.cache_host")
}

// DisableCaching :nodoc:
func DisableCaching() bool {
	return viper.GetBool("disable_caching")
}

// RedisLockHost :nodoc:
func RedisLockHost() string {
	return viper.GetString("redis.lock_host")
}

// RedisWorkerBrokerHost :nodoc:
func RedisWorkerBrokerHost() string {
	return viper.GetString("redis.worker_broker_host")
}

// CacheTTL :nodoc:
func CacheTTL() time.Duration {
	if !viper.IsSet("cache_ttl") {
		return time.Duration(DefaultCacheTTL) * time.Millisecond
	}

	return time.Duration(viper.GetInt("cache_ttl")) * time.Millisecond
}

// MobileWebRootURL :nodoc:
func MobileWebRootURL() string {
	switch Env() {
	case "production":
		return ProdKumparanMobileWebRootURL
	case "staging":
		return StagingKumparanMobileWebRootURL
	default:
		return DevKumparanMobileWebRootURL
	}
}

// DesktopWebRootURL :nodoc:
func DesktopWebRootURL() string {
	switch Env() {
	case "production":
		return ProdKumparanDesktopWebRootURL
	case "staging":
		return StagingKumparanDesktopWebRootURL
	default:
		return DevKumparanDesktopWebRootURL
	}
}

// AWSRegion :nodoc:
func AWSRegion() string {
	return viper.GetString("aws.region")
}

// AWSS3Bucket :nodoc:
func AWSS3Bucket() string {
	return viper.GetString("aws.s3_bucket")
}

// AWSS3Key :nodoc:
func AWSS3Key() string {
	return viper.GetString("aws.s3_key")
}

// AWSS3Secret :nodoc:
func AWSS3Secret() string {
	return viper.GetString("aws.s3_secret")
}

// UserServiceTarget :nodoc:
func UserServiceTarget() string {
	return viper.GetString("services.user_target")
}

//RPCClientTimeout :nodoc:
func RPCClientTimeout() time.Duration {
	if !viper.IsSet("rpc_client_timeout") {
		return time.Duration(DefaultRPCClientTimeout) * time.Millisecond
	}

	return time.Duration(viper.GetInt("rpc_client_timeout")) * time.Millisecond
}

//RPCServerTimeout :nodoc:
func RPCServerTimeout() time.Duration {
	if !viper.IsSet("rpc_server_timeout") {
		return time.Duration(DefaultRPCServerTimeout) * time.Millisecond
	}

	return time.Duration(viper.GetInt("rpc_server_timeout")) * time.Millisecond
}

// ServiceMaxConnPool :nodoc:
func ServiceMaxConnPool() int {
	if viper.GetInt("services.max_conn_pool") > 0 {
		return viper.GetInt("services.max_conn_pool")
	}
	return 200
}

// ServiceIdleConnPool :nodoc:
func ServiceIdleConnPool() int {
	if viper.GetInt("services.idle_conn_pool") > 0 {
		return viper.GetInt("services.idle_conn_pool")
	}
	return 100
}

// HTTPTimeout :nodoc:
func HTTPTimeout() time.Duration {
	if viper.GetInt("http_timeout") > 0 {
		return time.Duration(viper.GetInt("http_timeout")) * time.Millisecond
	}
	return DefaultHTTPTimeout
}

// SentryDSN :nodoc:
func SentryDSN() string {
	return viper.GetString("sentry_dsn")
}

// LogLevel :nodoc:
func LogLevel() string {
	return viper.GetString("log_level")
}

// GetConf :nodoc:
func GetConf() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./..")
	viper.AddConfigPath("./../..")
	viper.AddConfigPath("./../../..")
	viper.SetConfigName("config")
	viper.SetEnvPrefix("svc")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Warningf("%v", err)
	}

	return
}
