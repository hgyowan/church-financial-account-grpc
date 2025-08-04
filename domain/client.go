package domain

import "github.com/hgyowan/church-financial-account-grpc/domain/kakao"

type Client interface {
	kakao.KakaoClient
}
