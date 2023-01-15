package grpc

import (
	"fmt"

	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/jwtAuth"
	grpcPbV1 "github.com/vongdatcuong/music-streaming-protos/go/v1"
)

// Permissions
const permissionPermPrefix = "music_streaming.permission"

var permissionPermissions map[string]string = map[string]string{
	"READ":  permissionPermPrefix + ".read",
	"WRITE": permissionPermPrefix + ".write",
}

const userPermPrefix = "music_streaming.user"

var userPermissions map[string]string = map[string]string{
	"READ":  userPermPrefix + ".read",
	"WRITE": userPermPrefix + ".write",
}

// Endpoints
const permissionServicePath string = "/music_streaming.authentication.permission.PermissionService"
const userServicePath string = "/music_streaming.authentication.user.UserService"
const authPath string = "/music_streaming.authentication.user.UserService"

var EndPointPermissions map[string][]string = map[string][]string{
	// Permission
	permissionServicePath + "/GetPermissionList": {permissionPermissions["READ"]},
	permissionServicePath + "/CreatePermission":  {permissionPermissions["WRITE"]},
	permissionServicePath + "/PutPermission":     {permissionPermissions["WRITE"]},
	// User
	userServicePath + "/GetUserList":    {userPermissions["READ"]},
	userServicePath + "/GetUserDetails": {userPermissions["READ"]},
	userServicePath + "/CreateUser":     {userPermissions["WRITE"]},
	// TODO: Admin can update all users, but a particular user can only update his/her own profile
	userServicePath + "/PutUser":                 {userPermissions["WRITE"]},
	userServicePath + "/UpdateUserStatus":        {userPermissions["WRITE"]},
	userServicePath + "/UpdateUserPermissions":   {userPermissions["WRITE"]},
	userServicePath + "/GetUserListAutocomplete": {userPermissions["READ"]},
}

var EndPointNoAuthentication map[string]bool = map[string]bool{
	authPath + "/LogIn":               true,
	userServicePath + "/CreateUser":   true,
	userServicePath + "/Authenticate": true,
}

var EndPointPermissionFuncs map[string](func(jwtAuth.UserClaims, any) (bool, error)) = map[string](func(jwtAuth.UserClaims, any) (bool, error)){
	userServicePath + "/GetUserDetails":        CheckGetUserDetails,
	userServicePath + "/PutUser":               CheckPutUser,
	userServicePath + "/UpdateUserStatus":      CheckUpdateUserStatus,
	userServicePath + "/UpdateUserPermissions": CheckUpdateUserPermissions,
}

func CheckGetUserDetails(userClaims jwtAuth.UserClaims, req any) (bool, error) {
	assertedReq, ok := req.(*grpcPbV1.GetUserDetailsRequest)
	if !ok {
		return false, fmt.Errorf("want type *GetUserDetailsRequest;  got %T", req)
	}

	if userClaims.UserID == assertedReq.UserId {
		return true, nil
	}

	return false, nil
}

func CheckPutUser(userClaims jwtAuth.UserClaims, req any) (bool, error) {
	assertedReq, ok := req.(*grpcPbV1.PutUserRequest)
	if !ok {
		return false, fmt.Errorf("want type *PutUserRequest;  got %T", req)
	}

	if assertedReq.User != nil && userClaims.UserID == assertedReq.User.UserId {
		return true, nil
	}

	return false, nil
}

func CheckUpdateUserStatus(userClaims jwtAuth.UserClaims, req any) (bool, error) {
	assertedReq, ok := req.(*grpcPbV1.UpdateUserStatusRequest)
	if !ok {
		return false, fmt.Errorf("want type *UpdateUserStatusRequest;  got %T", req)
	}

	if userClaims.UserID == assertedReq.UserId {
		return true, nil
	}

	return false, nil
}

func CheckUpdateUserPermissions(userClaims jwtAuth.UserClaims, req any) (bool, error) {
	assertedReq, ok := req.(*grpcPbV1.UpdateUserPermissionsRequest)
	if !ok {
		return false, fmt.Errorf("want type *UpdateUserPermissionsRequest;  got %T", req)
	}

	if userClaims.UserID == assertedReq.UserId {
		return true, nil
	}

	return false, nil
}

// Http
/*const httpPath = "/api/gateway/v1"
const httpPermissionPath = httpPath + "/permission"
const httpUserPath = httpPath + "/user"
const httpAuthPath = httpPath + "/auth"

var HttpEndPointPermissions map[string][]string = map[string][]string{
	// Permission
	httpPermissionPath + "/list":                  {permissionPermissions["READ"]},
	httpPermissionPath + "/create_permission":     {permissionPermissions["WRITE"]},
	httpPermissionPath + "/put_permission":        {permissionPermissions["WRITE"]},
	httpPermissionPath + "/check_user_permission": {permissionPermissions["READ"]},
	// User
	httpUserPath + "/list":                   {userPermissions["READ"]},
	httpUserPath + "/details":                {userPermissions["READ"]},
	httpUserPath + "/put_user":               {userPermissions["WRITE"]},
	httpUserPath + "/update_user_status":     {userPermissions["WRITE"]},
	httpUserPath + "/update_user_permission": {userPermissions["WRITE"]},
}

var HttpEndPointNoAuthentication map[string]bool = map[string]bool{
	httpAuthPath + "/login":       true,
	httpUserPath + "/create_user": true,
}
*/
