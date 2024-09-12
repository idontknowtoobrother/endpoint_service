package utils

type ErrorResponse struct {
	status  int    `json:"status"`
	message string `json:"message"`
}

type SuccessResponse struct {
	status  int         `json:"status"`
	message string      `json:"message"`
	data    interface{} `json:"data"`
}

func NewErrorResponse(status int, message string) ErrorResponse {
	return ErrorResponse{
		status:  status,
		message: message,
	}
}

func NewSuccessResponse(status int, message string, data interface{}) SuccessResponse {
	return SuccessResponse{
		status:  status,
		message: message,
		data:    data,
	}
}
