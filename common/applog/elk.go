package applog

import (
	"crypto/tls"
	"fmt"

	esclient "common/es"

	"github.com/elastic/go-elasticsearch/v8"
)

var (
	globalEsClient *esclient.Client
)

func InitELK(elk *ELK) error {
	if elk == nil {
		return fmt.Errorf("elk配置不能为空")
	}

	// 初始化 Elasticsearch 客户端
	esCli, err := esclient.New(
		[]string{elk.Addr},
		&elasticsearch.Config{
			Username:      elk.Username,
			Password:      elk.Password,
			APIKey:        elk.APIKey,
			RetryOnStatus: elk.RetryOnStatus,
			MaxRetries:    elk.MaxRetries,
		},
		&tls.Config{
			InsecureSkipVerify: elk.InsecureSkipVerify,
		},
	)
	if err != nil {
		return fmt.Errorf("初始化 ES 客户端失败: %w", err)
	}
	globalEsClient = esCli

	// 提前创建 ES 索引
	if err := initEsLogIndex(esCli, elk.Index); err != nil {
		return fmt.Errorf("初始化 ES 日志索引失败: %w", err)
	}

	return nil
}

func GetEsClient() *esclient.Client {
	return globalEsClient
}

// 提前创建 ES 日志索引并定义映射
func initEsLogIndex(esCli *esclient.Client, indexName string) error {
	// 检查索引是否已存在，存在则直接返回
	exists, err := esCli.Index(indexName).ExistsIndex()
	if err != nil {
		return fmt.Errorf("检查索引是否存在失败: %w", err)
	}
	if exists {
		fmt.Printf("ES 索引 %s 已存在，无需重复创建\n", indexName)
		return nil
	}

	// 定义与 KafkaLog 对应的映射
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"dynamic": false, // 关闭动态映射
			"properties": map[string]interface{}{
				// 日志级别（精确匹配，支持聚合统计）
				"level": map[string]interface{}{
					"type": "keyword", // 不分词，适合精确筛选
				},

				"trace_id": map[string]interface{}{
					"type": "keyword", // 不分词，适合精确筛选
				},

				"span_id": map[string]interface{}{
					"type": "keyword", // 不分词，适合精确筛选
				},

				// 日志时间（支持常见时间格式）
				"time": map[string]interface{}{
					"type":   "date",
					"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd'T'HH:mm:ssZ||yyyy-MM-dd HH:mm:ss.SSS", // 支持带毫秒的格式
				},

				// 日志内容（中文分词，支持全文检索）
				"msg": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type":         "keyword",
							"ignore_above": 256,
						},
					},
				},

				// FileName：文件名
				"file": map[string]interface{}{
					"type": "keyword",
				},

				// Line：行号（数值类型，支持范围查询）
				"line": map[string]interface{}{
					"type": "integer", // 整数类型
				},

				// Request：请求数据
				"request": map[string]interface{}{
					"type":    "object",
					"dynamic": true,
				},

				// Response：响应数据
				"response": map[string]interface{}{
					"type":    "object",
					"dynamic": true,
				},

				// Field：额外扩展字段
				"field": map[string]interface{}{
					"type":    "object",
					"dynamic": true,
				},
			},
		},
	}

	// 创建索引
	if err := esCli.Index(indexName).CreateIndex(mapping); err != nil {
		return fmt.Errorf("创建 ES 索引失败: %w", err)
	}

	return nil
}
