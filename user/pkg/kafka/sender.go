package elk

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type LogData struct {
	Topic string
	//json数据
	Data []byte
}
type KafkaWriter struct {
	w    *kafka.Writer
	data chan LogData
}

func InitWriter(kkAddr string) *KafkaWriter {
	w := &kafka.Writer{
		Addr:     kafka.TCP(kkAddr),
		Balancer: &kafka.LeastBytes{},
	}
	k := &KafkaWriter{
		w:    w,
		data: make(chan LogData, 100),
	}
	go k.sendKafka()
	return k
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
					break
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
