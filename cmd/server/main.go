package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vongdatcuong/music-streaming-authentication/internal/database"
	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/jwtAuth"
	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/permission"
	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/user"
	"github.com/vongdatcuong/music-streaming-authentication/internal/transport/grpc"
	grpcTransport "github.com/vongdatcuong/music-streaming-authentication/internal/transport/grpc"
)

func Run() error {
	db, err := database.NewDatabase()

	if err != nil {
		return err
	}

	// Ping DB
	if err := db.Client.DB.Ping(); err != nil {
		return fmt.Errorf("could not ping the database: %w", err)
	}

	_, err = db.MigrateDB()

	if err != nil {
		return err
	}

	userService := user.NewService(db)
	permissionService := permission.NewService(db, userService)
	jwtAuthService := jwtAuth.NewService(os.Getenv("JWT_SECRET_KEY"), 6*time.Hour)
	authInterceptor := grpc.NewAuthInterceptor(jwtAuthService, userService, permissionService)
	grpcHandler := grpcTransport.NewHandler(permissionService, userService, authInterceptor)

	if err := grpcHandler.Server(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		logrus.Errorln(err)
	}
}
