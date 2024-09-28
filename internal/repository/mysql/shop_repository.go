package repository

import (
	"github.com/sigit14ap/shop-service/internal/domain"

	"gorm.io/gorm"
)

type ShopRepository interface {
	CreateShop(shop *domain.Shop) error
	GetShopByEmail(email string) (*domain.Shop, error)
	GetShopById(id uint64) (*domain.Shop, error)
}

type shopRepository struct {
	db *gorm.DB
}

func NewShopRepository(db *gorm.DB) ShopRepository {
	return &shopRepository{db: db}
}

func (repository *shopRepository) CreateShop(shop *domain.Shop) error {
	return repository.db.Create(shop).Error
}

func (repository *shopRepository) GetShopByEmail(email string) (*domain.Shop, error) {
	var shop domain.Shop
	err := repository.db.Where("email = ?", email).First(&shop).Error
	return &shop, err
}

func (repository *shopRepository) GetShopById(id uint64) (*domain.Shop, error) {
	var shop domain.Shop
	err := repository.db.Where("id = ?", id).First(&shop).Error
	return &shop, err
}
