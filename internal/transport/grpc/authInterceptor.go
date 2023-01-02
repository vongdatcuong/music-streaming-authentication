package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/jwtAuth"
	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/permission"
	common_utils "github.com/vongdatcuong/music-streaming-authentication/internal/modules/utils/common"
	grpcPbV1 "github.com/vongdatcuong/music-streaming-protos/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	jwtService            *jwtAuth.JwtService
	userService           UserServiceGrpc
	permissionService     PermissionServiceGrpc
	accessiblePermissions map[string][]string
}

func NewAuthInterceptor(jwtService *jwtAuth.JwtService, userService UserServiceGrpc, permissionService PermissionServiceGrpc) *AuthInterceptor {
	return &AuthInterceptor{jwtService: jwtService, userService: userService, permissionService: permissionService}
}

func (interceptor *AuthInterceptor) GrpcUnary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		injectedContext, err, _ := interceptor.authorize(ctx, req, md["authorization"], info.FullMethod, EndPointPermissions, EndPointNoAuthentication, EndPointPermissionFuncs)

		if err != nil {
			return getRespective403Response(info.FullMethod), nil
		}

		return handler(injectedContext, req)
	}
}

/*func (interceptor *AuthInterceptor) HttpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err, errCode := interceptor.authorize(r.Context(), r.Header["Authorization"], r.URL.Path, HttpEndPointPermissions)

		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, errCode, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}*/

func (interceptor *AuthInterceptor) authorize(ctx context.Context, req interface{}, authHeader []string,
	path string, permissionsMap map[string][]string, noAuthenMap map[string]bool,
	permissionFuncsMap map[string](func(jwtAuth.UserClaims, any) (bool, error))) (context.Context, error, uint32) {
	if noAuthenMap[path] {
		return ctx, nil, 0
	}
	accessToken, err := parseAuthorizationHeader(authHeader)
	if err != nil {
		return ctx, err, 1
	}

	claims, err := interceptor.jwtService.ValidateToken(accessToken)

	if err != nil {
		return ctx, err, 1
	}
	newCtx := interceptor.jwtService.InjectToken(ctx, accessToken)

	requiredPermFunc := permissionFuncsMap[path]

	if requiredPermFunc != nil {
		isAllowed, err := requiredPermFunc(*claims, req)

		if err != nil {
			return ctx, err, 1
		}

		// Success
		if isAllowed {
			return newCtx, nil, 0
		}

	}

	requiredPerm := permissionsMap[path]
	var firstRequiredPermName string

	if requiredPerm != nil && len(requiredPerm) > 0 {
		firstRequiredPermName = requiredPerm[0]
	}

	hasPerm, _, err := interceptor.permissionService.CheckUserPermission(ctx, claims.UserID, permission.Permission{
		// TODO: Check user has any permission in a list
		Name: firstRequiredPermName,
	})

	if err != nil {
		return newCtx, err, 1
	}

	if !hasPerm {
		return newCtx, fmt.Errorf("you have no permission to access this resource"), 403
	}

	// Success
	return newCtx, nil, 0
}

func parseAuthorizationHeader(values []string) (string, error) {
	if values == nil || len(values) == 0 {
		return "", fmt.Errorf("invalid authorization header")
	}
	authHeader := values[0]
	authHeaderParts := strings.Split(authHeader, " ")

	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("invalid authorization header")
	}

	return authHeaderParts[1], nil
}

// TODO: improve this
func getRespective403Response(path string) any {
	errCode, errMsg := common_utils.GetUInt32Pointer(403), common_utils.GetStringPointer("You have no permission to access this resource")

	if path == permissionServicePath+"/GetPermissionList" {
		return &grpcPbV1.GetPermissionListResponse{
			Error:    errCode,
			ErrorMsg: errMsg,
		}
	} else if path == permissionServicePath+"/CreatePermission" {
		return &grpcPbV1.CreatePermissionResponse{
			Error:    errCode,
			ErrorMsg: errMsg,
		}
	} else if path == permissionServicePath+"/PutPermission" {
		return &grpcPbV1.PutPermissionResponse{
			Error:    errCode,
			ErrorMsg: errMsg,
		}
	} else if path == userServicePath+"/GetUserList" {
		return &grpcPbV1.GetUserListResponse{
			Error:    errCode,
			ErrorMsg: errMsg,
		}
	} else if path == userServicePath+"/GetUserDetails" {
		return &grpcPbV1.GetUserDetailsResponse{
			Error:    errCode,
			ErrorMsg: errMsg,
		}
	} else if path == userServicePath+"/CreateUser" {
		return &grpcPbV1.CreateUserResponse{
			Error:    errCode,
			ErrorMsg: errMsg,
		}
	} else if path == userServicePath+"/PutUser" {
		return &grpcPbV1.PutUserResponse{
			Error:    errCode,
			ErrorMsg: errMsg,
		}
	} else if path == userServicePath+"/UpdateUserStatus" {
		return &grpcPbV1.UpdateUserStatusResponse{
			Error:    errCode,
			ErrorMsg: errMsg,
		}
	} else if path == userServicePath+"/UpdateUserPermissions" {
		return &grpcPbV1.UpdateUserPermissionsResponse{
			Error:    errCode,
			ErrorMsg: errMsg,
		}
	} else if path == userServicePath+"/Authenticate" {
		return &grpcPbV1.AuthenticateResponse{
			Error:    errCode,
			ErrorMsg: errMsg,
		}
	}
	return nil
}
