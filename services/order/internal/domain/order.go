package domain

import (
	"context"
	"encoding/json"
	"time"
)

type Order struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CustomerID    string    `gorm:"not null" json:"customer_id"`
	ProductID     string    `gorm:"not null" json:"product_id"`
	Quantity      int       `gorm:"not null" json:"quantity"`
	Status        string    `gorm:"not null;default:'PENDING'" json:"status"`
	CorrelationID string    `gorm:"not null;uniqueIndex" json:"correlation_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type OutboxEvent struct {
	ID            string          `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	EventType     string          `gorm:"not null"`
	Payload       json.RawMessage `gorm:"type:jsonb;not null"`
	CorrelationID string          `gorm:"not null"`
	Status        string          `gorm:"not null;default:'unprocessed'"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type OrderRepository interface {
	AutoMigrate() error
	CreateOrderWithOutbox(ctx context.Context, order *Order, event *OutboxEvent) error
	FindOrder(ctx context.Context, id string) (*Order, error)
}
