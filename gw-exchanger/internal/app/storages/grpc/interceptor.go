package grpc

import (
	"context"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type contextKey string

const dbContextKey = contextKey("db")

// Unary interceptor для добавления подключения к базе данных в контекст
func DBInterceptor(db *gorm.DB) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Добавляем DB в контекст
		ctx = context.WithValue(ctx, dbContextKey, db)
		return handler(ctx, req)
	}
}

// Вспомогательная функция для извлечения DB из контекста
func GetDB(ctx context.Context) *gorm.DB {
	db, ok := ctx.Value(dbContextKey).(*gorm.DB)
	if !ok {
		return nil
	}
	return db
}
