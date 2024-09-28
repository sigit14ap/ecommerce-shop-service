package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	helpers "github.com/sigit14ap/shop-service/helpers"
	"github.com/sigit14ap/shop-service/internal/services"
)

type ShopHandler struct {
	shopService services.ShopService
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

var validate *validator.Validate

func NewShopHandler(shopService services.ShopService) *ShopHandler {
	return &ShopHandler{shopService: shopService}
}

func (handler *ShopHandler) Register(context *gin.Context) {
	validate = validator.New()
	var registerRequest RegisterRequest
	err := context.BindJSON(&registerRequest)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	err = validate.Struct(registerRequest)
	if err != nil {
		helpers.ErrorValidationResponse(context, err)
		return
	}

	err = handler.shopService.Register(registerRequest.Email, registerRequest.Name, registerRequest.Password)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	helpers.SuccessResponse(context, gin.H{"message": "Shop created successfully"})
}

func (handler *ShopHandler) Login(context *gin.Context) {
	validate = validator.New()
	var loginRequest LoginRequest
	err := context.BindJSON(&loginRequest)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	err = validate.Struct(loginRequest)
	if err != nil {
		helpers.ErrorValidationResponse(context, err)
		return
	}

	token, err := handler.shopService.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusUnauthorized, "Email or password are incorrect")
		return
	}

	helpers.SuccessResponse(context, gin.H{"token": token})
}

func (handler *ShopHandler) Me(context *gin.Context) {
	shopId, exists := context.Get("shopID")

	if !exists {
		helpers.ErrorResponse(context, http.StatusUnauthorized, "Shop id not found")
		return
	}

	shopIdUint, valid := shopId.(uint64)
	if !valid {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Shop ID is not valid"})
		return
	}

	shop, err := handler.shopService.Me(shopIdUint)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusUnauthorized, "Shop data not found")
		return
	}

	helpers.SuccessResponse(context, gin.H{"shop": shop})
}
