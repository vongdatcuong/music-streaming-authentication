package main

import (
	"github.com/sirupsen/logrus"
	"github.com/vongdatcuong/music-streaming-authentication/internal/database"
)

func Run() error {
	_, err := database.NewDatabase()

	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		logrus.Errorln(err)
	}
}
