package esclient

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type Client struct {
	es    *elasticsearch.Client
	ctx   context.Context
	index string
}

func New(addr []string, cfg *elasticsearch.Config, tlsConfig *tls.Config) (*Client, error) {
	// 初始化默认配置
	defaultCfg := elasticsearch.Config{
		Addresses: addr,
		Transport: &http.Transport{
			ResponseHeaderTimeout: 30 * time.Second, // 默认超时设置
			TLSClientConfig:       tlsConfig,
		},
	}

	// 合并用户传入的配置
	if cfg != nil {
		if cfg.Username != "" {
			defaultCfg.Username = cfg.Username
		}
		if cfg.Password != "" {
			defaultCfg.Password = cfg.Password
		}
		if cfg.APIKey != "" {
			defaultCfg.APIKey = cfg.APIKey
		}

		if len(cfg.RetryOnStatus) > 0 {
			defaultCfg.RetryOnStatus = cfg.RetryOnStatus
		}
		if cfg.MaxRetries > 0 {
			defaultCfg.MaxRetries = cfg.MaxRetries
		}
	}

	es, err := elasticsearch.NewClient(defaultCfg)
	if err != nil {
		return nil, fmt.Errorf("初始化 ES 客户端失败: %w", err)
	}

	// 测试连接
	res, err := es.Ping(es.Ping.WithContext(context.Background()))
	if err != nil {
		return nil, fmt.Errorf("连接 ES 失败: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("ES 服务响应错误: %s", res.Status())
	}

	return &Client{
		es:  es,
		ctx: context.Background(),
	}, nil
}

// WithContext 设置上下文（如超时控制）
func (c *Client) WithContext(ctx context.Context) *Client {
	c.ctx = ctx
	return c
}

// Index 指定操作的索引名
func (c *Client) Index(index string) *Client {
	c.index = index
	return c
}

// Create 新增文档（指定 ID，为空则自动生成）
func (c *Client) Create(id string, doc interface{}) (string, error) {
	if c.index == "" {
		return "", fmt.Errorf("未指定索引名，请调用 Index() 方法")
	}

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return "", fmt.Errorf("文档序列化失败: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      c.index,
		DocumentID: id,
		Body:       strings.NewReader(string(docJSON)),
	}

	res, err := req.Do(c.ctx, c.es)
	if err != nil {
		return "", fmt.Errorf("创建文档失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return "", fmt.Errorf("创建文档响应错误: %s", res.Status())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}
	return result["_id"].(string), nil
}

// Get 获取文档
func (c *Client) Get(id string, out interface{}) error {
	if c.index == "" {
		return fmt.Errorf("未指定索引名")
	}

	req := esapi.GetRequest{
		Index:      c.index,
		DocumentID: id,
	}

	res, err := req.Do(c.ctx, c.es)
	if err != nil {
		return fmt.Errorf("获取文档失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("获取文档响应错误: %s", res.Status())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}
	source, ok := result["_source"]
	if !ok {
		return fmt.Errorf("文档不存在或无内容")
	}

	sourceJSON, err := json.Marshal(source)
	if err != nil {
		return fmt.Errorf("序列化文档内容失败: %w", err)
	}
	return json.Unmarshal(sourceJSON, out)
}

// Update 更新文档（部分更新）
func (c *Client) Update(id string, doc interface{}) error {
	if c.index == "" {
		return fmt.Errorf("未指定索引名")
	}

	updateBody := map[string]interface{}{"doc": doc}
	docJSON, err := json.Marshal(updateBody)
	if err != nil {
		return fmt.Errorf("更新内容序列化失败: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      c.index,
		DocumentID: id,
		Body:       strings.NewReader(string(docJSON)),
	}

	res, err := req.Do(c.ctx, c.es)
	if err != nil {
		return fmt.Errorf("更新文档失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("更新文档响应错误: %s", res.Status())
	}
	return nil
}

// Delete 删除文档
func (c *Client) Delete(id string) error {
	if c.index == "" {
		return fmt.Errorf("未指定索引名")
	}

	req := esapi.DeleteRequest{
		Index:      c.index,
		DocumentID: id,
	}

	res, err := req.Do(c.ctx, c.es)
	if err != nil {
		return fmt.Errorf("删除文档失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("删除文档响应错误: %s", res.Status())
	}
	return nil
}

// Bulk 批量操作
func (c *Client) Bulk(actions []esutil.BulkIndexerItem) error {
	if c.index == "" {
		return fmt.Errorf("未指定索引名")
	}

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         c.index,
		Client:        c.es,
		NumWorkers:    4,
		FlushBytes:    5e6,
		FlushInterval: 0,
	})
	if err != nil {
		return fmt.Errorf("创建批量索引器失败: %w", err)
	}

	for _, action := range actions {
		if err := bi.Add(c.ctx, action); err != nil {
			return fmt.Errorf("添加批量操作失败: %w", err)
		}
	}

	if err := bi.Close(c.ctx); err != nil {
		return fmt.Errorf("批量操作执行失败: %w", err)
	}

	stats := bi.Stats()
	if stats.NumFailed > 0 {
		return fmt.Errorf("批量操作失败，失败数: %d", stats.NumFailed)
	}
	return nil
}

// Search 执行查询
func (c *Client) Search(query interface{}, out interface{}) error {
	if c.index == "" {
		return fmt.Errorf("未指定索引名")
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("查询条件序列化失败: %w", err)
	}

	req := esapi.SearchRequest{
		Index: []string{c.index},
		Body:  strings.NewReader(string(queryJSON)),
	}

	res, err := req.Do(c.ctx, c.es)
	if err != nil {
		return fmt.Errorf("执行查询失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("查询响应错误: %s", res.Status())
	}

	return json.NewDecoder(res.Body).Decode(out)
}

// CreateIndex 创建索引
func (c *Client) CreateIndex(mapping interface{}) error {
	if c.index == "" {
		return fmt.Errorf("未指定索引名")
	}

	mappingJSON, err := json.Marshal(mapping)
	if err != nil {
		return fmt.Errorf("映射序列化失败: %w", err)
	}

	req := esapi.IndicesCreateRequest{
		Index: c.index,
		Body:  strings.NewReader(string(mappingJSON)),
	}

	res, err := req.Do(c.ctx, c.es)
	if err != nil {
		return fmt.Errorf("创建索引失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("创建索引响应错误: %s", res.Status())
	}
	return nil
}

// DeleteIndex 删除索引
func (c *Client) DeleteIndex() error {
	if c.index == "" {
		return fmt.Errorf("未指定索引名")
	}

	req := esapi.IndicesDeleteRequest{
		Index: []string{c.index},
	}

	res, err := req.Do(c.ctx, c.es)
	if err != nil {
		return fmt.Errorf("删除索引失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("删除索引响应错误: %s", res.Status())
	}
	return nil
}

// ExistsIndex 检查索引是否存在
func (c *Client) ExistsIndex() (bool, error) {
	if c.index == "" {
		return false, fmt.Errorf("未指定索引名")
	}

	req := esapi.IndicesExistsRequest{
		Index: []string{c.index},
	}

	res, err := req.Do(c.ctx, c.es)
	if err != nil {
		return false, fmt.Errorf("检查索引存在性失败: %w", err)
	}
	defer res.Body.Close()

	return !res.IsError(), nil
}
