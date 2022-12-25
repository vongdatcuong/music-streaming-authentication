package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/common"
	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/constants"
	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/user"
	time_utils "github.com/vongdatcuong/music-streaming-authentication/internal/modules/utils/time"
	validator_utils "github.com/vongdatcuong/music-streaming-authentication/internal/modules/utils/validator"
	"gorm.io/gorm"
)

func (db *Database) GetUserList(ctx context.Context, paginationInfo common.PaginationInfo, filter user.UserListFilter) ([]user.User, uint64, error) {
	var userSchemas []UserSchema
	var totalCount int64

	result := db.GormClient.WithContext(ctx).Preload("Permissions").Order("created_at desc, user_id").
		Where(UserSchema{UserID: filter.UserID, Email: filter.Email, Status: filter.Status}).Find(&userSchemas)

	if filter.CreatedTimeFrom != 0 {
		result = result.Where("created_at >= ?", filter.CreatedTimeFrom)
	}

	if filter.CreatedTimeTo != 0 {
		result = result.Where("created_at <= ?", filter.CreatedTimeTo)
	}

	result.Count(&totalCount)

	if result.Error != nil {
		return []user.User{}, 0, fmt.Errorf("could not count the total number of users: %w", result.Error)
	}

	result = result.Scopes(Paginate(paginationInfo)).Find(&userSchemas)

	if result.Error != nil {
		return []user.User{}, 0, fmt.Errorf("could not get the user list: %w", result.Error)
	}

	var users []user.User

	for _, schema := range userSchemas {
		users = append(users, convertUserSchemaToUser(schema, false))
	}

	return users, uint64(totalCount), nil
}

func (db *Database) GetUserDetails(ctx context.Context, id uint64) (user.User, error) {
	var record UserSchema

	result := db.GormClient.WithContext(ctx).Preload("Permissions").First(&record, id)
	if result.Error != nil {
		return user.User{}, fmt.Errorf("could not get user details: %w", result.Error)
	}

	return convertUserSchemaToUser(record, false), nil
}

func (db *Database) CreateUser(ctx context.Context, newUser user.User) (user.User, error) {
	userSchemaCreate := UserSchemaCreate{
		UserID:      newUser.UserID,
		Email:       newUser.Email,
		FirstName:   newUser.FirstName,
		LastName:    newUser.LastName,
		Status:      newUser.Status,
		Password:    newUser.Password,
		NewSongNoti: &newUser.NewSongNoti,
		CreatedAt:   time_utils.GetCurrentUnixTime(),
		UpdatedAt:   time_utils.GetCurrentUnixTime(),
	}

	err := validator_utils.ValidateStruct(userSchemaCreate)

	if err != nil {
		return user.User{}, fmt.Errorf("user is invalid: %w", err)
	}

	result := db.GormClient.WithContext(ctx).Create(&userSchemaCreate)

	if result.Error != nil {
		return user.User{}, fmt.Errorf("could not create user: %w", result.Error)
	}

	return convertUserSchemaCreateToUser(userSchemaCreate), nil
}

func (db *Database) PutUser(ctx context.Context, curUser user.User) (user.User, error) {
	userSchemaPut := UserSchemaPut{
		UserID:      curUser.UserID,
		Email:       curUser.Email,
		FirstName:   curUser.FirstName,
		LastName:    curUser.LastName,
		Status:      curUser.Status,
		NewSongNoti: &curUser.NewSongNoti,
		UpdatedAt:   time_utils.GetCurrentUnixTime(),
	}

	err := validator_utils.ValidateStruct(userSchemaPut)

	if err != nil {
		return user.User{}, fmt.Errorf("user is invalid: %w", err)
	}

	result := db.GormClient.WithContext(ctx).Updates(&userSchemaPut)

	if result.Error != nil {
		return user.User{}, fmt.Errorf("could not put user: %w", result.Error)
	}

	return convertUserSchemaPutToUser(userSchemaPut, curUser), nil
}

func (db *Database) UpdateUserStatus(ctx context.Context, userID uint64, status constants.ACTIVE_STATUS) error {
	schema := UpdateUserStatusSchema{
		UserID: userID,
		Status: status,
	}

	err := validator_utils.ValidateStruct(schema)

	if err != nil {
		return fmt.Errorf("params invalid: %w", err)
	}

	result := db.GormClient.WithContext(ctx).Updates(&schema)

	if result.Error != nil {
		return fmt.Errorf("could not update user status: %w", result.Error)
	}

	return nil
}

func (db *Database) DoesUserExist(ctx context.Context, id uint64) (bool, error) {
	var exists bool

	result := db.GormClient.Table(UserTableName).Select("count(*) > 0").Where("user_id = ?", id).Find(&exists)

	if result.Error != nil {
		return false, fmt.Errorf("could not check if user exists: %w", result.Error)
	}

	return exists, nil
}

func (db *Database) UpdateUserPermissions(ctx context.Context, userID uint64, addedPermissionIDs []uint64, deletedPermissionIDs []uint64) error {
	err := db.AddPermissionOfUser(ctx, userID, addedPermissionIDs)

	if err != nil {
		return err
	}

	err = db.DeletePermissionOfUser(ctx, userID, deletedPermissionIDs)

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) AddPermissionOfUser(ctx context.Context, userID uint64, permIDs []uint64) error {
	if permIDs == nil || len(permIDs) == 0 {
		return nil
	}

	var addedUserPermSchemas []UserPermissionSchema

	for _, addedID := range permIDs {
		addedUserPermSchemas = append(addedUserPermSchemas, UserPermissionSchema{UserID: userID, PermissionID: addedID,
			CreatedAt: time_utils.GetCurrentUnixTime(), UpdatedAt: time_utils.GetCurrentUnixTime()})
	}
	result := db.GormClient.WithContext(ctx).Create(addedUserPermSchemas)

	if result.Error != nil {
		return fmt.Errorf("could not add permissions to user with id %d: %w", userID, result.Error)
	}

	return nil
}

func (db *Database) DeletePermissionOfUser(ctx context.Context, userID uint64, permIDs []uint64) error {
	if permIDs == nil || len(permIDs) == 0 {
		return nil
	}

	var deletedUserPermSchemas []UserPermissionSchema

	for _, deletedID := range permIDs {
		deletedUserPermSchemas = append(deletedUserPermSchemas, UserPermissionSchema{UserID: userID, PermissionID: deletedID})
	}
	result := db.GormClient.WithContext(ctx).Delete(deletedUserPermSchemas)

	if result.Error != nil {
		return fmt.Errorf("could not delete permissions of user with id %d: %w", userID, result.Error)
	}

	return nil
}

func (db *Database) LogIn(ctx context.Context, loginUser user.User) (user.User, error) {
	var userSchema UserSchema

	result := db.GormClient.Model(&UserSchema{}).Where(UserSchema{Email: loginUser.Email}).First(&userSchema)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user.User{}, fmt.Errorf("could not login user: %w", result.Error)
	}

	return convertUserSchemaToUser(userSchema, true), nil
}
