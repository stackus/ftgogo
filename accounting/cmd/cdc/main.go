package main

import (
	"shared-go/applications"
)

func main() {
	cdc := applications.NewCDC(func(*applications.CDC) error { return nil })
	if err := cdc.Execute(); err != nil {
		panic(err)
	}
}
