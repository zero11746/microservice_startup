package database

import (
	"context"
	"errors"
	"gorm.io/gorm"
)

// TransactionManager 事务管理器
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager 创建事务管理器
func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// 事务上下文键
type txKey struct{}

// Execute 执行事务
func (m *TransactionManager) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		// 将事务对象 tx 存入 context
		txCtx := context.WithValue(ctx, txKey{}, tx)

		// 执行事务函数
		return fn(txCtx)
	})
}

// GetDBFromContext 从上下文中获取事务对象
func GetDBFromContext(ctx context.Context) (*gorm.DB, error) {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	if !ok {
		return nil, errors.New("transaction context not found")
	}
	return tx, nil
}
