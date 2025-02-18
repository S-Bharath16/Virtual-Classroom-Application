package response

import "net/http"

type APIResponse struct {
	ResponseCode int         `json:"responseCode"`
	ResponseBody ResponseBody `json:"responseBody"`
}

type ResponseBody struct {
	Message string      `json:"MESSAGE"`
	Data    interface{} `json:"DATA"`
}

func newResponse(code int, message string, data interface{}) APIResponse {
	return APIResponse{
		ResponseCode: code,
		ResponseBody: ResponseBody{
			Message: message,
			Data:    data,
		},
	}
}


func SetResponseOk(message string, data interface{}) APIResponse {
	return newResponse(http.StatusOK, message, data);
}


func SetResponseNotFound(message string, data interface{}) APIResponse {
	return newResponse(http.StatusNotFound, message, data);
}

func SetResponseBadRequest(message string, data interface{}) APIResponse {
	return newResponse(http.StatusBadRequest, message, data);
}

func SetResponseInternalError(data interface{}) APIResponse {
	return newResponse(http.StatusInternalServerError, "Internal Server Error occurred :(", data);
}

func SetResponseUnauth(data interface{}) APIResponse {
	return newResponse(http.StatusUnauthorized, "Unauthorized Access Denied !!", data);
}

func SetResponseTimedOut(message string, data interface{}) APIResponse {
	return newResponse(http.StatusRequestTimeout, message, data);
}

func SetResponseTransactionFailed(message string, data interface{}) APIResponse {
	return newResponse(http.StatusAccepted, message, data);
}