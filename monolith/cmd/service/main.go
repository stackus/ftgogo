package main

import (
	"github.com/stackus/ftgogo/accounting/acctmod"
	"github.com/stackus/ftgogo/consumer/consmod"
	"github.com/stackus/ftgogo/customer-web/cwebmod"
	"github.com/stackus/ftgogo/delivery/delvmod"
	"github.com/stackus/ftgogo/kitchen/kitcmod"
	"github.com/stackus/ftgogo/order-history/ohismod"
	"github.com/stackus/ftgogo/order/ordmod"
	"github.com/stackus/ftgogo/restaurant/restmod"
	"github.com/stackus/ftgogo/serviceapis"
	"shared-go/applications"
)

func main() {
	svc := applications.NewMonolith(initMonolith)
	if err := svc.Execute(); err != nil {
		panic(err)
	}
}

func initMonolith(svc *applications.Monolith) error {
	var err error

	serviceapis.RegisterTypes()

	err = acctmod.Setup(svc)
	if err != nil {
		return err
	}

	err = consmod.Setup(svc)
	if err != nil {
		return err
	}

	err = cwebmod.Setup(svc)
	if err != nil {
		return err
	}

	err = delvmod.Setup(svc)
	if err != nil {
		return err
	}

	err = kitcmod.Setup(svc)
	if err != nil {
		return err
	}

	err = ordmod.Setup(svc)
	if err != nil {
		return err
	}

	err = ohismod.Setup(svc)
	if err != nil {
		return err
	}

	err = restmod.Setup(svc)
	if err != nil {
		return err
	}

	return nil
}
