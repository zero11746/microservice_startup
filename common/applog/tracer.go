package applog

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/structs"
)

type FieldMap map[string]any

type Log struct {
	Level    string   `json:"level"`
	TraceID  string   `json:"trace_id"`
	SpanID   string   `json:"span_id"`
	Time     string   `json:"time"`
	Msg      string   `json:"msg"`
	FileName string   `json:"file"`
	Line     int      `json:"line"`
	Request  any      `json:"request"`
	Response any      `json:"response"`
	Field    FieldMap `json:"field"`
}

// ITracer 注入的内容
type ITracer interface {
	ID() string
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Close()
}

// Tracer 返回一个随机的uuid，作为trace trace_id
type Tracer struct {
	traceID string
	spanID  string
	req     any
	resp    any
	logger  *Logger
}

// NewTracer 返回一个tracer
func NewTracer(logger *Logger, traceID, spanID string) *Tracer {
	tracer := &Tracer{
		logger: logger,
	}
	tracer.traceID = traceID
	tracer.spanID = spanID

	return tracer
}

func getServicePos() (file string, line int) {
	var pcs [32]uintptr
	// 获取当前调用栈
	n := runtime.Callers(0, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := frames.Next()
		if strings.HasSuffix(frame.File, ".pb.go") &&
			!strings.Contains(frame.File, "protoc-gen-go") {
			return frame.File, frame.Line
		}
		if !more {
			break
		}
	}

	return "unknown", 0
}

// 组建输出
func (t *Tracer) makeMsg(level string, args ...interface{}) (map[string]interface{}, string) {
	message := ""
	for _, arg := range args {
		message += fmt.Sprintf(" %v", arg)
	}

	file, line := getServicePos()
	log := Log{
		Level:    level,
		TraceID:  t.traceID,
		SpanID:   t.spanID,
		Time:     time.Now().Format(time.DateTime),
		Msg:      message,
		FileName: file,
		Line:     line,
		Request:  t.req,
		Response: t.resp,
		Field:    FieldMap{},
	}

	s := structs.New(log)
	s.TagName = "json"
	result := s.Map()

	return result, message
}

func (t *Tracer) WithReq(req any) *Tracer {
	t.req = req
	return t
}

func (t *Tracer) WithResp(resp any) *Tracer {
	t.resp = resp
	return t
}

func (t *Tracer) ID() string {
	return t.traceID
}

func (t *Tracer) Debug(args ...interface{}) {
	t.logger.Debug(t.makeMsg(LoggerDebug, args))
}

func (t *Tracer) Info(args ...interface{}) {
	t.logger.Info(t.makeMsg(LoggerInfo, args))
}

func (t *Tracer) Warn(args ...interface{}) {
	t.logger.Warn(t.makeMsg(LoggerWarn, args))
}

func (t *Tracer) Error(args ...interface{}) {
	t.logger.Error(t.makeMsg("error", args))
}

func (t *Tracer) Fatal(args ...interface{}) {
	t.logger.Fatal(t.makeMsg("fatal", args))
}

func (t *Tracer) Close() {
	return
}

func (t *Tracer) Printf(v1 string, v2 ...interface{}) {
	msg := fmt.Sprintf(v1, v2...)

	msg = strings.Replace(msg, "\r", " ", -1)
	msg = strings.Replace(msg, "\n", " ", -1)

	// 直接输出
	t.Info(msg)
}
