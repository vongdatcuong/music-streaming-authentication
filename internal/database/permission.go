package database

import (
	"context"
	"fmt"

	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/constants"
	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/permission"
	time_utils "github.com/vongdatcuong/music-streaming-authentication/internal/modules/utils/time"
	validator_utils "github.com/vongdatcuong/music-streaming-authentication/internal/modules/utils/validator"
	"gorm.io/gorm"
)

type PermissionRowCreate struct {
	ID        uint64
	Name      string `validate:"required,max=256"`
	CreatedAt uint64
	UpdatedAt uint64
	Status    constants.ACTIVE_STATUS `validate:"required"`
}

type PermissionRowPut struct {
	Name      string `validate:"required,max=256"`
	UpdatedAt uint64
	Status    constants.ACTIVE_STATUS `validate:"required"`
}

func (db *Database) GetPermissionList(ctx context.Context) ([]permission.Permission, error) {
	var permissions []permission.Permission

	result := db.GormClient.Table(PermissionTableName).Find(&permissions)

	if result.Error != nil {
		return []permission.Permission{}, fmt.Errorf("could not get the permission list: %w", result.Error)
	}

	return permissions, nil

}

func (db *Database) CreatePermission(ctx context.Context, newPerm permission.Permission) (permission.Permission, error) {
	permRowCreate := PermissionRowCreate{
		Name:      newPerm.Name,
		Status:    newPerm.Status,
		CreatedAt: time_utils.GetCurrentUnixTime(),
		UpdatedAt: time_utils.GetCurrentUnixTime(),
	}

	if err := validator_utils.ValidateStruct(permRowCreate); err != nil {
		return permission.Permission{}, fmt.Errorf("permission is not valid: %w", err)
	}

	result := db.GormClient.Table(PermissionTableName).Create(&permRowCreate)

	if result.Error != nil {
		return permission.Permission{}, fmt.Errorf("could not create new permission: %w", result.Error)
	}

	return convertPermissionRowCreateToPermission(permRowCreate), nil
}

func (db *Database) PutPermission(ctx context.Context, existingPerm permission.Permission) (permission.Permission, error) {
	permissionRowPut := PermissionRowPut{
		Name:      existingPerm.Name,
		UpdatedAt: time_utils.GetCurrentUnixTime(),
		Status:    existingPerm.Status,
	}

	if err := validator_utils.ValidateStruct(permissionRowPut); err != nil {
		return permission.Permission{}, fmt.Errorf("permission is not valid: %w", err)
	}

	result := db.GormClient.Table(PermissionTableName).Where("permission_id = ?", existingPerm.PermissionID).Save(&permissionRowPut)

	if result.Error != nil {
		return permission.Permission{}, fmt.Errorf("could not update permission: %w", result.Error)
	}
	return convertPermissionRowPutToPermission(existingPerm, permissionRowPut), nil
}

// TODO: Test this after implementing User mmodule
func (db *Database) CheckUserPermission(ctx context.Context, userID uint64, perm permission.Permission) (bool, error) {
	if perm.PermissionID == 0 && perm.Name == "" {
		return false, fmt.Errorf("no permission is provided")
	}

	var result *gorm.DB
	var exists bool
	var permID uint64

	if perm.PermissionID == 0 {
		// Check by permission name
		permResult := db.GormClient.Table(PermissionTableName).Select("permission_id").Where("name = ?", perm.Name).Scan(&permID)
		if permResult.Error != nil || permID == 0 {
			return false, fmt.Errorf("could find permission by permission name: %w", permResult.Error)
		}
	} else {
		permID = perm.PermissionID
	}

	result = db.GormClient.Table(UserPermissionTableName).Select("count(*) > 0").Where("user_id = ? AND permission_id = ?", userID, permID).Find(&exists)

	if result.Error != nil {
		return false, fmt.Errorf("could not check user's permission: %w", result.Error)
	}

	return exists, nil
}

func (db *Database) DoesPermissionExist(ctx context.Context, id uint64) (bool, error) {
	var exists bool

	result := db.GormClient.Table(PermissionTableName).Select("count(*) > 0").Where("permission_id = ?", id).Find(&exists)

	if result.Error != nil {
		return false, fmt.Errorf("could not check if permission exists: %w", result.Error)
	}

	return exists, nil
}
