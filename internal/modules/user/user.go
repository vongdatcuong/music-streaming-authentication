package user

import (
	"context"
	"fmt"

	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/common"
	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/constants"
)

type UserStore interface {
	GetUserList(context.Context, common.PaginationInfo, UserListFilter) ([]User, uint64, error)
	GetUserDetails(context.Context, uint64) (User, error)
	CreateUser(context.Context, User) (User, error)
	PutUser(context.Context, User) (User, error)
	UpdateUserStatus(context.Context, uint64, constants.ACTIVE_STATUS) error
	DoesUserExist(context.Context, uint64) (bool, error)
	UpdateUserPermissions(context.Context, uint64, []uint64, []uint64) error
}

type UserService struct {
	store UserStore
}

type User struct {
	UserID      uint64
	Email       string
	FirstName   string
	LastName    string
	Status      constants.ACTIVE_STATUS
	Password    string
	NewSongNoti bool
	CreatedAt   uint64
	UpdatedAt   uint64
	Permissions []string
}

type UserListFilter struct {
	UserID          uint64
	Email           string
	Status          constants.ACTIVE_STATUS
	CreatedTimeFrom uint64
	CreatedTimeTo   uint64
}

func NewService(store UserStore) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) GetUserList(ctx context.Context, paginationInfo common.PaginationInfo, filter UserListFilter) ([]User, uint64, error) {
	userList, totalCount, err := s.store.GetUserList(ctx, paginationInfo, filter)

	if err != nil {
		return []User{}, 0, err
	}

	return userList, totalCount, nil
}

func (s *UserService) GetUserDetails(ctx context.Context, id uint64) (User, error) {
	user, err := s.store.GetUserDetails(ctx, id)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, newUser User) (User, error) {
	wrappedUser := User(newUser)
	wrappedUser.Status = constants.ACTIVE_STATUS_ACTIVE

	user, err := s.store.CreateUser(ctx, wrappedUser)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *UserService) PutUser(ctx context.Context, curUser User) (User, error) {
	userID := curUser.UserID
	doesExist, err := s.store.DoesUserExist(ctx, userID)

	if err != nil {
		return User{}, err
	}

	if !doesExist {
		return User{}, fmt.Errorf("could not find user with id %d", userID)
	}

	user, err := s.store.PutUser(ctx, curUser)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *UserService) UpdateUserStatus(ctx context.Context, userID uint64, status constants.ACTIVE_STATUS) error {
	doesExist, err := s.store.DoesUserExist(ctx, userID)

	if err != nil {
		return err
	}

	if !doesExist {
		return fmt.Errorf("could not find user with id %d", userID)
	}

	err = s.store.UpdateUserStatus(ctx, userID, status)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateUserPermissions(ctx context.Context, userID uint64, addedPermissionIDs []uint64, deletedPermissionIDs []uint64) error {
	doesExist, err := s.store.DoesUserExist(ctx, userID)

	if err != nil {
		return err
	}

	if !doesExist {
		return fmt.Errorf("could not find user with id %d", userID)
	}

	err = s.store.UpdateUserPermissions(ctx, userID, addedPermissionIDs, deletedPermissionIDs)

	if err != nil {
		return err
	}

	return nil
}
