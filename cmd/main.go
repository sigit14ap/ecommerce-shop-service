package main

import (
	"fmt"
	"os"

	"github.com/sigit14ap/shop-service/config"
	"github.com/sigit14ap/shop-service/helpers"
	delivery "github.com/sigit14ap/shop-service/internal/delivery/http"
	"github.com/sigit14ap/shop-service/internal/domain"
	repository "github.com/sigit14ap/shop-service/internal/repository/mysql"
	"github.com/sigit14ap/shop-service/internal/router"
	"github.com/sigit14ap/shop-service/internal/usecase"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()

	log := helpers.InitializeLogs()

	if log == nil {
		log.Fatal("Logger failed to started")
	}

	log.Info("Starting the shop service...")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&domain.Shop{})
	if err != nil {
		log.Fatalf("failed to auto-migrate Shop model: %v", err)
	}

	shopRepo := repository.NewShopRepository(db)
	shopService := usecase.NewShopUsecase(shopRepo)
	shopHandler := delivery.NewShopHandler(shopService)

	router := router.NewRouter(shopHandler)

	log.Info(router.Run(":" + os.Getenv("APP_PORT")))
}
