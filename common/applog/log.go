package applog

import (
	"common/applog/mq"
	"encoding/json"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path"
	"time"
)

// Logger 真正的执行器
type Logger struct {
	// 真正的执行器，之后可以复写executor，使用不同的日志执行器即可
	executor *logrus.Logger
	// 配置文件
	logConfig *LogConfig
}

// CustomFormatter 自定义日志格式化器
type CustomFormatter struct {
	// 是否使用JSON格式
	IsJSON bool
}

// Format 实现 logrus.Formatter 接口
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 获取时间戳
	timestamp := entry.Time.Format("2006-01-02 15:04:05")

	if f.IsJSON {
		// JSON格式处理
		jsonFormatter := &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		}
		data, err := jsonFormatter.Format(entry)
		if err != nil {
			return nil, err
		}

		// 在JSON前面加上时间戳
		formatted := fmt.Sprintf("%s: %s", timestamp, string(data))
		return []byte(formatted), nil
	} else {
		// 文本格式处理
		textFormatter := &logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		}
		data, err := textFormatter.Format(entry)
		if err != nil {
			return nil, err
		}

		// 在文本前面加上时间戳
		formatted := fmt.Sprintf("%s: %s", timestamp, string(data))
		return []byte(formatted), nil
	}
}

// DefaultLogger 返回一个默认实例
func DefaultLogger() *Logger {
	return &Logger{
		executor: logrus.StandardLogger(),
	}
}

// RotateLogs 设置日志切分
func RotateLogs(fileName string, split string, maxAge int) (*rotatelogs.RotateLogs, error) {
	var splitTime time.Duration
	var fileExt string

	// 根据配置设置分隔时间和文件扩展名
	switch split {
	case "h", "hour", "hourly":
		splitTime = time.Hour
		fileExt = ".%Y-%m-%d-%H"
	case "d", "day", "daily":
		splitTime = 24 * time.Hour
		fileExt = ".%Y-%m-%d"
	case "m", "month", "monthly":
		// 每月切分 30天
		splitTime = 30 * 24 * time.Hour
		fileExt = ".%Y-%m"
	case "w", "week", "weekly":
		splitTime = 7 * 24 * time.Hour
		fileExt = ".%Y-%m-%d" // 每周切分也按天命名，但切割间隔是7天
	default:
		// 默认按小时切分
		splitTime = time.Hour
		fileExt = ".%Y-%m-%d-%H"
	}

	maxAgeDuration := time.Duration(maxAge) * time.Hour

	// 生成一个切分的file
	rateLogs, err := rotatelogs.New(
		fileName+fileExt,
		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),
		// 文件最大保存时间，使用配置的max_age值
		rotatelogs.WithMaxAge(maxAgeDuration),
		// 日志切割时间间隔
		rotatelogs.WithRotationTime(splitTime),
	)

	if err != nil {
		return nil, err
	}

	return rateLogs, nil
}

// createFormatter 根据配置创建日志格式化器
func createFormatter(format string) logrus.Formatter {
	switch format {
	case "json", "JSON":
		return &CustomFormatter{
			IsJSON: true,
		}
	case "text", "TEXT":
		return &CustomFormatter{
			IsJSON: false,
		}
	default:
		// 默认使用文本格式
		return &CustomFormatter{
			IsJSON: false,
		}
	}
}

// NewLogger 用来初始化一个logger
func NewLogger(conf *LogConfig) *Logger {
	if len(conf.LogFile) == 0 {
		return nil
	}

	// 确定日志文件的地址
	fileName := conf.LogFile
	if len(conf.Path) != 0 {
		fileName = path.Join(conf.Path, conf.LogFile)
	}

	// 设置输出的level以及初始化对应的实例
	executor := logrus.New()
	switch conf.Level {
	case LoggerInfo:
		executor.SetLevel(logrus.InfoLevel)
		break
	case LoggerError:
		executor.SetLevel(logrus.ErrorLevel)
		break
	case LoggerWarn:
		executor.SetLevel(logrus.WarnLevel)
		break
	case LoggerFatal:
		executor.SetLevel(logrus.FatalLevel)
		break
	case LoggerDebug:
		executor.SetLevel(logrus.DebugLevel)
		break
	default:
		executor.SetLevel(logrus.InfoLevel)
		break
	}

	// 创建日志格式化器
	formatter := createFormatter(conf.Format)
	executor.SetFormatter(formatter)

	rotateInfoWriter, err := RotateLogs(fileName+".info", conf.Split, conf.MaxAge)
	if err != nil {
		panic(fmt.Sprintf("failed to create info log rotator: %v", err))
	}

	rotateWarnWriter, err := RotateLogs(fileName+".warn", conf.Split, conf.MaxAge)
	if err != nil {
		panic(fmt.Sprintf("failed to create warn log rotator: %v", err))
	}

	rotateErrorWriter, err := RotateLogs(fileName+".error", conf.Split, conf.MaxAge)
	if err != nil {
		panic(fmt.Sprintf("failed to create error log rotator: %v", err))
	}

	rotateDebugWriter, err := RotateLogs(fileName+".debug", conf.Split, conf.MaxAge)
	if err != nil {
		panic(fmt.Sprintf("failed to create debug log rotator: %v", err))
	}

	rotateFatalWriter, err := RotateLogs(fileName+".fatal", conf.Split, conf.MaxAge)
	if err != nil {
		panic(fmt.Sprintf("failed to create fatal log rotator: %v", err))
	}

	handler, err := os.OpenFile(os.DevNull, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModeAppend)
	if err != nil {
		panic(fmt.Sprintf("open file failed, %v", err))
	}

	var outWriter io.Writer
	if conf.Stdout {
		outWriter = io.MultiWriter(os.Stdout, handler)
	} else {
		outWriter = io.Writer(handler)
	}

	executor.SetOutput(outWriter)

	// 使用相同的格式化器创建
	ifHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: rotateDebugWriter,
		logrus.WarnLevel:  rotateWarnWriter,
		logrus.InfoLevel:  rotateInfoWriter,
		logrus.FatalLevel: rotateFatalWriter,
		logrus.ErrorLevel: rotateErrorWriter,
	}, formatter) // 使用配置的格式化器

	executor.AddHook(ifHook)

	if conf.ELK.IsSendELK {
		err = InitELK(&conf.ELK)
		if err != nil {
			log.Fatal(err)
			return nil
		}

		// 初始化 Kafka 生产者
		mq.InitLogWriter(conf.ELK.KafkaAddr)

		consumerGroup := fmt.Sprintf("log-consumer-%s", conf.ELK.KafkaTopic)
		mq.InitLogReader(
			conf.ELK.KafkaAddr,
			conf.ELK.KafkaTopic,
			consumerGroup,
			conf.ELK.Index,
			GetEsClient(),
		)

		mq.RegisterShutdownHook()
	}

	return &Logger{
		executor:  executor,
		logConfig: conf, // 保存配置引用
	}
}

// Info 日志
func (t *Logger) Info(field map[string]interface{}, msg string) {
	if t.executor.GetLevel() == logrus.DebugLevel {
		t.executor.WithFields(field).Debug(msg)
	}
	t.executor.WithFields(field).Info(msg)

	if t.logConfig.ELK.IsSendELK {
		data, _ := json.Marshal(field)
		mq.GetLogWriter().Send(mq.LogData{
			Data:  data,
			Topic: t.logConfig.ELK.KafkaTopic,
		})
	}
}

// Debug debug日志
func (t *Logger) Debug(field map[string]interface{}, msg string) {
	t.executor.WithFields(field).Debug(msg)

	if t.logConfig.ELK.IsSendELK {
		data, _ := json.Marshal(field)
		mq.GetLogWriter().Send(mq.LogData{
			Data:  data,
			Topic: t.logConfig.ELK.KafkaTopic,
		})
	}
}

// Error 错误日志
func (t *Logger) Error(field map[string]interface{}, msg string) {
	if t.executor.GetLevel() == logrus.DebugLevel {
		t.executor.WithFields(field).Debug(msg)
	}
	t.executor.WithFields(field).Error(msg)

	if t.logConfig.ELK.IsSendELK {
		data, _ := json.Marshal(field)
		mq.GetLogWriter().Send(mq.LogData{
			Data:  data,
			Topic: t.logConfig.ELK.KafkaTopic,
		})
	}
}

// Warn warn日志
func (t *Logger) Warn(field map[string]interface{}, msg string) {
	if t.executor.GetLevel() == logrus.DebugLevel {
		t.executor.WithFields(field).Debug(msg)
	}
	t.executor.WithFields(field).Warn(msg)

	if t.logConfig.ELK.IsSendELK {
		data, _ := json.Marshal(field)
		mq.GetLogWriter().Send(mq.LogData{
			Data:  data,
			Topic: t.logConfig.ELK.KafkaTopic,
		})
	}
}

// Fatal 输出fatal日志
func (t *Logger) Fatal(field map[string]interface{}, msg string) {
	if t.executor.GetLevel() == logrus.DebugLevel {
		t.executor.WithFields(field).Debug(msg)
	}
	t.executor.WithFields(field).Fatal(msg)

	if t.logConfig.ELK.IsSendELK {
		data, _ := json.Marshal(field)
		mq.GetLogWriter().Send(mq.LogData{
			Data:  data,
			Topic: t.logConfig.ELK.KafkaTopic,
		})
	}
}
