package main

import "github.com/sirupsen/logrus"

func Run() error {
	return nil
}

func main() {
	if err := Run(); err != nil {
		logrus.Errorln(err)
	}
}
