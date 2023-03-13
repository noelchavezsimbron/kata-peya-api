//go:build e2e
// +build e2e

package pets

import (
	"context"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var (
	feat = &feature{}
)

func TestFeatures(t *testing.T) {

	status := godog.TestSuite{
		Name:                 "pets",
		TestSuiteInitializer: InitializeSuite,
		ScenarioInitializer:  InitializeScenario,
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"features"},
			Output: colors.Colored(os.Stdout),
		},
	}.Run()

	if status == 2 {
		t.SkipNow()
	}

	if status != 0 {
		t.Fatal("test suite cannot run successfully")
	}
}

func InitializeSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		feat.StartSuite()
	})

	ctx.AfterSuite(func() {
		feat.Teardown()
	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		feat.Reset()
		return ctx, nil
	})

	ctx.Step(`^I send "([^"]*)" request to "([^"]*)"$`, feat.iSendRequestTo)
	ctx.Step(`^the are registered pets$`, feat.theAreRegisteredPets)
	ctx.Step(`^the response should match json:$`, feat.theResponseShouldMatchJson)
	ctx.Step(`^the response status code should be (\d+)$`, feat.theResponseStatusCodeShouldBe)
}
