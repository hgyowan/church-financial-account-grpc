package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/hgyowan/church-financial-account-grpc/domain/token"
	"github.com/hgyowan/church-financial-account-grpc/internal"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	"time"
)

func registerTokenService(s *service) {
	s.TokenService = &tokenService{s: s}
}

type tokenService struct {
	s *service
}

func (t *tokenService) generateClaim(request token.GenerateClaimRequest) token.JWTCustomClaims {
	now := time.Now().UTC()
	tokenClaim := token.JWTCustomClaims{
		UserID: request.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(internal.AccessTokenExpiredTime).Unix(),
			IssuedAt:  now.Unix(),
			Issuer:    envs.JwtIssuer,
		},
	}

	return tokenClaim
}

func (t *tokenService) IssueJWTToken(ctx context.Context, request token.IssueJWTTokenRequest) (*token.IssueJWTTokenResponse, error) {
	accessTokenClaim := t.generateClaim(token.GenerateClaimRequest{
		UserID:             request.UserID,
		ExpireTimeDuration: internal.AccessTokenExpiredTime,
		Secret:             envs.JwtAccessSecret,
	})
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaim)
	accessTokenString, err := accessToken.SignedString([]byte(envs.JwtAccessSecret))
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Create)
	}

	refreshTokenClaim := t.generateClaim(token.GenerateClaimRequest{
		UserID:             request.UserID,
		ExpireTimeDuration: internal.RefreshTokenExpiredTime,
		Secret:             envs.JwtRefreshSecret,
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaim)
	refreshTokenString, err := refreshToken.SignedString([]byte(envs.JwtRefreshSecret))
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Create)
	}

	tk := token.JWTToken{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}

	b, err := json.Marshal(tk)
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Create)
	}

	if err = t.s.externalRedisClient.Redis().Set(ctx, fmt.Sprintf("user_token:%s", request.UserID), b, internal.RefreshTokenExpiredTime).Err(); err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Create)
	}

	return &token.IssueJWTTokenResponse{
		JWTToken: tk,
	}, nil
}

func (t *tokenService) RefreshJWTToken(ctx context.Context, request token.RefreshJWTTokenRequest) (*token.RefreshJWTTokenResponse, error) {
	beforeToken, err := t.s.externalRedisClient.Redis().Get(ctx, fmt.Sprintf("user_token:%s", request.UserID)).Result()
	if err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Get)
	}

	var beforeJWTToken token.JWTToken
	if err = json.Unmarshal([]byte(beforeToken), &beforeJWTToken); err != nil {
		return nil, pkgError.WrapWithCode(err, pkgError.Get)
	}

	if request.RefreshToken != beforeJWTToken.RefreshToken {
		return nil, pkgError.WrapWithCode(pkgError.EmptyBusinessError(), pkgError.Expired)
	}

	tk, err := t.IssueJWTToken(ctx, token.IssueJWTTokenRequest{
		UserID: request.UserID,
	})
	if err != nil {
		return nil, pkgError.Wrap(err)
	}

	return &token.RefreshJWTTokenResponse{
		JWTToken: tk.JWTToken,
	}, nil
}
