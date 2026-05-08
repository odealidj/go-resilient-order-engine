package repository

import (
	"context"
	"go-resilient-order-engine/services/order/internal/domain"
	"go-resilient-order-engine/services/order/migrations"
	"log/slog"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

// AutoMigrate applies schema updates using production-grade Goose migrations
func (r *orderRepository) AutoMigrate() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	goose.SetBaseFS(migrations.EmbedFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	slog.Info("Running production-ready Goose database migrations on Primary DB...")
	return goose.Up(sqlDB, ".")
}

// CreateOrderWithOutbox saves the order and the outbox event atomically.
func (r *orderRepository) CreateOrderWithOutbox(ctx context.Context, order *domain.Order, event *domain.OutboxEvent) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		if err := tx.Create(event).Error; err != nil {
			return err
		}
		return nil
	})
}

// FindOrder is a read operation. GORM dbresolver will route this to the Replica DB.
func (r *orderRepository) FindOrder(ctx context.Context, id string) (*domain.Order, error) {
	var order domain.Order
	if err := r.db.WithContext(ctx).First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}
