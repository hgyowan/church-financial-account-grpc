package workspace

import (
	pkgCrypto "github.com/hgyowan/go-pkg-library/crypto"
	"gorm.io/gorm"
	"time"
)

type Workspace struct {
	ID                      string     `gorm:"column:id;type:varchar(36);primaryKey"`
	Name                    string     `gorm:"column:name;type:varchar(64);not null"`
	Description             *string    `gorm:"column:description;type:varchar(256);default:''"`
	OwnerID                 string     `gorm:"column:owner_id;type:varchar(36);not null"`
	RepresentativeName      string     `gorm:"column:representative_name;type:varchar(256);not null" crypto:"type:fixed_cbc;context:representative_name"`
	RepresentativePhone     string     `gorm:"column:representative_phone;type:varchar(256);not null" crypto:"type:fixed_cbc;context:representative_phone"`
	RepresentativeEmail     string     `gorm:"column:representative_email;type:varchar(254);not null" crypto:"type:fixed_cbc;context:representative_email"`
	Address1                string     `gorm:"column:address1;type:varchar(256);not null" crypto:"type:fixed_cbc;context:address1"`
	Address2                string     `gorm:"column:address2;type:varchar(256);not null" crypto:"type:fixed_cbc;context:address2"`
	ZipCode                 string     `gorm:"column:zip_code;type:varchar(256);not null" crypto:"type:fixed_cbc;context:zip_code"`
	ThumbnailURL            *string    `gorm:"column:thumbnail_url;type:varchar(512);not null"`
	SealURL                 *string    `gorm:"column:seal_url;type:varchar(512);default:''"`
	BusinessRegistrationURL *string    `gorm:"column:business_registration_url;type:varchar(512);default:''"`
	BusinessRegistrationNum *string    `gorm:"column:business_registration_num;type:varchar(256);default:''" crypto:"type:fixed_cbc;context:business_registration_num"`
	CreatedAt               time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt               *time.Time `gorm:"column:updated_at"`
}

func (w *Workspace) TableName() string {
	return "workspaces"
}

func (w *Workspace) BeforeCreate(*gorm.DB) error {
	return pkgCrypto.EncryptScheme(w)
}

func (w *Workspace) AfterCreate(*gorm.DB) error {
	return pkgCrypto.DecryptScheme(w)
}

func (w *Workspace) BeforeUpdate(*gorm.DB) error {
	return pkgCrypto.EncryptScheme(w)
}

func (w *Workspace) AfterUpdate(*gorm.DB) error {
	return pkgCrypto.DecryptScheme(w)
}

func (w *Workspace) AfterFind(*gorm.DB) error {
	return pkgCrypto.DecryptScheme(w)
}

type WorkspaceUser struct {
	WorkspaceID string     `gorm:"column:workspace_id;type:varchar(36);primaryKey"`
	UserID      string     `gorm:"column:user_id;type:varchar(36);primaryKey"`
	IsAdmin     bool       `gorm:"column:is_admin;type:bool;not null;default:false"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
}

func (wu *WorkspaceUser) TableName() string {
	return "workspace_users"
}

type WorkspaceInvite struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement"`
	WorkspaceID string    `gorm:"column:workspace_id;type:varchar(36);not null"`
	UserID      string    `gorm:"column:user_id;type:varchar(36);not null"`
	Message     string    `gorm:"column:message;type:varchar(256);default:''"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
}

func (wi *WorkspaceInvite) TableName() string {
	return "workspace_invites"
}

type WorkspaceSimple struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	IsOwner  bool      `json:"isOwner"`
	IsAdmin  bool      `json:"isAdmin"`
	JoinedAt time.Time `json:"joinedAt"`
}
