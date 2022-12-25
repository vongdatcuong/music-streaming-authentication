package grpc

import (
	"fmt"
	"net"
	"os"

	grpcPbV1 "github.com/vongdatcuong/music-streaming-protos/go/v1"
	"google.golang.org/grpc"
)

type Handler struct {
	grpcPbV1.UnimplementedPermissionServiceServer
	grpcPbV1.UnimplementedUserServiceServer
	permissionService PermissionServiceGrpc
	userService       UserServiceGrpc
	authInterceptor   *AuthInterceptor
}

func NewHandler(permissionService PermissionServiceGrpc, userService UserServiceGrpc, authInterceptor *AuthInterceptor) *Handler {
	h := &Handler{permissionService: permissionService, userService: userService, authInterceptor: authInterceptor}

	return h
}

func (h *Handler) RunGrpcServer(port string, channel chan error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		channel <- fmt.Errorf("could not listen on port %s: %w", port, err)
		return
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(h.authInterceptor.GrpcUnary()))
	grpcPbV1.RegisterPermissionServiceServer(grpcServer, h)
	grpcPbV1.RegisterUserServiceServer(grpcServer, h)

	if err := grpcServer.Serve(lis); err != nil {
		channel <- fmt.Errorf("could not server Grpc server on port %s: %w", port, err)
	}
}

/*func (h *Handler) RunRestServer(port string, channel chan error) {
	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true, // Rest Server to return the same fields as protobuf
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)
	muxCtx, cancelMuxCtx := context.WithCancel(context.Background())
	defer cancelMuxCtx()
	err := grpcPbV1.RegisterPermissionServiceHandlerServer(muxCtx, gwmux, h)

	if err != nil {
		channel <- fmt.Errorf("Failed to register Permission Rest endpoints: %w", err)
		return
	}

	err = grpcPbV1.RegisterUserServiceHandlerServer(muxCtx, gwmux, h)

	if err != nil {
		channel <- fmt.Errorf("Failed to register User Rest endpoints: %w", err)
		return
	}

	restLis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		channel <- fmt.Errorf("could not listen on port %s: %w", port, err)
		return
	}
	httpMux := http.NewServeMux()
	httpMux.Handle("/", h.authInterceptor.HttpMiddleware(gwmux))

	if err := http.Serve(restLis, httpMux); err != nil {
		channel <- fmt.Errorf("could not serve Rest server on port %s: %w", port, err)
	}
}*/

func (h *Handler) Server() error {
	grpcChannel := make(chan error)
	go h.RunGrpcServer(os.Getenv("GRPC_PORT"), grpcChannel)
	//go h.RunRestServer(os.Getenv("REST_PORT"), restChannel)

	select {
	case grpcError := <-grpcChannel:
		return grpcError
		//case restError := <-restChannel:
		//	return restError
	}
}
