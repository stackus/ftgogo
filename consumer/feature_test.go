package accounting

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	flag "github.com/spf13/pflag"

	"github.com/stackus/ftgogo/consumer/features/steps"
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
	state.RegisterRegisterConsumerSteps(ctx)
	state.RegisterUpdateConsumerSteps(ctx)
	state.RegisterValidateOrderByConsumerSteps(ctx)
	state.RegisterAddAddressSteps(ctx)
	state.RegisterUpdateAddressSteps(ctx)
	state.RegisterRemoveAddressSteps(ctx)
	state.RegisterGetAddressSteps(ctx)
	state.RegisterGetConsumerSteps(ctx)
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
		Name:                "consumer features",
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
