package token

type TokenService interface {
	IssueJWTToken()
	RefreshJWTToken()
}
