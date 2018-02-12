package models

type Token struct {
	AccessToken  string
	RefreshToken string
}

type OauthJumpResponse struct {
	Token       *Token
	QueryString string
}
