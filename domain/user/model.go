package user

import (
	pkgCrypto "github.com/hgyowan/go-pkg-library/crypto"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID          string     `gorm:"column:id;type:varchar(36);primaryKey;not null"`
	Email       string     `gorm:"column:email;type:varchar(254);unique;not null" crypto:"type:fixed_cbc;context:email"`
	Name        string     `gorm:"column:name;type:varchar(256);not null" crypto:"type:fixed_cbc;context:name"`
	Nickname    string     `gorm:"column:nickname;type:varchar(64);not null;default:''"`
	Password    string     `gorm:"column:password;type:varchar(64);"`
	PhoneNumber string     `gorm:"column:phone_number;type:varchar(256);not null" crypto:"type:fixed_cbc;context:phone_number"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(*gorm.DB) error {
	return pkgCrypto.EncryptScheme(u)
}

func (u *User) AfterCreate(*gorm.DB) error {
	return pkgCrypto.DecryptScheme(u)
}

func (u *User) BeforeUpdate(*gorm.DB) error {
	return pkgCrypto.EncryptScheme(u)
}

func (u *User) AfterUpdate(*gorm.DB) error {
	return pkgCrypto.DecryptScheme(u)
}

func (u *User) AfterFind(*gorm.DB) error {
	return pkgCrypto.DecryptScheme(u)
}

type UserLoginLog struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement;not null"`
	UserID    string    `gorm:"column:user_id;type:varchar(36);not null"`
	IP        *string   `gorm:"column:ip;type:varchar(39)"`
	Browser   *string   `gorm:"column:browser;type:varchar(64)"`
	OS        *string   `gorm:"column:os;type:varchar(64)"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
}

func (ull *UserLoginLog) TableName() string {
	return "user_login_logs"
}

type UserSSO struct {
	UserID         string    `gorm:"column:user_id;type:varchar(36);primaryKey;not null"`
	Provider       string    `gorm:"column:provider;type:varchar(64);not null"`
	ProviderUserID string    `gorm:"column:provider_user_id;type:varchar(128);not null"`
	Email          string    `gorm:"column:email;type:varchar(254);not null" crypto:"type:fixed_cbc;context:email"`
	CreatedAt      time.Time `gorm:"column:created_at;not null"`
}

func (us *UserSSO) TableName() string {
	return "user_sso"
}

func (us *UserSSO) BeforeCreate(*gorm.DB) error {
	return pkgCrypto.EncryptScheme(us)
}

func (us *UserSSO) AfterCreate(*gorm.DB) error {
	return pkgCrypto.DecryptScheme(us)
}

func (us *UserSSO) BeforeUpdate(*gorm.DB) error {
	return pkgCrypto.EncryptScheme(us)
}

func (us *UserSSO) AfterUpdate(*gorm.DB) error {
	return pkgCrypto.DecryptScheme(us)
}

func (us *UserSSO) AfterFind(*gorm.DB) error {
	return pkgCrypto.DecryptScheme(us)
}

type UserConsent struct {
	UserID            string     `gorm:"column:user_id;type:varchar(36);primaryKey;not null"`
	IsTermsAgreed     bool       `gorm:"column:is_terms_agreed;not null;default:true"`
	IsMarketingAgreed *bool      `gorm:"column:is_marketing_agreed;not null;default:false"`
	CreatedAt         time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt         *time.Time `gorm:"column:updated_at"`
}

func (uc *UserConsent) TableName() string {
	return "user_consents"
}

type SSOUser struct {
	SSOUserID string `json:"sso_user_id"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
}
