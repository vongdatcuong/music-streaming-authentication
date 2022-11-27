package database

import "github.com/vongdatcuong/music-streaming-authentication/internal/modules/permission"

func convertPermissionSchemaToPermission(schema PermissionSchema) permission.Permission {
	return permission.Permission{
		PermissionID: schema.PermissionID,
		Name:         schema.Name,
		CreatedAt:    schema.CreatedAt,
		UpdatedAt:    schema.UpdatedAt,
		Status:       schema.Status,
	}
}

func convertPermissionSchemaCreateToPermission(schema PermissionSchemaCreate) permission.Permission {
	return permission.Permission{
		PermissionID: schema.PermissionID,
		Name:         schema.Name,
		CreatedAt:    schema.CreatedAt,
		UpdatedAt:    schema.UpdatedAt,
		Status:       schema.Status,
	}
}

func convertPermissionSchemaPutToPermission(existingPerm permission.Permission, schema PermissionSchemaPut) permission.Permission {
	return permission.Permission{
		PermissionID: existingPerm.PermissionID,
		Name:         schema.Name,
		CreatedAt:    existingPerm.CreatedAt,
		UpdatedAt:    schema.UpdatedAt,
		Status:       schema.Status,
	}
}
