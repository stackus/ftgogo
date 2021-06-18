package accounting

import (
	"os"
	"testing"

	"github.com/cucumber/godog"

	"github.com/stackus/ftgogo/accounting/features/steps"
)

type featureTests struct{}

func (f featureTests) initScenario(ctx *godog.ScenarioContext) {
	state := steps.NewFeatureState()

	ctx.BeforeScenario(func(*godog.Scenario) {
		state.Reset()
	})

	state.RegisterCommonSteps(ctx)
	state.RegisterCreateAccountSteps(ctx)
	state.RegisterDisableAccountSteps(ctx)
	state.RegisterEnableAccountSteps(ctx)
	state.RegisterGetAccountSteps(ctx)
	state.RegisterAuthorizeOrderSteps(ctx)
}

func TestFeatures(t *testing.T) {
	format := "progress"
	for _, arg := range os.Args[1:] {
		if arg == "-test.v=true" { // go test transforms -v option
			format = "pretty"
			break
		}
	}

	opts := godog.Options{
		Format: format,
	}

	features := featureTests{}

	status := godog.TestSuite{
		Name:                "accounting features",
		ScenarioInitializer: features.initScenario,
		Options:             &opts,
	}.Run()

	if status != 0 {
		t.Fail()
	}
}
