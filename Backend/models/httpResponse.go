package models

type APIResponse struct {
	ResponseCode int         `json:"responseCode"`
	ResponseBody ResponseBody `json:"responseBody"`
}

type ResponseBody struct {
	MESSAGE string      `json:"MESSAGE"`
	DATA    interface{} `json:"DATA"`
}