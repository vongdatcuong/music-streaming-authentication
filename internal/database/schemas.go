package database

import "github.com/vongdatcuong/music-streaming-authentication/internal/modules/constants"

// PERMISSION
type PermissionSchema struct {
	PermissionID uint64                  `gorm:"column:permission_id;primaryKey"`
	Name         string                  `gorm:"column:name"`
	CreatedAt    uint64                  `gorm:"column:created_at"`
	UpdatedAt    uint64                  `gorm:"column:updated_at"`
	Status       constants.ACTIVE_STATUS `gorm:"column:status"`
}

type PermissionSchemaCreate struct {
	PermissionID uint64 `gorm:"column:permission_id;primaryKey"`
	Name         string `gorm:"column:name" validate:"required,max=256"`
	CreatedAt    uint64
	UpdatedAt    uint64
	Status       constants.ACTIVE_STATUS `gorm:"column:status" validate:"required"`
}

type PermissionSchemaPut struct {
	PermissionID uint64 `gorm:"column:permission_id;primaryKey"`
	Name         string `gorm:"column:name" validate:"required,max=256"`
	UpdatedAt    uint64
	Status       constants.ACTIVE_STATUS `gorm:"column:status" validate:"required"`
}

// USER
type UserSchema struct {
	UserID      uint64                  `gorm:"column:user_id;primaryKey"`
	Email       string                  `gorm:"column:email"`
	FirstName   string                  `gorm:"column:first_name"`
	LastName    string                  `gorm:"column:last_name"`
	Status      constants.ACTIVE_STATUS `gorm:"column:status"`
	Password    string                  `gorm:"column:password"`
	NewSongNoti string                  `gorm:"column:new_song_noti"`
	CreatedAt   uint64                  `gorm:"column:created_at"`
	UpdatedAt   uint64                  `gorm:"column:updated_at"`
}

// USER_PERMISSION
type UserPermissionSchema struct {
	PermissionID uint64 `gorm:"column:permission_id;primaryKey"`
	UserID       uint64 `gorm:"column:user_id;primaryKey"`
	CreatedAt    uint64 `gorm:"column:created_at"`
	UpdatedAt    uint64 `gorm:"column:updated_at"`
}
