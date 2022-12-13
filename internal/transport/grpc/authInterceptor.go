package grpc

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/jwtAuth"
	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/permission"
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
		err, _ := interceptor.authorize(ctx, md["authorization"], info.FullMethod, EndPointPermissions, EndPointNoAuthentication)

		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (interceptor *AuthInterceptor) HttpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err, errCode := interceptor.authorize(r.Context(), r.Header["Authorization"], r.URL.Path, HttpEndPointPermissions, HttpEndPointNoAuthentication)

		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, errCode, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, authHeader []string, path string, permissionsMap map[string][]string, noAuthenMap map[string]bool) (error, uint32) {
	if noAuthenMap[path] {
		return nil, 0
	}

	accessToken, err := parseAuthorizationHeader(authHeader)

	if err != nil {
		return err, 1
	}

	claims, err := interceptor.jwtService.ValidateToken(accessToken)

	if err != nil {
		return err, 1
	}

	doesExist, err := interceptor.userService.DoesUserExist(ctx, claims.UserID)

	if err != nil {
		return err, 1
	}

	if !doesExist {
		return fmt.Errorf("invalid token"), 403
	}

	requiredPerm := permissionsMap[path]

	if requiredPerm != nil {
		hasPerm, err := interceptor.permissionService.CheckUserPermission(ctx, claims.UserID, permission.Permission{
			// TODO: Check user has any permission in a list
			Name: requiredPerm[0],
		})

		if err != nil {
			return fmt.Errorf("could not check user permission: %w", err), 1
		}

		if !hasPerm {
			return fmt.Errorf("you have no permission to access this resource"), 403
		}
	}

	return nil, 0
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
