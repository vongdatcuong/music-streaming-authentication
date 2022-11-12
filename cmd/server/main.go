package main

import (
	"github.com/sirupsen/logrus"
	"github.com/vongdatcuong/music-streaming-authentication/internal/database"
)

func Run() error {
	db, err := database.NewDatabase()

	if err != nil {
		return err
	}

	// Ping DB
	if err := db.Client.DB.Ping(); err != nil {
		//return fmt.Errorf("could not ping the database: %w", err)
	}

	_, err = db.MigrateDB()

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
