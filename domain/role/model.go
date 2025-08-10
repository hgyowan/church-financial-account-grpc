package role

import "time"

type AuthRole struct {
	ID          string     `gorm:"column:id;type:varchar(36);primaryKey"`
	Name        string     `gorm:"column:name;type:varchar(64);not null"`
	Description string     `gorm:"column:description;type:varchar(256);default:''"`
	URL         string     `gorm:"column:url;type:varchar(512);not null"`
	Method      string     `gorm:"column:method;type:varchar(7)"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
}

func (ar *AuthRole) TableName() string {
	return "auth_roles"
}

type AuthGroup struct {
	ID          string     `gorm:"column:id;type:varchar(36);primaryKey"`
	Name        string     `gorm:"column:name;type:varchar(64);not null"`
	Description string     `gorm:"column:description;type:varchar(256);default:''"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
}

func (ag *AuthGroup) TableName() string {
	return "auth_groups"
}

type AuthGroupRole struct {
	GroupID   string    `gorm:"column:group_id;type:varchar(36);primaryKey"`
	RoleID    string    `gorm:"column:role_id;type:varchar(36);primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
}

func (agr *AuthGroupRole) TableName() string {
	return "auth_group_roles"
}

type UserRoleGroup struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement"`
	GroupID     string    `gorm:"column:group_id;type:varchar(36);not null"`
	WorkspaceID string    `gorm:"column:workspace_id;type:varchar(36);not null"`
	UserID      string    `gorm:"column:user_id;type:varchar(36);not null"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
}

func (urg *UserRoleGroup) TableName() string {
	return "user_role_groups"
}
