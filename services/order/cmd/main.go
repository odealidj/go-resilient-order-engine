package main

import (
	"go-resilient-order-engine/common/logger"
	config "go-resilient-order-engine/services/analytic/confg"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func main() {
	logger.Init()
	cfg := config.LoadConfig()

	// 1. setup gorm
	db, err := gorm.Open(postgres.Open(cfg.PrimaryDBURL), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to Primary DB", "error", err)
		panic(err)
	}

	// 2. Setup Read/Write Split with dbresolver plugin
	err = db.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{postgres.Open(cfg.ReplicaDBURL)},
		Policy:   dbresolver.RandomPolicy{},
	}).SetMaxIdleConns(10).SetConnMaxLifetime(time.Hour))
	if err != nil {
		slog.Error("Failed to setup DB resolver for GORM", "error", err)
		panic(err)
	}

	slog.Info("Database connected with Read/Write Split (Primary/Replica) via GORM dbresolver")

}
