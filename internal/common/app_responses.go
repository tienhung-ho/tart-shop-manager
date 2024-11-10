package common

import "net/http"

type successRes struct {
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
	Paging     interface{} `json:"paging"`
	Filter     interface{} `json:"filter"`
}

func NewSuccesResponse(data, paging, filter interface{}) *successRes {
	return &successRes{StatusCode: http.StatusOK, Data: data, Paging: paging, Filter: filter}
}

func SimpleSuccesResponse(data interface{}) *successRes {
	return NewSuccesResponse(data, nil, nil)
}

type userResponesToken struct {
	StatusCode   int         `json:"status_code"`
	AccessToken  interface{} `json:"accesstoken"`
	RefreshToken interface{} `json:"refreshtoken"`
	Data         interface{} `json:"data"`
}

func NewReponseUserToken(accessToken, refreshToken string, data interface{}) *userResponesToken {
	return &userResponesToken{
		StatusCode:   http.StatusOK,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Data:         data,
	}
}

type errTokenRespone struct {
	Error interface{} `json:"error"`
	Type  interface{} `json:"type"`
}

func NewReponseErrToken(errToken, tokenType string) *errTokenRespone {
	return &errTokenRespone{
		Error: errToken,
		Type:  tokenType,
	}
}

type dataRes struct {
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
	Message    string      `json:"message"`
}

func NewDataResponse(data interface{}, message string) *dataRes {
	return &dataRes{StatusCode: http.StatusOK, Data: data, Message: message}
}
