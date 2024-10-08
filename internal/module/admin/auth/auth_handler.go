package auth

import (
	"encoding/json"
	"net/http"
	"ticket-app/config"
	"ticket-app/internal/module/customer/auth/request"
	"ticket-app/internal/module/customer/auth/response"
	"ticket-app/pkg/constant"
	"ticket-app/pkg/helpers"
	"ticket-app/pkg/jwt"
	middleware "ticket-app/pkg/middleware"
	"ticket-app/pkg/rabbitmq"
	baseResponse "ticket-app/pkg/response"

	validator "github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService AuthService
	validate    *validator.Validate
	config      config.AppConfig
}

func NewAuthHandler(authService AuthService, config config.AppConfig) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
		config:      config,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(constant.INVALID_REQUEST))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(helpers.FormatValidationError(err)))
		return
	}

	data, token, err := h.authService.Login(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	var loginResponse response.LoginResponse
	loginResponse.UserentityToLoginResponse(data, token)

	rabbitmqResponse, err := rabbitmq.MappingJsonToRabbitMQMessage(loginResponse)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	go rabbitmq.PublishMessage((*rabbitmq.RabbitMQConnection)(&h.config.RabbitMQConn), rabbitmqResponse, "admin_login")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.LOGIN_SUCCESS, loginResponse))
}

func (h *AuthHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req request.RefreshTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(constant.INVALID_REQUEST))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(helpers.FormatValidationError(err)))
		return
	}

	userID := r.Context().Value(middleware.UserContextKey("userID")).(string)
	valid, err := jwt.ValidateRefreshToken(userID, req.RefreshToken)
	if !valid || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(constant.INVALID_TOKEN))
		return
	}

	token, err := h.authService.RefreshToken(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	var refreshTokenResponse response.RefreshTokenResponse
	refreshTokenResponse.UserentityToRefreshTokenResponse(token)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.GENERATE_TOKEN_SUCCESS, refreshTokenResponse))
}
