package database

import "github.com/vongdatcuong/music-streaming-authentication/internal/modules/user"

func convertUserSchemaToUser(schema UserSchema) user.User {
	var perms []string

	for _, permSchema := range schema.Permissions {
		perms = append(perms, permSchema.Name)
	}

	return user.User{
		UserID:      schema.UserID,
		Email:       schema.Email,
		FirstName:   schema.FirstName,
		LastName:    schema.LastName,
		Status:      schema.Status,
		NewSongNoti: schema.NewSongNoti,
		CreatedAt:   schema.CreatedAt,
		UpdatedAt:   schema.UpdatedAt,
		Permissions: perms,
	}
}

func convertUserSchemaCreateToUser(schema UserSchemaCreate) user.User {
	return user.User{
		UserID:      schema.UserID,
		Email:       schema.Email,
		FirstName:   schema.FirstName,
		LastName:    schema.LastName,
		Status:      schema.Status,
		NewSongNoti: *schema.NewSongNoti,
		CreatedAt:   schema.CreatedAt,
		UpdatedAt:   schema.UpdatedAt,
	}
}

func convertUserSchemaPutToUser(schema UserSchemaPut, curUser user.User) user.User {
	return user.User{
		UserID:      schema.UserID,
		Email:       schema.Email,
		FirstName:   schema.FirstName,
		LastName:    schema.LastName,
		Status:      schema.Status,
		NewSongNoti: *schema.NewSongNoti,
		CreatedAt:   curUser.CreatedAt,
		UpdatedAt:   schema.UpdatedAt,
	}
}
