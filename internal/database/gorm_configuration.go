package database

const PermissionTableName = "permission"
const UserTableName = "user"
const UserPermissionTableName = "user_permission"

type Tabler interface {
	TableName() string
}

func (PermissionSchema) TableName() string {
	return PermissionTableName
}

func (PermissionSchemaCreate) TableName() string {
	return PermissionTableName
}

func (PermissionSchemaPut) TableName() string {
	return PermissionTableName
}

func (UserSchema) TableName() string {
	return UserTableName
}

func (UserSchemaCreate) TableName() string {
	return UserTableName
}

func (UserSchemaPut) TableName() string {
	return UserTableName
}

func (UpdateUserStatusSchema) TableName() string {
	return UserTableName
}

func (UserPermissionSchema) TableName() string {
	return UserPermissionTableName
}
