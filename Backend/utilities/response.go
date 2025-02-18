package utilities

import (
	"net/http"

	"Backend/models"
)

func newResponse(code int, message string, data interface{}) models.APIResponse {
	return models.APIResponse{
		ResponseCode: code,
		ResponseBody: models.ResponseBody{
			MESSAGE: message,
			DATA:    data,
		},
	}
}

func SetResponseOk(message string, data ...interface{}) models.APIResponse {
	var d interface{}
	if len(data) == 0 {
		d = map[string]interface{}{}
	} else {
		d = data[0]
	}
	return newResponse(http.StatusOK, message, d)
}

func SetResponseNotFound(message string, data ...interface{}) models.APIResponse {
	var d interface{}
	if len(data) == 0 {
		d = map[string]interface{}{}
	} else {
		d = data[0]
	}
	return newResponse(http.StatusNotFound, message, d)
}

func SetResponseBadRequest(message string, data ...interface{}) models.APIResponse {
	var d interface{}
	if len(data) == 0 {
		d = map[string]interface{}{}
	} else {
		d = data[0]
	}
	return newResponse(http.StatusBadRequest, message, d)
}

func SetResponseInternalError(message string, data ...interface{}) models.APIResponse {
	var d interface{}
	if len(data) == 0 {
		d = map[string]interface{}{}
	} else {
		d = data[0]
	}
	return newResponse(http.StatusInternalServerError, message, d)
}

func SetResponseUnauth(data ...interface{}) models.APIResponse {
	var d interface{}
	if len(data) == 0 {
		d = map[string]interface{}{}
	} else {
		d = data[0]
	}
	return newResponse(http.StatusUnauthorized, "Unauthorized Access Denied !!", d)
}

func SetResponseTimedOut(message string, data ...interface{}) models.APIResponse {
	var d interface{}
	if len(data) == 0 {
		d = map[string]interface{}{}
	} else {
		d = data[0]
	}
	return newResponse(http.StatusRequestTimeout, message, d)
}

func SetResponseTransactionFailed(message string, data ...interface{}) models.APIResponse {
	var d interface{}
	if len(data) == 0 {
		d = map[string]interface{}{}
	} else {
		d = data[0]
	}
	return newResponse(http.StatusAccepted, message, d)
}