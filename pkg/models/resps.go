package models

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PositiveAPIMapResponse struct {
	Status int                    `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

type PositiveAPIResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type PositiveObjectResponse struct {
	Status int           `json:"status"`
	Data   []interface{} `json:"data"`
}

type NegativeResponse struct {
	Status      int    `json:"status"`
	Description string `json:"description"`
	Subcode     int    `json:"subcode"`
}

func DefaultPosObjectResponse(data []interface{}) PositiveObjectResponse {
	return PositiveObjectResponse{
		Status: 200,
		Data:   data,
	}
}

func DefaultPosMapResponse(data map[string]interface{}) PositiveAPIMapResponse {
	return PositiveAPIMapResponse{
		Status: 200,
		Data:   data,
	}
}

func DefaultPosResponse(data interface{}) PositiveAPIResponse {
	return PositiveAPIResponse{
		Status: 200,
		Data:   data,
	}
}

func DefaultNegResponse(status int, description string, subcode int) NegativeResponse {
	return NegativeResponse{
		Status:      status,
		Description: description,
		Subcode:     subcode,
	}
}

func RespHandler(w http.ResponseWriter, r *http.Request, resp interface{}, err error, status int, action string) {
	if err != nil {
		w.WriteHeader(status)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Printf("[ CRIT ][ %s ] could not encode response: %v\n", action, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
