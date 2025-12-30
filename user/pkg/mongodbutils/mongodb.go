package mongodbutils

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// InsertOne 插入单条文档
func InsertOne(ctx context.Context, collectionName string, document interface{}, dbName ...string) (primitive.ObjectID, error) {
	coll := GetCollection(collectionName, dbName...)
	result, err := coll.InsertOne(ctx, document)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("insert document failed: %w", err)
	}

	// 检查插入ID是否为ObjectID类型
	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		return id, nil
	}
	return primitive.NilObjectID, errors.New("inserted ID is not a valid ObjectID")
}

// InsertMany 插入多条文档
func InsertMany(ctx context.Context, collectionName string, documents []interface{}, dbName ...string) ([]interface{}, error) {
	coll := GetCollection(collectionName, dbName...)
	result, err := coll.InsertMany(ctx, documents)
	if err != nil {
		return nil, fmt.Errorf("insert many documents failed: %w", err)
	}
	return result.InsertedIDs, nil
}

// FindOne 查询单条文档
func FindOne(ctx context.Context, collectionName string, filter interface{}, result interface{}, dbName ...string) error {
	coll := GetCollection(collectionName, dbName...)
	err := coll.FindOne(ctx, filter).Decode(result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil // 未找到文档不视为错误，返回nil
		}
		return fmt.Errorf("find one document failed: %w", err)
	}
	return nil
}

// Find 查询多条文档（支持分页、排序、投影）
func Find(ctx context.Context, collectionName string, filter interface{}, results interface{},
	opts ...*options.FindOptions) error {

	coll := GetCollection(collectionName)
	cursor, err := coll.Find(ctx, filter, opts...)
	if err != nil {
		return fmt.Errorf("find documents failed: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, results); err != nil {
		return fmt.Errorf("decode documents failed: %w", err)
	}
	return nil
}

// UpdateOne 更新单条文档
func UpdateOne(ctx context.Context, collectionName string, filter interface{}, update interface{}, dbName ...string) (*mongo.UpdateResult, error) {
	coll := GetCollection(collectionName, dbName...)
	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("update one document failed: %w", err)
	}
	return result, nil
}

// UpdateMany 更新多条文档
func UpdateMany(ctx context.Context, collectionName string, filter interface{}, update interface{}, dbName ...string) (*mongo.UpdateResult, error) {
	coll := GetCollection(collectionName, dbName...)
	result, err := coll.UpdateMany(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("update many documents failed: %w", err)
	}
	return result, nil
}

// UpsertOne 更新单条文档，如果不存在则插入
func UpsertOne(ctx context.Context, collectionName string, filter interface{}, update interface{}, dbName ...string) (*mongo.UpdateResult, error) {
	coll := GetCollection(collectionName, dbName...)

	// 设置 upsert 选项
	opts := options.Update().SetUpsert(true)

	result, err := coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, fmt.Errorf("upsert one document failed: %w", err)
	}
	return result, nil
}

// DeleteOne 删除单条文档
func DeleteOne(ctx context.Context, collectionName string, filter interface{}, dbName ...string) (*mongo.DeleteResult, error) {
	coll := GetCollection(collectionName, dbName...)
	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("delete one document failed: %w", err)
	}
	return result, nil
}

// DeleteMany 删除多条文档
func DeleteMany(ctx context.Context, collectionName string, filter interface{}, dbName ...string) (*mongo.DeleteResult, error) {
	coll := GetCollection(collectionName, dbName...)
	result, err := coll.DeleteMany(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("delete many documents failed: %w", err)
	}
	return result, nil
}

// CountDocuments 统计文档数量
func CountDocuments(ctx context.Context, collectionName string, filter interface{}, dbName ...string) (int64, error) {
	coll := GetCollection(collectionName, dbName...)
	count, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("count documents failed: %w", err)
	}
	return count, nil
}

// FindById 通过ID查询文档（MongoDB默认ID字段为 "_id"）
func FindById(ctx context.Context, collectionName string, id primitive.ObjectID, result interface{}, dbName ...string) error {
	return FindOne(ctx, collectionName, bson.M{"_id": id}, result, dbName...)
}

// UpdateById 通过ID更新文档
func UpdateById(ctx context.Context, collectionName string, id primitive.ObjectID, update interface{}, dbName ...string) (*mongo.UpdateResult, error) {
	return UpdateOne(ctx, collectionName, bson.M{"_id": id}, update, dbName...)
}

// DeleteById 通过ID删除文档
func DeleteById(ctx context.Context, collectionName string, id primitive.ObjectID, dbName ...string) (*mongo.DeleteResult, error) {
	return DeleteOne(ctx, collectionName, bson.M{"_id": id}, dbName...)
}

type TransactionFunc func(ctx mongo.SessionContext) (interface{}, error)

// ExecuteTransaction 执行事务（仅使用源码中存在的函数）
func ExecuteTransaction(ctx context.Context, fn TransactionFunc) (interface{}, error) {
	// 1. 启动会话
	session, err := client.StartSession()
	if err != nil {
		return nil, fmt.Errorf("启动会话失败: %w", err)
	}
	defer session.EndSession(ctx) // 确保会话关闭

	// 2. 将外部上下文转换为会话上下文
	sessionCtx := mongo.NewSessionContext(ctx, session)

	// 3. 配置事务选项（读关注、写关注）
	txOpts := options.Transaction().
		SetReadConcern(readconcern.Majority()).
		SetWriteConcern(writeconcern.Majority())

	// 4. 开始事务
	if err := session.StartTransaction(txOpts); err != nil {
		return nil, fmt.Errorf("开始事务失败: %w", err)
	}

	// 5. 执行事务处理函数
	result, err := fn(sessionCtx)
	if err != nil {
		// 失败时回滚事务
		if abortErr := session.AbortTransaction(sessionCtx); abortErr != nil {
			return nil, err
		}
		return nil, err
	}

	// 6. 提交事务
	if err := session.CommitTransaction(sessionCtx); err != nil {
		return nil, err
	}

	return result, nil
}
