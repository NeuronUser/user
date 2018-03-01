package models

type Token struct {
	AccessToken  string
	RefreshToken string
}

type OauthJumpResponse struct {
	UserID      string
	Token       *Token
	QueryString string
}
