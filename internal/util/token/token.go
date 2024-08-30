package tokenutil

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewToken(accessToken, refreshToken string) *Token {
	return &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
