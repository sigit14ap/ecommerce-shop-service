package services

import (
	"errors"

	"github.com/sigit14ap/shop-service/helpers"
	"github.com/sigit14ap/shop-service/internal/domain"
	repository "github.com/sigit14ap/shop-service/internal/repository/mysql"
	"golang.org/x/crypto/bcrypt"
)

type ShopService interface {
	Register(email string, name string, password string) error
	Login(email string, password string) (string, error)
	Me(id uint64) (*domain.Shop, error)
}

type shopService struct {
	shopRepository repository.ShopRepository
}

func NewShopService(shopRepo repository.ShopRepository) ShopService {
	return &shopService{
		shopRepository: shopRepo,
	}
}

func (service *shopService) Register(email string, name string, password string) error {

	_, err := service.shopRepository.GetShopByEmail(email)
	if err == nil {
		return errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	shop := &domain.Shop{
		Email:    email,
		Name:     name,
		Password: string(hashedPassword),
	}

	return service.shopRepository.CreateShop(shop)
}

func (service *shopService) Login(email string, password string) (string, error) {
	var shop *domain.Shop
	var err error

	shop, err = service.shopRepository.GetShopByEmail(email)

	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(shop.Password), []byte(password))
	if err != nil {
		return "", err
	}

	token, err := helpers.GenerateJWT(shop.Email, shop.ID)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (service *shopService) Me(id uint64) (*domain.Shop, error) {
	shop, err := service.shopRepository.GetShopById(id)

	if err != nil {
		return nil, err
	}

	return shop, nil
}
