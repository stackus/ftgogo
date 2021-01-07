package main

import (
	"shared-go/applications"
)

func main() {
	svc := applications.NewService(initApplication)
	if err := svc.Execute(); err != nil {
		panic(err)
	}
}
