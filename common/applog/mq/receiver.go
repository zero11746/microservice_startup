package mq

import (
	esclient "common/es"
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

var globalKafkaReader *KafkaReader

type KafkaReader struct {
	R       *kafka.Reader
	esCli   *esclient.Client
	ctx     context.Context
	cancel  context.CancelFunc
	esIndex string
}

func InitLogReader(kkAddr, logTopic, groupID, esIndex string, esCli *esclient.Client) {
	// 配置 Kafka 阅读器，只消费 "log" 主题
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kkAddr},
		Topic:    logTopic,
		GroupID:  groupID,
		MinBytes: 1024,
		MaxBytes: 1024 * 1024,
	})

	ctx, cancel := context.WithCancel(context.Background())
	globalKafkaReader = &KafkaReader{
		R:       reader,
		esCli:   esCli,
		ctx:     ctx,
		cancel:  cancel,
		esIndex: esIndex,
	}
	go globalKafkaReader.readMsg()
}

// 消费消息并写入 ES
func (r *KafkaReader) readMsg() {
	defer func() {
		// 退出时关闭 Kafka 阅读器
		if err := r.R.Close(); err != nil {
			log.Printf("关闭 Kafka 阅读器失败: %v", err)
		}
		r.cancel() // 取消上下文
		log.Println("消费循环已退出")
	}()

	for {
		select {
		case <-r.ctx.Done():
			return
		default:
			// 读取 Kafka 消息
			m, err := r.R.ReadMessage(r.ctx)
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					log.Printf("Kafka 读取消息失败: %v，5秒后重试", err)
					time.Sleep(5 * time.Second)
				}
				continue
			}

			log.Printf("收到消息: 主题=%s, 分区=%d, 偏移量=%d", m.Topic, m.Partition, m.Offset)

			doc := m.Value

			// 写入 Elasticsearch
			_, err = r.esCli.Index(r.esIndex).Create("", doc)
			if err != nil {
				log.Printf("写入 ES 失败: 消息值=%s, 错误=%v", string(doc), err)
				if retryErr := r.retryWriteES("logs", doc, 3); retryErr != nil {
					log.Printf("多次重试后仍写入失败: %v", retryErr)
				}
				continue
			}

			// 消费成功，提交偏移量
			if err := r.R.CommitMessages(r.ctx, m); err != nil {
				log.Printf("提交偏移量失败: %v", err)
			} else {
				log.Printf("消息处理成功: 偏移量=%d", m.Offset)
			}
		}
	}
}

// 重试写入 ES
func (r *KafkaReader) retryWriteES(index string, doc []byte, retries int) error {
	var lastErr error
	for i := 0; i < retries; i++ {
		_, err := r.esCli.Index(index).Create("", doc)
		if err == nil {
			return nil
		}
		lastErr = err
		delay := time.Duration(1<<i) * time.Second
		log.Printf("第 %d 次重试写入 ES 失败，%v 后重试: %v", i+1, delay, err)
		time.Sleep(delay)
	}
	return fmt.Errorf("达到最大重试次数(%d)，最后错误: %w", retries, lastErr)
}

func (r *KafkaReader) Close() {
	r.cancel()
}
