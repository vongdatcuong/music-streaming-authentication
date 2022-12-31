package permission

import (
	"context"
	"fmt"

	constants "github.com/vongdatcuong/music-streaming-authentication/internal/modules/constants"
)

type PermissionStore interface {
	GetPermissionList(context.Context) ([]Permission, error)
	CreatePermission(context.Context, Permission) (Permission, error)
	PutPermission(context.Context, Permission) (Permission, error)
	CheckUserPermission(context.Context, uint64, Permission) (bool, error)
	DoesPermissionExist(context.Context, uint64) (bool, error)
}

type UserService interface {
	DoesUserExist(context.Context, uint64) (bool, error)
}

type PermissionService struct {
	store       PermissionStore
	userService UserService
}

type Permission struct {
	PermissionID uint64
	Name         string
	CreatedAt    uint64
	UpdatedAt    uint64
	Status       constants.ACTIVE_STATUS
}

func NewService(store PermissionStore, userService UserService) *PermissionService {
	return &PermissionService{
		store:       store,
		userService: userService,
	}
}

func (s *PermissionService) GetPermissionList(ctx context.Context) ([]Permission, error) {
	permissionList, err := s.store.GetPermissionList(ctx)

	if err != nil {
		return []Permission{}, err
	}

	return permissionList, nil
}

func (s *PermissionService) CreatePermission(ctx context.Context, newPerm Permission) (Permission, error) {
	wrappedPerm := Permission(newPerm)
	wrappedPerm.Status = constants.ACTIVE_STATUS_ACTIVE

	permission, err := s.store.CreatePermission(ctx, wrappedPerm)

	if err != nil {
		return Permission{}, err
	}

	return permission, nil
}

func (s *PermissionService) PutPermission(ctx context.Context, existingPerm Permission) (Permission, error) {
	permID := existingPerm.PermissionID
	doesExist, err := s.store.DoesPermissionExist(ctx, permID)

	if err != nil {
		return Permission{}, err
	}

	if !doesExist {
		return Permission{}, fmt.Errorf("could not find permission with id %d", permID)
	}

	permission, err := s.store.PutPermission(ctx, existingPerm)

	if err != nil {
		return Permission{}, err
	}

	return permission, nil
}

// TODO: Remove check user exist completely
func (s *PermissionService) CheckUserPermission(ctx context.Context, userID uint64, perm Permission) (bool, bool, error) {
	/*doesUserExist, err := s.userService.DoesUserExist(ctx, userID)

	if err != nil {
		return false, false, err
	}

	if !doesUserExist {
		return false, false, fmt.Errorf("user does not exist")
	}*/
	// perm empty => endpoint doesn't require any permission
	if perm.Name == "" && perm.PermissionID == 0 {
		return true, true, nil
	}

	hasPerm, err := s.store.CheckUserPermission(ctx, userID, perm)

	if err != nil {
		return false, false, err
	}

	return hasPerm, true, nil
}
