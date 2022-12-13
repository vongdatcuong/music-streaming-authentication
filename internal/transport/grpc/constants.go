package grpc

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
const authPath string = "/music_streaming.authentication.user.LogIn"

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
	userServicePath + "/PutUser":               {userPermissions["WRITE"]},
	userServicePath + "/UpdateUserStatus":      {userPermissions["WRITE"]},
	userServicePath + "/UpdateUserPermissions": {userPermissions["WRITE"]},
}

var EndPointNoAuthentication map[string]bool = map[string]bool{
	authPath + "/LogIn":             true,
	userServicePath + "/CreateUser": true,
}

// Http
const httpPath = "/api/gateway/v1"
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
