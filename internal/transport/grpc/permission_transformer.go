package grpc

import (
	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/permission"
	grpcPbV1 "github.com/vongdatcuong/music-streaming-protos/go/v1"
)

func convertPermissionToGrpcPermission(permission permission.Permission) *grpcPbV1.Permission {
	return &grpcPbV1.Permission{
		PermissionId: permission.PermissionID,
		Name:         permission.Name,
		CreatedAt:    permission.CreatedAt,
		UpdatedAt:    permission.UpdatedAt,
		Status:       uint32(permission.Status),
	}
}
