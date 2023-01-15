package grpc

import (
	"context"

	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/common"
	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/constants"
	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/user"
	common_utils "github.com/vongdatcuong/music-streaming-authentication/internal/modules/utils/common"
	grpcPbV1 "github.com/vongdatcuong/music-streaming-protos/go/v1"
)

type UserServiceGrpc interface {
	GetUserList(context.Context, common.PaginationInfo, user.UserListFilter) ([]user.User, uint64, error)
	GetUserDetails(context.Context, uint64) (user.User, error)
	CreateUser(context.Context, user.User) (user.User, error)
	PutUser(context.Context, user.User) (user.User, error)
	UpdateUserStatus(context.Context, uint64, constants.ACTIVE_STATUS) error
	UpdateUserPermissions(context.Context, uint64, []uint64, []uint64) error
	DoesUserExist(context.Context, uint64) (bool, error)
	LogIn(context.Context, user.User) (user.User, error)
	GetUserListAutocomplete(context.Context, common.PaginationInfo, user.UserListAutocompleteFilter) ([]user.User, uint64, error)
}

func (h *Handler) GetUserList(ctx context.Context, req *grpcPbV1.GetUserListRequest) (*grpcPbV1.GetUserListResponse, error) {
	var pagination common.PaginationInfo = common.PaginationInfo{}
	var filter user.UserListFilter = user.UserListFilter{}

	if req.PaginationInfo != nil {
		pagination = common.PaginationInfo{
			Offset: req.PaginationInfo.Offset,
			Limit:  req.PaginationInfo.Limit,
		}
	}

	if req.Filter != nil {
		filter.UserID = req.Filter.UserId
		filter.Email = req.Filter.Email
		filter.Status = constants.ACTIVE_STATUS(req.Filter.Status)
		filter.CreatedTimeFrom = req.Filter.CreatedTimeFrom
		filter.CreatedTimeTo = req.Filter.CreatedTimeTo
	}

	users, totalCount, err := h.userService.GetUserList(ctx, pagination, filter)

	if err != nil {
		return &grpcPbV1.GetUserListResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	var grpcUsers [](*grpcPbV1.User)

	for _, item := range users {
		grpcUsers = append(grpcUsers, convertUserToGrpcUser(item))
	}

	return &grpcPbV1.GetUserListResponse{
		Data: &grpcPbV1.GetUserListResponseData{
			Users:      grpcUsers,
			TotalCount: &totalCount,
		},
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}

func (h *Handler) GetUserDetails(ctx context.Context, req *grpcPbV1.GetUserDetailsRequest) (*grpcPbV1.GetUserDetailsResponse, error) {
	fetchUser, err := h.userService.GetUserDetails(ctx, req.UserId)

	if err != nil {
		return &grpcPbV1.GetUserDetailsResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	return &grpcPbV1.GetUserDetailsResponse{
		Data: &grpcPbV1.GetUserDetailsResponseData{
			User: convertUserToGrpcUser(fetchUser),
		},
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}

func (h *Handler) CreateUser(ctx context.Context, req *grpcPbV1.CreateUserRequest) (*grpcPbV1.CreateUserResponse, error) {
	if req.User == nil {
		return &grpcPbV1.CreateUserResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer("user must not be empty"),
		}, nil
	}

	newUser := user.User{
		Email:       req.User.Email,
		FirstName:   req.User.FirstName,
		LastName:    req.User.LastName,
		Status:      constants.ACTIVE_STATUS(req.User.Status),
		Password:    req.User.Password,
		NewSongNoti: req.User.NewSongNoti,
	}
	_, err := h.userService.CreateUser(ctx, newUser)

	if err != nil {
		return &grpcPbV1.CreateUserResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	return &grpcPbV1.CreateUserResponse{
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}

func (h *Handler) PutUser(ctx context.Context, req *grpcPbV1.PutUserRequest) (*grpcPbV1.PutUserResponse, error) {
	if req.User == nil {
		return &grpcPbV1.PutUserResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer("user must not be empty"),
		}, nil
	}

	curUser := user.User{
		UserID:      req.User.UserId,
		Email:       req.User.Email,
		FirstName:   req.User.FirstName,
		LastName:    req.User.LastName,
		Status:      constants.ACTIVE_STATUS(req.User.Status),
		NewSongNoti: req.User.NewSongNoti,
	}
	_, err := h.userService.PutUser(ctx, curUser)

	if err != nil {
		return &grpcPbV1.PutUserResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	return &grpcPbV1.PutUserResponse{
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}

func (h *Handler) UpdateUserStatus(ctx context.Context, req *grpcPbV1.UpdateUserStatusRequest) (*grpcPbV1.UpdateUserStatusResponse, error) {
	err := h.userService.UpdateUserStatus(ctx, req.UserId, constants.ACTIVE_STATUS(req.Status))

	if err != nil {
		return &grpcPbV1.UpdateUserStatusResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	return &grpcPbV1.UpdateUserStatusResponse{
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}

func (h *Handler) UpdateUserPermissions(ctx context.Context, req *grpcPbV1.UpdateUserPermissionsRequest) (*grpcPbV1.UpdateUserPermissionsResponse, error) {
	err := h.userService.UpdateUserPermissions(ctx, req.UserId, req.AddedPermissionIds, req.DeletedPermissionIds)

	if err != nil {
		return &grpcPbV1.UpdateUserPermissionsResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	return &grpcPbV1.UpdateUserPermissionsResponse{
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}

func (h *Handler) LogIn(ctx context.Context, req *grpcPbV1.LogInRequest) (*grpcPbV1.LogInResponse, error) {
	fetchedUser, err := h.userService.LogIn(ctx, user.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		return &grpcPbV1.LogInResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	token, err := h.authInterceptor.jwtService.GenerateToken(fetchedUser.UserID)

	if err != nil {
		return &grpcPbV1.LogInResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	return &grpcPbV1.LogInResponse{
		Token:    token,
		User:     convertUserToGrpcUser(fetchedUser),
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}

func (h *Handler) Authenticate(ctx context.Context, req *grpcPbV1.AuthenticateRequest) (*grpcPbV1.AuthenticateResponse, error) {
	doesExist, err := h.userService.DoesUserExist(ctx, req.UserId)

	if err != nil {
		return &grpcPbV1.AuthenticateResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	return &grpcPbV1.AuthenticateResponse{
		IsAuthenticated: &doesExist,
		Error:           nil,
		ErrorMsg:        nil,
	}, nil
}

func (h *Handler) GetUserListAutocomplete(ctx context.Context, req *grpcPbV1.GetUserListAutocompleteRequest) (*grpcPbV1.GetUserListAutocompleteResponse, error) {
	var pagination common.PaginationInfo = common.PaginationInfo{}
	var filter user.UserListAutocompleteFilter = user.UserListAutocompleteFilter{}

	if req.PaginationInfo != nil {
		pagination = common.PaginationInfo{
			Offset: req.PaginationInfo.Offset,
			Limit:  req.PaginationInfo.Limit,
		}
	}

	if req.Filter != nil {
		filter.Email = req.Filter.Email
	}

	users, totalCount, err := h.userService.GetUserListAutocomplete(ctx, pagination, filter)

	if err != nil {
		return &grpcPbV1.GetUserListAutocompleteResponse{
			Error:    common_utils.GetUInt32Pointer(1),
			ErrorMsg: common_utils.GetStringPointer(err.Error()),
		}, nil
	}

	var grpcUsers [](*grpcPbV1.User)

	for _, item := range users {
		grpcUsers = append(grpcUsers, convertUserToGrpcUserAutocomplete(item))
	}

	return &grpcPbV1.GetUserListAutocompleteResponse{
		Data: &grpcPbV1.GetUserListAutocompleteResponseData{
			Users:      grpcUsers,
			TotalCount: &totalCount,
		},
		Error:    nil,
		ErrorMsg: nil,
	}, nil
}
