package bdd

import (
	"flag"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"

	"knative.dev/eventing/test/pkg/injection"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "knative.dev/pkg/system/testing"
)

var opt = godog.Options{
	Output: colors.Colored(os.Stdout),
}

func TestMain(m *testing.M) {
	flag.Parse()
	ctx := injection.InjectionEnabled()

	if len(flag.Args()) > 0 {
		opt.Paths = flag.Args()
	} else {
		opt.Paths = []string{
			"./features/",
		}
	}

	opt.Format = "pretty"

	status := godog.RunWithOptions("CloudEvents", func(s *godog.Suite) {
		CloudEventsFeatureContext(ctx, s)
	}, opt)

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
