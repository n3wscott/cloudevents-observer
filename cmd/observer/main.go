package main

import (
	"flag"
	"knative.dev/eventing/test/pkg/injection"
	"os"

	"knative.dev/eventing/test/pkg/eventlog/cloudout"
	"knative.dev/eventing/test/pkg/eventlog/recorder"
	"knative.dev/eventing/test/pkg/eventlog/writer"
	"knative.dev/eventing/test/pkg/observer"
)

func main() {
	flag.Parse()
	ctx := injection.InjectionEnabled()

	obs := observer.New(
		writer.NewEventLog(ctx, os.Stdout),
		recorder.NewFromEnv(ctx),
		cloudout.NewFromEnv(ctx),
	)

	if err := obs.Start(ctx); err != nil {
		panic(err)
	}
}
