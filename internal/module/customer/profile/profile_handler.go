package profile

import (
	"encoding/json"
	"net/http"
	"ticket-app/internal/module/customer/profile/request"
	"ticket-app/internal/module/customer/profile/response"
	"ticket-app/pkg/constant"
	"ticket-app/pkg/helpers"
	middleware "ticket-app/pkg/middleware"
	baseResponse "ticket-app/pkg/response"

	validator "github.com/go-playground/validator/v10"
)

type ProfileHandler struct {
	ProfileService ProfileService
	validate       *validator.Validate
}

func NewProfileHandler(ProfileService ProfileService) *ProfileHandler {
	return &ProfileHandler{ProfileService: ProfileService,
		validate: validator.New()}
}

func (h *ProfileHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserContextKey("userID")).(string)

	data, err := h.ProfileService.GetMe(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(err.Error()))
		return
	}

	var userResponse response.UserResponse
	userResponse.UserentityToUserResponse(data)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.SUCCESS_GET, userResponse))
}

func (h *ProfileHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserContextKey("userID")).(string)

	var userUpdateRequest request.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&userUpdateRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(constant.INVALID_REQUEST))
		return
	}

	if err := h.validate.Struct(userUpdateRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(baseResponse.ErrorResponse(helpers.FormatValidationError(err)))
		return
	}

	data, err := h.ProfileService.UpdateUser(userID, userUpdateRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(baseResponse.InternalServerErrorResponse(constant.FAILED_UPDATE_PROFILE))
		return
	}

	var userResponse response.UserResponse
	userResponse.UserentityToUserResponse(data)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.SUCCESS_UPDATE, userResponse))
}

func (h *ProfileHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserContextKey("userID")).(string)
	deletedBy := userID

	err := h.ProfileService.DeleteUser(userID, deletedBy)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(baseResponse.InternalServerErrorResponse(constant.DELETED_FAILED))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(baseResponse.SuccessResponse(constant.SUCCESS_DELETED, nil))
}
