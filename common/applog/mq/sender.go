package mq

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

var globalKafkaWriter *KafkaWriter

type LogData struct {
	Topic string
	//json数据
	Data []byte
}
type KafkaWriter struct {
	w    *kafka.Writer
	data chan LogData
}

func InitLogWriter(kkAddr string) {
	w := &kafka.Writer{
		Addr:     kafka.TCP(kkAddr),
		Balancer: &kafka.LeastBytes{},
	}
	k := &KafkaWriter{
		w:    w,
		data: make(chan LogData, 100),
	}

	go k.sendKafka()
	globalKafkaWriter = k
	return
}

func GetLogWriter() *KafkaWriter {
	return globalKafkaWriter
}

func (w *KafkaWriter) Send(data LogData) {
	w.data <- data
}

func (w *KafkaWriter) Close() {
	if w.w != nil {
		w.w.Close()
	}
}

func (w *KafkaWriter) sendKafka() {
	for {
		select {
		case data := <-w.data:
			messages := []kafka.Message{
				{
					Topic: data.Topic,
					Value: data.Data,
				},
			}
			var err error
			const retries = 3
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			for i := 0; i < retries; i++ {
				// attempt to create topic prior to publishing the message
				err = w.w.WriteMessages(ctx, messages...)
				if err == nil {
					log.Printf("kafka send writemessage success \n")
					break
				} else {
					log.Printf("kafka send writemessage err %s \n", err.Error())
				}
				if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
					time.Sleep(time.Millisecond * 250)
					continue
				}
				if err != nil {
					log.Printf("kafka send writemessage err %s \n", err.Error())
				}
			}
			cancel()
		}
	}

}

// RegisterShutdownHook 注册程序退出钩子，释放资源
func RegisterShutdownHook() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigChan
		log.Println("收到退出信号，开始释放资源...")

		// 关闭 Kafka 生产者
		if globalKafkaWriter != nil {
			globalKafkaWriter.Close()
			log.Println("Kafka 生产者已关闭")
		}

		// 关闭 Kafka 消费者
		if globalKafkaReader != nil {
			globalKafkaReader.Close()
			log.Println("Kafka 消费者已关闭")
		}

		os.Exit(0)
	}()
}
