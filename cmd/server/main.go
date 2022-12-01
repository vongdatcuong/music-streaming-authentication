package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/vongdatcuong/music-streaming-authentication/internal/database"
	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/permission"
	"github.com/vongdatcuong/music-streaming-authentication/internal/modules/user"
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

	//res, _ := db.CheckUserPermission(context.Background(), 1, permission.Permission{Name: "musc_streaming.song.read"})
	//res2, _ := db.GetPermissionList(context.Background())
	//logrus.Info(res)

	permissionService := permission.NewService(db)
	userService := user.NewService(db)
	grpcHandler := grpcTransport.NewHandler(permissionService, userService)

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
