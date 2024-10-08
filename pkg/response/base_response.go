package response

import "net/http"

type BaseResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(message string, data interface{}) BaseResponse {
	return BaseResponse{
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string) BaseResponse {
	return BaseResponse{
		Status:  http.StatusBadRequest,
		Message: message,
		Data:    nil,
	}
}

func UnauthorizedResponse(message string) BaseResponse {
	return BaseResponse{
		Status:  http.StatusUnauthorized,
		Message: message,
		Data:    nil,
	}
}

func InternalServerErrorResponse(message string) BaseResponse {
	return BaseResponse{
		Status:  http.StatusInternalServerError,
		Message: message,
		Data:    nil,
	}
}
