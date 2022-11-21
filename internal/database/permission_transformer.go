package database

import "github.com/vongdatcuong/music-streaming-authentication/internal/modules/permission"

func convertPermissionRowCreateToPermission(permRowCreate PermissionRowCreate) permission.Permission {
	return permission.Permission{
		PermissionID: permRowCreate.ID,
		Name:         permRowCreate.Name,
		CreatedAt:    permRowCreate.CreatedAt,
		UpdatedAt:    permRowCreate.UpdatedAt,
		Status:       permRowCreate.Status,
	}
}

func convertPermissionRowPutToPermission(existingPerm permission.Permission, permRowPut PermissionRowPut) permission.Permission {
	return permission.Permission{
		PermissionID: existingPerm.PermissionID,
		Name:         permRowPut.Name,
		CreatedAt:    existingPerm.CreatedAt,
		UpdatedAt:    permRowPut.UpdatedAt,
		Status:       permRowPut.Status,
	}
}
