package grpc

import (
	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/user"
	grpcPbV1 "github.com/vongdatcuong/music-streaming-protos/go/v1"
)

func convertUserToGrpcUser(user user.User) *grpcPbV1.User {
	return &grpcPbV1.User{
		UserId:      user.UserID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Status:      uint32(user.Status),
		NewSongNoti: user.NewSongNoti,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Permissions: user.Permissions,
	}
}
