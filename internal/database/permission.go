package database

import (
	"context"
	"fmt"

	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/permission"
	time_utils "github.com/vongdatcuong/music-streaming-authentication/internal/modules/utils/time"
	validator_utils "github.com/vongdatcuong/music-streaming-authentication/internal/modules/utils/validator"
	"gorm.io/gorm"
)

func (db *Database) GetPermissionList(ctx context.Context) ([]permission.Permission, error) {
	var permissionSchemas []PermissionSchema

	result := db.GormClient.WithContext(ctx).Find(&permissionSchemas)

	if result.Error != nil {
		return []permission.Permission{}, fmt.Errorf("could not get the permission list: %w", result.Error)
	}

	var permissions []permission.Permission

	for _, schema := range permissionSchemas {
		permissions = append(permissions, convertPermissionSchemaToPermission(schema))
	}

	return permissions, nil

}

func (db *Database) CreatePermission(ctx context.Context, newPerm permission.Permission) (permission.Permission, error) {
	permSchemaCreate := PermissionSchemaCreate{
		Name:      newPerm.Name,
		Status:    newPerm.Status,
		CreatedAt: time_utils.GetCurrentUnixTime(),
		UpdatedAt: time_utils.GetCurrentUnixTime(),
	}

	if err := validator_utils.ValidateStruct(permSchemaCreate); err != nil {
		return permission.Permission{}, fmt.Errorf("permission is invalid: %w", err)
	}

	result := db.GormClient.WithContext(ctx).Create(&permSchemaCreate)

	if result.Error != nil {
		return permission.Permission{}, fmt.Errorf("could not create new permission: %w", result.Error)
	}

	return convertPermissionSchemaCreateToPermission(permSchemaCreate), nil
}

func (db *Database) PutPermission(ctx context.Context, existingPerm permission.Permission) (permission.Permission, error) {
	permSchemaPut := PermissionSchemaPut{
		PermissionID: existingPerm.PermissionID,
		Name:         existingPerm.Name,
		UpdatedAt:    time_utils.GetCurrentUnixTime(),
		Status:       existingPerm.Status,
	}

	if err := validator_utils.ValidateStruct(permSchemaPut); err != nil {
		return permission.Permission{}, fmt.Errorf("permission is not valid: %w", err)
	}

	result := db.GormClient.WithContext(ctx).Updates(permSchemaPut)

	if result.Error != nil {
		return permission.Permission{}, fmt.Errorf("could not put permission: %w", result.Error)
	}

	return convertPermissionSchemaPutToPermission(existingPerm, permSchemaPut), nil
}

func (db *Database) CheckUserPermission(ctx context.Context, userID uint64, perm permission.Permission) (bool, error) {
	if perm.PermissionID == 0 && perm.Name == "" {
		return false, fmt.Errorf("no permission is provided")
	}

	var result *gorm.DB
	var exists bool
	var permID uint64

	if perm.PermissionID == 0 {
		// Check by permission name
		permResult := db.GormClient.Model(PermissionSchema{}).Select("permission_id").Where(PermissionSchema{Name: perm.Name}).Scan(&permID)
		if permResult.Error != nil {
			return false, fmt.Errorf("could find permission by permission name: %w", permResult.Error)
		} else if permID == 0 {
			return false, fmt.Errorf("permission does not exist")
		}
	} else {
		permID = perm.PermissionID
	}

	result = db.GormClient.Model(UserPermissionSchema{}).Select("count(*) > 0").Where(UserPermissionSchema{PermissionID: permID, UserID: userID}).Find(&exists)

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
