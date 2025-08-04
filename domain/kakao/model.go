package kakao

type Profile struct {
	Nickname          string `json:"nickname"`
	ThumbnailImageURL string `json:"thumbnail_image_url"`
	ProfileImageURL   string `json:"profile_image_url"`
}

type KakaoAccount struct {
	Profile         Profile `json:"profile"`
	IsEmailVerified bool    `json:"is_email_verified"` // 이메일 인증 여부
	Email           string  `json:"email"`
}
