package common

type successRes struct {
	Data   interface{} `json:"data"`
	Paging interface{} `json:"paging"`
	Filter interface{} `json:"filter"`
}

func NewSuccesResponse(data, paging, filter interface{}) *successRes {
	return &successRes{Data: data, Paging: paging, Filter: filter}
}

func SimpleSuccesResponse(data interface{}) *successRes {
	return NewSuccesResponse(data, nil, nil)
}

type userResponesToken struct {
	AccessToken  interface{} `json:"accesstoken"`
	RefreshToken interface{} `json:"refreshtoken"`
	Data         interface{} `json:"data"`
}

func NewReponseUserToken(accessToken, refreshToken string, data interface{}) *userResponesToken {
	return &userResponesToken{
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
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func NewDataResponse(data interface{}, message string) *dataRes {
	return &dataRes{Data: data, Message: message}
}
