package main

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	flag "github.com/spf13/pflag"

	"github.com/stackus/ftgogo/restaurant/features/steps"
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
	state.RegisterCreateRestaurantSteps(ctx)
}

func TestFeatures(t *testing.T) {
	flag.Parse()
	for _, arg := range os.Args[1:] {
		if arg == "-test.v=true" {
			opts.Format = "pretty"
			break
		}
	}

	status := godog.TestSuite{
		Name:                "restaurant features",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	if status != 0 {
		t.Fail()
	}
}
