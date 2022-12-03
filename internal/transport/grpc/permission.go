package grpc

import (
	"context"

	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/constants"
	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/permission"
	common_utils "github.com/vongdatcuong/music-streaming-authentication/internal/modules/utils/common"
	grpcPbV1 "github.com/vongdatcuong/music-streaming-authentication/protos/v1/pb"
)

type PermissionServiceGrpc interface {
	GetPermissionList(context.Context) ([]permission.Permission, error)
	CreatePermission(context.Context, permission.Permission) (permission.Permission, error)
	PutPermission(context.Context, permission.Permission) (permission.Permission, error)
	CheckUserPermission(context.Context, uint64, permission.Permission) (bool, error)
}

func (h *Handler) GetPermissionList(ctx context.Context, req *grpcPbV1.GetPermissionListRequest) (*grpcPbV1.GetPermissionListResponse, error) {
	permissions, err := h.permissionService.GetPermissionList(ctx)

	if err != nil {
		return &grpcPbV1.GetPermissionListResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	grpcPermissions := [](*grpcPbV1.Permission){}

	for i := 0; i < len(permissions); i++ {
		grpcPermissions = append(grpcPermissions, convertPermissionToGrpcPermission(permissions[i]))
	}
	return &grpcPbV1.GetPermissionListResponse{
		Data: &grpcPbV1.GetPermissionListResponseData{
			Permissions: grpcPermissions,
		},
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}

func (h *Handler) CreatePermission(ctx context.Context, req *grpcPbV1.CreatePermissionRequest) (*grpcPbV1.CreatePermissionResponse, error) {
	if req.Permission == nil {
		return &grpcPbV1.CreatePermissionResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer("permission must not be empty"),
		}, nil
	}

	newPermission := permission.Permission{
		Name: req.Permission.Name,
	}
	_, err := h.permissionService.CreatePermission(ctx, newPermission)

	if err != nil {
		return &grpcPbV1.CreatePermissionResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	return &grpcPbV1.CreatePermissionResponse{
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}

func (h *Handler) PutPermission(ctx context.Context, req *grpcPbV1.PutPermissionRequest) (*grpcPbV1.PutPermissionResponse, error) {
	if req.Permission == nil {
		return &grpcPbV1.PutPermissionResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer("permission must not be empty"),
		}, nil
	}

	existingPermission := permission.Permission{
		PermissionID: req.Permission.PermissionId,
		Name:         req.Permission.Name,
		Status:       constants.ACTIVE_STATUS(req.Permission.Status),
	}
	_, err := h.permissionService.PutPermission(ctx, existingPermission)

	if err != nil {
		return &grpcPbV1.PutPermissionResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	return &grpcPbV1.PutPermissionResponse{
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}

func (h *Handler) CheckUserPermission(ctx context.Context, req *grpcPbV1.CheckUserPermissionRequest) (*grpcPbV1.CheckUserPermissionResponse, error) {
	hasPerm, err := h.permissionService.CheckUserPermission(ctx, req.UserId, permission.Permission{PermissionID: req.PermissionId, Name: req.PermissionName})

	if err != nil {
		return &grpcPbV1.CheckUserPermissionResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	return &grpcPbV1.CheckUserPermissionResponse{
		HasPermission: &hasPerm,
		Error:         nil,
		ErrorMsg:      nil,
	}, nil
}
