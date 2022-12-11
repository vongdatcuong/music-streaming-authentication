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

var EndPointPermissions map[string][]string = map[string][]string{
	// Permission
	permissionServicePath + "/GetPermissionList": {permissionPermissions["READ"]},
	permissionServicePath + "/CreatePermission":  {permissionPermissions["WRITE"]},
	permissionServicePath + "/PutPermission":     {permissionPermissions["WRITE"]},
	// User
	userServicePath + "/GetUserList":           {userPermissions["READ"]},
	userServicePath + "/GetUserDetails":        {userPermissions["READ"]},
	userServicePath + "/CreateUser":            {userPermissions["WRITE"]},
	userServicePath + "/PutUser":               {userPermissions["WRITE"]},
	userServicePath + "/UpdateUserStatus":      {userPermissions["WRITE"]},
	userServicePath + "/UpdateUserPermissions": {userPermissions["WRITE"]},
}

// Http
const httpPath = "/api/gateway/v1"
const httpPermissionPath = httpPath + "/permission"
const httpUserPath = httpPath + "/user"

var HttpEndPointPermissions map[string][]string = map[string][]string{
	// Permission
	httpPermissionPath + "/list":                  {permissionPermissions["READ"]},
	httpPermissionPath + "/create_permission":     {permissionPermissions["WRITE"]},
	httpPermissionPath + "/put_permission":        {permissionPermissions["WRITE"]},
	httpPermissionPath + "/check_user_permission": {permissionPermissions["READ"]},
	// User
	httpUserPath + "/list":                   {userPermissions["READ"]},
	httpUserPath + "/details":                {userPermissions["READ"]},
	httpUserPath + "/create_user":            {userPermissions["WRITE"]},
	httpUserPath + "/put_user":               {userPermissions["WRITE"]},
	httpUserPath + "/update_user_status":     {userPermissions["WRITE"]},
	httpUserPath + "/update_user_permission": {userPermissions["WRITE"]},
}
