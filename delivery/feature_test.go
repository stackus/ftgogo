package main

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	flag "github.com/spf13/pflag"

	"github.com/stackus/ftgogo/delivery/features/steps"
)

var opts = godog.Options{
	Format:   "progress",
	NoColors: true,
}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	state := steps.NewFeatureState()

	ctx.BeforeScenario(func(*godog.Scenario) {
		state.Reset()
	})

	state.RegisterCommonSteps(ctx)
	state.RegisterCancelDeliverySteps(ctx)
	state.RegisterCreateDeliverySteps(ctx)
	state.RegisterCreateRestaurantSteps(ctx)
	state.RegisterScheduleDeliverySteps(ctx)
	state.RegisterSetCourierAvailabilitySteps(ctx)
	state.RegisterGetDeliverySteps(ctx)
	state.RegisterGetCourierSteps(ctx)
}

func TestMain(m *testing.M) {
	flag.Parse()
	for _, arg := range os.Args[1:] {
		if arg == "-test.v=true" {
			opts.Format = "pretty"
			break
		}
	}

	status := godog.TestSuite{
		Name:                "delivery features",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	if st := m.Run(); st != 0 {
		os.Exit(st)
	}

	if status != 0 {
		os.Exit(status)
	}
}
