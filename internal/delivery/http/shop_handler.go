package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	helpers "github.com/sigit14ap/shop-service/helpers"
	"github.com/sigit14ap/shop-service/internal/delivery/dto"
	"github.com/sigit14ap/shop-service/internal/usecase"
)

type ShopHandler struct {
	shopUsecase usecase.ShopUsecase
}

var validate *validator.Validate

func NewShopHandler(shopUsecase usecase.ShopUsecase) *ShopHandler {
	return &ShopHandler{shopUsecase: shopUsecase}
}

func (handler *ShopHandler) Register(context *gin.Context) {
	validate = validator.New()
	var registerRequest dto.RegisterRequest
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

	err = handler.shopUsecase.Register(registerRequest.Email, registerRequest.Name, registerRequest.Password)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	helpers.SuccessResponse(context, gin.H{"message": "Shop created successfully"})
}

func (handler *ShopHandler) Login(context *gin.Context) {
	validate = validator.New()
	var loginRequest dto.LoginRequest
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

	token, err := handler.shopUsecase.Login(loginRequest.Email, loginRequest.Password)
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

	shop, err := handler.shopUsecase.Me(shopIdUint)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusUnauthorized, "Shop data not found")
		return
	}

	helpers.SuccessResponse(context, gin.H{"shop": shop})
}
