package applog

// LogConfig 日志的配置相关
type LogConfig struct {
	Path    string `toml:"path"`
	LogFile string `toml:"log_file"`
	Split   string `toml:"split"`
	Level   string `toml:"level"`
	Stdout  bool   `toml:"stdout"`
	MaxAge  int    `toml:"max_age"`
	Format  string `toml:"format"`
	ELK     ELK
}

type ELK struct {
	IsSendELK          bool   `toml:"is_send_elk"`          // 是否发送日志到 ELK
	KafkaAddr          string `toml:"kafka_addr"`           // Kafka 地址
	KafkaTopic         string `toml:"kafka_topic"`          // Kafka 日志主题
	Addr               string `toml:"addr"`                 // Elasticsearch 地址
	Index              string `toml:"index"`                // Elasticsearch 索引名（主配置）
	Username           string `toml:"username"`             // Elasticsearch 认证用户名
	Password           string `toml:"password"`             // Elasticsearch 认证密码
	APIKey             string `toml:"api_key"`              // Elasticsearch API 密钥（可选）
	InsecureSkipVerify bool   `toml:"insecure_skip_verify"` // 是否跳过 SSL 证书验证（测试用）
	RetryOnStatus      []int  `toml:"retry_on_status"`      // 需要重试的 HTTP 状态码
	MaxRetries         int    `toml:"max_retries"`          // 最大重试次数
}

const (
	LoggerDebug = "debug"
	LoggerWarn  = "warn"
	LoggerInfo  = "info"
	LoggerError = "error"
	LoggerFatal = "fatal"
)
